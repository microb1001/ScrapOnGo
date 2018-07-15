package main

import (
	"encoding/xml"
	"fmt"
	"os"
	"io"
	"log"
)

type DescriptionType struct {
	About   string     `xml:"RDF:about,attr"`
	Id      string     `xml:"NS9:id,attr"`
	Type    string     `xml:"NS9:type,attr"`
	Title   string     `xml:"NS9:title,attr"`
	Chars   string     `xml:"NS9:chars,attr"`
	Comment string     `xml:"NS9:comment,attr"`
	Icon    string     `xml:"NS9:icon,attr"`
	Source  string     `xml:"NS9:source,attr"`
}

type UrnType string
var DictDesc map[UrnType] DescriptionType = make(map[UrnType] DescriptionType)
var DictSeq map[UrnType] []UrnType = make(map[UrnType] []UrnType)
type Name_x2 struct{old xml.Name;new xml.Name;attricnt int}

var urnSeq UrnType = "ERR"

// Пространства имен
const SpaceNS1="http://amb.vis.ne.jp/mozilla/scrapbook-rdf#"
const SpaceNC="http://home.netscape.com/NC-rdf#"
const SpaceRDF="http://www.w3.org/1999/02/22-rdf-syntax-ns#"
// Токены
var TnNil=xml.Name{"",""}
var TnRDF=xml.Name{SpaceRDF,"RDF"}
var TnDe=xml.Name{SpaceRDF,"Description"}
var TnBo=xml.Name{SpaceNC,"BookmarkSeparator"}
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
var ProcInst xml.ProcInst
func main() {
	xmlFile, err := os.Open("scrapbook.rdf")
	if err != nil { fmt.Println(err)}

	defer xmlFile.Close()

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
			case Name_x2{TnRDF,TnDe,8},Name_x2{TnRDF,TnBo,8}:
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
				if (cs.new==TnBo) && (ee.Type!="separator"){panic("Тип неверен separator "+string(urn))}
				if (cs.new!=TnBo) && (ee.Type=="separator"){panic("Сепаратор без сепаратора separator "+string(urn))}
				DictDesc[urn]=ee
			case Name_x2{TnRDF,TnSe,1}:
				if tt.Attr[0].Name!=TnAb {fmt.Println("Ошибка Seq: ", tt)}
				urnSeq= UrnType(tt.Attr[0].Value)
				if len(DictSeq[urnSeq])!=0 {fmt.Println("Дублирование Seq: ", tt)}
				//var ur []UrnType // почему то нельзя объединить со следующей строкой
				DictSeq[urnSeq]=make([]UrnType, 0, 10)
			case Name_x2{TnSe,TnLi,1}:
				if tt.Attr[0].Name!=TnRe {fmt.Println("Ошибка Li: ", tt)}
				DictSeq[urnSeq]=append(DictSeq[urnSeq], UrnType(tt.Attr[0].Value))
			default:
				fmt.Println("Нераспознан : ", tt)
			}
		case xml.EndElement:
			if up[len(up)-1] != tt.Name {fmt.Println("EndElementОшибка", up[len(up)],t)} // Вообще-то закрытие проверает сам декодер
			up=up[:len(up)-1]
		case xml.CharData:
			//fmt.Println("TEXT", tt) // Тут нужна проверка что это только концы строк и пробелы
		case xml.ProcInst:
			ProcInst=tt.Copy()
			//fmt.Println("ProcInst", tt, string(tt.Inst))  // Тоже проверка

		case xml.Comment:
			fmt.Println("ERR-COMMENT", tt)
		case xml.Directive:
			fmt.Println("ERR-Directive", tt)
		default:
			fmt.Println("ERR", t)
		}


	}
	//fmt.Printf("%+v\n", DictSeq)
	//fmt.Printf("%+v\n", DictDesc)
	xmlFile2, err := os.Create("scrapbookout.rdf")
	if err != nil { fmt.Println(err)}

	defer xmlFile2.Close()

	encoder := xml.NewEncoder(xmlFile2)
	encoder.Indent("", "  ")
//============================типы только для оформления вывода в XML=================================
	type li2Type struct {
		Urn UrnType `xml:"RDF:resource,attr"`
	}
	type Seq1Type struct {
		Urn UrnType `xml:"RDF:about,attr"`
		Li  [] li2Type    `xml:"RDF:li"`
	}
	type Scrap0Type struct {
		XMLName      xml.Name                `xml:"RDF:RDF"`
		Descriptions []DescriptionType `xml:"RDF:Description"`
		Separators []DescriptionType   `xml:"NC:BookmarkSeparator"`
		Seq []Seq1Type                       `xml:"RDF:Seq"`
		NameSpaceNS1 string                  `xml:"xmlns:NS9,attr"`
		NameSpaceNC string                   `xml:"xmlns:NC,attr"`
		NameSpaceRDF string                  `xml:"xmlns:RDF,attr"`
	}
	var res Scrap0Type
//==========================================++++++++++++++++++=======================================
	res.NameSpaceNS1=SpaceNS1
	res.NameSpaceNC=SpaceNC
	res.NameSpaceRDF=SpaceRDF

	for _,v :=range DictDesc {
		if v.Type=="separator" {res.Separators=append(res.Separators,v)
		} else {
		res.Descriptions=append(res.Descriptions,v)
		}
	}
	for k,v :=range DictSeq {
		var u Seq1Type
		u.Urn=k
		for _,m:=range v {
			u.Li=append(u.Li,li2Type{Urn:m})
		}
		res.Seq=append(res.Seq,u)
	}
	err=encoder.EncodeToken(ProcInst)
	if err != nil {
		log.Fatal(err)
	}
	err=encoder.Encode(&res)
	if err != nil {
		log.Fatal(err)
	}
	buf, err := xml.MarshalIndent(&res,"", "   ")
	if err != nil {
		log.Fatal(err)
	}
	_=buf
	//fmt.Println(string(buf))
	//fmt.Printf("%+v\n", DictDesc)
}
