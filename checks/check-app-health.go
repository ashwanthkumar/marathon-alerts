package checks

import (
	"fmt"
	"strings"
	"time"
	"net/http"
	"io/ioutil"
    "strconv"
    "log"

	maps "github.com/ashwanthkumar/golang-utils/maps"
	"github.com/gambol99/go-marathon"
)

type AppHealth struct{}

func (h *AppHealth) Name() string {
	return "apphealth"
}

// Checks App health status of App Endpoint by reading router.hosts 
//  and Check for 200 status code
func (h *AppHealth) Check(app marathon.Application) AppCheck {
	host := maps.GetString(app.Labels, "router.hosts", "notDefined")
    host_path  := app.HealthChecks[0].Path
	host_url := strings.Join([]string{"http://", host, host_path}, "")
    result := Pass
    message := fmt.Sprintf("HTTP Status OK!!")

	resp, err := http.Get(host_url)
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	responseString := string(body)


	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		result = Pass
	} else {
		result = Critical
		message = fmt.Sprintf("HTTP Response Status: ,  " + strconv.Itoa(resp.StatusCode), responseString)
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
