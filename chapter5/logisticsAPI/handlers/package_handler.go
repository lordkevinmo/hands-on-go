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
)

// GetPackage fetch the package infos from the db collection given the id
func (db *DB) GetPackage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var pack data.Package
	objectID, _ := primitive.ObjectIDFromHex(vars["id"])
	filter := bson.M{"_id": objectID}
	err := db.collection.FindOne(context.TODO(), filter).Decode(&pack)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(pack)
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

// AddPackage add new package to the db collection
func (db *DB) AddPackage(w http.ResponseWriter, r *http.Request) {
	var pack data.Package
	postBody, err := ioutil.ReadAll(r.Body)
	json.Unmarshal(postBody, &pack)
	result, err := db.collection.InsertOne(context.TODO(), pack)
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

// UpdatePackage modifies an existing package infos given it id
func (db *DB) UpdatePackage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var pack data.Package
	putBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(putBody, &pack)

	objectID, _ := primitive.ObjectIDFromHex(vars["id"])
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": pack}
	_, err := db.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Updated successfully!"))
	}
}

// RemovePackage deletes a package from a db collection
func (db *DB) RemovePackage(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	objectID, _ := primitive.ObjectIDFromHex(vars["id"])
	filter := bson.M{"_id": objectID}

	_, err := db.collection.DeleteOne(context.TODO(), filter)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write([]byte("Deleted Successfully"))
	}
}
