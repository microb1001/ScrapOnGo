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
	Other   []string     `xml:",any"`
}

type Rdf_RdfOut struct {
	XMLName xml.Name `xml:"RDF:RDF"`
	RDF_Descriptions []RDF_DescriptionOut `xml:"RDF:Description"`
	Other   []string     `xml:",any"`
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
	//Other   string     `xml:",any"`
}
type RDF_DescriptionOut struct {
	//XMLName xml.Name `xml:"http://www.w3.org/1999/02/22-rdf-syntax-ns# Description"`
	About   string     `xml:"RDF:about,attr"`
	Id      string     `xml:"NS2:id,attr"`
	Type    string     `xml:"NS2:type,attr"`
	Title   string     `xml:"NS2:title,attr"`
	Chars   string     `xml:"NS2:chars,attr"`
	Comment string     `xml:"NS2:comment,attr"`
	Icon    string     `xml:"NS2:icon,attr"`
	Source  string     `xml:"NS2:source,attr"`
	//Other   string     `xml:",any"`
}
type Redirect struct {
	Title string `xml:"title,attr"`
}

type Page struct {
	Title string `xml:"title"`
	Redir Redirect `xml:"redirect"`
	Text string `xml:"revision>text"`
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
	var res2 Rdf_RdfOut
	err := xml.Unmarshal(data, &res)
	if err != nil {
		panic(err)
	}
	err2 := xml.Unmarshal(byteValue, &res1)
	if err2 != nil {
		panic(err2)
	}
	for _, v := range res1.RDF_Descriptions {
		res2.RDF_Descriptions=append(res2.RDF_Descriptions,RDF_DescriptionOut(v))
	}
	//res2.RDF_Descriptions[1]=RDF_DescriptionOut(res1.RDF_Descriptions[1])
	buf, err := xml.MarshalIndent(&res2,"", "   ")
	if err != nil {
		log.Fatal(err)
	}
	_=buf
//	fmt.Println(string(buf))

	xmlFile.Close()
	xmlFile, err3 := os.Open("scrapbook.rdf")
	if err3 != nil {
		fmt.Println("Error opening file:", err)
		return
	}
	decoder1 := xml.NewDecoder(xmlFile)
	for {
		// Read tokens from the XML document in a stream.
		t, _ := decoder1.Token()
		if t == nil {
			break
		}
		// Inspect the type of the token just read.

		switch se := t.(type) {
		case xml.StartElement:
			fmt.Println("StartElement", se.Name)
			// If we just read a StartElement token
			// ...and its name is "page"
			if se.Name.Local == "page" {
				var p Page
				// decode a whole chunk of following XML into the
				// variable p which is a Page (se above)
				decoder1.DecodeElement(&p, &se)
				// Do some stuff with the page.
				//p.Title = CanonicalizeTitle(p.Title)
				//...
			}
			//...
		case xml.EndElement:
			fmt.Println("EndElement", t)
		}
	}
	//fmt.Printf("%+v\n", res)
	//fmt.Printf("%+v\n", res1)
}
