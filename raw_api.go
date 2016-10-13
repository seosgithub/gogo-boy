package gogo_boy

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

/*
	----------------------------------------------------------------------
	Raw requests for direct API usage
	----------------------------------------------------------------------
*/

const (
	TrackEndpoint           = "https://api.appboy.com/users/track"
	DeletePushTokenEndpoint = "https://api.appboy.com/push_notification/remove"
	CampaignTriggerEndpoint = "https://api.appboy.com/campaigns/trigger/send"

	timeoutDuration = time.Second * time.Duration(5)
)

// A track request is used to track purchases, user events, etc. It's configured
// so you can batch requests which won't count against your API limit.
type RawTrackRequest struct {
	AppGroupId string              `json:"app_group_id"`
	Attributes []RawAttributesInfo `json:"attributes,omitempty"` // Attributes are per-user information
	Purchases  []RawPurchaseInfo   `json:"purchases,omitempty"`  // Purchases are special events, each bound to a user
	Events     []RawEventInfo      `json:"events,omitempty"`     // Events
}

// App boy seperates deletion of push tokens into a seperate endpoint that's different
// than track request.  Why?  no idea.
type RawPushTokenDeleteRequest struct {
	AppGroupId string             `json:"app_group_id"`
	PushTokens []RawPushTokenInfo `json:"push_tokens"`
}

// Trigger a campaign
type RawCampaignTriggerRequest struct {
	Recipients []RawCampaignRecipient `json:"recipients"`
	CampaignId string                 `json:"campaign_id"`
	AppGroupId string                 `json:"app_group_id"`
}

// A campaign trigger has many user recipients
type RawCampaignRecipient struct {
	ExternalId        string                 `json:"external_user_id"`
	TriggerProperties map[string]interface{} `json:"trigger_properties"`
}

type RawAttributesInfo struct {
	ExternalId string `json:"external_id"` // The id of your user in your database

	//PushTokenImport bool               `json:"push_token_import,omitempty"` // Are you importing a push token?
	PushTokens []RawPushTokenInfo `json:"push_tokens,omitempty"` // A list of push tokens

	FirstName string `json:"first_name,omitempty"` // User's first name
	LastName  string `json:"last_name,omitempty"`  // User's last name
	Email     string `json:"email,omitempty"`      // User's email

	// These don't really exist, we have to dynamically place them in the JSON
	// at a later point
	CustomAttributes map[string]interface{} `json:"-"`
}

// You may upload a push token via the API but most people
type RawPushTokenInfo struct {
	AppId string `json:"app_id"`
	Token string `json:"token"`
}

type RawPurchaseInfo struct {
	ExternalId string  `json:"external_id"`
	ProductId  string  `json:"product_id"`
	Currency   string  `json:"currency"`
	Price      float32 `json:"price"`
	Quantity   int     `json:"quantity"`
	Time       string  `json:"time"`
}

type RawEventInfo struct {
	ExternalId string `json:"external_id"`
	//AppId      string `json:"app_id,omitempty"`
	Name               string `json:"name"`
	Time               string `json:"time"` // Time is in ISO 8601 format
	UpdateExistingOnly bool   `json:"_update_existing_only,omitempty"`

	// This library doesn't support these (yet)
	//Properties         string `json:"properties"`
}

