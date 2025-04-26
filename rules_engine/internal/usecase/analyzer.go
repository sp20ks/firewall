package usecase

import (
	"fmt"
	"net"
	"net/url"
	"regexp"
	"rules-engine/internal/entity"
	"rules-engine/internal/repository"
	"strings"
)

const sqlInjectionPattern = `(?i)(\b(select|insert|update|delete|drop|union|join|cast|create|alter|truncate|grant|revoke|nullif|execute)\b[\s\S]*?['";\-\+=])|(\b(or|and)\b\s+('[^']*'|\d+)\s*=\s*('[^']*'|\d+))|(--|#)|(;[\s]*(select|insert|update|delete|drop|create|alter|truncate))|(%27|%2D%2D|%23)`
const xssPattern = `(?i)<script.*?>.*?</script>`

type AnalyzerUseCase struct {
	ruleRepo   repository.RuleRepository
	ipListRepo repository.IPListRepository

	sqlPattern *regexp.Regexp
	xssPattern *regexp.Regexp
}

func NewAnalyzerUseCase(ruleRepo repository.RuleRepository, ipListRepo repository.IPListRepository) *AnalyzerUseCase {
	sqlRegex := regexp.MustCompile(sqlInjectionPattern)
	xssRegex := regexp.MustCompile(xssPattern)

	return &AnalyzerUseCase{
		ruleRepo:   ruleRepo,
		ipListRepo: ipListRepo,
		sqlPattern: sqlRegex,
		xssPattern: xssRegex,
	}
}

