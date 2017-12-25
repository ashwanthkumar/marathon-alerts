package checks

import (
	"testing"

	"github.com/gambol99/go-marathon"
	"github.com/stretchr/testify/assert"
)

// === AppHealth ===
func TestAppHealthGetHosts(t *testing.T) {
	check := AppHealth{}
	appLabels := make(map[string]string)
	appLabels["router.hosts"] = "zeus.prod.indix.tv"
	s := marathon.HealthCheck{
		GracePeriodSeconds:     180,
		IntervalSeconds:        90,
		TimeoutSeconds:         30,
		MaxConsecutiveFailures: 3,
		PortIndex:              0,
		Path:                   "/health",
		Protocol:               "HTTP",
	}
	var appHC []*marathon.HealthCheck
	appHC = append(appHC, &s)

	app := marathon.Application{
		ID:           "/production.zeus",
		Instances:    10,
		TasksHealthy: 5,
		Labels:       appLabels,
		HealthChecks: appHC,
	}

	appCheck := check.Check(app)
	assert.Equal(t, Pass, appCheck.Result)
	assert.Equal(t, "apphealth", appCheck.CheckName)
	assert.Equal(t, "/production.zeus", appCheck.App)
}

func TestAppHealthGetSingleHost(t *testing.T) {
	check := AppHealth{}
	appLabels := make(map[string]string)
	appLabels["router.hosts"] = "developer.indix.com:spectre.indix.tv"
	s := marathon.HealthCheck{
		GracePeriodSeconds:     180,
		IntervalSeconds:        90,
		TimeoutSeconds:         30,
		MaxConsecutiveFailures: 3,
		PortIndex:              0,
		Path:                   "/",
		Protocol:               "HTTP",
	}
	var appHC []*marathon.HealthCheck
	appHC = append(appHC, &s)

	app := marathon.Application{
		ID:           "/production.spectre",
		Instances:    10,
		TasksHealthy: 5,
		Labels:       appLabels,
		HealthChecks: appHC,
	}

	appCheck := check.Check(app)
	assert.Equal(t, Pass, appCheck.Result)
	assert.Equal(t, "apphealth", appCheck.CheckName)
	assert.Equal(t, "/production.spectre", appCheck.App)
}

func TestAppHealthProtocol(t *testing.T) {
	check := AppHealth{}
	appLabels := make(map[string]string)
	s := marathon.HealthCheck{
		GracePeriodSeconds:     240,
		IntervalSeconds:        60,
		TimeoutSeconds:         20,
		MaxConsecutiveFailures: 3,
		Protocol:               "COMMAND",
	}
	var appHC []*marathon.HealthCheck
	appHC = append(appHC, &s)

	app := marathon.Application{
		ID:           "/production.autobot",
		Instances:    10,
		TasksHealthy: 5,
		Labels:       appLabels,
		HealthChecks: appHC,
	}

	appCheck := check.Check(app)
	assert.Equal(t, Warning, appCheck.Result)
	assert.Equal(t, "apphealth", appCheck.CheckName)
	assert.Equal(t, "/production.autobot", appCheck.App)
	assert.Equal(t, "The healtcheck can be run for an App with HTTP endpoint!!!", appCheck.Message)
}

func TestAppHTTP(t *testing.T) {
	check := AppHealth{}
	appLabels := make(map[string]string)
	appLabels["router.hosts"] = "zeus.prod.indix.tv"
	s := marathon.HealthCheck{
		GracePeriodSeconds:     180,
		IntervalSeconds:        90,
		TimeoutSeconds:         30,
		MaxConsecutiveFailures: 3,
		PortIndex:              0,
		Path:                   "/health",
		Protocol:               "HTTP",
	}
	var appHC []*marathon.HealthCheck
	appHC = append(appHC, &s)

	app := marathon.Application{
		ID:           "/production.zeus",
		Instances:    10,
		TasksHealthy: 5,
		Labels:       appLabels,
		HealthChecks: appHC,
	}

	appCheck := check.Check(app)
	assert.Equal(t, Pass, appCheck.Result)
	assert.Equal(t, "HTTP Status OK!!", appCheck.Message)
	assert.Equal(t, "apphealth", appCheck.CheckName)
	assert.Equal(t, "/production.zeus", appCheck.App)
}
