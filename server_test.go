package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"
)

type TestCase struct {
	Url  string
	Want []*City
}

func TestCityHandler(t *testing.T) {
	var testcases = []TestCase{
		{
			Url: "/city/NoRealCityName",
		},
		{
			Url: "/city/Aachen",
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
		req := httptest.NewRequest(http.MethodGet, testcase.Url, nil)
		w := httptest.NewRecorder()

		log.Println(req)

		FetchCity(w, req)
		res := w.Result()
		log.Println(res)
		log.Println(json.NewEncoder(w).Encode(res))
	}

	/*if !want.MatchString(msg) || err != nil {
	    t.Fatalf(`Hello("Gladys") = %q, %v, want match for %#q, nil`, msg, err, want)
	}*/
}
