package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

func TestCityHandler(t *testing.T) {
	//define test struct
	type test struct {
		name string
		url  string
		city string
		want []*City
	}

	//define test cases
	testcases := []test{
		{
			name: "Test with no real city name",
			url:  "/city",
			city: "NoRealCityName",
		},
		{
			name: "Call with correct city name",
			url:  "/city",
			city: "Aachen",
			want: []*City{
				{
					Id:          3097,
					Name:        "Aachen",
					Countrycode: "DEU",
					District:    "Nordrhein-Westfalen",
					Population:  243825,
				},
			},
		},
	}

	for _, test := range testcases {
		//for each case run test
		t.Run(test.name, func(t *testing.T) {
			//mock request
			req, _ := http.NewRequest(http.MethodGet, test.url, nil)
			w := httptest.NewRecorder()

			//populate mux vars
			vars := map[string]string{
				"name": test.city,
			}
			r := mux.SetURLVars(req, vars)

			//execute function to test
			FetchCity(w, r)
			var target []*City
			var city City
			res := w.Result()

			if len(test.want) == 0 {
				//if no result expected, check if result is returned
				if res.StatusCode != 404 {
					t.Fatal("Found resource, but did not expect to")
				}
			} else {
				//if result is expected, decode json response
				if err := json.NewDecoder(res.Body).Decode(&city); err != nil {
					t.Fatal(err)
				}
				//check if it matches the expected result
				target = append(target, &city)
				for index, result := range target {
					if !result.Compare(test.want[index]) {
						t.Fatalf("Result for index %v does not match expected state", index)
					}
				}
			}
		})
	}
}

func TestCitiesHandler(t *testing.T) {
	//define test struct
	type test struct {
		name      string
		url       string
		wantCount int
		want      []*City
	}

	//define test cases
	testcases := []test{
		{
			name:      "Call without limit",
			url:       "/city",
			wantCount: 4079,
		},
		{
			name:      "Call with limit of 2",
			url:       "/city?limit=2",
			wantCount: 2,
			want: []*City{
				{
					Id:          1,
					Name:        "Kabul",
					Countrycode: "AFG",
					District:    "Kabol",
					Population:  1780000,
				},
				{
					Id:          2,
					Name:        "Qandahar",
					Countrycode: "AFG",
					District:    "Qandahar",
					Population:  237500,
				},
			},
		},
	}

	for _, test := range testcases {
		//for each case run test
		t.Run(test.name, func(t *testing.T) {
			//mock request
			req, _ := http.NewRequest(http.MethodGet, test.url, nil)
			w := httptest.NewRecorder()

			//execute function to test
			FetchCities(w, req)
			var target []*City
			res := w.Result()

			//decode json response
			if err := json.NewDecoder(res.Body).Decode(&target); err != nil {
				t.Fatal(err)
			}

			if test.wantCount != len(target) {
				//if number of results does not match expected number of results, fail
				t.Fatalf("Number of returned resources (%v) does not match expected value %v", len(target), test.wantCount)
			}
			if len(test.want) != 0 {
				//check if results match expected results
				for index, result := range target {
					if !result.Compare(test.want[index]) {
						t.Fatalf("Result for index %v does not match expected state", index)
					}
				}
			}
		})
	}
}
