package main

import (
	"os"
	"fmt"
	"encoding/xml"
	"log"
	"io"
)

func LoadAll (fileName string) DictType {
	var ScrapLoad DictType = DictType {
		make(map[UrnType] DescriptionType),
		make(map[UrnType] []UrnType),
		make([]xml.ProcInst,0)}
	var urnSeq UrnType = "ERR"
	xmlFile, err := os.Open(fileName)
	if err != nil { fmt.Println(err); panic(err)}

	defer xmlFile.Close()

	decoder := xml.NewDecoder(xmlFile)
	var list = []xml.Name{TnNil}
	for {
		currentToken, err := decoder.Token()
		if err != nil && err != io.EOF {fmt.Println(err); panic(err)}
		if currentToken == nil {break} // Конец файла

		switch elem := currentToken.(type) {
		case xml.StartElement:
			var tokensPair = Name_x2{list[len(list)-1], elem.Name,len(elem.Attr)}
			list =append(list, elem.Name)

			switch tokensPair {
			case Name_x2{TnNil,TnRDF,3}: // Начало

			case Name_x2{TnRDF,TnDe,8},  // Дескриптор или сепаратор
				 Name_x2{TnRDF,TnBo,8}:
				var urn UrnType
				var newDesc DescriptionType
				for _, a := range elem.Attr{
					switch a.Name {
					case TnAb: newDesc.About= a.Value; urn = UrnType(a.Value)
					case TnId: newDesc.Id= a.Value
					case TnTy: newDesc.Type= a.Value
					case TnTi: newDesc.Title= a.Value
					case TnCh: newDesc.Chars= a.Value
					case TnCo: newDesc.Comment= a.Value
					case TnIc: newDesc.Icon= a.Value
					case TnSo: newDesc.Source= a.Value
					default: fmt.Println(a.Name);panic("Неопознанный тег в Description")
					}
				}
				// Проверки ошибок
				if _,ok:= ScrapLoad.Desc[urn];ok {panic("Дубликат description "+string(urn))}
				if (tokensPair.new==TnBo) && (newDesc.Type!="separator"){panic("Тип неверен, должен быть separator "+string(urn))}
				if (tokensPair.new!=TnBo) && (newDesc.Type=="separator"){panic("Дескриптор, с типом separator "+string(urn))}
				ScrapLoad.Desc[urn]= newDesc

			case Name_x2{TnRDF,TnSe,1}: // Seq - Последовательность дескрипторов
				if elem.Attr[0].Name!=TnAb {fmt.Println(elem); panic("Ошибка Seq: должен быть атрибут about")}
				urnSeq= UrnType(elem.Attr[0].Value)
				if len(ScrapLoad.Seq[urnSeq])!=0 {fmt.Println(elem); panic("Дублирование Seq:")}
				ScrapLoad.Seq[urnSeq]=make([]UrnType, 0, 10)

			case Name_x2{TnSe,TnLi,1}: // Элемент последовательности
				if elem.Attr[0].Name!=TnRe {fmt.Println("Ошибка Li: ", elem);panic("Li должен иметь атрибут resource")}
				ScrapLoad.Seq[urnSeq]=append(ScrapLoad.Seq[urnSeq], UrnType(elem.Attr[0].Value))

				default:
				fmt.Println(elem);panic("Нераспознан");
			}
		case xml.EndElement:
			if list[len(list)-1] != elem.Name {fmt.Println("EndElementОшибка", list[len(list)], currentToken)} // Вообще-то закрытие проверает сам декодер
			list = list[:len(list)-1]
		case xml.CharData:
			//fmt.Println("TEXT", elem) // Тут нужна проверка что это только концы строк и пробелы
		case xml.ProcInst:
			ScrapLoad.ProcInst = append(ScrapLoad.ProcInst,elem.Copy())

		case xml.Comment:
			fmt.Println("ERR-COMMENT", elem)
		case xml.Directive:
			fmt.Println("ERR-Directive", elem)
		default:
			fmt.Println("ERR", currentToken)
		}


	}
	return ScrapLoad
}

func SaveAll (fileName string, ScrapSave DictType ){
	xmlFile2, err := os.Create(fileName)
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

	for _,tk:=range (ScrapSave.ProcInst){
		err=encoder.EncodeToken(tk)
		if err != nil {
			log.Fatal(err)
		}
		}
	for _,v :=range ScrapSave.Desc {
		if v.Type=="separator" {res.Separators=append(res.Separators,v)
		} else {
			res.Descriptions=append(res.Descriptions,v)
		}
	}
	for k,v :=range ScrapSave.Seq {
		var u Seq1Type
		u.Urn=k
		for _,m:=range v {
			u.Li=append(u.Li,li2Type{Urn:m})
		}
		res.Seq=append(res.Seq,u)
	}

	err=encoder.Encode(&res)
	if err != nil {
		log.Fatal(err)
	}
}
