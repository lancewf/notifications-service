This is for users that only have a few nodes and want to be notified when InSpec runs fail. 
The user will be notified of failures through slack and IFTTT. With IFTTT one can get a notification on their phone, an email, or change their lights to flashing red. 
One extra feature is the service does not over notify. When the system is in the failed state a user is notified when the failure starts, ends, and configurable daily reminders. 

To set this up, point the Effortless InSpec Automate `server_url` config to this server's "/inspec_reports" URL. Below is an example of the Effortless InSpec config. This forwards all the InSpec run reports to this service. It is also configurable to forward all messages to an external Automate. 
```
[automate]
enable = true
server_url = "http://localhost:8095/inspec_reports"
token = 'none'
user = 'none'
```


Below is a example config.
```
[service]
host = "localhost"
port = 8095

[ifttt_webhook]
url = "https://maker.ifttt.com/trigger/inspec/with/key/fake"

[slack_webhook]
url = "https://hooks.slack.com/services/"

[inspec]
min_impact_to_notify = 0.7
```
