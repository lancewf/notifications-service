package inspec

import (
	"encoding/json"
	"fmt"
	"github.com/buger/jsonparser"
	"github.com/lancewf/notifications-service/pkg/util"
	log "github.com/sirupsen/logrus"
)

type Report struct {
	MinImpact float32
	ID        string
	NodeName  string
	Profiles  []Profile
}

type Profile struct {
	Name     string
	Controls []Control
}

type Control struct {
	ID      string
	Title   string
	Impact  float32
	Results []Result
}

type Result struct {
	Status   string //passed, skipped, failed
	CodeDesc string
}

type SlackMessage struct {
	Username    string            `json:"username"`
	Text        string            `json:"text"`
	IconURL     string            `json:"icon_url"`
	Attachments []SlackAttachment `json:"attachments"`
}

type SlackAttachment struct {
	Text   string       `json:"text"`
	Color  string       `json:"color"`
	Fields []SlackField `json:"fields"`
}

type SlackField struct {
	Title string `json:"title"`
	Value string `json:"value"`
	Short bool   `json:"short"`
}

type IFTTTMessage struct {
	Value1 string `json:"value1"`
	Value2 string `json:"value2"`
	Value3 string `json:"value3"`
}

func ParseReport(rawReport []byte, minImpact float32) Report {
	profiles := getProfiles(rawReport)
	return Report{
		MinImpact: minImpact,
		NodeName:  util.GetStringIfExists("node_name", rawReport),
		ID:        util.GetStringIfExists("node_uuid", rawReport),
		Profiles:  profiles,
	}
}

func (report Report) HasNotificationToSend() bool {
	return len(report.failedProfiles(0.0)) > 0
}

func (report Report) WebHookMessage() string {
	return "Failed InSpec Report!"
}

func (report Report) IFTTTWebHookMessage() string {
	failedProfilesName := "<None>"
	failedControlInfo := "<None>"
	failedProfiles := report.failedProfiles(0.0)
	if len(failedProfiles) > 0 {
		failedProfile := failedProfiles[0]
		failedProfilesName = failedProfile.Name

		failedControls := failedProfile.failedControls(0.0)
		if len(failedControls) > 0 {
			failedControl := failedControls[0]
			failedControlInfo = fmt.Sprintf("%s:%s", failedControl.ID, failedControl.Title)
		}
	}

	msg := IFTTTMessage{
		Value1: failedProfilesName,
		Value2: failedControlInfo,
		Value3: fmt.Sprintf("%d", report.numberOfFailedTests()),
	}

	JSONRaw, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Error parsing message %v", err)
		return ""
	}

	return string(JSONRaw)
}

func (report Report) SlackWebhookMessage() string {
	failedProfilesName := "<None>"
	failedControlInfo := "<None>"
	failedProfiles := report.failedProfiles(0.0)
	if len(failedProfiles) > 0 {
		failedProfile := failedProfiles[0]
		failedProfilesName = failedProfile.Name

		failedControls := failedProfile.failedControls(0.0)
		if len(failedControls) > 0 {
			failedControl := failedControls[0]
			failedControlInfo = fmt.Sprintf("%s:%s", failedControl.ID, failedControl.Title)
		}
	}

	msg := SlackMessage{
		Username: "Notification Service",
		Text:     fmt.Sprintf("InSpec found a critical control failure on node %q", report.NodeName),
		IconURL:  "https://docs.chef.io/_static/chef_logo_v2.png",
		Attachments: []SlackAttachment{
			{
				Text: fmt.Sprintf("%d tests failed. Rerun the test locally for full details.",
					report.numberOfFailedTests()),
				Color: "warning",
				Fields: []SlackField{
					{
						Title: "Control ID::Title",
						Value: failedControlInfo,
						Short: false,
					},
					{
						Title: "Profile",
						Value: failedProfilesName,
						Short: false,
					},
					{
						Title: "Node",
						Value: report.NodeName,
						Short: false,
					},
					{
						Title: "Highest Failed Impact",
						Value: fmt.Sprintf("%.1f", report.highestFailedImpact()),
						Short: false,
					},
				},
			},
		},
	}

	JSONRaw, err := json.Marshal(msg)
	if err != nil {
		log.Errorf("Error parsing message %v", err)
		return fmt.Sprintf("{\"text\": \"%s\"}", "Error")
	}

	return string(JSONRaw)
}

