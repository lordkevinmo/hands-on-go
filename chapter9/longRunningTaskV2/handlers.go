package main

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"github.com/go-redis/redis"
	"github.com/lordkevinmo/hands-on-go/chapter9/longRunningTaskV1/models"

	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

// JobServer holds handler function
type JobServer struct {
	Queue       amqp.Queue
	Channel     *amqp.Channel
	Conn        *amqp.Connection
	redisClient *redis.Client
}

func (s *JobServer) publish(jsonBody []byte) error {
	message := amqp.Publishing{
		ContentType: "application/json",
		Body:        jsonBody,
	}

	err := s.Channel.Publish(
		"",
		queueName,
		false,
		false,
		message,
	)

	handleError(err, "Error while generating UUID")
	return err
}

func (s *JobServer) asyncDBHandler(w http.ResponseWriter, r *http.Request) {
	jobID, err := uuid.NewRandom()
	queryParams := r.URL.Query()
	// Ex: client_time: 1569174071
	unixTime, err := strconv.ParseInt(queryParams.Get("client_time"), 10, 64)
	clientTime := time.Unix(unixTime, 0)
	handleError(err, "Error while converting client time")

	jsonBody, err := json.Marshal(models.Job{ID: jobID,
		Type:      "A",
		ExtraData: models.Log{ClientTime: clientTime},
	})
	handleError(err, "JSON Body creation failed")

	if s.publish(jsonBody) == nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonBody)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *JobServer) asyncCallbackHandler(w http.ResponseWriter, r *http.Request) {
	jobID, err := uuid.NewRandom()

	jsonBody, err := json.Marshal(models.Job{ID: jobID,
		Type:      "B",
		ExtraData: "",
	})
	handleError(err, "JSON body creation failed")

	if s.publish(jsonBody) == nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonBody)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *JobServer) asyncMailHandler(w http.ResponseWriter, r *http.Request) {
	jobID, err := uuid.NewRandom()

	jsonBody, err := json.Marshal(models.Job{ID: jobID,
		Type:      "C",
		ExtraData: "",
	})
	handleError(err, "JSON body creation failed")

	err = s.publish(jsonBody)

	if err == nil {
		w.WriteHeader(http.StatusOK)
		w.Header().Set("Content-Type", "application/json")
		w.Write(jsonBody)
	} else {
		w.WriteHeader(http.StatusInternalServerError)
	}
}

func (s *JobServer) statusHandler(w http.ResponseWriter, r *http.Request) {
	queryParams := r.URL.Query()
	// fetch UUID from query
	uuid := queryParams.Get("uuid")
	w.Header().Set("Content-Type", "application/json")
	jobStatus := s.redisClient.Get(uuid)
	status := map[string]string{"uuid": uuid, "status": jobStatus.Val()}
	response, err := json.Marshal(status)
	handleError(err, "Cannot create response for client")
	w.Write(response)
}
