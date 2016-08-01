package gogo_boy

import (
	"fmt"
	"time"
)

/*
	----------------------------------------------------------------------
  Much prettier API interfaces that build on raw_api
	----------------------------------------------------------------------
*/

type ConfigureInfo struct {
	AppId      string
	AppGroupId string
}

var appId string
var appGroupId string

func Configure(_appGroupId string, _appId string) {
	appGroupId = _appGroupId
	appId = _appId
}

type TrackRequest struct {
	AppGroupId string
	AppId      string
	ExternalId string

	// Unlike the RawTrackRequest, we're only
	// considering one user which will only ever
	// use one attribute
	Attributes          map[string]interface{}
	PushTokenAttributes []string

	PurchaseEvents []*PurchaseEvent
	Events         []*Event
}

func NewTrackRequest(externalId string) *TrackRequest {
	return &TrackRequest{
		AppGroupId:          appGroupId,
		AppId:               appId,
		Attributes:          map[string]interface{}{},
		PushTokenAttributes: []string{},
		ExternalId:          externalId,
	}
}

func (tr *TrackRequest) SetFirstName(name string) {
	tr.Attributes["first_name"] = name
}

func (tr *TrackRequest) SetEmail(email string) {
	tr.Attributes["email"] = email
}

func (tr *TrackRequest) AddEvent(event interface{}) {
	switch event := event.(type) {
	case *PurchaseEvent:
		tr._addPurchaseEvent(event)
	case *Event:
		tr._addEvent(event)
	default:
		panic(fmt.Errorf("Unknow event type %T", event))
	}
}

func (tr *TrackRequest) _addPurchaseEvent(event *PurchaseEvent) {
	tr.PurchaseEvents = append(tr.PurchaseEvents, event)
}

func (tr *TrackRequest) _addEvent(event *Event) {
	tr.Events = append(tr.Events, event)
}

func (tr *TrackRequest) AddPushToken(token string) {
	tr.PushTokenAttributes = append(tr.PushTokenAttributes, token)
}

func (tr *TrackRequest) SetCustomValueAttribute(name string, value interface{}) {
	tr.Attributes[name] = value
}

func (tr *TrackRequest) Post() error {
	rt := &RawTrackRequest{
		AppGroupId: tr.AppGroupId,
		Attributes: []RawAttributesInfo{
			RawAttributesInfo{
				CustomAttributes: map[string]interface{}{},
				PushTokens:       []RawPushTokenInfo{},
			},
		},
		Purchases: []RawPurchaseInfo{},
	}

	rt.Attributes[0].ExternalId = tr.ExternalId

	for k, v := range tr.Attributes {
		switch k {
		case "first_name":
			rt.Attributes[0].FirstName = v.(string)
		case "email":
			rt.Attributes[0].Email = v.(string)
		default:
			rt.Attributes[0].CustomAttributes[k] = v
		}
	}

	for _, pt := range tr.PushTokenAttributes {
		rt.Attributes[0].PushTokenImport = true
		rt.Attributes[0].PushTokens = append(rt.Attributes[0].PushTokens, RawPushTokenInfo{
			Token: pt,
			AppId: tr.AppId,
		})
	}

	for _, pt := range tr.PurchaseEvents {
		rpi := pt.RawPurchaseInfo
		rpi.ExternalId = tr.ExternalId
		rt.Purchases = append(rt.Purchases, rpi)
	}

	for _, pt := range tr.Events {
		rpi := pt.RawEventInfo
		rt.Events = append(rt.Events, rpi)
	}

	return RawPostTrackRequest(rt)
}

type PurchaseEvent struct {
	RawPurchaseInfo
}

func NewPurchaseEvent() *PurchaseEvent {
	return &PurchaseEvent{}
}

func (e *PurchaseEvent) SetProductId(val string) {
	e.ProductId = val
}

func (e *PurchaseEvent) SetCurrencyUSD() {
	e.Currency = "USD"
}

func (e *PurchaseEvent) SetPrice(val float32) {
	e.Price = val
}

func (e *PurchaseEvent) SetQuantity(val int) {
	e.Quantity = val
}

func (e *PurchaseEvent) SetTime(time time.Time) {
	e.Time = time.Format("2006-02-01")
}

type Event struct {
	RawEventInfo
}

func NewEvent() *Event {
	return &Event{}
}

func (e *Event) SetName(val string) {
	e.Name = val
}

func (e *Event) SetTime(time time.Time) {
	e.Time = time.Format("2006-02-01")
}
