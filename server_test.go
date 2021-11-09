package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

type test struct {
	name string
	url  string
	city string
	want []*City
}

func TestCityHandler(t *testing.T) {
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
		/*{
			Url: "/city?limit=1",
			Want: []*City{
				{
					Id:          1,
					Name:        "Kabul",
					Countrycode: "AFG",
					District:    "Kabol",
					Population:  1780000,
				},
			},
		},*/
	}

	for _, test := range testcases {
		t.Run(test.name, func(t *testing.T) {
			req, _ := http.NewRequest(http.MethodGet, test.url, nil)
			w := httptest.NewRecorder()

			vars := map[string]string{
				"name": test.city,
			}
			r := mux.SetURLVars(req, vars)

			FetchCity(w, r)
			var target []*City
			var city City
			res := w.Result()

			if len(test.want) == 0 {
				if res.StatusCode != 404 {
					t.Fatal("Found resource, but did not want to")
				}
			} else {
				if err := json.NewDecoder(res.Body).Decode(&city); err != nil {
					t.Fatal(err)
				}
				target = append(target, &city)
				for index, result := range target {
					if !result.Compare(test.want[index]) {
						t.Fatal("Result does not match expected state")
					}
				}
			}
		})
	}
}
