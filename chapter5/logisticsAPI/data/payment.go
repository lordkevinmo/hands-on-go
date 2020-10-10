package data

import "time"

// Payment holds the information related to the receiver
type Payment struct {
	ID             interface{}    `json:"id" bson:"_id,omitempty"`
	InitiatedOn    time.Time      `json:"initiated_on" bson:"initiated_on"`
	SuccessfulOn   time.Time      `json:"successful_on" bson:"successful_on"`
	MerchantID     interface{}    `json:"merchant_id" bson:"merchant_id"`
	ModeOfPayment  string         `json:"mode_of_payment" bson:"mode_of_payment"`
	PaymentDetails paymentDetails `json:"payment_details" bson:"payment_details"`
}

// paymentdetails holds the payment details information
type paymentDetails struct {
	TransactionToken string `json:"transaction_token" bson:"transaction_token"`
}
