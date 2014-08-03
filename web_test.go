package main

import (
	// "io/ioutil"
	"bytes"
	"fmt"
	"net/http"
	"net/http/httptest"
	"path/filepath"
	"reflect"
	"runtime"
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
		out: "<root/>",
	},
	{
		in: []byte(`{"errorCode": "E0000001","errorSummary": "Api validation failed",
                 "errorLink": "E0000001", "errorId": "oaeHfmOAx1iRLa0H10DeMz5fQ","errorCauses": [
        {"errorSummary": "login: An object with this field already exists in the current organization"}]}`),

		out: "<errorId>oaeHfmOAx1iRLa0H10DeMz5fQ</errorId>",
	},
}

var jsonTestCasesFailures = []struct {
	in          []byte
	description string
	out         string
}{
	{description: "Invalid JSON",
		in:  []byte(`jibberish`),
		out: "",
	},
	{description: "Invalid JSON 2",
		in: []byte("{]"), out: "",
	},
}

// {`{"root_element":{"foo":"bar1"}}`,"<root_element><foo>bar1</foo></root_element>" }}

// assert fails the test if the condition is false.
func assert(tb testing.TB, condition bool, msg string, v ...interface{}) {
	if !condition {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: "+msg+"\033[39m\n\n", append([]interface{}{filepath.Base(file), line}, v...)...)
		tb.FailNow()
	}
}

// ok fails the test if an err is not nil.
func ok(tb testing.TB, err error) {
	if err != nil {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d: unexpected error: %s\033[39m\n\n", filepath.Base(file), line, err.Error())
		tb.FailNow()
	}
}

// equals fails the test if exp is not equal to act.
func equals(tb testing.TB, exp, act interface{}) {
	if !reflect.DeepEqual(exp, act) {
		_, file, line, _ := runtime.Caller(1)
		fmt.Printf("\033[31m%s:%d:\n\n\texp: %#v\n\n\tgot: %#v\033[39m\n\n", filepath.Base(file), line, exp, act)
		tb.FailNow()
	}
}

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

		assert(t, (response.Code == http.StatusOK), "Response body did not contain expected %v:\n\t: %v", "200", response.Code)
		assert(t, strings.Contains(body, test.out), "Did not get expected response. ------ \n\t Recieved: \n\t %v:\n\t Expected: \n\t %v", body, test.out)

	}

}

func TestJson2xmlFailTestCases(t *testing.T) {

	for _, test := range jsonTestCasesFailures {

		request, _ := http.NewRequest("POST", "/json2xml", bytes.NewReader(test.in))

		request.Header.Add("Content-Type", "application/json")

		response := httptest.NewRecorder()

		json2Xml(response, request)

		assert(t, response.Code == http.StatusNotAcceptable, "Expected response code 406 but received %v", response.Code)

	}

}
