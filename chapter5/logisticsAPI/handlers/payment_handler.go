package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gorilla/mux"
	"github.com/lordkevinmo/hands-on-go/chapter5/logisticsAPI/data"
)

// GetPayment fetch the payment info from DB using the ID
func (db *DB) GetPayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var payment data.Payment
	objectID, pError := primitive.ObjectIDFromHex(vars["id"])
	if pError != nil {
		panic(pError)
	}
	filter := bson.M{"_id": objectID}
	err := db.collection.FindOne(context.TODO(), filter).Decode(&payment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		response, jError := json.Marshal(payment)
		if jError != nil {
			panic(jError)
		}
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

// AddPayment add the payment to the db collection
func (db *DB) AddPayment(w http.ResponseWriter, r *http.Request) {
	var payment data.Payment
	body, readError := ioutil.ReadAll(r.Body)
	if readError != nil {
		panic(readError)
	}
	json.Unmarshal(body, &payment)
	result, err := db.collection.InsertOne(context.TODO(), payment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		response, jError := json.Marshal(result)
		if jError != nil {
			log.Println("Error when marshalling JSON")
		}
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

// UpdatePayment modifies the payment info from the db collection
func (db *DB) UpdatePayment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var payment data.Payment
	body, readError := ioutil.ReadAll(r.Body)
	if readError != nil {
		panic(readError)
	}
	json.Unmarshal(body, &payment)

	objectID, _ := primitive.ObjectIDFromHex(vars["id"])
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": &payment}
	_, err := db.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Updated Successfully!"))
	}
}
