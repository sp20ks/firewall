package usecase

import (
	"net/url"
	"regexp"
	"rules-engine/internal/entity"
	"rules-engine/internal/repository"
	"strings"
)

type AnalizerUseCase struct {
	ruleRepo repository.RuleRepository
}

func NewAnalizerUseCase(ruleRepo repository.RuleRepository) *AnalizerUseCase {
	return &AnalizerUseCase{ruleRepo: ruleRepo}
}

// TODO: добавить ip листы сюда
func (a *AnalizerUseCase) AnalyzeRequest(request *entity.Request) (*entity.ScanResult, error) {
	rules, err := a.ruleRepo.GetRulesByURL(extractPath(request.URL), request.Method)
	if err != nil {
	}

	result := &entity.ScanResult{
		Action: entity.ActionAllow,
		Reason: "Request passed all checks.",
	}

	for _, rule := range rules {
		if rule.IsActive == nil || !*rule.IsActive {
			continue
		}

		var tempResult *entity.ScanResult
		switch rule.AttackType {
		case "xss":
			tempResult = a.applyXSSRule(request, &rule)
		case "csrf":
			tempResult = a.applyCSRFRule(request, &rule)
		case "sqli":
			tempResult = a.applySQLIRule(request, &rule)
		default:
			continue
		}

		if tempResult != nil && tempResult.Action == "block" {
			return tempResult, nil
		}

		if tempResult != nil {
			result = tempResult
		}
	}

	return result, nil
}

func (a *AnalizerUseCase) applyXSSRule(request *entity.Request, rule *entity.Rule) *entity.ScanResult {
	xssDetected := false
	modifiedBody := request.Body
	modifiedURL := request.URL
	if matched, _ := regexp.MatchString("(?i)<script.*?>.*?</script>", request.Body); matched {
		xssDetected = true
		switch rule.ActionType {
		case entity.ActionSanitize:
			modifiedBody = sanitizeXSS(request.Body)
		case entity.ActionEscape:
			modifiedBody = escapeHTML(request.Body)
		}
	}

	decodedURL := decodeURL(request.URL)
	if matched, _ := regexp.MatchString("(?i)<script.*?>.*?</script>", decodedURL); matched {
		xssDetected = true
		switch rule.ActionType {
		case entity.ActionSanitize:
			modifiedURL = sanitizeXSS(decodedURL)
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
			if modifiedBody != request.Body {
				result.ModifiedBody = modifiedBody
			}
			if modifiedURL != request.URL {
				result.ModifiedURL = modifiedURL
			}
		}

		return result
	}

	return nil
}

func (a *AnalizerUseCase) applyCSRFRule(request *entity.Request, rule *entity.Rule) *entity.ScanResult {
	if request.Headers["X-CSRF-Token"] == "" {
		return &entity.ScanResult{
			Action: entity.ActionBlock,
			Reason: "Missing CSRF token.",
		}
	}
	return nil
}

func (a *AnalizerUseCase) applySQLIRule(request *entity.Request, rule *entity.Rule) *entity.ScanResult {
	sqlDetected := false
	modifiedBody := request.Body
	modifiedURL := request.URL

	pattern := regexp.MustCompile(`(?i)'?\s*OR\s+1=1\s*--`)

	if pattern.MatchString(request.Body) {
		sqlDetected = true
		switch rule.ActionType {
		case entity.ActionSanitize:
			modifiedBody = pattern.ReplaceAllString(request.Body, "")
		case entity.ActionEscape:
			modifiedBody = escapeSQL(request.Body)
		}
	}

	if pattern.MatchString(request.URL) {
		sqlDetected = true
		switch rule.ActionType {
		case entity.ActionSanitize:
			modifiedURL = pattern.ReplaceAllString(request.URL, "")
		case entity.ActionEscape:
			modifiedURL = escapeSQL(request.URL)
		}
	}

	if sqlDetected {
		result := &entity.ScanResult{
			Action: rule.ActionType,
			Reason: "SQL injection detected in request.",
		}

		if rule.ActionType != entity.ActionBlock {
			if modifiedBody != request.Body {
				result.ModifiedBody = modifiedBody
			}
			if modifiedURL != request.URL {
				result.ModifiedURL = modifiedURL
			}
		}

		return result
	}

	return nil
}

func sanitizeXSS(input string) string {
	re := regexp.MustCompile("(?i)<script.*?>.*?</script>")
	return re.ReplaceAllString(input, "")
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
