package gogo_boy

import (
	"fmt"
	"strings"
	"testing"

	"github.com/jarcoal/httpmock"
	. "github.com/smartystreets/goconvey/convey"
)

func TestRawAPI(t *testing.T) {
	before := func() {
		httpmock.Activate()
	}

	after := func() {
		httpmock.DeactivateAndReset()
	}

	Convey("Can execute a track request to app-boy", t, func() {
		before()
		defer after()

		// Mock request to app-boy
		var request map[string]interface{}
		MockTrackSuccess(func(_request map[string]interface{}) {
			request = _request
		})

		// Construct a track request
		trackRequest := &RawTrackRequest{
			AppGroupId: "foo",
			Attributes: []RawAttributesInfo{
				RawAttributesInfo{
					FirstName:  "bar",
					ExternalId: "foo",
					Email:      "test@test.com",
					CustomAttributes: map[string]interface{}{
						"custom_value_attribute":  "custom_value_attribute_value",
						"custom_value_attribute2": "custom_value_attribute_value2",
					},
					PushTokens: []RawPushTokenInfo{
						RawPushTokenInfo{
							AppId: "blah",
							Token: "apple-token",
						},
					},
				},
			},
			Purchases: []RawPurchaseInfo{
				RawPurchaseInfo{
					ExternalId: "foo",
					ProductId:  "baz",
					Currency:   "USD",
					Price:      4.92,
					Quantity:   1,
					Time:       "Z070000",
				},
			},
			Events: []RawEventInfo{
				RawEventInfo{
					ExternalId: "foo",
					Name:       "bak",
					Time:       "Z070000",
				},
			},
		}

		err := RawPostTrackRequest(trackRequest)
		checkErr(err)
		So(err, ShouldEqual, nil)

		// Root
		So(request["app_group_id"], ShouldEqual, "foo")

		// Check attributes
		attributes := request["attributes"].([]interface{})
		attribute := attributes[0].(map[string]interface{})
		So(attribute["first_name"], ShouldEqual, "bar")
		So(attribute["external_id"], ShouldEqual, "foo")
		So(attribute["email"], ShouldEqual, "test@test.com")
		So(attribute["custom_value_attribute"], ShouldEqual, "custom_value_attribute_value")
		So(attribute["custom_value_attribute2"], ShouldEqual, "custom_value_attribute_value2")

		pushTokenAttributes := attribute["push_tokens"].([]interface{})
		pushTokenAttribute := pushTokenAttributes[0].(map[string]interface{})
		So(pushTokenAttribute["app_id"], ShouldEqual, "blah")
		So(pushTokenAttribute["token"], ShouldEqual, "apple-token")

		// Check purchases
		purchases := request["purchases"].([]interface{})
		purchase := purchases[0].(map[string]interface{})
		So(purchase["external_id"], ShouldEqual, "foo")
		So(purchase["product_id"], ShouldEqual, "baz")
		So(purchase["currency"], ShouldEqual, "USD")
		So(purchase["price"], ShouldEqual, 4.92)
		So(purchase["quantity"], ShouldEqual, 1)
		So(purchase["time"], ShouldEqual, "Z070000")

		// Check events
		events := request["events"].([]interface{})
		event := events[0].(map[string]interface{})
		So(event["external_id"], ShouldEqual, "foo")
		So(event["name"], ShouldEqual, "bak")
		So(event["time"], ShouldEqual, "Z070000")
	})

	Convey("Can execute a track request missing some fields", t, func() {
		before()
		defer after()

		// Mock request to app-boy
		var request map[string]interface{}
		MockTrackSuccess(func(_request map[string]interface{}) {
			request = _request
		})

		// Construct a track request
		trackRequest := &RawTrackRequest{
			AppGroupId: "foo",
			Attributes: []RawAttributesInfo{
				RawAttributesInfo{
					PushTokens: []RawPushTokenInfo{},
					ExternalId: "foo",
				},
			},
		}

		err := RawPostTrackRequest(trackRequest)
		checkErr(err)
		So(err, ShouldEqual, nil)

		// Root
		So(request["app_group_id"], ShouldEqual, "foo")

		// Check attributes
		attributes := request["attributes"].([]interface{})

		attribute := attributes[0].(map[string]interface{})

		So(attribute["external_id"], ShouldEqual, "foo")

		if _, ok := attribute["first_name"]; ok {
			panic("first_name existed, should have not")
		}

		if _, ok := attribute["email"]; ok {
			panic("email existed, should have not")
		}

		if _, ok := attribute["email"]; ok {
			panic("email existed, should have not")
		}

		if _, ok := attribute["push_tokens"]; ok {
			panic("push_tokens existed, should have not")
		}

	})

	Convey("Does fail gracefully with non succesful request", t, func() {
		before()
		defer after()

		// Mock request to app-boy
		var request map[string]interface{}
		MockTrackSuccess(func(_request map[string]interface{}) {
			request = _request
		})

		// Construct a track request
		trackRequest := &RawTrackRequest{
			AppGroupId: "foo",
			Attributes: []RawAttributesInfo{
				RawAttributesInfo{
					FirstName:  "bar",
					ExternalId: "foo",
					Email:      "test@test.com",
					CustomAttributes: map[string]interface{}{
						"custom_value_attribute":  "custom_value_attribute_value",
						"custom_value_attribute2": "custom_value_attribute_value2",
					},
					PushTokens: []RawPushTokenInfo{
						RawPushTokenInfo{
							AppId: "blah",
							Token: "apple-token",
						},
					},
				},
			},
			Purchases: []RawPurchaseInfo{
				RawPurchaseInfo{
					ExternalId: "foo",
					ProductId:  "baz",
					Currency:   "USD",
					Price:      4.92,
					Quantity:   1,
					Time:       "Z070000",
				},
			},
			Events: []RawEventInfo{
				RawEventInfo{
					ExternalId: "foo",
					Name:       "bak",
					Time:       "Z070000",
				},
			},
		}

		err := RawPostTrackRequest(trackRequest)
		checkErr(err)
		So(err, ShouldEqual, nil)

		// Root
		So(request["app_group_id"], ShouldEqual, "foo")

		// Check attributes
		attributes := request["attributes"].([]interface{})
		attribute := attributes[0].(map[string]interface{})
		So(attribute["first_name"], ShouldEqual, "bar")
		So(attribute["external_id"], ShouldEqual, "foo")
		So(attribute["email"], ShouldEqual, "test@test.com")
		So(attribute["custom_value_attribute"], ShouldEqual, "custom_value_attribute_value")
		So(attribute["custom_value_attribute2"], ShouldEqual, "custom_value_attribute_value2")

		pushTokenAttributes := attribute["push_tokens"].([]interface{})
		pushTokenAttribute := pushTokenAttributes[0].(map[string]interface{})
		So(pushTokenAttribute["app_id"], ShouldEqual, "blah")
		So(pushTokenAttribute["token"], ShouldEqual, "apple-token")

		// Check purchases
		purchases := request["purchases"].([]interface{})
		purchase := purchases[0].(map[string]interface{})
		So(purchase["external_id"], ShouldEqual, "foo")
		So(purchase["product_id"], ShouldEqual, "baz")
		So(purchase["currency"], ShouldEqual, "USD")
		So(purchase["price"], ShouldEqual, 4.92)
		So(purchase["quantity"], ShouldEqual, 1)
		So(purchase["time"], ShouldEqual, "Z070000")

		// Check events
		events := request["events"].([]interface{})
		event := events[0].(map[string]interface{})
		So(event["external_id"], ShouldEqual, "foo")
		So(event["name"], ShouldEqual, "bak")
		So(event["time"], ShouldEqual, "Z070000")
	})

	Convey("Can execute a campaign trigger request", t, func() {
		before()
		defer after()

		// Mock request to app-boy
		var request map[string]interface{}
		MockCampaignTriggerSuccess(func(_request map[string]interface{}) {
			request = _request
		})

		// Construct a push token delete request
		dr := &RawCampaignTriggerRequest{
			AppGroupId: "xxx",
			CampaignId: "yyy",
			Recipients: []RawCampaignRecipient{
				RawCampaignRecipient{
					ExternalId: "4900",
					TriggerProperties: map[string]interface{}{
						"like_count": 31,
					},
				},
			},
		}

		err := RawPostCampaignTriggerRequest(dr)
		checkErr(err)
		So(err, ShouldEqual, nil)

		// Root
		So(request["app_group_id"], ShouldEqual, "xxx")
		So(request["campaign_id"], ShouldEqual, "yyy")

		// Get first recipient
		recs := request["recipients"].([]interface{})
		rec := recs[0].(map[string]interface{})
		So(rec["external_user_id"], ShouldEqual, "4900")
		triggerProperties := rec["trigger_properties"].(map[string]interface{})
		So(triggerProperties["like_count"], ShouldEqual, 31)
	})

	Convey("Gracefully handles error from campaign trigger", t, func() {
		before()
		defer after()

		// Mock request to app-boy
		var request map[string]interface{}
		MockCampaignTriggerFailure(func(_request map[string]interface{}) {
			request = _request
		})

		dr := &RawCampaignTriggerRequest{
			AppGroupId: "xxx",
			CampaignId: "yyy",
			Recipients: []RawCampaignRecipient{
				RawCampaignRecipient{
					ExternalId: "4900",
					TriggerProperties: map[string]interface{}{
						"like_count": 31,
					},
				},
			},
		}

		err := RawPostCampaignTriggerRequest(dr)
		So(err, ShouldNotEqual, nil)
		errStr := fmt.Sprintf("%s", err)
		So(strings.Contains(errStr, "An error message"), ShouldEqual, true)
		So(strings.Contains(errStr, "400"), ShouldEqual, true)
	})

	Convey("Can execute a delete push notification token request", t, func() {
		before()
		defer after()

		// Mock request to app-boy
		var request map[string]interface{}
		MockDeletPushTokenSuccess(func(_request map[string]interface{}) {
			request = _request
		})

		// Construct a push token delete request
		dr := &RawPushTokenDeleteRequest{
			AppGroupId: "foo",
			PushTokens: []RawPushTokenInfo{
				RawPushTokenInfo{
					AppId: "blah",
					Token: "apple-token",
				},
			},
		}

		err := RawPostDeletePushTokenRequest(dr)
		checkErr(err)
		So(err, ShouldEqual, nil)

		// Root
		So(request["app_group_id"], ShouldEqual, "foo")

		pushTokenAttributes := request["push_tokens"].([]interface{})
		pushTokenAttribute := pushTokenAttributes[0].(map[string]interface{})
		So(pushTokenAttribute["app_id"], ShouldEqual, "blah")
		So(pushTokenAttribute["token"], ShouldEqual, "apple-token")
	})

	Convey("Gracefully handles error from push token deletion", t, func() {
		before()
		defer after()

		// Mock request to app-boy
		var request map[string]interface{}
		MockDeletPushTokenFailure(func(_request map[string]interface{}) {
			request = _request
		})

		// Construct a push token delete request
		dr := &RawPushTokenDeleteRequest{
			AppGroupId: "foo",
			PushTokens: []RawPushTokenInfo{
				RawPushTokenInfo{
					AppId: "blah",
					Token: "apple-token",
				},
			},
		}

		err := RawPostDeletePushTokenRequest(dr)
		So(err, ShouldNotEqual, nil)
		errStr := fmt.Sprintf("%s", err)
		So(strings.Contains(errStr, "An error message"), ShouldEqual, true)
		So(strings.Contains(errStr, "400"), ShouldEqual, true)
	})

}
