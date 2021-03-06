[![Build Status](https://travis-ci.org/ashwanthkumar/marathon-alerts.svg?branch=master)](https://travis-ci.org/ashwanthkumar/marathon-alerts)
[![codecov.io](https://codecov.io/github/ashwanthkumar/marathon-alerts/coverage.svg?branch=master)](https://codecov.io/github/ashwanthkumar/marathon-alerts?branch=master)

# marathon-alerts

Marathon Alerts is a tool for monitoring the apps running on marathon. Inspired from [kubernetes-alerts](https://github.com/AcalephStorage/kubernetes-alerts) and [consul-alerts](https://github.com/AcalephStorage/consul-alerts).

This was initially built for Marathon 0.8.0, hence we don't use the event bus.

## Usage
```
$ marathon-alerts --help
Usage of marathon-alerts:
      --alerts-suppress-duration duration              Suppress alerts for this duration once notified (default 30m0s)
      --check-interval duration                        Check runs periodically on this interval (default 30s)
      --check-min-healthy-critical-threshold value     Min Healthy instances check fail threshold (default 0.5)
      --check-min-healthy-warn-threshold value         Min Healthy instances check warning threshold (default 0.75)
      --check-min-instances-critical-threshold value   Min Instances check fail threshold (default 0.5)
      --check-min-instances-warn-threshold value       Min Instances check warning threshold (default 0.75)
      --debug                                          Enable debug mode. More counters for now.
      --pid string                                     File to write PID file (default "PID")
      --slack-channel string                           #Channel / @User to post the alert (defaults to webhook configuration)
      --slack-owner string                             Comma list of owners who should be alerted on the post
      --slack-webhook string                           Comma list of Slack webhooks to post the alert
      --uri string                                     Marathon URI to connect
```

Example invocation would be like the following
```
$ marathon-alerts --uri http://marathon1:8080,marathon2:8080 \
                  --slack-webhook https://hooks.slack.com/services/..../ \
                  --slack-owner ashwanthkumar,slackbot
```

## App Labels
Apart from the flags that are used while starting up, the functionality can be controlled at an app level using labels in the app specification. The following table explains the properties and it's usage.

| Property  | Description  |  Example  |
|  ---    |   ---      |  ---    |
| alerts.enabled  | Controls if the alerts for the app should be enabled or disabled. Defaults - true | false |
| alerts.checks.subscribe  | Comma separated list of checks that needs to be run. Defaults - all | all |
| alerts.routes  | Ability to route different checks to different notifiers based on their level. See the section below on Routes to understand how you can add routes to your apps. Defaults - `*/resolved/*;*/warning/*;*/critical/*` | min-healthy/critical/pagerduty;min-healthy/warning/slack |
| alerts.min-healthy.critical.threshold  | Failure threshold for min-healthy check. Defaults - `--check-min-healthy-critical-threshold` | 0.5 |
| alerts.min-healthy.warn.threshold  | Warning threshold for min-healthy check. Defaults - `--check-min-healthy-warn-threshold` | 0.4 |
| alerts.min-instances.critical.threshold  | Failure threshold for min-instances check. Defaults - `--check-min-instances-critical-threshold` | 0.5 |
| alerts.min-instances.warn.threshold  | Warning threshold for min-instances check. Defaults - `--check-min-instances-warn-threshold` | 0.4 |
| alerts.slack.webhook  | Comma separated list of Slack webhooks to send slack notifications. Overrides - `--slack-webhook` | http://hooks.slack.com/.../ |
| alerts.slack.channel  | #Channel / @User to post the alert into. Overrides - `--slack-channel`  | z_development |
| alerts.slack.owners  | Comma separated list of users who should be tagged in the alert. Overrides - `--slack-owner`  | ashwanthkumar,slackbot |

## Metrics
We collect some metrics internally in marathon-alerts. They're dumped periodically to STDERR. You can find the list of metrics and it's usage in the following table

| Metric | Description |
| :------------- | :------------- |
| alerts-suppressed-cleaned | Number of alerts we cleaned up because they got expired from suppress duration. |
| marathon-all-apps-response-time | Response time of marathon's /v2/apps API call |
| notifications-total | Total number of notifications we sent from AlertManager to NotificationManager |
| notifications-warning | Number of Warning check notifications we sent from AlertManager to NotificationManager |
| notifications-critical | Number of Critical check notifications we sent from AlertManager to NotificationManager |
| notifications-resolved | Number of Pass (aka Resolved) check notifications we sent from AlertManager to NotificationManager |
| notifications-rate | Meter metric that denotes the rate at which notifications are being sent |

## Debug Metrics
Apart from the standard metrics above, we also collect quite a few other metrics, mostly for debugging purposes. You can enable these metrics if run `marathon-alerts` with a `--debug` flag.

| Metric | Description |
| :------------- | :------------- |
| alerts-suppressed-called | Number of times we called AlertManager.cleanUpSupressedAlerts() |
| alerts-process-check-called | Number of times we called AlertManager.processCheck() |
| alerts-manager-stopped | Number of times we called AlertManager.Stop() |
| apps-checker-stopped | Number of times we called AppChecker.Stop() |
| apps-checker-marathon-all-apps-api | Number of times we called Marathon's /v2/apps API |
| apps-checker-alerts-sent | Number of checks we sent to AlertManager from AppChecker |
| apps-checker-check-&lt;name&gt; | Number of checks identified by &lt;name&gt; we sent to AlertManager |
| apps-checker-app-&lt;id&gt; | Number of checks for an app identified by &lt;id&gt; we sent to AlertManager |
| apps-checker-&lt;id&gt;-&lt;name&gt; | Number of checks identified by &lt;name&gt; for an app identified by &lt;id&gt; we sent to AlertManager |
| notifications-warning-rate | Meter metric that denotes the rate at which warning notifications are being sent |
| notifications-critical-rate | Meter metric that denotes the rate at which critical notifications are being sent |
| notifications-resolved-rate | Meter metric that denotes the rate at which resolved notifications are being sent |

## Routes
From v0.3.0-RC7 onwards we've an ability to route different check alerts to different notifiers. On a per-app basis you can control the routes using `alerts.routes` label. The format of the value should be as following -
```
<check-name>/<check-level>/<notifier-name>;[<check-name>/<check-level>/<notifier-name>]
```

### Rules
1. Check name and Notifier names can be glob patterns. No complicated regex allowed as of now.
2. Check level has to be one of warning / pass / critical / resolved.
3. Multiple routes can be defined by separating them using `;`.

Default routes if none specified is -  `"*/warning/*;*/critical/*;*/resolved/*"`. It means we'll route all check's warning / critical / resolved notifications to all available notifiers.

## Releases
Binaries are available [here](https://github.com/ashwanthkumar/marathon-alerts/releases).

## Deployment
We've a sample `marathon.json.conf` that we use in our organization along with [`marathonctl deploy`](https://github.com/ashwanthkumar/marathonctl).

## Building
To build from source, you need [`glide`](http://glide.sh/) tool in `$PATH`.

```
$ cd $GOPATH/src
$ mkdir -p github.com/ashwanthkumar/marathon-alerts
$ git clone https://github.com/ashwanthkumar/marathon-alerts.git github.com/ashwanthkumar/marathon-alerts
$ cd github.com/ashwanthkumar/marathon-alerts
$ make setup  # Downloads the required dependencies
$ make test   # Runs the test
$ make build  # Builds the distribution specific binary
```

## Available Checks
- [x] `min-healthy` - Minimum % of Task instances that should be healthy else this check is fired.
- [x] `min-instances` - Minimum % of Task instances that should be healthy or staged, else this check is fired.
- [ ] `max-instances` - If the number of instances goes beyond some % of the pre-defined max limit
- [x] `suspended` - If the service was suspended by mistake or unintentionally. `min-healthy` doesn't catch suspended services today.

## Notifiers
- [x] Slack
- [ ] Influx
- [ ] Pager Duty
- [ ] Email

## Contribute
If you've any feature requests or issues, please open a Github issue. We accept PRs. Fork away!

## License
http://www.apache.org/licenses/LICENSE-2.0
