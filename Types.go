package main

import "encoding/xml"

type DescriptionType struct {
	About   string     `xml:"RDF:about,attr"`
	Id      string     `xml:"NS1:id,attr"`
	Type    string     `xml:"NS1:type,attr"`
	Title   string     `xml:"NS1:title,attr"`
	Chars   string     `xml:"NS1:chars,attr"`
	Comment string     `xml:"NS1:comment,attr"`
	Icon    string     `xml:"NS1:icon,attr"`
	Source  string     `xml:"NS1:source,attr"`
}

type UrnType string
type DictType struct {
	File string
	Desc map[UrnType] DescriptionType    // = make(map[UrnType] DescriptionType)
	Seq map[UrnType] []UrnType   // = make(map[UrnType] []UrnType)
	ProcInst []xml.ProcInst // <? xml version?>
	}

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