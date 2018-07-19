package main

import (
	"encoding/xml"
	"encoding/csv"
	"log"
	"fmt"
	"os"
)

type Name_x2 struct{old xml.Name;new xml.Name;attricnt int}
//var DictDesc map[UrnType] DescriptionType = make(map[UrnType] DescriptionType)
//var DictSeq map[UrnType] []UrnType = make(map[UrnType] []UrnType)
//var Dict DictType = DictType {
//	make(map[UrnType] DescriptionType),
//	make(map[UrnType] []UrnType),
//	make([]xml.ProcInst,0)}
//
//

var Scrap []DictType
func main() {
	//Dict:=LoadAll("E:/test/SB-Note-part3-0618/scrapbook.rdf",)
 	//SaveAll("scrapbookout.rdf",Dict)

	csvFile, err := os.Open("list.csv")
	if err != nil { fmt.Println(err); panic(err)}
	defer csvFile.Close()

	r := csv.NewReader(csvFile)
	ScrapFolders, err := r.ReadAll()
		if err != nil {	log.Fatal(err)}

	for _,a := range(ScrapFolders){
		fmt.Println("Обрабатывается: "+a[0])
	Scrap=append(Scrap, LoadAll(a[0]))
	}
	for _,a := range(ScrapFolders){
		fmt.Println("Проверка целостности : "+a[0])
	Integrity(Scrap[0])
	}
	for _,a := range(Scrap){
		fmt.Println("Сохраняется: "+a.File)
		SaveAll(a.File,a)
	}

}
