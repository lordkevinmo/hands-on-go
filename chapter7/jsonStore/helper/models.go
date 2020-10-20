package helper

import (
	"fmt"

	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const (
	host     = "127.0.0.1"
	port     = 5432
	user     = "lordkevinmo"
	password = "zyh5d117l77fq0a/AAA7eiakjzoJBKAE7NNlgPYGa?dl=0"
	dbname   = "mydb"
)

// Shipment holds the data structure of the shipment
type Shipment struct {
	gorm.Model
	Packages []Package
	Data     string `sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB" json:"-"`
}

// Package holds the data structure of the package
type Package struct {
	gorm.Model
	Data string `sql:"type:JSONB NOT NULL DEFAULT '{}'::JSONB"`
}

// TableName is used to rename the table name since gorm creates the table with plural nums
func (Shipment) TableName() string {
	return "Shipment"
}

// TableName is used here as the same purpose as the previous one
func (Package) TableName() string {
	return "Package"
}

// InitDB is used to init the database
func InitDB() (*gorm.DB, error) {
	var connectionString = fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	var err error

	db, err := gorm.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}
	db.AutoMigrate(&Shipment{}, &Package{})
	return db, nil
}