// Post to track request endpoint
func RawPostTrackRequest(trackRequest *RawTrackRequest) error {
	// Marshal into a JSON string
	__json, err := json.Marshal(trackRequest)
	if err != nil {
		return fmt.Errorf("PostTrackRequest failed: %s", err)
	}

	// Then un-marshal because we need to place our custom attributes
	var _json map[string]interface{}
	if err := json.Unmarshal([]byte(__json), &_json); err != nil {
		return fmt.Errorf("PostTrackRequest failed to unmarshal _json: %s", err)
	}

	// For each user attribute, find custom attributes
	if _json["attributes"] != nil {
		_jsonAttributes := _json["attributes"].([]interface{})
		attributes := trackRequest.Attributes
		for i, attribute := range attributes {
			for k, v := range attribute.CustomAttributes {
				_jsonAttribute := _jsonAttributes[i].(map[string]interface{})
				_jsonAttribute[k] = v
			}
		}
	}

	json, err := json.Marshal(_json)
	if err != nil {
		return fmt.Errorf("PostTrackRequest failed to add custom attributes: %s", err)
	}

	// Create post request and make sure you set the content type
	req, err := http.NewRequest("POST", TrackEndpoint, bytes.NewBuffer(json))
	if err != nil {
		return fmt.Errorf("PostTrackRequest failed: %s", err)
	}
	if err != nil {
		return fmt.Errorf("PostTrackRequest failed: %s", err)
	}
	req.Header.Add("Content-Type", "application/json")

	// Our HTTP client
	client := &http.Client{
		Timeout: timeoutDuration,
	}

	// Execute request
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("PostTrackRequest failed: %s", err)
	}

	// Read body
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("PostTrackRequest failed: %s", err)
	}

	//fmt.Printf("%v", string(json))
	//fmt.Printf("%v", string(body))

	// App-Boy returns a 201 if this is successful
	if resp.StatusCode != 201 {
		return fmt.Errorf("PostTrackRequest failed: Expected status code from app boy to be a 201 but we received a: %d with the payload: '%s'\n", resp.StatusCode, body)
	}

	return nil
}

// Post push token request endpoint
func RawPostDeletePushTokenRequest(rawReq *RawPushTokenDeleteRequest) error {
	jsonStr, err := json.Marshal(rawReq)
	if err != nil {
		return fmt.Errorf("RawPostDeletPushTokenRequest failed: %s", err)
	}

	// Create post request and make sure you set the content type
	req, err := http.NewRequest("POST", DeletePushTokenEndpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return fmt.Errorf("RawPostDeletPushTokenRequest failed: %s", err)
	}
	if err != nil {
		return fmt.Errorf("RawPostDeletPushTokenRequest failed: %s", err)
	}
	req.Header.Add("Content-Type", "application/json")

	// Our HTTP client
	client := &http.Client{
		Timeout: timeoutDuration,
	}

	// Execute request
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("RawPostDeletPushTokenRequest failed: %s", err)
	}

	// Read body
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("RawPostDeletPushTokenRequest failed: %s", err)
	}

	//fmt.Printf("%v", string(json))
	//fmt.Printf("%v", string(body))

	// App-Boy returns a 201 if this is successful
	if resp.StatusCode != 201 {
		return fmt.Errorf("RawPostDeletPushTokenRequest failed: Expected status code from app boy to be a 201 but we received a: %d with the payload: '%s'\n", resp.StatusCode, body)
	}

	return nil
}

func RawPostCampaignTriggerRequest(rawReq *RawCampaignTriggerRequest) error {
	jsonStr, err := json.Marshal(rawReq)
	if err != nil {
		return fmt.Errorf("RawCampaignTriggerRequest failed: %s", err)
	}

	// Create post request and make sure you set the content type
	req, err := http.NewRequest("POST", CampaignTriggerEndpoint, bytes.NewBuffer(jsonStr))
	if err != nil {
		return fmt.Errorf("RawCampaignTriggerRequest failed: %s", err)
	}
	if err != nil {
		return fmt.Errorf("RawCampaignTriggerRequest failed: %s", err)
	}
	req.Header.Add("Content-Type", "application/json")

	// Our HTTP client
	client := &http.Client{
		Timeout: timeoutDuration,
	}

	// Execute request
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return fmt.Errorf("RawCampaignTriggerRequest failed: %s", err)
	}

	// Read body
	var body []byte
	body, err = ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil {
		return fmt.Errorf("RawCampaignTriggerRequest failed: %s", err)
	}

	//fmt.Printf("%v", string(json))
	//fmt.Printf("%v", string(body))

	// App-Boy returns a 201 if this is successful
	if resp.StatusCode != 201 {
		return fmt.Errorf("RawCampaignTriggerRequest failed: Expected status code from app boy to be a 201 but we received a: %d with the payload: '%s'\n", resp.StatusCode, body)
	}

	return nil
}
