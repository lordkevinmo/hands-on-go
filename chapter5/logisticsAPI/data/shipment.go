package data

import "time"

// Shipment holds all the information about the shipment and the logistics
type Shipment struct {
	ID         interface{} `json:"id" bson:"_id,omitempty"`
	Sender     interface{} `json:"sender" bson:"sender"`
	Receiver   interface{} `json:"receiver" bson:"receiver"`
	Package    interface{} `json:"package" bson:"package"`
	Payment    interface{} `json:"payment" bson:"payment"`
	Carrier    interface{} `json:"carrier" bson:"carrier"`
	PromisedOn time.Time   `json:"promised_on" bson:"promised_on"`
}
