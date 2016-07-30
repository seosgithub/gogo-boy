package gogo_boy

import (
	"bytes"
	"encoding/json"
	"net/http"

	"github.com/jarcoal/httpmock"
)

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

func StopMocks() {
	httpmock.DeactivateAndReset()
}
