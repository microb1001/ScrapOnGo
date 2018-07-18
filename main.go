package main

import (
	"encoding/xml"
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

func main() {
	Dict:=LoadAll("scrapbook.rdf",)
 	SaveAll("scrapbookout.rdf",Dict)


}
