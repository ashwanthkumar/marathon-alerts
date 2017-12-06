package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/ashwanthkumar/marathon-alerts/checks"
	"github.com/ashwanthkumar/marathon-alerts/notifiers"
	flag "github.com/spf13/pflag"

	marathon "github.com/gambol99/go-marathon"
	"github.com/rcrowley/go-metrics"
)

var appChecker AppChecker
var alertManager AlertManager

// Check settings
var minHealthyWarningThreshold float32
var minHealthyCriticalThreshold float32
var minInstancesWarningThreshold float32
var minInstancesCriticalThreshold float32

// Required flags
var marathonURI string
var checkInterval time.Duration
var alertSuppressDuration time.Duration
var debugMode bool
var pidFile string

// Slack flags
var slackWebhooks string
var slackChannel string
var slackOwners string

// DebugMetricsRegistry is used for pushing debug level metrics by rest of the app
var DebugMetricsRegistry metrics.Registry

func main() {
	log.SetFlags(log.Ldate | log.Ltime | log.Lmicroseconds | log.LUTC | log.Lshortfile)
	log.SetOutput(os.Stdout)
	os.Args[0] = "marathon-alerts"
	defineFlags()
	flag.Parse()
	pid := []byte(fmt.Sprintf("%d\n", os.Getpid()))
	err := ioutil.WriteFile(pidFile, pid, 0644)
	if err != nil {
		fmt.Println("Unable to write pid file. ")
		log.Fatalf("Error - %v\n", err)
	}

	client, err := marathonClient(marathonURI)
	if err != nil {
		fmt.Printf("%v\n", err)
		os.Exit(1)
	}
	DebugMetricsRegistry = metrics.NewPrefixedRegistry("debug")

	minHealthyTasks := &checks.MinHealthyTasks{
		DefaultCriticalThreshold: minHealthyCriticalThreshold,
		DefaultWarningThreshold:  minHealthyWarningThreshold,
	}
	minInstances := &checks.MinInstances{
		DefaultCriticalThreshold: minHealthyCriticalThreshold,
		DefaultWarningThreshold:  minHealthyWarningThreshold,
	}
	suspendedCheck := &checks.SuspendedCheck{}
	appHealth := &checks.AppHealth{
      DefaultCriticalThreshold: minHealthyCriticalThreshold,
	  DefaultWarningThreshold:  minHealthyWarningThreshold,
	}
	checks := []checks.Checker{minHealthyTasks, minInstances, suspendedCheck, appHealth}

	appChecker = AppChecker{
		Client:        client,
		CheckInterval: checkInterval,
		Checks:        checks,
	}
	appChecker.Start()

	var allNotifiers []notifiers.Notifier
	slack := notifiers.Slack{
		Webhook: slackWebhooks,
		Channel: slackChannel,
		Owners:  slackOwners,
	}
	allNotifiers = append(allNotifiers, &slack)

	alertManager = AlertManager{
		CheckerChan:      appChecker.AlertsChannel,
		SuppressDuration: alertSuppressDuration,
		Notifiers:        allNotifiers,
	}
	alertManager.Start()

	metrics.RegisterDebugGCStats(DebugMetricsRegistry)
	metrics.RegisterRuntimeMemStats(DebugMetricsRegistry)
	go metrics.CaptureDebugGCStats(DebugMetricsRegistry, 15*time.Minute)
	go metrics.CaptureRuntimeMemStats(DebugMetricsRegistry, 5*time.Minute)
	go metrics.Log(metrics.DefaultRegistry, 60*time.Second, log.New(os.Stderr, "metrics: ", log.Lmicroseconds))
	if debugMode {
		go metrics.Log(DebugMetricsRegistry, 300*time.Second, log.New(os.Stderr, "debug-metrics: ", log.Lmicroseconds))
	}
	appChecker.RunWaitGroup.Wait()
	// Handle signals and cleanup all routines
}

func marathonClient(uri string) (marathon.Marathon, error) {
	config := marathon.NewDefaultConfig()
	config.URL = uri
	config.HTTPClient = &http.Client{
		Timeout: (30 * time.Second),
	}

	return marathon.NewClient(config)
}

func defineFlags() {
	flag.StringVar(&marathonURI, "uri", "", "Marathon URI to connect")
	flag.StringVar(&pidFile, "pid", "PID", "File to write PID file")
	flag.BoolVar(&debugMode, "debug", false, "Enable debug mode. More counters for now.")
	flag.DurationVar(&checkInterval, "check-interval", 60*time.Second, "Check runs periodically on this interval")
	flag.DurationVar(&alertSuppressDuration, "alerts-suppress-duration", 30*time.Minute, "Suppress alerts for this duration once notified")

	// Check flags
	flag.Float32Var(&minHealthyWarningThreshold, "check-min-healthy-warn-threshold", 0.75, "Min Healthy instances check warning threshold")
	flag.Float32Var(&minHealthyCriticalThreshold, "check-min-healthy-critical-threshold", 0.5, "Min Healthy instances check fail threshold")
	flag.Float32Var(&minInstancesWarningThreshold, "check-min-instances-warn-threshold", 0.75, "Min Instances check warning threshold")
	flag.Float32Var(&minInstancesCriticalThreshold, "check-min-instances-critical-threshold", 0.5, "Min Instances check fail threshold")

	// Slack flags
	flag.StringVar(&slackWebhooks, "slack-webhook", "", "Comma list of Slack webhooks to post the alert")
	flag.StringVar(&slackChannel, "slack-channel", "", "#Channel / @User to post the alert (defaults to webhook configuration)")
	flag.StringVar(&slackOwners, "slack-owner", "", "Comma list of owners who should be alerted on the post")
}
