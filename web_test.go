package main

import (
	// "io/ioutil"
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

var jsonTestCases = []struct {
	in  []byte
	out string
}{
	{
		in:  []byte(`{"root_element":{"foo":"bar1"}}`),
		out: "<root_element><foo>bar1</foo></root_element>",
	},
	{
		in:  []byte("{}"),
		out: "<doc/>",
	},
	{
		in: []byte(`{"errorCode": "E0000001","errorSummary": "Api validation failed",
                 "errorLink": "E0000001", "errorId": "oaeHfmOAx1iRLa0H10DeMz5fQ","errorCauses": [
        {"errorSummary": "login: An object with this field already exists in the current organization"}]}`),

		out: "<errorId>oaeHfmOAx1iRLa0H10DeMz5fQ</errorId>",
	},
}

// {`{"root_element":{"foo":"bar1"}}`,"<root_element><foo>bar1</foo></root_element>" }}

func TestHome(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(RootHandler))
	defer ts.Close()
	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	RootHandler(response, request)

	if response.Code != http.StatusOK {
		t.Fatalf("Response body did not contain expected %v:\n\tbody: %v", "200", response.Code)
	}

}

func TestJson2xmlTestCases(t *testing.T) {

	for _, test := range jsonTestCases {

		request, _ := http.NewRequest("POST", "/json2xml", bytes.NewReader(test.in))

		request.Header.Add("Content-Type", "application/json")

		response := httptest.NewRecorder()

		json2Xml(response, request)

		body := response.Body.String()

		if response.Code != http.StatusOK {
			t.Errorf("Response body did not contain expected %v:\n\tbody: %v", "200", response.Code)
		}

		if !strings.Contains(body, test.out) {
			t.Errorf("Response body did not contain expected %v:\n\t %v", body, test.out)
		}

	}

}
