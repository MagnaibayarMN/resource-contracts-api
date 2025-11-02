// Package sql provides PostgreSQL database operations for administrative units.
// It manages province and district data used for geographic filtering.
package sql

import (
	"context"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx/v5/pgxpool"
)

// Unit represents a basic administrative unit with ID and name.
type Unit struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// BaseProvince contains common fields for provinces and districts.
type BaseProvince struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Note     string `json:"-"`
	Type     int    `json:"type"`
	Location string `json:"-"`
}

// Province represents a top-level administrative division.
type Province struct {
	BaseProvince
	ParentID int `json:"parentId"`
}

// District represents a sub-division within a province.
type District struct {
	BaseProvince
}

// Pgsql is the global PostgreSQL connection pool.
var Pgsql *pgxpool.Pool

// EstablishPgSQL initializes the PostgreSQL connection pool.
// It reads the connection string from PGSQL_URL environment variable.
func EstablishPgSQL() {
	poolConfig, err := pgxpool.ParseConfig(os.Getenv("PGSQL_URL"))
	poolConfig.MaxConnIdleTime = 0
	poolConfig.MinConns = 2

	if err != nil {
		log.Fatal(err)
	}
	conn, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to connect to database: %v\n", err)
		os.Exit(1)
	}
	Pgsql = conn
}

// GetProvincesAllUnits retrieves all provinces and districts as a map.
// Returns a map with unit ID as key and name as value.
//
// Returns:
//   - map[int]string: Map of unit ID to name
//   - error: Error if query fails
func GetProvincesAllUnits() (map[int]string, error) {
	query := "select id, name from mongolian_provinces order by name asc"

	rows, err := Pgsql.Query(context.Background(), query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var provinces []Unit
	for rows.Next() {
		var province Unit
		if err := rows.Scan(&province.ID, &province.Name); err != nil {
			return nil, err
		}
		provinces = append(provinces, province)
	}

	result := make(map[int]string)
	for _, item := range provinces {
		result[item.ID] = item.Name
	}

	return result, err
}

// GetProvinces retrieves provinces or districts based on the provided ID.
// If no provinceId is provided, returns all provinces.
// If provinceId is provided, returns districts within that province.
//
// Parameters:
//   - provinceId: Optional province ID to filter districts
//
// Returns:
//   - []Province: List of provinces or districts
//   - error: Error if query fails
func GetProvinces(provinceId string) (provinces []Province, err error) {

	var parentId *int

	if provinceId != "" {
		_provinceId, err := strconv.Atoi(provinceId)
		if err != nil {
			return nil, fmt.Errorf("%s is not a instance of int", provinceId)
		}

		parentId = &_provinceId
	}

	params := []any{}
	_type := 1
	query := "select * from mongolian_provinces where type = $1 order by name asc"

	if parentId != nil {
		_type = 2
		query = "select * from mongolian_provinces where type = $1 and parent_id = $2 order by name asc"
		params = append(params, _type)
		params = append(params, parentId)
	} else {
		params = append(params, _type)
	}

	rows, err := Pgsql.Query(context.Background(), query, params...)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	for rows.Next() {
		var province Province
		if err := rows.Scan(&province.ID, &province.Name, &province.Note, &province.ParentID, &province.Type, &province.Location); err != nil {
			return nil, err
		}
		provinces = append(provinces, province)
	}

	return provinces, nil
}
