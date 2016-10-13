package gogo_boy

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/jarcoal/httpmock"
)

func MockCampaignTriggerFailure(requestChecker func(map[string]interface{})) {
	httpmock.Activate()
	httpmock.RegisterResponder("POST", CampaignTriggerEndpoint,
		func(req *http.Request) (*http.Response, error) {
			buf := new(bytes.Buffer)
			buf.ReadFrom(req.Body)
			_request := buf.String()

			var request map[string]interface{}
			err := json.Unmarshal([]byte(_request), &request)
			checkErr(err)
			requestChecker(request)

			response := getFixtureWithPath("triggered_campaign_res_err.json")
			resp := httpmock.NewStringResponse(400, response)
			return resp, nil
		},
	)
}

func MockCampaignTriggerSuccess(requestChecker func(map[string]interface{})) {
	httpmock.Activate()
	httpmock.RegisterResponder("POST", CampaignTriggerEndpoint,
		func(req *http.Request) (*http.Response, error) {
			buf := new(bytes.Buffer)
			buf.ReadFrom(req.Body)
			_request := buf.String()

			var request map[string]interface{}
			err := json.Unmarshal([]byte(_request), &request)
			checkErr(err)
			requestChecker(request)

			response := getFixtureWithPath("triggered_campaign_res.json")
			resp := httpmock.NewStringResponse(201, response)
			return resp, nil
		},
	)
}

func MockDeletPushTokenFailure(requestChecker func(map[string]interface{})) {
	httpmock.Activate()
	httpmock.RegisterResponder("POST", DeletePushTokenEndpoint,
		func(req *http.Request) (*http.Response, error) {
			buf := new(bytes.Buffer)
			buf.ReadFrom(req.Body)
			_request := buf.String()

			var request map[string]interface{}
			err := json.Unmarshal([]byte(_request), &request)
			checkErr(err)
			requestChecker(request)

			response := getFixtureWithPath("push_removal_res_err.json")
			resp := httpmock.NewStringResponse(400, response)
			return resp, nil
		},
	)
}

func MockDeletPushTokenSuccess(requestChecker func(map[string]interface{})) {
	httpmock.Activate()
	httpmock.RegisterResponder("POST", DeletePushTokenEndpoint,
		func(req *http.Request) (*http.Response, error) {
			buf := new(bytes.Buffer)
			buf.ReadFrom(req.Body)
			_request := buf.String()

			var request map[string]interface{}
			err := json.Unmarshal([]byte(_request), &request)
			checkErr(err)
			requestChecker(request)

			response := getFixtureWithPath("push_removal_res.json")
			resp := httpmock.NewStringResponse(201, response)
			return resp, nil
		},
	)
}

func MockTrackSuccess(requestChecker func(map[string]interface{})) {
	httpmock.Activate()
	httpmock.RegisterResponder("POST", TrackEndpoint,
		func(req *http.Request) (*http.Response, error) {
			buf := new(bytes.Buffer)
			buf.ReadFrom(req.Body)
			_request := buf.String()

			var request map[string]interface{}
			err := json.Unmarshal([]byte(_request), &request)
			checkErr(err)
			requestChecker(request)

			response := getFixtureWithPath("track_success_res.json")
			resp := httpmock.NewStringResponse(201, response)
			return resp, nil
		},
	)
}

func MockTrackFailure(requestChecker func(map[string]interface{})) {
	httpmock.Activate()
	httpmock.RegisterResponder("POST", TrackEndpoint,
		func(req *http.Request) (*http.Response, error) {
			buf := new(bytes.Buffer)
			buf.ReadFrom(req.Body)
			_request := buf.String()

			var request map[string]interface{}
			err := json.Unmarshal([]byte(_request), &request)
			checkErr(err)
			requestChecker(request)

			resp := httpmock.NewStringResponse(500, "error")
			return resp, nil
		},
	)
}

func StopMocks() {
	httpmock.DeactivateAndReset()
}
