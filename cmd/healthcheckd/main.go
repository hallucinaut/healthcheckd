package main

import (
	"fmt"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/fatih/color"
)

type HealthCheck struct {
	Name     string
	URL      string
	Method   string
	Timeout  time.Duration
	Status   string
	Response time.Duration
}

func main() {
	if len(os.Args) < 2 {
		fmt.Println(color.CyanString("healthcheckd - Multi-Service Health Aggregator"))
		fmt.Println()
		fmt.Println("Usage: healthcheckd <service1> <service2> ...")
		fmt.Println("Format: name=url[method]")
		fmt.Println()
		fmt.Println("Examples:")
		fmt.Println("  healthcheckd api=http://localhost:8080/health")
		fmt.Println("  web=http://localhost:3000[GET]")
		os.Exit(1)
	}

	fmt.Println(color.CyanString("\n=== SERVICE HEALTH CHECK ===\n"))

	var checks []HealthCheck
	for _, arg := range os.Args[1:] {
		parts := parseServiceConfig(arg)
		if len(parts) >= 2 {
			checks = append(checks, HealthCheck{
				Name:    parts[0],
				URL:     parts[1],
				Method:  "GET",
				Timeout: 5 * time.Second,
			})
		}
	}

	runHealthChecks(checks)
	generateGrafanaConfig()
}

func parseServiceConfig(config string) []string {
	parts := strings.Split(config, "=")
	if len(parts) != 2 {
		return nil
	}

	name := parts[0]
	url := parts[1]

	// Check for method override in URL
	if strings.Contains(url, "[") && strings.HasSuffix(url, "]") {
		parts = strings.SplitN(url, "[", 2)
		if len(parts) == 2 {
			url = parts[0]
			// method can be extracted if needed
		}
	}

	return []string{name, url}
}

func runHealthChecks(checks []HealthCheck) {
	var up, down int

	for _, check := range checks {
		start := time.Now()
		status, _ := checkEndpoint(check.URL, check.Method, check.Timeout)
		elapsed := time.Since(start)

		check.Response = elapsed
		check.Status = status

		if status == "UP" {
			up++
			fmt.Printf("%-20s %s (%s)\n", check.Name, color.GreenString("UP"), formatDuration(elapsed))
		} else {
			down++
			fmt.Printf("%-20s %s (%s)\n", check.Name, color.RedString("DOWN"), formatDuration(elapsed))
		}
	}

	fmt.Println()
	fmt.Printf("Summary: %d UP, %d DOWN\n", up, down)
}

func checkEndpoint(url, method string, timeout time.Duration) (string, time.Duration) {
	client := &http.Client{Timeout: timeout}
	req, _ := http.NewRequest(method, url, nil)
	resp, err := client.Do(req)

	if err != nil {
		return "DOWN", 0
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		return "UP", 0
	}

	return "DOWN", 0
}

func formatDuration(d time.Duration) string {
	if d < time.Millisecond {
		return fmt.Sprintf("%dµs", d.Microseconds())
	}
	return fmt.Sprintf("%dms", d.Milliseconds())
}

func generateGrafanaConfig() {
	fmt.Print(color.YellowString("\n=== GRAFANA DASHBOARD CONFIG ==="))
	fmt.Print(`
{
  "dashboard": {
    "title": "Service Health Dashboard",
    "panels": [
      {
        "title": "Service Status",
        "type": "stat",
        "targets": [
          {
            "expr": "up{job=\"services\"}"
          }
        ]
      }
    ]
  }
}
`)
}