package data

// Package holds the information about the object
type Package struct {
	ID         interface{}      `json:"id" bson:"_id,omitempty"`
	Dimensions packageDimension `json:"dimensions" bson:"dimensions"`
	Weight     uint32           `json:"weight" bson:"weight"`
	IsDamaged  bool             `json:"is_damaged" bson:"is_damaged"`
	Status     string           `json:"status" bson:"status"`
}

type packageDimension struct {
	Width  uint32 `json:"width" bson:"width"`
	Height uint32 `json:"height" bson:"height"`
}
