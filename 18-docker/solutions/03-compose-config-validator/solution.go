package composevalidator

import (
	"bufio"
	"fmt"
	"strings"
)

type ComposeViolation struct {
	RuleID  string
	Service string
	Message string
}

const (
	RuleMissingService      = "MISSING_REQUIRED_SERVICE"
	RuleHardcodedSecret     = "HARDCODED_SECRET_DETECTED"
	RuleDatabaseHealthcheck = "DATABASE_HEALTHCHECK_REQUIRED"
)

func isDatabaseService(name string) bool {
	lower := strings.ToLower(name)
	dbKeywords := []string{"postgres", "mysql", "mariadb", "mongo", "redis", "db"}
	for _, kw := range dbKeywords {
		if strings.Contains(lower, kw) {
			return true
		}
	}
	return false
}

type serviceMeta struct {
	hasHealthcheck bool
}

func ValidateCompose(content string, requiredServices []string) ([]ComposeViolation, error) {
	var violations []ComposeViolation

	servicesFound := make(map[string]*serviceMeta)
	scanner := bufio.NewScanner(strings.NewReader(content))

	inServices := false
	var currentService string

	for scanner.Scan() {
		line := scanner.Text()
		trimmed := strings.TrimSpace(line)
		if trimmed == "" || strings.HasPrefix(trimmed, "#") {
			continue
		}

		// Calculate indentation
		indent := len(line) - len(strings.TrimLeft(line, " "))

		if indent == 0 && strings.HasPrefix(trimmed, "services:") {
			inServices = true
			continue
		}

		if indent == 0 && inServices && !strings.HasPrefix(trimmed, "services:") {
			// Left services block (e.g. volumes: or networks:)
			inServices = false
			currentService = ""
		}

		if inServices && indent == 2 && strings.HasSuffix(trimmed, ":") {
			// Found service declaration
			currentService = strings.TrimSuffix(trimmed, ":")
			servicesFound[currentService] = &serviceMeta{}
			continue
		}

		if inServices && currentService != "" {
			if strings.HasPrefix(trimmed, "healthcheck:") {
				servicesFound[currentService].hasHealthcheck = true
			}

			// Check environment secrets
			lower := strings.ToLower(trimmed)
			if (strings.Contains(lower, "password") || strings.Contains(lower, "secret")) && strings.Contains(trimmed, ":") {
				parts := strings.SplitN(trimmed, ":", 2)
				if len(parts) == 2 {
					val := strings.TrimSpace(parts[1])
					if val != "" && !strings.HasPrefix(val, "${") {
						violations = append(violations, ComposeViolation{
							RuleID:  RuleHardcodedSecret,
							Service: currentService,
							Message: fmt.Sprintf("Environment secret %s in service %s must use variable interpolation (${VAR})", strings.TrimSpace(parts[0]), currentService),
						})
					}
				}
			}
		}
	}

	// 1. Check required services
	for _, reqSvc := range requiredServices {
		if _, found := servicesFound[reqSvc]; !found {
			violations = append(violations, ComposeViolation{
				RuleID:  RuleMissingService,
				Service: reqSvc,
				Message: fmt.Sprintf("Required service '%s' is missing from docker-compose.yml", reqSvc),
			})
		}
	}

	// 2. Check database healthchecks
	for svcName, meta := range servicesFound {
		if isDatabaseService(svcName) && !meta.hasHealthcheck {
			violations = append(violations, ComposeViolation{
				RuleID:  RuleDatabaseHealthcheck,
				Service: svcName,
				Message: fmt.Sprintf("Database service '%s' must declare a healthcheck block", svcName),
			})
		}
	}

	return violations, nil
}
