package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      Storage
}

func NewAPIServer(listenAddr string, store Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/weather", makeHTTPHandleFunc(s.handleWeatherReport))

	router.HandleFunc("/weather/{id}", makeHTTPHandleFunc(s.handleGetWeatherReportById))

	log.Println("Json API running om port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

// Catch all function for CRUD Operations
func (s *APIServer) handleWeatherReport(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetWeatherReports(w, r)
	}
	if r.Method == "POST" {
		return s.handleCreateWeatherReport(w, r)
	}
	if r.Method == "UPDATE" {
		return s.handleCreateWeatherReport(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteWeatherReport(w, r)
	}
	return fmt.Errorf("method not allowed %s", r.Method)
}

// Get Weather Reports
func (s *APIServer) handleGetWeatherReports(w http.ResponseWriter, r *http.Request) error {
	weatherReport, err := s.store.GetWeatherReports()
	if err != nil {
		return err
	}

	return writeJson(w, http.StatusOK, weatherReport)
}

// Get singular weather report based on id
func (s *APIServer) handleGetWeatherReportById(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	println(id)
	weatherReport := NewWeatherReport("Sunny!", 2, 0.1)
	return writeJson(w, http.StatusOK, weatherReport)
}

func (s *APIServer) handleCreateWeatherReport(w http.ResponseWriter, r *http.Request) error {
	createRequest := new(CreateWeatherReportRequest)
	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		return err
	}
	WeatherReport := NewWeatherReport(createRequest.Description, createRequest.Temperature, createRequest.RainChance)
	if err := s.store.CreateWeatherReport(WeatherReport); err != nil {
		return err
	}

	return writeJson(w, http.StatusOK, createRequest)
}

func (s *APIServer) handleUpdateWeatherReport(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteWeatherReport(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func writeJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunction func(http.ResponseWriter, *http.Request) error

type apiError struct {
	Error string
}

// Acts as a wrapper function for HTTP calls
func makeHTTPHandleFunc(f apiFunction) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if err := f(w, r); err != nil {
			//Handle error
			writeJson(w, http.StatusBadRequest, apiError{Error: err.Error()})
		}
	}
}
