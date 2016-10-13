package gogo_boy

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAPI(t *testing.T) {
	var client *Client
	var appClient *AppClient
	before := func() {
		client = NewClient("foo")
		appClient = client.NewAppClient("blah")
	}

	after := func() {
		StopMocks()
	}

	Convey("Can execute a track request to app-boy", t, func() {
		before()
		defer after()

		// Mock request to app-boy
		var request map[string]interface{}
		MockTrackSuccess(func(_request map[string]interface{}) { request = _request })

		externalId := "holah"
		a := appClient.NewTrackRequest(externalId)
		err := a.Post()

		err = a.Post()
		checkErr(err)
		So(err, ShouldEqual, nil)

		// Root
		So(request["app_group_id"], ShouldEqual, "foo")

		attributes := request["attributes"].([]interface{})
		attribute := attributes[0].(map[string]interface{})
		So(attribute["push_token_import"], ShouldEqual, nil)
		So(attribute["external_id"], ShouldEqual, externalId)
	})

	Convey("Can execute a track request to app-boy with attributes", t, func() {
		before()
		defer after()

		// Mock request to app-boy
		var request map[string]interface{}
		MockTrackSuccess(func(_request map[string]interface{}) { request = _request })

		a := appClient.NewTrackRequest("holah")
		a.SetFirstName("foo")
		a.SetEmail("test@test.com")
		a.SetCustomValueAttribute("baz", "bar")
		a.AddPushToken("apple-token")
		err := a.Post()

		err = a.Post()
		checkErr(err)
		So(err, ShouldEqual, nil)

		// Check attributes
		attributes := request["attributes"].([]interface{})
		attribute := attributes[0].(map[string]interface{})
		So(attribute["first_name"], ShouldEqual, "foo")
		So(attribute["email"], ShouldEqual, "test@test.com")
		So(attribute["baz"], ShouldEqual, "bar")
		//So(attribute["push_token_import"], ShouldEqual, true)

		pushTokenAttributes := attribute["push_tokens"].([]interface{})
		pushTokenAttribute := pushTokenAttributes[0].(map[string]interface{})
		So(pushTokenAttribute["app_id"], ShouldEqual, "blah")
		So(pushTokenAttribute["token"], ShouldEqual, "apple-token")
	})

	Convey("Can execute a track request for a purchase", t, func() {
		before()
		defer after()

		// Mock request to app-boy
		var request map[string]interface{}
		MockTrackSuccess(func(_request map[string]interface{}) { request = _request })

		a := appClient.NewTrackRequest("holah")
		pEvent := NewPurchaseEvent()
		pEvent.SetProductId("blah")
		pEvent.SetCurrencyUSD()
		pEvent.SetPrice(4.29)
		pEvent.SetQuantity(1)
		pEvent.SetTime(time.Unix(0, 0))
		a.AddEvent(pEvent)

		err := a.Post()
		checkErr(err)
		So(err, ShouldEqual, nil)

		// Check purchases
		purchases := request["purchases"].([]interface{})
		purchase := purchases[0].(map[string]interface{})
		So(purchase["external_id"], ShouldEqual, "holah")
		So(purchase["product_id"], ShouldEqual, "blah")
		So(purchase["currency"], ShouldEqual, "USD")
		So(purchase["price"], ShouldEqual, 4.29)
		So(purchase["quantity"], ShouldEqual, 1)
		So(purchase["time"], ShouldEqual, "1970-01-01T00:00:00") // Epoch 0
	})

	Convey("Can execute a track request for a event", t, func() {
		before()
		defer after()

		// Mock request to app-boy
		var request map[string]interface{}
		MockTrackSuccess(func(_request map[string]interface{}) { request = _request })

		a := appClient.NewTrackRequest("holah")
		event := NewEvent()
		event.SetName("blah")
		event.SetTime(time.Unix(0, 0))
		a.AddEvent(event)

		eventB := NewEvent()
		eventB.SetName("foo")
		eventB.SetTime(time.Unix(0, 0))
		a.AddEvent(eventB)

		err := a.Post()
		checkErr(err)
		So(err, ShouldEqual, nil)

		// Check events
		events := request["events"].([]interface{})
		So(len(events), ShouldEqual, 2)
		_event := events[0].(map[string]interface{})
		So(_event["name"], ShouldEqual, "blah")
		So(_event["time"], ShouldEqual, "1970-01-01T00:00:00") // Epoch 0
		So(_event["external_id"], ShouldEqual, "holah")
	})

	Convey("Can marshal a track request", t, func() {
		before()
		defer after()

		// Mock request to app-boy
		var request map[string]interface{}
		MockTrackSuccess(func(_request map[string]interface{}) { request = _request })

		a := appClient.NewTrackRequest("holah")
		a.SetEmail("test@test.com")
		a.SetFirstName("foo")
		a.SetCustomValueAttribute("foo", "bar")
		event := NewEvent()
		event.SetName("blah")
		event.SetTime(time.Unix(86399, 0))
		a.AddEvent(event)

		eventB := NewEvent()
		eventB.SetName("foo")
		eventB.SetTime(time.Unix(0, 0))
		a.AddEvent(eventB)

		eventC := NewPurchaseEvent()
		eventC.SetTime(time.Unix(86400, 0))
		eventC.SetQuantity(1)
		eventC.SetProductId("foo")
		eventC.SetPrice(1)
		eventC.SetCurrencyUSD()
		a.AddEvent(eventC)

		checkErr(a.Post())
		requestA := request
		request = nil

		res, err := json.Marshal(a)
		checkErr(err)

		var req TrackRequest
		err = json.Unmarshal(res, &req)
		checkErr(err)

		checkErr(req.Post())
		requestB := request

		So(requestA, ShouldNotEqual, nil)
		So(reflect.DeepEqual(requestA, requestB), ShouldEqual, true)

		So(len(req.Attributes), ShouldEqual, 3)
		So(len(req.Events), ShouldEqual, 2)
		So(len(req.PurchaseEvents), ShouldEqual, 1)
		So(req.PurchaseEvents[0].Price, ShouldEqual, 1)
		So(req.PurchaseEvents[0].Time, ShouldEqual, "1970-01-02T00:00:00")
		So(req.Events[0].Name, ShouldEqual, "blah")
		So(req.Events[0].Time, ShouldEqual, "1970-01-01T23:59:59")
	})

	Convey("Disabling mocks won't work", t, func() {
		before()
		defer after()

		// Mock request to app-boy
		var request map[string]interface{}
		MockTrackSuccess(func(_request map[string]interface{}) { request = _request })

		StopMocks()

		a := appClient.NewTrackRequest("holah")
		a.SetEmail("test@test.com")
		a.SetFirstName("foo")
		a.SetCustomValueAttribute("foo", "bar")

		// This will fail because it hits app boys servers
		err := a.Post()
		So(err, ShouldNotEqual, nil)
	})

	Convey("Using the failure mock will result in an error", t, func() {
		before()
		defer after()

		// Mock request to app-boy
		var request map[string]interface{}
		MockTrackFailure(func(_request map[string]interface{}) { request = _request })

		StopMocks()

		a := appClient.NewTrackRequest("holah")
		a.SetEmail("test@test.com")
		a.SetFirstName("foo")
		a.SetCustomValueAttribute("foo", "bar")

		// This will fail because it hits app boys servers
		err := a.Post()
		So(err, ShouldNotEqual, nil)
	})

	Convey("Can execute a track request to app-boy with attributes and delete push token request", t, func() {
		before()
		defer after()

		// Mock request to app-boy
		var trackRequest map[string]interface{}
		var deletePushTokenRequest map[string]interface{}
		MockTrackSuccess(func(_request map[string]interface{}) { trackRequest = _request })
		MockDeletPushTokenSuccess(func(_request map[string]interface{}) { deletePushTokenRequest = _request })

		a := appClient.NewTrackRequest("holah")
		a.SetFirstName("foo")
		a.SetEmail("test@test.com")
		a.SetCustomValueAttribute("baz", "bar")
		a.AddPushToken("apple-token")
		a.RemovePushToken("apple-token2")
		err := a.Post()

		err = a.Post()
		checkErr(err)
		So(err, ShouldEqual, nil)

		// Check attributes
		attributes := trackRequest["attributes"].([]interface{})
		attribute := attributes[0].(map[string]interface{})
		So(attribute["first_name"], ShouldEqual, "foo")
		So(attribute["email"], ShouldEqual, "test@test.com")
		So(attribute["baz"], ShouldEqual, "bar")
		//So(attribute["push_token_import"], ShouldEqual, true)

		pushTokenAttributes := attribute["push_tokens"].([]interface{})
		pushTokenAttribute := pushTokenAttributes[0].(map[string]interface{})
		So(pushTokenAttribute["app_id"], ShouldEqual, "blah")
		So(pushTokenAttribute["token"], ShouldEqual, "apple-token")

		// Delete push tokens
		deletePushTokenAttributes := deletePushTokenRequest["push_tokens"].([]interface{})
		deletePushTokenAttribute := deletePushTokenAttributes[0].(map[string]interface{})
		So(deletePushTokenAttribute["app_id"], ShouldEqual, "blah")
		So(deletePushTokenAttribute["token"], ShouldEqual, "apple-token2")

	})

	Convey("Can execute a campaign trigger", t, func() {
		before()
		defer after()

		// Mock request to app-boy
		var request map[string]interface{}
		MockCampaignTriggerSuccess(func(_request map[string]interface{}) { request = _request })

		a := client.NewCampaignTriggerRequest("my-campaign-id")

		// Batch 50 because it's an API upper limit before it needs to be split
		// into multiple requests
		for i := 0; i < 50; i++ {
			a.addRecipient("4900", map[string]interface{}{
				"like_count": 31,
			})
		}
		err := a.Post()

		checkErr(err)
		So(err, ShouldEqual, nil)

		// Check attributes
		appGroupId := request["app_group_id"].(string)
		campaignId := request["campaign_id"].(string)
		recs := request["recipients"].([]interface{})
		rec := recs[0].(map[string]interface{})

		triggerProperties := rec["trigger_properties"].(map[string]interface{})
		likeCount := int(triggerProperties["like_count"].(float64))

		So(appGroupId, ShouldEqual, "foo")
		So(campaignId, ShouldEqual, "my-campaign-id")
		So(len(recs), ShouldEqual, 50)
		So(rec["external_user_id"], ShouldEqual, "4900")
		So(likeCount, ShouldEqual, 31)
	})

	Convey("Does error if you go over 50 recipients", t, func() {
		before()
		defer after()

		// Mock request to app-boy
		var request map[string]interface{}
		MockCampaignTriggerSuccess(func(_request map[string]interface{}) { request = _request })

		a := client.NewCampaignTriggerRequest("my-campaign-id")

		// Batch 51 because it's over the api limit of 50
		for i := 0; i < 51; i++ {
			a.addRecipient("4900", map[string]interface{}{
				"like_count": 31,
			})
		}
		// This will fail because we have over 50 recipients
		err := a.Post()
		So(err, ShouldNotEqual, nil)
		So(strings.Contains(fmt.Sprintf("%s", err), "50"), ShouldEqual, true)
		So(strings.Contains(fmt.Sprintf("%s", err), "51"), ShouldEqual, true)
		So(strings.Contains(fmt.Sprintf("%s", err), "my-campaign-id"), ShouldEqual, true)
	})

	Convey("Does batch campaign trigger into two requests for 51 entries", t, func() {
		before()
		defer after()

		// Mock request to app-boy
		var request map[string]interface{}
		MockCampaignTriggerSuccess(func(_request map[string]interface{}) { request = _request })

		a := client.NewCampaignTriggerRequest("my-campaign-id")
		a.addRecipient("4900", map[string]interface{}{
			"like_count": 31,
		})
		a.addRecipient("3333", map[string]interface{}{
			"like_count": 33,
		})
		err := a.Post()

		checkErr(err)
		So(err, ShouldEqual, nil)

		// Check attributes
		appGroupId := request["app_group_id"].(string)
		campaignId := request["campaign_id"].(string)
		recs := request["recipients"].([]interface{})
		rec := recs[0].(map[string]interface{})

		triggerProperties := rec["trigger_properties"].(map[string]interface{})
		likeCount := int(triggerProperties["like_count"].(float64))

		So(appGroupId, ShouldEqual, "foo")
		So(campaignId, ShouldEqual, "my-campaign-id")
		So(len(recs), ShouldEqual, 2)
		So(rec["external_user_id"], ShouldEqual, "4900")
		So(likeCount, ShouldEqual, 31)
	})

}
