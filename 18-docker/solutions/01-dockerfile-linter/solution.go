package dockerlinter

import (
	"bufio"
	"strings"
)

type LintViolation struct {
	RuleID  string
	Message string
}

const (
	RuleMultiStage      = "MULTI_STAGE_REQUIRED"
	RuleCGODisabled     = "CGO_DISABLED_REQUIRED"
	RuleDependencyCache = "DEPENDENCY_CACHE_ORDER"
	RuleNonRootUser     = "NON_ROOT_USER_REQUIRED"
)

func LintDockerfile(content string) ([]LintViolation, error) {
	var violations []LintViolation

	scanner := bufio.NewScanner(strings.NewReader(content))
	fromCount := 0
	hasCGO := false
	hasNonRoot := false

	copyModLine := -1
	copyAllLine := -1
	currentLine := 0

	for scanner.Scan() {
		currentLine++
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		upper := strings.ToUpper(line)

		if strings.HasPrefix(upper, "FROM ") {
			fromCount++
		}

		if strings.Contains(line, "CGO_ENABLED=0") {
			hasCGO = true
		}

		if strings.HasPrefix(upper, "USER ") {
			userVal := strings.TrimSpace(line[5:])
			if userVal != "root" && userVal != "0" {
				hasNonRoot = true
			}
		}

		if strings.HasPrefix(upper, "COPY ") {
			if strings.Contains(line, "go.mod") {
				if copyModLine == -1 {
					copyModLine = currentLine
				}
			}
			if strings.Contains(line, ". .") || strings.Contains(line, "./ ./") {
				if copyAllLine == -1 {
					copyAllLine = currentLine
				}
			}
		}
	}

	if fromCount < 2 {
		violations = append(violations, LintViolation{
			RuleID:  RuleMultiStage,
			Message: "Dockerfile must use multi-stage builds (at least 2 FROM instructions)",
		})
	}

	if !hasCGO {
		violations = append(violations, LintViolation{
			RuleID:  RuleCGODisabled,
			Message: "Static compilation requires CGO_ENABLED=0",
		})
	}

	if copyAllLine != -1 && copyModLine != -1 && copyModLine > copyAllLine {
		violations = append(violations, LintViolation{
			RuleID:  RuleDependencyCache,
			Message: "COPY go.mod must precede COPY . . to leverage layer caching",
		})
	}

	if !hasNonRoot {
		violations = append(violations, LintViolation{
			RuleID:  RuleNonRootUser,
			Message: "Container must run as a non-root user (USER <non-root>)",
		})
	}

	return violations, nil
}
