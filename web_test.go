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

func TestPostJSON_1(t *testing.T) {

	var jsondata = []byte(`{"root_element":{"foo":"bar1"}}`)

	request, _ := http.NewRequest("POST", "/json2xml", bytes.NewReader(jsondata))

	request.Header.Add("Content-Type", "application/json")

	response := httptest.NewRecorder()

	json2Xml(response, request)

	body := response.Body.String()
	if response.Code != http.StatusOK {
		t.Errorf("Response body did not contain expected %v:\n\tbody: %v", "200", response.Code)
	}

	if !strings.Contains(body, "<root_element><foo>bar1</foo></root_element>") {
		t.Errorf("Response body did not contain expected %v:\n\t", body)
	}

}

func TestJson2xml_2(t *testing.T) {

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
			t.Errorf("Response body did not contain expected %v:\n\t", body)
		}

		// actual := Translate(test.in)
		// assert.Equal(t, test.out, actual)
	}

	// t.Errorf("failed dude")
}