func (report Report) numberOfFailedTests() int {
	numberOfFailedTest := 0
	failedProfiles := report.failedProfiles(report.MinImpact)
	for _, failedProfile := range failedProfiles {
		numberOfFailedTest = failedProfile.numberOfFailedTests(0.0)
	}

	return numberOfFailedTest
}

func (report Report) highestFailedImpact() float32 {
	var highestFailedImpact float32 = -1.0
	for _, failedProfile := range report.failedProfiles(report.MinImpact) {
		profileHighestFailedImpact := failedProfile.highestFailedImpact(report.MinImpact)
		if highestFailedImpact < profileHighestFailedImpact {
			highestFailedImpact = profileHighestFailedImpact
		}
	}

	return highestFailedImpact
}

func (report Report) failedProfiles(minImpact float32) []Profile {
	profiles := make([]Profile, 0)
	for _, profile := range report.Profiles {
		failedControls := profile.failedControls(minImpact)

		if len(failedControls) > 0 {
			profiles = append(profiles, profile)
		}
	}

	return profiles
}

func (profile Profile) highestFailedImpact(minImpact float32) float32 {
	var highestFailedImpact float32 = -1.0
	for _, failedControl := range profile.failedControls(minImpact) {
		if highestFailedImpact < failedControl.Impact {
			highestFailedImpact = failedControl.Impact
		}
	}

	return highestFailedImpact
}

func (profile Profile) numberOfFailedTests(minImpact float32) int {
	numberOfFailedTest := 0
	for _, failedControl := range profile.failedControls(minImpact) {
		numberOfFailedTest = numberOfFailedTest + failedControl.numberOfFailedTests()
	}
	return numberOfFailedTest
}

func (profile Profile) failedControls(minImpact float32) []Control {
	controls := make([]Control, 0)
	for _, control := range profile.Controls {
		if control.Impact >= minImpact {
			failedResults := control.failedResults()

			if len(failedResults) > 0 {
				controls = append(controls, control)
			}
		}
	}

	return controls
}

func (control Control) failedResults() []Result {
	results := make([]Result, 0)
	for _, result := range control.Results {
		if result.Status == "failed" {
			results = append(results, result)
		}
	}

	return results
}

func (control Control) numberOfFailedTests() int {
	return len(control.failedResults())
}

func getProfiles(rawReport []byte) []Profile {
	profiles := make([]Profile, 0)
	jsonparser.ArrayEach(rawReport, func(rawProfile []byte, _ jsonparser.ValueType, _ int, err error) {
		controls := getControls(rawProfile)
		profile := Profile{
			Name:     util.GetStringIfExists("name", rawProfile),
			Controls: controls,
		}
		profiles = append(profiles, profile)
	}, "profiles")

	return profiles
}

func getControls(rawProfile []byte) []Control {
	controls := make([]Control, 0)
	jsonparser.ArrayEach(rawProfile, func(rawControl []byte, _ jsonparser.ValueType, _ int, err error) {
		results := getResults(rawControl)
		control := Control{
			ID:      util.GetStringIfExists("id", rawControl),
			Title:   util.GetStringIfExists("title", rawControl),
			Impact:  util.GetFloat32IfExists("impact", rawControl),
			Results: results,
		}
		controls = append(controls, control)
	}, "controls")

	return controls
}

func getResults(rawControl []byte) []Result {
	results := make([]Result, 0)
	jsonparser.ArrayEach(rawControl, func(rawResult []byte, _ jsonparser.ValueType, _ int, err error) {
		result := Result{
			Status:   util.GetStringIfExists("status", rawResult),
			CodeDesc: util.GetStringIfExists("code_desc", rawResult),
		}
		results = append(results, result)
	}, "results")

	return results
}
