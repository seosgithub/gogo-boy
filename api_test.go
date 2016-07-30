package gogo_boy

import (
	"encoding/json"
	"reflect"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func TestAPI(t *testing.T) {
	before := func() {
		Configure("foo", "blah")
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
		a := NewTrackRequest(externalId)
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

		a := NewTrackRequest("holah")
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
		So(attribute["push_token_import"], ShouldEqual, true)

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

		a := NewTrackRequest("holah")
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
		So(purchase["time"], ShouldEqual, "1969-31-12") // Epoch 0
	})

	Convey("Can execute a track request for a event", t, func() {
		before()
		defer after()

		// Mock request to app-boy
		var request map[string]interface{}
		MockTrackSuccess(func(_request map[string]interface{}) { request = _request })

		a := NewTrackRequest("holah")
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
		So(_event["time"], ShouldEqual, "1969-31-12") // Epoch 0
	})

	Convey("Can marshal a track request", t, func() {
		before()
		defer after()

		// Mock request to app-boy
		var request map[string]interface{}
		MockTrackSuccess(func(_request map[string]interface{}) { request = _request })

		a := NewTrackRequest("holah")
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
		So(req.PurchaseEvents[0].Time, ShouldEqual, "1970-01-01")
		So(req.Events[0].Name, ShouldEqual, "blah")
		So(req.Events[1].Time, ShouldEqual, "1969-31-12")
	})

	Convey("Disabling mocks won't work", t, func() {
		before()
		defer after()

		// Mock request to app-boy
		var request map[string]interface{}
		MockTrackSuccess(func(_request map[string]interface{}) { request = _request })

		StopMocks()

		a := NewTrackRequest("holah")
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

		a := NewTrackRequest("holah")
		a.SetEmail("test@test.com")
		a.SetFirstName("foo")
		a.SetCustomValueAttribute("foo", "bar")

		// This will fail because it hits app boys servers
		err := a.Post()
		So(err, ShouldNotEqual, nil)
	})

}
