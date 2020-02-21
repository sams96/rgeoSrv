package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/go-test/deep"
)

func TestRootHandler(t *testing.T) {
	tests := []struct {
		name           string
		in             *http.Request
		out            *httptest.ResponseRecorder
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Algeria",
			in:             httptest.NewRequest("GET", "/?1.880273&31.787305", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"Country":"Algeria","CountryLong":"People's Democratic Republic of Algeria","CountryCode2":"DZ","CountryCode3":"DZA","Continent":"Africa","Region":"Africa","SubRegion":"Northern Africa"}`,
		},
		{
			name:           "Madagascar",
			in:             httptest.NewRequest("GET", "/?47.478275&-17.530126", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"Country":"Madagascar","CountryLong":"Republic of Madagascar","CountryCode2":"MG","CountryCode3":"MDG","Continent":"Africa","Region":"Africa","SubRegion":"Eastern Africa"}`,
		},
		{
			name:           "Zimbabwe",
			in:             httptest.NewRequest("GET", "/?29.832875&-19.948725", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"Country":"Zimbabwe","CountryLong":"Republic of Zimbabwe","CountryCode2":"ZW","CountryCode3":"ZWE","Continent":"Africa","Region":"Africa","SubRegion":"Eastern Africa"}`,
		},
		{
			name:           "Ocean",
			in:             httptest.NewRequest("GET", "/?0&0", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "country not found\n",
		},
		{
			name:           "North Pole",
			in:             httptest.NewRequest("GET", "/?-135&90", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "country not found\n",
		},
		{
			name:           "South Pole",
			in:             httptest.NewRequest("GET", "/?44.99&-89.99", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"Country":"Antarctica","CountryLong":"","CountryCode2":"AQ","CountryCode3":"ATA","Continent":"Antarctica","Region":"Antarctica","SubRegion":"Antarctica"}`,
		},
		{
			name:           "Alaska",
			in:             httptest.NewRequest("GET", "/?-150.542&66.3", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"Country":"United States of America","CountryLong":"United States of America","CountryCode2":"US","CountryCode3":"USA","Continent":"North America","Region":"Americas","SubRegion":"Northern America"}`,
		},
		{
			name:           "UK",
			in:             httptest.NewRequest("GET", "/?0&52", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"Country":"United Kingdom","CountryLong":"United Kingdom of Great Britain and Northern Ireland","CountryCode2":"GB","CountryCode3":"GBR","Continent":"Europe","Region":"Europe","SubRegion":"Northern Europe"}`,
		},
		{
			name:           "Libya",
			in:             httptest.NewRequest("GET", "/?24.994611&25.860750", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"Country":"Libya","CountryLong":"Libya","CountryCode2":"LY","CountryCode3":"LBY","Continent":"Africa","Region":"Africa","SubRegion":"Northern Africa"}`,
		},
		{
			name:           "Egypt",
			in:             httptest.NewRequest("GET", "/?25.005187&25.855963", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"Country":"Egypt","CountryLong":"Arab Republic of Egypt","CountryCode2":"EG","CountryCode3":"EGY","Continent":"Africa","Region":"Africa","SubRegion":"Northern Africa"}`,
		},
		{
			name:           "US Border",
			in:             httptest.NewRequest("GET", "/?-102.560616&48.998074", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"Country":"United States of America","CountryLong":"United States of America","CountryCode2":"US","CountryCode3":"USA","Continent":"North America","Region":"Americas","SubRegion":"Northern America"}`,
		},
		{
			name:           "Canada Border",
			in:             httptest.NewRequest("GET", "/?-102.560616&49.0", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"Country":"Canada","CountryLong":"Canada","CountryCode2":"CA","CountryCode3":"CAN","Continent":"North America","Region":"Americas","SubRegion":"Northern America"}`,
		},
	}

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			rootHandler(test.out, test.in)
			if test.out.Code != test.expectedStatus {
				t.Errorf("expected: %d\ngot: %d\n",
					test.expectedStatus, test.out.Code)
				t.Fail()
			}

			body := test.out.Body.String()
			if diff := deep.Equal(test.expectedBody, body); diff != nil {
				t.Error(diff)
			}
		})
	}
}
