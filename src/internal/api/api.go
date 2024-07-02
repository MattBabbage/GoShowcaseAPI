package api

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/MattBabbage/GoShowcaseAPI/internal/storage"
	"github.com/MattBabbage/GoShowcaseAPI/internal/types"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type APIServer struct {
	listenAddr string
	store      storage.Storage
}

func NewAPIServer(listenAddr string, store storage.Storage) *APIServer {
	return &APIServer{
		listenAddr: listenAddr,
		store:      store,
	}
}

func (s *APIServer) Run() {
	router := mux.NewRouter()
	router.HandleFunc("/weather", makeHTTPHandleFunc(s.handleWeatherReport))

	router.HandleFunc("/weather/{id}", makeHTTPHandleFunc(s.handleWeatherReportById))

	log.Println("Json API running on port: ", s.listenAddr)
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
func (s *APIServer) handleWeatherReportById(w http.ResponseWriter, r *http.Request) error {
	if r.Method == "GET" {
		return s.handleGetWeatherReportById(w, r)
	}
	if r.Method == "DELETE" {
		return s.handleDeleteWeatherReport(w, r)
	}
	return fmt.Errorf("invalid request")
}

// Get singular weather report based on id
func (s *APIServer) handleGetWeatherReportById(w http.ResponseWriter, r *http.Request) error {
	uid, err := getUUIDFromString(mux.Vars(r)["id"])
	if err != nil {
		return err
	}
	weatherReport, err := s.store.GetWeatherReportByID(uid)
	if err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, weatherReport)
}

func (s *APIServer) handleCreateWeatherReport(w http.ResponseWriter, r *http.Request) error {
	createRequest := new(types.CreateWeatherReportRequest)
	if err := json.NewDecoder(r.Body).Decode(&createRequest); err != nil {
		return err
	}
	WeatherReport := types.NewWeatherReport(createRequest.Description, createRequest.Temperature, createRequest.RainChance)
	if err := s.store.CreateWeatherReport(WeatherReport); err != nil {
		return err
	}

	return writeJson(w, http.StatusOK, createRequest)
}

func (s *APIServer) handleUpdateWeatherReport(w http.ResponseWriter, r *http.Request) error {
	weatherReport := new(types.WeatherReport)
	if err := json.NewDecoder(r.Body).Decode(&weatherReport); err != nil {
		return err
	}
	if err := s.store.UpdateWeatherReport(weatherReport); err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, weatherReport)
}

func (s *APIServer) handleDeleteWeatherReport(w http.ResponseWriter, r *http.Request) error {
	uid, err := getUUIDFromString(mux.Vars(r)["id"])
	if err != nil {
		return err
	}

	if err := s.store.DeleteWeatherReport(uid); err != nil {
		return err
	}
	return writeJson(w, http.StatusOK, map[string]uuid.UUID{"deleted": uid})
}

func writeJson(w http.ResponseWriter, status int, v any) error {
	w.Header().Add("Content-Type", "application/json")
	w.WriteHeader(status)
	return json.NewEncoder(w).Encode(v)
}

type apiFunction func(http.ResponseWriter, *http.Request) error

type apiError struct {
	Error string `json:"error"`
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

// Helper function to return an appropriate error
func getUUIDFromString(idString string) (uuid.UUID, error) {
	uid, err := uuid.Parse(idString)
	if err != nil {
		return uuid.Nil, fmt.Errorf("invalid uuid given '%s'", idString)
	}
	return uid, nil
}
