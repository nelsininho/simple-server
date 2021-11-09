package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"

	"database/sql"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

const getCity = `SELECT id, name, countrycode, district, population FROM city WHERE name = $1`
const getCities = `SELECT id, name, countrycode, district, population FROM city`
const getCitiesWithLimit = `SELECT id, name, countrycode, district, population FROM city LIMIT $1`

type City struct {
	Id          int
	Name        string
	Countrycode string
	District    string
	Population  int
}

func logging(f http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Called route %v`\n", req.URL.Path)
		f(w, req)
	}
}

func FetchCity(w http.ResponseWriter, req *http.Request) {
	vars := mux.Vars(req)
	cityname := vars["name"]

	log.Println("Name:", vars)

	db, err := sql.Open("postgres", "postgresql://docker:docker@localhost:5432/world?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	city := City{}

	if err = db.QueryRow(getCity, cityname).Scan(&city.Id, &city.Name, &city.Countrycode, &city.District, &city.Population); err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			log.Printf(" %v: unknown city", cityname)
			return
		}
	}

	log.Println(city)
	json.NewEncoder(w).Encode(city)
}

func FetchCities(w http.ResponseWriter, req *http.Request) {
	m, _ := url.ParseQuery(req.URL.RawQuery)
	limit := m["limit"]

	db, err := sql.Open("postgres", "postgresql://docker:docker@localhost:5432/world?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
	}

	cities := []City{}
	var rows *sql.Rows
	if len(limit) == 0 {
		rows, err = db.Query(getCities)
	} else {
		rows, err = db.Query(getCitiesWithLimit, limit[0])
	}

	if rows != nil {
		for rows.Next() {
			var city City
			if err := rows.Scan(&city.Id, &city.Name, &city.Countrycode, &city.District, &city.Population); err != nil {
				log.Println("Error")
			}
			cities = append(cities, city)
		}
		if err = rows.Err(); err != nil {
			log.Println("Error")
		}
	}

	json.NewEncoder(w).Encode(cities)
}

func main() {
	r := mux.NewRouter()

	cityrouter := r.PathPrefix("/city").Subrouter()
	cityrouter.HandleFunc("/{name}", logging(FetchCity)).Methods("GET")
	cityrouter.HandleFunc("", logging(FetchCities)).Methods("GET")

	http.ListenAndServe(":8089", r)
}
