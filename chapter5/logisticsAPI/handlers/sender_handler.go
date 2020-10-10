package handlers

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lordkevinmo/hands-on-go/chapter5/logisticsAPI/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// DB holds the mongo db driver collection
type DB struct {
	collection *mongo.Collection
}

// GetSender fetch sender with a given ID
func (db *DB) GetSender(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var sender data.Sender
	objectID, _ := primitive.ObjectIDFromHex(vars["id"])
	filter := bson.M{"_id": objectID}
	err := db.collection.FindOne(context.TODO(), filter).Decode(&sender)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(sender)
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

// AddSender add new sender to our collection DB
func (db *DB) AddSender(w http.ResponseWriter, r *http.Request) {
	var sender data.Sender
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &sender)
	result, err := db.collection.InsertOne(context.TODO(), sender)
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

// UpdateSender modified the sender given the ID
func (db *DB) UpdateSender(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var sender data.Sender
	putBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(putBody, &sender)

	objectID, _ := primitive.ObjectIDFromHex(vars["id"])
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": &sender}
	_, err := db.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Updated successfully"))
	}
}

// RemoveSender deleted the sender from the db collection given the ID
func (db *DB) RemoveSender(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	objectID, _ := primitive.ObjectIDFromHex(vars["id"])
	filter := bson.M{"_id": objectID}

	_, err := db.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Deleted Successfuly"))
	}
}
