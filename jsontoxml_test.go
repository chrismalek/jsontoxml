package main

import (
	// "io/ioutil"
	"bytes"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func TestHome(t *testing.T) {

	request, _ := http.NewRequest("GET", "/", nil)
	response := httptest.NewRecorder()

	rootHandler(response, request)

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
		t.Errorf("Response body did not contain expected %v:\n\tbody: %v", "San Francisco", body)
	}

}

func Testjson2xml_1(t *testing.T) {
	t.Errorf("failed dude")
}
