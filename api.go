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

type Client struct {
	appGroupId string
}

type AppClient struct {
	*Client

	appId string
}

func NewClient(appGroupId string) *Client {
	client := &Client{
		appGroupId: appGroupId,
	}

	return client
}

type TrackRequest struct {
	AppGroupId string
	AppId      string
	ExternalId string

	// Unlike the RawTrackRequest, we're only
	// considering one user which will only ever
	// use one attribute
	Attributes                map[string]interface{}
	PushTokenAttributes       []string
	DeletePushTokenAttributes []string // Tokens that you want deleted

	PurchaseEvents []*PurchaseEvent
	Events         []*Event
}

type CampaignTriggerRequest struct {
	AppGroupId string
	CampaignId string

	Recipients []RawCampaignRecipient
}

func (c *Client) NewAppClient(appId string) *AppClient {
	return &AppClient{
		Client: c,
		appId:  appId,
	}
}

func (c *AppClient) NewTrackRequest(externalId string) *TrackRequest {
	return &TrackRequest{
		AppGroupId:          c.appGroupId,
		AppId:               c.appId,
		Attributes:          map[string]interface{}{},
		PushTokenAttributes: []string{},
		ExternalId:          externalId,
	}
}

func (c *Client) NewCampaignTriggerRequest(campaignId string) *CampaignTriggerRequest {
	return &CampaignTriggerRequest{
		AppGroupId: c.appGroupId,
		CampaignId: campaignId,
		Recipients: []RawCampaignRecipient{},
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

// This will tell app-boy to remove the push token included
func (tr *TrackRequest) RemovePushToken(token string) {
	tr.DeletePushTokenAttributes = append(tr.DeletePushTokenAttributes, token)
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
		//rt.Attributes[0].PushTokenImport = true
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
		rpi.ExternalId = tr.ExternalId
		rt.Events = append(rt.Events, rpi)
	}

	// Run the regular track requests first
	err := RawPostTrackRequest(rt)
	if err != nil {
		return err
	}

	// Now run the push token deletions if there are any
	if len(tr.DeletePushTokenAttributes) > 0 {
		dr := &RawPushTokenDeleteRequest{
			PushTokens: []RawPushTokenInfo{},
		}

		for _, pt := range tr.DeletePushTokenAttributes {
			dr.PushTokens = append(dr.PushTokens, RawPushTokenInfo{
				Token: pt,
				AppId: tr.AppId,
			})
		}
		err = RawPostDeletePushTokenRequest(dr)
		if err != nil {
			return err
		}
	}

	return nil
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
	e.Time = time.UTC().Format("2006-01-02T15:04:05")
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
	e.Time = time.UTC().Format("2006-01-02T15:04:05")
}

func (ctr *CampaignTriggerRequest) addRecipient(externalId string, triggerProperties map[string]interface{}) {
	rec := RawCampaignRecipient{
		ExternalId:        externalId,
		TriggerProperties: triggerProperties,
	}

	ctr.Recipients = append(ctr.Recipients, rec)
}

func (ctr *CampaignTriggerRequest) Post() error {
	rt := &RawCampaignTriggerRequest{
		AppGroupId: ctr.AppGroupId,
		CampaignId: ctr.CampaignId,
		Recipients: ctr.Recipients,
	}

	if lr := len(rt.Recipients); lr > 50 {
		return fmt.Errorf("Tried to post a CampaignTriggerRequest for the [AppBoyCampaign](campaign_id: %s) but there were %d recipients which exceeds the maximum of 50 per request.  You will need to break your campaign trigger requests up into multiple requests in order to send more than 50 recipients", rt.CampaignId, lr)
	}

	err := RawPostCampaignTriggerRequest(rt)
	return err
}
