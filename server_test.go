package main

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gorilla/mux"
)

type TestCase struct {
	Url  string
	Name string
	Want []*City
}

func TestCityHandler(t *testing.T) {
	testcases := []TestCase{
		/*{
			Url:  "/city",
			Name: "NoRealCityName",
		},*/
		{
			Url:  "/city",
			Name: "Aachen",
			Want: []*City{
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

	for _, testcase := range testcases {
		req, _ := http.NewRequest(http.MethodGet, testcase.Url, nil)
		w := httptest.NewRecorder()

		vars := map[string]string{
			"name": testcase.Name,
		}
		r := mux.SetURLVars(req, vars)

		FetchCity(w, r)
		var target []*City
		var city City
		res := w.Result()

		json.NewDecoder(res.Body).Decode(&city)
		target = append(target, &city)
		/*if reflect.TypeOf(target).Kind() != reflect.TypeOf(testcase.Want).Kind() {
			target = []*City{target}
		}*/
		for index, result := range target {
			if !result.Compare(testcase.Want[index]) {
				t.Fatal("Result does not match expected state")
			}
		}

		/*if len(testcase.Want) != len(json.NewEncoder(w).Encode(res)) || err != nil {
			t.Fatalf(`Return value of call to route %q does not match %#q in length`, testcase.Url, testcase.Want)
		}*/
	}
}
