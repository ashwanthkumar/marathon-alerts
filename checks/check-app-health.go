package checks

import (
	"fmt"
	"strings"
	"time"

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
	host_path := app.HealthChecks[0].path
	host_url = strings.Join("http://", host, host_path)

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
		result =: Pass
		fmt.Println("HTTP Status OK!")
	} else {
		result = Critical
		message = fmt.Println("HTTP Response Status: " + strconv.Itoa(resp.StatusCode))
	}

	return AppCheck{
		App:       app.ID,
		Labels:    app.Labels,
		CheckName: n.Name(),
		Result:    result,
		Message:   message,
		Timestamp: time.Now(),
	}
}
