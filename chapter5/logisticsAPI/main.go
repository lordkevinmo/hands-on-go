package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/mux"
	"github.com/lordkevinmo/hands-on-go/chapter5/logisticsAPI/data"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// DB holds the mongo db driver collection
type DB struct {
	collection *mongo.Collection
}

func main() {
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")
	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		panic(err)
	}
	defer client.Disconnect(context.TODO())

	senderCollection := client.Database("logisticsDB").Collection("sender")
	senderDB := &DB{collection: senderCollection}
	receiverCollection := client.Database("logisticsDB").Collection("receiver")
	receiverDB := &DB{collection: receiverCollection}
	carrierCollection := client.Database("logisticsDB").Collection("carrier")
	carrierDB := &DB{collection: carrierCollection}
	packageCollection := client.Database("logisticsDB").Collection("package")
	packageDB := &DB{collection: packageCollection}
	paymentCollection := client.Database("logisticsDB").Collection("payment")
	payDB := &DB{collection: paymentCollection}
	shipmentCollection := client.Database("logisticsDB").Collection("shipment")
	shipDB := &DB{collection: shipmentCollection}

	r := mux.NewRouter()
	// Sender Endpoint
	r.HandleFunc("/v1/sender/{id:[a-zA-Z0-9]*}", senderDB.GetSender).Methods(http.MethodGet)
	r.HandleFunc("/v1/sender", senderDB.AddSender).Methods(http.MethodPost)
	r.HandleFunc("/v1/sender/{id:[a-zA-Z0-9]*}", senderDB.UpdateSender).Methods(http.MethodPut)
	r.HandleFunc("/v1/sender/{id:[a-zA-Z0-9]*}", senderDB.RemoveSender).Methods(http.MethodDelete)
	// Receiver Endpoint
	r.HandleFunc("/v1/receiver/{id:[a-zA-Z0-9]*}", receiverDB.GetReceiver).Methods(http.MethodGet)
	r.HandleFunc("/v1/receiver/", receiverDB.AddReceiver).Methods(http.MethodPost)
	r.HandleFunc("/v1/receiver/{id:[a-zA-Z0-9]*}", receiverDB.UpdateReceiver).Methods(http.MethodPut)
	r.HandleFunc("/v1/receiver/{id:[a-zA-Z0-9]*}", receiverDB.RemoveReceiver).Methods(http.MethodDelete)
	// Carrier Endpoint
	r.HandleFunc("/v1/carrier/{id:[a-zA-Z0-9]*}", carrierDB.GetCarrier).Methods(http.MethodGet)
	r.HandleFunc("/v1/carrier/", carrierDB.AddCarrier).Methods(http.MethodPost)
	r.HandleFunc("/v1/carrier/{id:[a-zA-Z0-9]*}", carrierDB.UpdateCarrier).Methods(http.MethodPut)
	r.HandleFunc("/v1/carrier/{id:[a-zA-Z0-9]*}", carrierDB.RemoveCarrier).Methods(http.MethodDelete)
	// Package Endpoint
	r.HandleFunc("/v1/package/{id:[a-zA-Z0-9]*}", packageDB.GetPackage).Methods(http.MethodGet)
	r.HandleFunc("/v1/package/", packageDB.AddPackage).Methods(http.MethodPost)
	r.HandleFunc("/v1/package/{id:[a-zA-Z0-9]*}", packageDB.UpdatePackage).Methods(http.MethodPut)
	r.HandleFunc("/v1/package/{id:[a-zA-Z0-9]*}", packageDB.RemovePackage).Methods(http.MethodDelete)
	// Payment Endpoint
	r.HandleFunc("/v1/payment/{id:[a-zA-Z0-9]*}", payDB.GetPayment).Methods(http.MethodGet)
	r.HandleFunc("/v1/payment/", payDB.AddPayment).Methods(http.MethodPost)
	r.HandleFunc("/v1/payment/{id:[a-zA-Z0-9]*}", payDB.UpdatePayment).Methods(http.MethodPut)
	// Shipment Endpoint
	r.HandleFunc("/v1/shipment/{id:[a-zA-Z0-9]*}", shipDB.GetShipment).Methods(http.MethodGet)
	r.HandleFunc("/v1/shipment/", shipDB.AddShipment).Methods(http.MethodPost)
	r.HandleFunc("/v1/shipment/{id:[a-zA-Z0-9]*}", shipDB.UpdateShipment).Methods(http.MethodPut)

	srv := &http.Server{
		Handler:      r,
		Addr:         "127.0.0.1:8000",
		WriteTimeout: 15 * time.Second,
		ReadTimeout:  15 * time.Second,
	}
	log.Fatal(srv.ListenAndServe())
}

/********************************** CARRIER ENDPOINT ***************************************/

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

/********************************** PACKAGE ENDPOINT ***************************************/

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

/********************************** PAYMENT ENDPOINT ***************************************/

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

/********************************** RECEIVER ENDPOINT ***************************************/

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

// RemoveReceiver deleted the receiver from the db collection given the ID
func (db *DB) RemoveReceiver(w http.ResponseWriter, r *http.Request) {
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

/********************************** SENDER ENDPOINT ***************************************/

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

/********************************** SHIPMENT ENDPOINT ***************************************/

// GetShipment fetch the shipment infos from the db collection given the ID
func (db *DB) GetShipment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var shipment data.Shipment
	objectID, _ := primitive.ObjectIDFromHex(vars["id"])
	filter := bson.M{"_id": objectID}
	err := db.collection.FindOne(context.TODO(), filter).Decode(&shipment)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "application/json")
		response, _ := json.Marshal(shipment)
		w.WriteHeader(http.StatusOK)
		w.Write(response)
	}
}

//AddShipment add new shipment to the db collection
func (db *DB) AddShipment(w http.ResponseWriter, r *http.Request) {
	var shipment data.Shipment
	postBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(postBody, &shipment)
	result, err := db.collection.InsertOne(context.TODO(), shipment)
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

// UpdateShipment modifies the shipement given the id
func (db *DB) UpdateShipment(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	var shipment data.Shipment
	putBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(putBody, &shipment)

	objectID, _ := primitive.ObjectIDFromHex(vars["id"])
	filter := bson.M{"_id": objectID}
	update := bson.M{"$set": shipment}
	_, err := db.collection.UpdateOne(context.TODO(), filter, update)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(err.Error()))
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.Write([]byte("Updated Successfully!"))
	}
}
