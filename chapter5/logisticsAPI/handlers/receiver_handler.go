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

// GetReceiver fetch Receiver with a given ID
func (db *DB) GetReceiver(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var receiver data.Receiver
	objectID, _ := primitive.ObjectIDFromHex(vars["id"])
	filter := bson.M{"_id": objectID}
	err := db.collection.FindOne(context.TODO(), filter).Decode(&receiver)

	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(receiver)
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

// AddReceiver add new receiver to our collection DB
func (db *DB) AddReceiver(w http.ResponseWriter, r *http.Request) {
	var receiver data.Receiver
	body, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(body, &receiver)
	result, err := db.collection.InsertOne(context.TODO(), receiver)
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

// UpdateReceiver modified the receiever given the ID
func (db *DB) UpdateReceiver(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var receiver data.Receiver
	putBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(putBody, &receiver)

	objectID, _ := primitive.ObjectIDFromHex(vars["id"])
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": &receiver}
	_, err := db.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Updated successfully"))
	}
}

// RemoveReceiever deleted the receiver from the db collection given the ID
func (db *DB) RemoveReceiever(w http.ResponseWriter, r *http.Request) {
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
