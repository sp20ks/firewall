package usecase

import (
	"fmt"
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

func (a *AnalizerUseCase) AnalyzeRequest(request *entity.Request) (*entity.ScanResult, error) {
	rules, err := a.ruleRepo.GetRulesByURL(extractPath(request.URL), request.Method)
	if err != nil {
		return nil, fmt.Errorf("error fetching rules for resource %s %s: %w", request.Method, extractPath(request.URL), err)
	}

	for _, rule := range rules {
		if rule.IsActive == nil || !*rule.IsActive {
			continue
		}

		var result *entity.ScanResult
		switch rule.AttackType {
		case "xss":
			result = a.applyXSSRule(request, &rule)
		case "csrf":
			result = a.applyCSRFRule(request, &rule)
		case "sqli":
			result = a.applySQLIRule(request, &rule)
		default:
			continue
		}

		if result != nil {
			return result, nil
		}
	}

	return &entity.ScanResult{
		Action: entity.ActionAllow,
		Reason: "Request passed all checks.",
	}, nil
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

	if matched, _ := regexp.MatchString("(?i)<script.*?>.*?</script>", request.URL); matched {
		xssDetected = true
		switch rule.ActionType {
		case entity.ActionSanitize:
			modifiedURL = sanitizeXSS(request.URL)
		case entity.ActionEscape:
			modifiedURL = escapeHTML(request.URL)
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
