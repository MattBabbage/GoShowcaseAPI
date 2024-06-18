package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func writeJson(w http.ResponseWriter, status int, v any) error {
	w.WriteHeader(status)
	w.Header().Set("Content-Type", "application/json")
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

type APIServer struct {
	listenAddr string
}

func NewAPIServer(listenAddr string) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/weather", makeHTTPHandleFunc(s.handleWeatherReport))

	router.HandleFunc("/weather/{id}", makeHTTPHandleFunc(s.handleGetWeatherReport))

	log.Println("Json API running om port: ", s.listenAddr)
	http.ListenAndServe(s.listenAddr, router)
}

// Catch all function for CRUD Operations
func (s *APIServer) handleWeatherReport(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetWeatherReport(w, r)
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

func (s *APIServer) handleGetWeatherReport(w http.ResponseWriter, r *http.Request) error {
	id := mux.Vars(r)["id"]
	println(id)
	weatherReport := NewWeatherReport("Sunny!", 2, 0.1)
	return writeJson(w, http.StatusOK, weatherReport)
}

func (s *APIServer) handleCreateWeatherReport(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleUpdateWeatherReport(w http.ResponseWriter, r *http.Request) error {
	return nil
}

func (s *APIServer) handleDeleteWeatherReport(w http.ResponseWriter, r *http.Request) error {
	return nil
}
