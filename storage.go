package main

import (
	"database/sql"
	"fmt"

	"github.com/google/uuid"
	_ "github.com/lib/pq"
)

type Storage interface {
	CreateWeatherReport(*WeatherReport) error
	DeleteWeatherReport(uuid.UUID) error
	UpdateWeatherReport(*WeatherReport) error
	GetWeatherReports() ([]*WeatherReport, error)
	GetWeatherReportByID(id uuid.UUID) (*WeatherReport, error)
}

type PostgresStore struct {
	db *sql.DB
}

func NewPostgressStore() (*PostgresStore, error) {
	connStr := "user=postgres dbname=postgres password=mysecretpassword sslmode=disable"
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}
	if err := db.Ping(); err != nil {
		return nil, err
	}
	return &PostgresStore{
		db: db,
	}, nil
}

func (s *PostgresStore) Init() error {
	return s.createWeatherReportTable()
}

func (s *PostgresStore) createWeatherReportTable() error {
	query := `create table if not exists weatherreports
				(
					id uuid primary key,
					description varchar(255),
					temperature int,
					chance_rain float(24),
					created_at timestamp
				)`
	_, err := s.db.Exec(query)
	return err
}

func (s *PostgresStore) CreateWeatherReport(w *WeatherReport) error {
	query :=
		`insert into weatherreports
		(id, description, temperature, chance_rain, created_at)
		values ($1, $2, $3, $4, $5)`
	resp, err := s.db.Query(query, w.ID, w.Description, w.Temperature, w.RainChance, w.CreatedAt)
	if err != nil {
		return err
	}
	println(resp)
	return nil
}

func (s *PostgresStore) GetWeatherReportByID(id uuid.UUID) (*WeatherReport, error) {
	rows, err := s.db.Query("select * from weatherreports where id = $1", id)
	if err != nil {
		return nil, err
	}
	for rows.Next() {
		return scanIntoWeatherReport(rows)
	}
	return nil, fmt.Errorf("weather report '%s' not found", id.String())
}

func (s *PostgresStore) UpdateWeatherReport(*WeatherReport) error {

	return nil
}

func (s *PostgresStore) DeleteWeatherReport(id uuid.UUID) error {
	_, err := s.db.Query("delete from weatherreports where id = $1", id)
	return err
}

func (s *PostgresStore) GetWeatherReports() ([]*WeatherReport, error) {
	rows, err := s.db.Query("select * from weatherreports")

	if err != nil {
		return nil, err
	}

	println(rows)
	reports, err := scanIntoWeatherReports(rows)

	if err != nil {
		return nil, err
	}

	return reports, nil
}

func scanIntoWeatherReports(rows *sql.Rows) ([]*WeatherReport, error) {
	reports := []*WeatherReport{}
	for rows.Next() {
		report := new(WeatherReport)
		err := rows.Scan(&report.ID, &report.Description, &report.Temperature, &report.RainChance, &report.CreatedAt)
		if err != nil {
			return nil, err
		}
		reports = append(reports, report)
	}
	return reports, nil
}

func scanIntoWeatherReport(rows *sql.Rows) (*WeatherReport, error) {
	report := new(WeatherReport)
	err := rows.Scan(&report.ID, &report.Description, &report.Temperature, &report.RainChance, &report.CreatedAt)
	return report, err
}
