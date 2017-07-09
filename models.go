package gostorages

import "encoding/xml"

// metadata response

type MetadataRoot struct {
	GHPayload *GHPayload `xml:"http://www.vizrt.com/types payload,omitempty" json:"payload,omitempty"`
}

type GHPayload struct {
	AttrModel string `xml:" model,attr"  json:",omitempty"`
	AttrXmlns string `xml:" xmlns,attr"  json:",omitempty"`
	GHField []*GHField `xml:"http://www.vizrt.com/types field,omitempty" json:"field,omitempty"`
	XMLName  xml.Name `xml:"http://www.vizrt.com/types payload,omitempty" json:"payload,omitempty"`
}

type GHField struct {
	AttrName string `xml:" name,attr"  json:",omitempty"`
	GHField []*GHField `xml:"http://www.vizrt.com/types field,omitempty" json:"field,omitempty"`
	GHList *GHList `xml:"http://www.vizrt.com/types list,omitempty" json:"list,omitempty"`
	GHValue *GHValue `xml:"http://www.vizrt.com/types value,omitempty" json:"value,omitempty"`
	XMLName  xml.Name `xml:"http://www.vizrt.com/types field,omitempty" json:"field,omitempty"`
}

type GHValue struct {
	Text string `xml:",chardata" json:",omitempty"`
	XMLName  xml.Name `xml:"http://www.vizrt.com/types value,omitempty" json:"value,omitempty"`
}

type GHList struct {
	GHPayload *GHPayload `xml:"http://www.vizrt.com/types payload,omitempty" json:"payload,omitempty"`
	XMLName  xml.Name `xml:"http://www.vizrt.com/types list,omitempty" json:"list,omitempty"`
}

// Image create response

type GHChidleyRoot314159 struct {
	GHFeed *GHFeed `xml:"http://www.w3.org/2005/Atom feed,omitempty" json:"feed,omitempty"`
}

type GHFeed struct {
	AttrXmlnsMedia string `xml:"xmlns media,attr"  json:",omitempty"`
	AttrXmlns string `xml:" xmlns,attr"  json:",omitempty"`
	GHAuthor *GHAuthor `xml:"http://www.w3.org/2005/Atom author,omitempty" json:"author,omitempty"`
	GHEntry *GHEntry `xml:"http://www.w3.org/2005/Atom entry,omitempty" json:"entry,omitempty"`
	GHId *GHId `xml:"http://www.w3.org/2005/Atom id,omitempty" json:"id,omitempty"`
	GHTitle *GHTitle `xml:"http://www.w3.org/2005/Atom title,omitempty" json:"title,omitempty"`
	GHUpdated *GHUpdated `xml:"http://www.w3.org/2005/Atom updated,omitempty" json:"updated,omitempty"`
	XMLName  xml.Name `xml:"http://www.w3.org/2005/Atom feed,omitempty" json:"feed,omitempty"`
}

type GHTitle struct {
	Text string `xml:",chardata" json:",omitempty"`
	XMLName  xml.Name `xml:"http://www.w3.org/2005/Atom title,omitempty" json:"title,omitempty"`
}

type GHUpdated struct {
	Text string `xml:",chardata" json:",omitempty"`
	XMLName  xml.Name `xml:"http://www.w3.org/2005/Atom updated,omitempty" json:"updated,omitempty"`
}

type GHAuthor struct {
	GHName *GHName `xml:"http://www.w3.org/2005/Atom name,omitempty" json:"name,omitempty"`
	XMLName  xml.Name `xml:"http://www.w3.org/2005/Atom author,omitempty" json:"author,omitempty"`
}

type GHName struct {
	Text string `xml:",chardata" json:",omitempty"`
	XMLName  xml.Name `xml:"http://www.w3.org/2005/Atom name,omitempty" json:"name,omitempty"`
}

type GHId struct {
	Text string `xml:",chardata" json:",omitempty"`
	XMLName  xml.Name `xml:"http://www.w3.org/2005/Atom id,omitempty" json:"id,omitempty"`
}

type GHEntry struct {
	AttrXmlns string `xml:" xmlns,attr"  json:",omitempty"`
	GHCategory *GHCategory `xml:"http://www.w3.org/2005/Atom category,omitempty" json:"category,omitempty"`
	GHContent *GHContent `xml:"http://search.yahoo.com/mrss/ content,omitempty" json:"content,omitempty"`
	GHId *GHId `xml:"http://www.w3.org/2005/Atom id,omitempty" json:"id,omitempty"`
	GHLink []*GHLink `xml:"http://www.w3.org/2005/Atom link,omitempty" json:"link,omitempty"`
	GHPublished *GHPublished `xml:"http://www.w3.org/2005/Atom published,omitempty" json:"published,omitempty"`
	GHThumbnail *GHThumbnail `xml:"http://search.yahoo.com/mrss/ thumbnail,omitempty" json:"thumbnail,omitempty"`
	GHTitle *GHTitle `xml:"http://www.w3.org/2005/Atom title,omitempty" json:"title,omitempty"`
	GHUpdated *GHUpdated `xml:"http://www.w3.org/2005/Atom updated,omitempty" json:"updated,omitempty"`
	XMLName  xml.Name `xml:"http://www.w3.org/2005/Atom entry,omitempty" json:"entry,omitempty"`
}

type GHPublished struct {
	Text string `xml:",chardata" json:",omitempty"`
	XMLName  xml.Name `xml:"http://www.w3.org/2005/Atom published,omitempty" json:"published,omitempty"`
}

type GHCategory struct {
	AttrScheme string `xml:" scheme,attr"  json:",omitempty"`
	AttrTerm string `xml:" term,attr"  json:",omitempty"`
	XMLName  xml.Name `xml:"http://www.w3.org/2005/Atom category,omitempty" json:"category,omitempty"`
}

type GHLink struct {
	AttrHref string `xml:" href,attr"  json:",omitempty"`
	AttrRel string `xml:" rel,attr"  json:",omitempty"`
	AttrType string `xml:" type,attr"  json:",omitempty"`
	XMLName  xml.Name `xml:"http://www.w3.org/2005/Atom link,omitempty" json:"link,omitempty"`
}

type GHContent struct {
	AttrType string `xml:" type,attr"  json:",omitempty"`
	AttrUrl string `xml:" url,attr"  json:",omitempty"`
	XMLName  xml.Name `xml:"http://search.yahoo.com/mrss/ content,omitempty" json:"content,omitempty"`
}

type GHThumbnail struct {
	AttrHeight string `xml:" height,attr"  json:",omitempty"`
	AttrUrl string `xml:" url,attr"  json:",omitempty"`
	AttrWidth string `xml:" width,attr"  json:",omitempty"`
	XMLName  xml.Name `xml:"http://search.yahoo.com/mrss/ thumbnail,omitempty" json:"thumbnail,omitempty"`
}
