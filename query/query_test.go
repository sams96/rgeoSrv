/*
Copyright 2020 Sam Smith

Licensed under the Apache License, Version 2.0 (the "License"); you may not use
this file except in compliance with the License.  You may obtain a copy of the
License at

   http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software distributed
under the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR
CONDITIONS OF ANY KIND, either express or implied.  See the License for the
specific language governing permissions and limitations under the License.
*/

package query

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/sergi/go-diff/diffmatchpatch"
)

func TestQuery(t *testing.T) {
	tests := []struct {
		name           string
		in             *http.Request
		out            *httptest.ResponseRecorder
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Algeria",
			in:             httptest.NewRequest("GET", "/query?1.880273&31.787305", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"country":"Algeria","country_long":"People's Democratic Republic of Algeria","country_code_2":"DZ","country_code_3":"DZA","continent":"Africa","region":"Africa","subregion":"Northern Africa","province":"El Bayadh","province_code":"DZ-32"}`,
		},
		{
			name:           "Madagascar",
			in:             httptest.NewRequest("GET", "/query?47.523836&-18.905691", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"country":"Madagascar","country_long":"Republic of Madagascar","country_code_2":"MG","country_code_3":"MDG","continent":"Africa","region":"Africa","subregion":"Eastern Africa","province":"Analamanga","province_code":"MG-T","city":"Antananarivo"}`,
		},
		{
			name:           "Zimbabwe",
			in:             httptest.NewRequest("GET", "/query?29.832875&-19.948725", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"country":"Zimbabwe","country_long":"Republic of Zimbabwe","country_code_2":"ZW","country_code_3":"ZWE","continent":"Africa","region":"Africa","subregion":"Eastern Africa","province":"Midlands","province_code":"ZW-MI"}`,
		},
		{
			name:           "Ocean",
			in:             httptest.NewRequest("GET", "/query?0&0", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "country not found\n",
		},
		{
			name:           "North Pole",
			in:             httptest.NewRequest("GET", "/query?-135&90", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusInternalServerError,
			expectedBody:   "country not found\n",
		},
		{
			name:           "South Pole",
			in:             httptest.NewRequest("GET", "/query?44.99&-89.99", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"country":"Antarctica","country_code_2":"AQ","country_code_3":"ATA","continent":"Antarctica","region":"Antarctica","subregion":"Antarctica","province":"Antarctica","province_code":"AQ-X01~"}`,
		},
		{
			name:           "Alaska",
			in:             httptest.NewRequest("GET", "/query?-149.901785&61.199134", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"country":"United States of America","country_long":"United States of America","country_code_2":"US","country_code_3":"USA","continent":"North America","region":"Americas","subregion":"Northern America","province":"Alaska","province_code":"US-AK","city":"Anchorage"}`,
		},
		{
			name:           "UK",
			in:             httptest.NewRequest("GET", "/query?0&51.5045", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"country":"United Kingdom","country_long":"United Kingdom of Great Britain and Northern Ireland","country_code_2":"GB","country_code_3":"GBR","continent":"Europe","region":"Europe","subregion":"Northern Europe","province":"Tower Hamlets","province_code":"GB-TWH","city":"London"}`,
		},
		{
			name:           "Libya",
			in:             httptest.NewRequest("GET", "/query?24.98&25.86", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"country":"Libya","country_long":"Libya","country_code_2":"LY","country_code_3":"LBY","continent":"Africa","region":"Africa","subregion":"Northern Africa","province":"Al Kufrah","province_code":"LY-KF"}`,
		},
		{
			name:           "Egypt",
			in:             httptest.NewRequest("GET", "/query?25.005187&25.855963", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"country":"Egypt","country_long":"Arab Republic of Egypt","country_code_2":"EG","country_code_3":"EGY","continent":"Africa","region":"Africa","subregion":"Northern Africa","province":"Al Wadi at Jadid","province_code":"EG-WAD"}`,
		},
		{
			name:           "US Border",
			in:             httptest.NewRequest("GET", "/query?-102.560616&48.992073", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"country":"United States of America","country_long":"United States of America","country_code_2":"US","country_code_3":"USA","continent":"North America","region":"Americas","subregion":"Northern America","province":"North Dakota","province_code":"US-ND"}`,
		},
		{
			name:           "Canada Border",
			in:             httptest.NewRequest("GET", "/query?-102.560616&49.0", nil),
			out:            httptest.NewRecorder(),
			expectedStatus: http.StatusOK,
			expectedBody:   `{"country":"Canada","country_long":"Canada","country_code_2":"CA","country_code_3":"CAN","continent":"North America","region":"Americas","subregion":"Northern America","province":"Saskatchewan","province_code":"CA-SK"}`,
		},
	}

	q, err := NewHandlers(nil)
	if err != nil {
		t.Error(err)
	}

	dmp := diffmatchpatch.New()

	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			q.query(test.out, test.in)
			if test.out.Code != test.expectedStatus {
				t.Errorf("expected: %d\ngot: %d\n",
					test.expectedStatus, test.out.Code)
				t.Fail()
			}

			body := test.out.Body.String()
			if test.expectedBody != body {
				diff := dmp.DiffMain(test.expectedBody, body, false)
				t.Error(dmp.DiffPrettyText(diff))
			}
		})
	}
}
