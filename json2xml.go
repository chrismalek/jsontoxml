package main

import (
	"github.com/clbanning/mxj/j2x"
	"github.com/clbanning/mxj/x2j"
	"io/ioutil"
	"log"
	"net/http"
)

func json2Xml(rw http.ResponseWriter, req *http.Request) {

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		panic(err)
	}

	log.Println(string(body))

	if err != nil {
		panic(err)
	}
	var xmloutput []byte

	xmloutput, err = j2x.JsonToXml(body)

	log.Println(string(xmloutput))

	if err != nil {
		panic(err)
	}

	rw.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
	rw.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
	rw.Header().Set("Expires", "0")                                         // Proxies
	rw.Header().Set("Content-Type", "application/xml")

	rw.Write(xmloutput)

}

func xml2Json(rw http.ResponseWriter, req *http.Request) {

	body, err := ioutil.ReadAll(req.Body)

	if err != nil {
		panic(err)
	}

	log.Println(string(body))

	if err != nil {
		panic(err)
	}
	var jsonoutput []byte

	jsonoutput, err = x2j.XmlToJson(body)

	log.Println(string(jsonoutput))

	if err != nil {
		panic(err)
	}

	rw.Header().Set("Cache-Control", "no-cache, no-store, must-revalidate") // HTTP 1.1.
	rw.Header().Set("Pragma", "no-cache")                                   // HTTP 1.0.
	rw.Header().Set("Expires", "0")                                         // Proxies
	rw.Header().Set("Content-Type", "application/json")

	rw.Write(jsonoutput)

}

func main() {
	http.HandleFunc("/json2xml", json2Xml)
	http.HandleFunc("/xml2json", xml2Json)

	log.Fatal(http.ListenAndServe(":8082", nil))
}
