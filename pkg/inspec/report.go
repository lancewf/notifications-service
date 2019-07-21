package inspec

import (
	"github.com/buger/jsonparser"
	"github.com/lancewf/notifications-service/pkg/util"
)

type Report struct {
	ID       string
	Profiles []Profile
}

type Profile struct {
	Name     string
	Controls []Control
}

type Control struct {
	Title   string
	Impact  float32
	Results []Result
}

type Result struct {
	Status   string //passed, skipped, failed
	CodeDesc string
}

func (report Report) HasNotificationToSend() bool {
	return len(report.failedProfiles(0.0)) > 0
}

func (report Report) WebHookMessage() string {
	return "Failed InSpec Report!"
}

func (report Report) IFTTTWebHookMessage() string {
	return "{\"value1\" : \"InSpec\", \"value2\" : \"profile names\", \"value3\" : \"whale-server\"}"
}

func (report Report) SlackWebhookMessage() string {
	return "{\"text\": \"Failed InSpec Report!\"}"
}

func (report Report) failedProfiles(impact float32) []Profile {
	profiles := make([]Profile, 0)
	for _, profile := range report.Profiles {
		failedControls := profile.failedControls(impact)

		if len(failedControls) > 0 {
			profiles = append(profiles, profile)
		}
	}

	return profiles
}

func (profile Profile) failedControls(impact float32) []Control {
	controls := make([]Control, 0)
	for _, control := range profile.Controls {
		if control.Impact >= impact {
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

func ParseReport(rawReport []byte) Report {
	profiles := getProfiles(rawReport)
	return Report{
		ID:       util.GetStringIfExists("node_uuid", rawReport),
		Profiles: profiles,
	}
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
