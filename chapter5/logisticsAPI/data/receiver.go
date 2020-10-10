package data

// Receiver holds the user who received the shipment information
type Receiver struct {
	ID        interface{} `json:"id" bson:"_id,omitempty"`
	FirstName string      `json:"first_name" bson:"first_name"`
	LastName  string      `json:"last_name" bson:"last_name"`
	Address   Address     `json:"address" bson:"address"`
	Phone     string      `json:"phone" bson:"phone"`
}