func (a *AnalyzerUseCase) AnalyzeRequest(request *entity.Request) (*entity.ScanResult, error) {
	result, err := a.applyIPLists(request)
	if err != nil {
		return nil, err
	}

	// TODO: приоритеты у правил??
	if result.Action == entity.ActionBlock {
		return result, nil
	}

	result, err = a.applyRules(request)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func (a *AnalyzerUseCase) applyRules(request *entity.Request) (*entity.ScanResult, error) {
	rules, err := a.ruleRepo.GetRulesByURL(extractPath(request.URL), request.Method)
	if err != nil {
		return nil, fmt.Errorf("error while loading rules for resource")
	}

	result := &entity.ScanResult{
		Action:       entity.ActionAllow,
		Reason:       "Request passed all checks.",
		ModifiedURL:  request.URL,
		ModifiedBody: request.Body,
	}

	for _, rule := range rules {
		if rule.IsActive == nil || !*rule.IsActive {
			continue
		}

		// передаем в apply функции модифицированные url и body, чтобы в случае нескольких правил с sanitize или escape применились все действия
		var tempResult *entity.ScanResult
		switch rule.AttackType {
		case "xss":
			tempResult = a.applyXSSRule(result.ModifiedURL, result.ModifiedBody, rule)
		case "csrf":
			tempResult = a.applyCSRFRule(request, rule)
		case "sqli":
			tempResult = a.applySQLIRule(result.ModifiedURL, result.ModifiedBody, rule)
		default:
			continue
		}

		if tempResult != nil && tempResult.Action == "block" {
			return tempResult, nil
		}

		if tempResult != nil {
			result.ModifiedBody = tempResult.ModifiedBody
			result.ModifiedURL = tempResult.ModifiedURL
		}
	}

	// важно: мы отдаем только allow или block. не даем информацию о escape или sanitize
	return result, nil
}

func (a *AnalyzerUseCase) applyXSSRule(url, body string, rule entity.Rule) *entity.ScanResult {
	xssDetected := false
	modifiedBody := body
	modifiedURL := url
	if a.xssPattern.MatchString(body) {
		xssDetected = true
		switch rule.ActionType {
		case entity.ActionSanitize:
			modifiedBody = a.sqlPattern.ReplaceAllString(body, "")
		case entity.ActionEscape:
			modifiedBody = escapeHTML(body)
		}
	}

	decodedURL := decodeURL(url)
	if a.xssPattern.MatchString(decodedURL) {
		xssDetected = true
		switch rule.ActionType {
		case entity.ActionSanitize:
			modifiedURL = a.sqlPattern.ReplaceAllString(decodedURL, "")
		case entity.ActionEscape:
			modifiedURL = escapeHTML(decodedURL)
		}
	}

	if xssDetected {
		result := &entity.ScanResult{
			Action: rule.ActionType,
			Reason: "XSS detected in request.",
		}

		if rule.ActionType != entity.ActionBlock {
			if modifiedBody != body {
				result.ModifiedBody = modifiedBody
			}
			if modifiedURL != url {
				result.ModifiedURL = modifiedURL
			}
		}

		return result
	}

	return nil
}

func (a *AnalyzerUseCase) applyCSRFRule(request *entity.Request, rule entity.Rule) *entity.ScanResult {
	if request.Headers["X-Csrf-Token"] == "" {
		return &entity.ScanResult{
			Action: entity.ActionBlock,
			Reason: "Missing CSRF token.",
		}
	}
	return nil
}

func (a *AnalyzerUseCase) applySQLIRule(url, body string, rule entity.Rule) *entity.ScanResult {
	sqlDetected := false
	modifiedBody := body
	modifiedURL := url

	if a.sqlPattern.MatchString(body) {
		sqlDetected = true
		switch rule.ActionType {
		case entity.ActionSanitize:
			modifiedBody = a.sqlPattern.ReplaceAllString(body, "")
		case entity.ActionEscape:
			modifiedBody = escapeSQL(body)
		}
	}

	decodedURL := decodeURL(url)
	if a.sqlPattern.MatchString(decodedURL) {
		sqlDetected = true
		switch rule.ActionType {
		case entity.ActionSanitize:
			modifiedURL = a.sqlPattern.ReplaceAllString(decodedURL, "")
		case entity.ActionEscape:
			modifiedURL = escapeSQL(decodedURL)
		}
	}

	if sqlDetected {
		result := &entity.ScanResult{
			Action: rule.ActionType,
			Reason: "SQL injection detected in request.",
		}

		if rule.ActionType != entity.ActionBlock {
			if modifiedBody != body {
				result.ModifiedBody = modifiedBody
			}
			if modifiedURL != url {
				result.ModifiedURL = modifiedURL
			}
		}

		return result
	}

	return nil
}

func escapeHTML(input string) string {
	replacer := strings.NewReplacer(
		"<", "&lt;",
		">", "&gt;",
		"\"", "&quot;",
		"'", "&#39;",
		"&", "&amp;",
	)
	return replacer.Replace(input)
}

func escapeSQL(input string) string {
	return strings.ReplaceAll(input, "'", "''")
}

func extractPath(fullURL string) string {
	fullURL = strings.TrimSpace(fullURL)

	if idx := strings.IndexAny(fullURL, "?&"); idx != -1 {
		return fullURL[:idx]
	}

	return fullURL
}

func decodeURL(raw string) string {
	decoded, err := url.QueryUnescape(raw)
	if err != nil {
		return raw
	}
	return decoded
}

func (a *AnalyzerUseCase) applyIPLists(request *entity.Request) (*entity.ScanResult, error) {
	lists, err := a.ipListRepo.GetIPListsByURL(extractPath(request.URL), request.Method)
	if err != nil {
		return nil, fmt.Errorf("error while loading ip lists for resource")
	}

	var result *entity.ScanResult

	for _, list := range lists {
		result = a.checkIP(request, list)
		if result.Action == entity.ActionBlock {
			return result, nil
		}
	}
	return &entity.ScanResult{
		Action: entity.ActionAllow,
		Reason: "Requester IP was detected in ip list.",
	}, nil
}

func (a *AnalyzerUseCase) checkIP(request *entity.Request, iPList entity.IPList) *entity.ScanResult {
	ip := net.ParseIP(strings.Split(request.IP, ":")[0])
	if ip == nil {
		return &entity.ScanResult{
			Action: entity.ActionBlock,
			Reason: "Invalid IP address format.",
		}
	}

	result := &entity.ScanResult{
		Action: entity.ActionAllow,
		Reason: "Requester IP was detected in ip list.",
	}

	if iPList.IP.Contains(ip) && iPList.ListType == "blacklist" || !iPList.IP.Contains(ip) {
		result.Action = entity.ActionBlock
	}

	return result
}
