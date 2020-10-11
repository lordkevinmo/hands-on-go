package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"

	"github.com/gorilla/mux"
	"github.com/lordkevinmo/hands-on-go/chapter5/logisticsAPI/data"
)

// GetCarrier fetch a carrier infos from db collection
func (db *DB) GetCarrier(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var carrier data.Carrier
	objectID, _ := primitive.ObjectIDFromHex(vars["id"])
	filter := bson.M{"_id": objectID}
	err := db.collection.FindOne(context.TODO(), filter).Decode(&carrier)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(carrier)
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

// AddCarrier add a new carrier infos to the db collection
func (db *DB) AddCarrier(w http.ResponseWriter, r *http.Request) {
	var carrier data.Carrier
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &carrier)
	result, err := db.collection.InsertOne(context.TODO(), carrier)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(result)
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

// UpdateCarrier modifies the carrier infos given the id
func (db *DB) UpdateCarrier(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var carrier data.Carrier
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &carrier)

	objectID, _ := primitive.ObjectIDFromHex(vars["id"])
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": carrier}
	_, err := db.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Updated Successfully!"))
	}
}

// RemoveCarrier deletes the carrier infos given the id
func (db *DB) RemoveCarrier(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	objectID, _ := primitive.ObjectIDFromHex(vars["id"])
	filter := bson.M{"_id": objectID}

	_, err := db.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Deleted Successfully!"))
	}
}
