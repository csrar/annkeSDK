package models

import "encoding/xml"

type LoginResponse struct {
	XMLName                  xml.Name `xml:"SessionLoginCap"`
	Version                  string   `xml:"version,attr"`
	Xmlns                    string   `xml:"xmlns,attr"`
	SessionID                string   `xml:"sessionID"`
	Challenge                string   `xml:"challenge"`
	Iterations               int      `xml:"iterations"`
	IsIrreversible           bool     `xml:"isIrreversible"`
	Salt                     string   `xml:"salt"`
	IsSessionIDValidLongTerm struct {
		Text bool   `xml:",chardata"`
		Opt  string `xml:"opt,attr"`
	} `xml:"isSessionIDValidLongTerm"`
	SessionIDVersion int `xml:"sessionIDVersion"`
}

type Session struct {
	XMLName                  xml.Name `xml:"SessionLogin"`
	Text                     string   `xml:",chardata"`
	UserName                 string   `xml:"userName"`
	Password                 string   `xml:"password"`
	SessionID                string   `xml:"sessionID"`
	IsSessionIDValidLongTerm bool     `xml:"isSessionIDValidLongTerm"`
	SessionIDVersion         int      `xml:"sessionIDVersion"`
}
