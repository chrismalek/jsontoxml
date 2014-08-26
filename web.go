package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"github.com/clbanning/anyxml"
	"github.com/clbanning/mxj"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func setResponseHeaders(rw http.ResponseWriter) {

	rw.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
	rw.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
	rw.Header().Set("Expires", "0")                                         // Proxies
	rw.Header().Set("Content-Type", "application/xml")

}

func json2Xml(rw http.ResponseWriter, req *http.Request) {

	switch req.Method {

	case "POST":
		switch req.Header.Get("Content-Type") {
		case "application/json":
			body, _ := ioutil.ReadAll(req.Body)

			log.Println(req.Header.Get("x-dbname"))
			var parsedJson interface{}

			if err := json.Unmarshal(body, &parsedJson); err != nil {
				http.Error(rw, "Error converting reading JSON", 406)
				log.Println("invalid JSON found \n", string(body))
				return
			}

			var xmlout []byte

			xmlout, err2 := anyxml.Xml(parsedJson, "root")
			if err2 != nil {
				http.Error(rw, "Error converting to XML", 406)
				return
			}

			setResponseHeaders(rw)
			output := []byte(xml.Header)
			output = append(output, xmlout...)
			rw.Write(output)

		default:
			http.Error(rw, "Please send along the content-type header of application/json", 406)

		}

	default:
		http.Error(rw, "Invalid request method.", 405)

	}
}

func xml2Json(rw http.ResponseWriter, req *http.Request) {

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		panic(err)
	}

	if req.Method == "POST" {
		switch req.Header.Get("Content-Type") {
		case "application/xml":

			mapVal, merr := mxj.NewMapXml(body)
			if merr != nil {
				http.Error(rw, "Error Reading XML", 406)
				log.Println("invalid xml \n", string(body))

				return
			}
			jVal, jerr := mapVal.Json()
			if jerr != nil {
				http.Error(rw, "Error converting to JSON", 406) // handle error
				return
			}

			setResponseHeaders(rw)

			rw.Write(jVal)

		default:
			http.Error(rw, "Please send along the content-type header of application/xml", 406)

		}

	} else {
		http.Error(rw, "Invalid request method.", 405)
	}

}

func RootHandler(res http.ResponseWriter, req *http.Request) {
	var tmpl = template.Must(template.ParseFiles("templates/about.html"))
	err := tmpl.Execute(res, nil)

	if err != nil {
		http.Error(res, "Internal server error", 500)
		return
	}
}

func main() {

	http.HandleFunc("/json2xml", json2Xml)
	http.HandleFunc("/xml2json", xml2Json)
	http.HandleFunc("/", RootHandler)

	var portNumber string = os.Getenv("PORT")

	if portNumber == "" {
		portNumber = "9000"
	}

	fmt.Println("listening on port:" + portNumber)
	log.Fatal(http.ListenAndServe(":"+portNumber, nil))

}
