package stellplatzanzahl

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
	"time"
)

type ApiAntwort struct {
	Raststaetten []Raststaette `json:"parking_lorry"`
}

type Raststaette struct {
	Extent                   string
	Identifier               string
	Title                    string
	Point                    string
	Subtitle                 string                     `json:"subtitle"`
	Coordinates              Coordinates                `json:"coordinate"`
	LorryParkingFeatureIcons []LorryParkingFeatureIcons `json:"lorryParkingFeatureIcons"`
	Description              []string                   `json:"description"`
}

type Coordinates struct {
	Latitude  string `json:"lat"`
	Longitude string `json:"long"`
}

type LorryParkingFeatureIcons struct {
	Zusatzinfos1 Picknickmoeglichkeit `json:"0"`
	Zusatzinfos2 Toiletten            `json:"1"`
}

type Picknickmoeglichkeit struct {
	DescriptionP string `json:"description"`
}

type Toiletten struct {
	DescriptionT string `json:"description"`
}
type Parkingslots struct {
	PKW int
	LKW int
}

type error interface {
	Error() string
}

func (p *Parkingslots) Sum() int {
	return p.LKW + p.PKW
}

func ParkinglorrySum(HighwayName string) (*Parkingslots, error) {

	hostname := "https://verkehr.autobahn.de/o/autobahn/"
	path1 := "//"

	path3 := "//services/parking_lorry"

	result, err := url.JoinPath(hostname, path1, HighwayName, path3)
	if err != nil {
		return nil, fmt.Errorf("error joining urlpath: %v", err)
	}

	client := http.Client{
		Timeout: 6 * time.Second,
	}
	resp, err := client.Get(result)
	if err != nil {
		return nil, fmt.Errorf("error in ApiConnection: %v", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("error Http Status: %v", err)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("error reading body: %v", err)
	}

	var Liste ApiAntwort
	err = json.Unmarshal(body, &Liste)
	if err != nil {
		return nil, fmt.Errorf("error unmarsheling json: %v", err)
	}

	totalPKWParkingSpaces := 0

	for _, p := range Liste.Raststaetten {
		for _, d := range p.Description {
			re := regexp.MustCompile(`PKW.Stellplätze:\s*([0-9]+)`)
			matches := re.FindStringSubmatch(d)
			if len(matches) > 1 {
				countStr := matches[1]
				count, err := strconv.Atoi(countStr)
				if err != nil {
					fmt.Println("Error converting string to integer:", err)
					continue
				}
				totalPKWParkingSpaces += count
			}
		}
	}

	totalLKWParkingSpaces := 0

	for _, p := range Liste.Raststaetten {
		for _, d := range p.Description {
			re := regexp.MustCompile(`LKW.Stellplätze:\s*([0-9]+)`)
			matches := re.FindStringSubmatch(d)
			if len(matches) > 1 {
				countStr := matches[1]
				count, err := strconv.Atoi(countStr)
				if err != nil {
					fmt.Println("Error converting string to integer:", err)
					continue
				}
				totalLKWParkingSpaces += count
			}
		}
	}

	return &Parkingslots{LKW: totalLKWParkingSpaces, PKW: totalPKWParkingSpaces}, nil
}
