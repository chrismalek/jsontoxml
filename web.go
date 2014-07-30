package main

import (
	"fmt"
	"github.com/clbanning/mxj"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"regexp"
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
			body, err := ioutil.ReadAll(req.Body)

			match, _ := regexp.MatchString(`\B\[`, string(body))
			if match {
				body = []byte(`{"top-level-array":` + string(body) + `}`)
				log.Println("\n----\n matched top-level-array\n")
				log.Println("new Body \n", string(body))
			}
			// log.Println(string(body))

			mapVal, merr := mxj.NewMapJson(body)
			if merr != nil {
				http.Error(rw, "Error converting reading JSON", 406)
				return
			}
			xmlVal, xerr := mapVal.Xml()
			if xerr != nil {
				http.Error(rw, "Error converting to XML", 406)
				return
			}

			// xmloutput, err = j2x.JsonToXml(body)

			// log.Println(string(xmlVal))

			if err != nil {
				log.Println(err)
				http.Error(rw, "Could not convert to xml", 400)
				return
			}
			setResponseHeaders(rw)
			// rw.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
			// rw.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
			// rw.Header().Set("Expires", "0")                                         // Proxies
			// rw.Header().Set("Content-Type", "application/xml")

			rw.Write(xmlVal)

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

			// log.Println(string(body))

			// var jsonoutput []byte

			// jsonoutput, err = x2j.XmlToJson(body)

			// log.Println(string(jsonoutput))

			mapVal, merr := mxj.NewMapXml(body)
			if merr != nil {
				http.Error(rw, "Error Reading XML", 406)
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

func rootHandler(res http.ResponseWriter, req *http.Request) {
	var tmpl = template.Must(template.ParseFiles("templates/about.html"))
	err := tmpl.Execute(res, nil)

	if err != nil {
		http.Error(res, "Internal server error", 500)
		return
	}

	// fmt.Fprintln(res, "This page intentionally left blank")
}

func main() {

	http.HandleFunc("/json2xml", json2Xml)
	http.HandleFunc("/xml2json", xml2Json)
	http.HandleFunc("/", rootHandler)

	var portNumber string = os.Getenv("PORT")

	if portNumber == "" {
		portNumber = "9000"
	}

	fmt.Println("listening on port:" + portNumber)
	log.Fatal(http.ListenAndServe(":"+portNumber, nil))

}
