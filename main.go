package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"io/ioutil"
	"log"
)

type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:",chardata"`
}

type Object struct {
	Properties []Property `xml:"PROPERTY"`
	Basetype   string     `xml:"basetype,attr"`
	Name       string     `xml:"name,attr"`
	Oid        int        `xml:"oid,attr"`
}

type Response struct {
	XMLName xml.Name `xml:"RESPONSE"`
	Objects []Object `xml:"OBJECT"`
}

type Rdf_Rdf struct {
	XMLName xml.Name `xml:"http://www.w3.org/1999/02/22-rdf-syntax-ns# RDF"`
	RDF_Descriptions []RDF_Description `xml:"http://www.w3.org/1999/02/22-rdf-syntax-ns# Description"`
}

type RDF_Description struct {
	//XMLName xml.Name `xml:"http://www.w3.org/1999/02/22-rdf-syntax-ns# Description"`
	About   string     `xml:"http://www.w3.org/1999/02/22-rdf-syntax-ns# about,attr"`
	Id      string     `xml:"http://amb.vis.ne.jp/mozilla/scrapbook-rdf# id,attr"`
	Type    string     `xml:"http://amb.vis.ne.jp/mozilla/scrapbook-rdf# type,attr"`
	Title   string     `xml:"http://amb.vis.ne.jp/mozilla/scrapbook-rdf# title,attr"`
	Chars   string     `xml:"http://amb.vis.ne.jp/mozilla/scrapbook-rdf# chars,attr"`
	Comment string     `xml:"http://amb.vis.ne.jp/mozilla/scrapbook-rdf# comment,attr"`
	Icon    string     `xml:"http://amb.vis.ne.jp/mozilla/scrapbook-rdf# icon,attr"`
	Source  string     `xml:"http://amb.vis.ne.jp/mozilla/scrapbook-rdf# source,attr"`
}
func main() {
	// Open our xmlFile
	xmlFile, err1 := os.Open("scrapbook.rdf")
	// if we os.Open returns an error then handle it
	if err1 != nil {
		fmt.Println(err1)
	}
	defer xmlFile.Close()
	byteValue, _ := ioutil.ReadAll(xmlFile)

	data := []byte(`<RESPONSE>
    <OBJECT basetype="status" name="status" oid="1">
      <PROPERTY name="response-type">success</PROPERTY>
      <PROPERTY name="response-type-numeric">0</PROPERTY>
      <PROPERTY name="response">0f738648db95bb1f6ca37f6b8b5aafa8</PROPERTY>
      <PROPERTY name="return-code">1</PROPERTY>
    </OBJECT>
</RESPONSE>`)

	var res Response
	var res1 Rdf_Rdf
	err := xml.Unmarshal(data, &res)
	if err != nil {
		panic(err)
	}
	err2 := xml.Unmarshal(byteValue, &res1)
	if err2 != nil {
		panic(err2)
	}
	buf, err := xml.MarshalIndent(&res1,"", "   ")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println(string(buf))

	//fmt.Printf("%+v\n", res)
	//fmt.Printf("%+v\n", res1)
}
