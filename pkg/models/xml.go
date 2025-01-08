package models

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
)

type Envelope struct {
	XMLName xml.Name `xml:"Envelope"`
	Cube    Cube     `xml:"Cube"`
}

type Cube struct {
	XMLName xml.Name   `xml:"Cube"`
	Cubes   []DateCube `xml:"Cube"`
}

type DateCube struct {
	XMLName xml.Name       `xml:"Cube"`
	Time    string         `xml:"time,attr"`
	Cubes   []CurrencyCube `xml:"Cube"`
}

type CurrencyCube struct {
	XMLName  xml.Name `xml:"Cube"`
	Currency string   `xml:"currency,attr"`
	Rate     string   `xml:"rate,attr"`
}

func FetchXML(url string) (*Envelope, error) {
	var envelope Envelope
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch XML: %s", resp.Status)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if err := xml.Unmarshal(body, &envelope); err != nil {
		return nil, fmt.Errorf("unmarshal : %w", err)
	}

	return &envelope, nil
}
