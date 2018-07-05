package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"io/ioutil"
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
	XMLName xml.Name `xml:"RDF"`
	//NS1   string     `xml:"xmlns:NS1,attr"`
	//NC   string     `xml:"xmlns:NC,attr"`
	//RDF   string     `xml:"xmlns:RDF,attr"`
	RDF_Descriptions []RDF_Description `xml:"Description"`
}

type RDF_Description struct {
	Properties []Property `xml:"PROPERTY"`
	About   string     `xml:"about,attr"`
	Id      string     `xml:"id,attr"`
	Type    string     `xml:"type,attr"`
	Title   string     `xml:"title,attr"`
	Chars   string     `xml:"chars,attr"`
	Comment string     `xml:"comment,attr"`
	Icon    string     `xml:"icon,attr"`
	Source  string     `xml:"source,attr"`
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
	fmt.Printf("%+v\n", res)
	fmt.Printf("%+v\n", res1)
}
