package data

// Carrier holds the carrier information
type Carrier struct {
	ID          interface{} `json:"id" bson:"_id,omitempty"`
	Name        string      `json:"name" bson:"name"`
	CarrierCode int32       `json:"carrier_code" bson:"carrier_code"`
	IsPartner   bool        `json:"is_partner" bson:"is_partner"`
}
