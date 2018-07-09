package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"log"
	"io"
)

type RdfType struct {
	XMLName      xml.Name          `xml:"RDF:RDF"`
	Descriptions []DescriptionType `xml:"RDF:Description"`
}

type DescriptionType struct {
	//XMLName xml.Name `xml:"http://www.w3.org/1999/02/22-rdf-syntax-ns# Description"`
	About   string     `xml:"RDF:about,attr"`
	Id      string     `xml:"NS2:id,attr"`
	Type    string     `xml:"NS2:type,attr"`
	Title   string     `xml:"NS2:title,attr"`
	Chars   string     `xml:"NS2:chars,attr"`
	Comment string     `xml:"NS2:comment,attr"`
	Icon    string     `xml:"NS2:icon,attr"`
	Source  string     `xml:"NS2:source,attr"`
}
type UrnType string

var DictDesc map[UrnType] DescriptionType = make(map[UrnType] DescriptionType)
var DictSeq map[UrnType] []UrnType = make(map[UrnType] []UrnType)

type Name_x2 struct{old xml.Name;new xml.Name;attricnt int}

var res RdfType
var urnSeq UrnType = "ERR"

// Пространства имен
const SpaceNS1="http://amb.vis.ne.jp/mozilla/scrapbook-rdf#"
const SpaceNC="http://home.netscape.com/NC-rdf#"
const SpaceRDF="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
// Токены
var TnNil=xml.Name{"",""}
var TnRDF=xml.Name{SpaceRDF,"RDF"}
var TnDe=xml.Name{SpaceRDF,"Description"}
var TnSe=xml.Name{SpaceRDF,"Seq"}
var TnLi=xml.Name{SpaceRDF,"li"}
var TnAb=xml.Name{SpaceRDF,"about"}
var TnId=xml.Name{SpaceNS1,"id"}
var TnTy=xml.Name{SpaceNS1,"type"}
var TnTi=xml.Name{SpaceNS1,"title"}
var TnCh=xml.Name{SpaceNS1,"chars"}
var TnCo=xml.Name{SpaceNS1,"comment"}
var TnIc=xml.Name{SpaceNS1,"icon"}
var TnSo=xml.Name{SpaceNS1,"source"}
var TnRe=xml.Name{SpaceRDF,"resource"}

func main() {
	xmlFile, errOs := os.Open("scrapbook.rdf")
	if errOs != nil { fmt.Println(errOs)}
	defer xmlFile.Close()
	//byteValue, _ := ioutil.ReadAll(xmlFile)

	decoder := xml.NewDecoder(xmlFile)
	var up  = []xml.Name{TnNil}
	for {
		t, err := decoder.Token()
		if err != nil && err != io.EOF {fmt.Println(err); panic(err)}
		if t== nil {break} // Конец файла

		switch tt := t.(type) {
		case xml.StartElement:
			var cs  = Name_x2{up[len(up)-1],tt.Name,len(tt.Attr)}
			up=append(up,tt.Name)

			switch cs {
			case Name_x2{TnNil,TnRDF,3}:
			case Name_x2{TnRDF,TnDe,8}:
				var urn UrnType
				var ee DescriptionType
				for _, a := range tt.Attr{
					switch a.Name {
					case TnAb: ee.About= a.Value; urn = UrnType(a.Value)
					case TnId: ee.Id= a.Value
					case TnTy: ee.Type= a.Value
					case TnTi: ee.Title= a.Value
					case TnCh: ee.Chars= a.Value
					case TnCo: ee.Comment= a.Value
					case TnIc: ee.Icon= a.Value
					case TnSo: ee.Source= a.Value
					default: fmt.Println(a.Name);panic("Неопознанный тег в Description")
					}
				}
				if _,ok:=DictDesc[urn];ok {panic("Дубликат description "+string(urn))}
				DictDesc[urn]=ee
			case Name_x2{TnRDF,TnSe,1}:
				if tt.Attr[0].Name!=TnAb {fmt.Println("Ошибка Seq: ", tt)}
				urnSeq=UrnType(tt.Attr[0].Value)
				if len(DictSeq[urnSeq])!=0 {fmt.Println("Дублирование Seq: ", tt)}

			case Name_x2{TnSe,TnLi,1}:
				if tt.Attr[0].Name!=TnRe {fmt.Println("Ошибка Li: ", tt)}
				DictSeq[urnSeq]=append(DictSeq[urnSeq], UrnType(tt.Attr[0].Value))
			default:
			fmt.Println("Нераспознан : ", tt)
			}

			//...
		case xml.EndElement:
			if up[len(up)-1] != tt.Name {fmt.Println("EndElementОшибка", up[len(up)],t)} // Вообще-то закрытие проверает сам декодер
			up=up[:len(up)-1]
			//fmt.Println("EndElement", t,up)
		case xml.CharData:
			//fmt.Println("TEXT", tt) // Тут нужна проверка что это только концы строк и пробелы
		case xml.ProcInst:
			//fmt.Println("ProcInst", tt, string(tt.Inst))  // Тоже проверка

		case xml.Comment:
			fmt.Println("ERR-COMMENT", tt)
		case xml.Directive:
			fmt.Println("ERR-Directive", tt)
		default:
			fmt.Println("ERR", t)
		}

		buf, err := xml.MarshalIndent(&res,"", "   ")
		if err != nil {
			log.Fatal(err)
		}
		_=buf
		//	fmt.Println(string(buf))
	}
	fmt.Printf("%+v\n", DictSeq)
	fmt.Printf("%+v\n", DictDesc)
}
