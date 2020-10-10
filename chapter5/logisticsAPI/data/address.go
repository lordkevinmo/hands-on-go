package data

// Address holds the user address informaton
type Address struct {
	AddressType string `json:"addressType" bson:"addressType"`
	Street      string `json:"street" bson:"street"`
	City        string `json:"city" bson:"city"`
	State       string `json:"state" bson:"state"`
	PinCode     uint16 `json:"pinCode" bson:"pinCode"`
	Country     string `json:"country" bson:"country"`
}
