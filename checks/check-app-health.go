package checks

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	maps "github.com/ashwanthkumar/golang-utils/maps"
	"github.com/gambol99/go-marathon"
)

// AppHealth to check for app's health
type AppHealth struct{}

func (h *AppHealth) Name() string {
	return "apphealth"
}

// Check App health status of App Endpoint by reading router.hosts
func (h *AppHealth) Check(app marathon.Application) AppCheck {
	hostList := maps.GetString(app.Labels, "router.hosts", "notDefined")
	hostSlice := strings.Split(hostList, ":")
	host := hostSlice[0]
	healthCheckPath := app.HealthChecks[0].Path
	healthCheckProtocol := app.HealthChecks[0].Protocol
	hostURL := strings.Join([]string{"http://", host, healthCheckPath}, "")
	result := Pass
	message := fmt.Sprintf("HTTP Status OK!!")

	if healthCheckProtocol == "HTTP" {
		resp, err := http.Get(hostURL)
		if err != nil {
			log.Println(log.Ldate|log.Ltime, "ERROR:", err)
		} else {
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.Println(log.Ldate|log.Ltime, "ERROR:", err)
			}
			responseString := string(body)

			if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
				result = Pass
			} else {
				result = Critical
				message = fmt.Sprintf("HTTP Response Status: %s  %s ", strconv.Itoa(resp.StatusCode), responseString)
			}
		}
	} else {
		result = Warning
		message = fmt.Sprintf("The healtcheck can be run for an App with HTTP endpoint!!!")
	}
	return AppCheck{
		App:       app.ID,
		Labels:    app.Labels,
		CheckName: h.Name(),
		Result:    result,
		Message:   message,
		Timestamp: time.Now(),
	}
}
