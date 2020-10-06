package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"time"

	"github.com/emicklei/go-restful"
	"github.com/lordkevinmo/hands-on-go/src/chapter4/dbutils"
	_ "github.com/mattn/go-sqlite3"
)

// DB driver visible to whole program
var DB *sql.DB

// TrainResource is the model for holding rail information
type TrainResource struct {
	ID              int
	DriverName      string
	OperatingStatus bool
}

// StationResource holds information about locations
type StationResource struct {
	ID          int
	Name        string
	OpeningTime time.Time
	ClosingTime time.Time
}

// ScheduleResource links both trains and stations
type ScheduleResource struct {
	ID          int
	TrainID     int
	StationID   int
	ArrivalTime time.Time
}

// Register adds paths and routes to a new service instance for train
func (t *TrainResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/v1/train").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/{train-id}").To(t.getTrain))
	ws.Route(ws.POST("").To(t.createTrain))
	ws.Route(ws.DELETE("/{train-id}").To(t.removeTrain))

	container.Add(ws)
}

// Register adds paths and routes to a new service instance for station
func (s *StationResource) Register(container *restful.Container) {
	ws := new(restful.WebService)
	ws.Path("/v1/station").Consumes(restful.MIME_JSON).Produces(restful.MIME_JSON)

	ws.Route(ws.GET("/{station-id}").To(s.getStation))
	ws.Route(ws.POST("").To(s.createStation))
	ws.Route(ws.DELETE("/{station-id}").To(s.removeStation))

	container.Add(ws)
}

// GET http://localhost:8000/v1/station/1
func (s StationResource) getStation(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("station-id")
	err := DB.QueryRow("SELECT ID, NAME, OPENING_TIME, CLOSING_TIME FROM station where id=?", id).Scan(&s.ID, &s.Name, &s.OpeningTime, &s.ClosingTime)
	if err != nil {
		log.Println(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Station could not be found")
	} else {
		response.WriteEntity(s)
	}
}

// POST http://localhost:8000/v1/station
func (s StationResource) createStation(request *restful.Request, response *restful.Response) {
	log.Println(request.Request.Body)
	decoder := json.NewDecoder(request.Request.Body)
	var tempStation StationResource
	err := decoder.Decode(&tempStation)
	log.Println(tempStation.Name)
	if err != nil {
		log.Println(err)
	}
	statement, _ := DB.Prepare("INSERT INTO station (NAME, OPENING_TIME, CLOSING_TIME) values (?, ?, ?)")
	result, err := statement.Exec(tempStation.Name, tempStation.OpeningTime, tempStation.ClosingTime)
	if err == nil {
		newID, _ := result.LastInsertId()
		tempStation.ID = int(newID)
		response.WriteHeaderAndEntity(http.StatusCreated, tempStation)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

// DELETE http://localhost:8000/station/1
func (s StationResource) removeStation(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("station-id")
	statement, _ := DB.Prepare("DELETE FROM station where id=?")
	_, err := statement.Exec(id)
	if err == nil {
		response.WriteHeader(http.StatusOK)
	} else {
		response.AddHeader("Content-Type", "application/json")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

// GET http://localhost:8000/v1/trains/1
func (t TrainResource) getTrain(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("train-id")
	err := DB.QueryRow("SELECT ID, DRIVER_NAME, OPERATING_STATUS FROM train where id=?", id).Scan(&t.ID, &t.DriverName, &t.OperatingStatus)
	if err != nil {
		log.Println(err)
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusNotFound, "Train could not be found.")
	} else {
		response.WriteEntity(t)
	}
}

// POST http://localhost:8000/v1/trains
func (t TrainResource) createTrain(request *restful.Request, response *restful.Response) {
	log.Println(request.Request.Body)
	decoder := json.NewDecoder(request.Request.Body)
	var b TrainResource
	err := decoder.Decode(&b)
	log.Println(b.DriverName, b.OperatingStatus)
	if err != nil {
		log.Println(err)
	}
	statement, _ := DB.Prepare("insert into train (DRIVER_NAME, OPERATING_STATUS) values (?, ?)")
	result, err := statement.Exec(b.DriverName, b.OperatingStatus)
	if err == nil {
		newID, _ := result.LastInsertId()
		b.ID = int(newID)
		response.WriteHeaderAndEntity(http.StatusCreated, b)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

// DELETE http://localhost:8000/v1/trains/1
func (t TrainResource) removeTrain(request *restful.Request, response *restful.Response) {
	id := request.PathParameter("train-id")
	statement, _ := DB.Prepare("delete from train where id=?")
	_, err := statement.Exec(id)
	if err == nil {
		response.WriteHeader(http.StatusOK)
	} else {
		response.AddHeader("Content-Type", "text/plain")
		response.WriteErrorString(http.StatusInternalServerError, err.Error())
	}
}

func main() {
	log.Println("Time", time.Now())
	var err error
	// Connect to database
	DB, err = sql.Open("sqlite3", "./railAPI.db")
	if err != nil {
		log.Println("Driver Creation failed!")
	}
	// Create tables
	dbutils.Initialize(DB)
	// Instantiate container
	wsContainer := restful.NewContainer()
	wsContainer.Router(restful.CurlyRouter{})
	// train endpoint processing
	t := TrainResource{}
	t.Register(wsContainer)
	// station endpoint processing
	s := StationResource{}
	s.Register(wsContainer)

	log.Printf("Start listening on localhost:8000")
	server := &http.Server{
		Addr:    ":8000",
		Handler: wsContainer,
	}
	log.Fatal(server.ListenAndServe())
}
