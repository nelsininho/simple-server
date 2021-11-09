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

const getCities = `SELECT id, name, countrycode, district, population FROM city`

type City struct {
	Id          int
	Name        string
	Countrycode string
	District    string
	Population  int
}

func (c *City) Compare(c2 *City) bool {
	return c.Id == c2.Id &&
		c.Name == c2.Name &&
		c.Countrycode == c2.Countrycode &&
		c.District == c2.District &&
		c.Population == c2.Population
}

func logging(f http.HandlerFunc) http.HandlerFunc {
	//logging middleware for displaying the url path
	return func(w http.ResponseWriter, req *http.Request) {
		log.Printf("Called route %v`\n", req.URL.Path)
		f(w, req)
	}
}

func FetchCity(w http.ResponseWriter, req *http.Request) {
	//read requested city name
	vars := mux.Vars(req)
	cityname := vars["name"]

	//connect do db
	db, err := sql.Open("postgres", "postgresql://docker:docker@postgres-db/world?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Internal Server Error!"))
		return
	}

	//get city information
	city := City{}
	if err = db.QueryRow(getCities+` WHERE name = $1`, cityname).Scan(&city.Id, &city.Name, &city.Countrycode, &city.District, &city.Population); err != nil {
		log.Println(err)
		if err == sql.ErrNoRows {
			//if no result returned, log unknown city
			log.Printf(" %v: unknown city", cityname)
			w.WriteHeader(http.StatusNotFound)
			w.Write([]byte("404 - City not found!"))
			return
		}
	}

	//return result in JSON format
	log.Println(city)
	json.NewEncoder(w).Encode(city)
}

func FetchCities(w http.ResponseWriter, req *http.Request) {
	//read possible limit parameter
	m, _ := url.ParseQuery(req.URL.RawQuery)
	limit := m["limit"]

	//connect to db
	db, err := sql.Open("postgres", "postgresql://docker:docker@postgres-db/world?sslmode=disable")
	if err != nil {
		log.Fatalln(err)
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("500 - Internal Server Error!"))
		return
	}

	//read results from db with or without limit, depending on whether it was requested
	cities := []City{}
	var rows *sql.Rows
	if len(limit) == 0 {
		rows, err = db.Query(getCities)
	} else {
		query := getCities + ` LIMIT $1`
		rows, err = db.Query(query, limit[0])
	}

	//scan result rows and append them to return slice
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

	//return result
	json.NewEncoder(w).Encode(cities)
}

func main() {
	r := mux.NewRouter()

	//create endpoint for returning city by name and for returning multiple cities
	cityrouter := r.PathPrefix("/city").Subrouter()
	cityrouter.HandleFunc("/{name}", logging(FetchCity)).Methods("GET")
	cityrouter.HandleFunc("", logging(FetchCities)).Methods("GET")

	http.ListenAndServe(":8080", r)
}
