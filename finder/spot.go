package finder

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

const (
	AWSSpotPriceUrl = "https://website.spot.ec2.aws.a2z.com/spot.js"
)

type SpotPrice struct {
	Linux *float64
	MSWin *float64
}

type SpotValueColumn struct {
	Name     string    `json:"name"`
	RawPrice *RawPrice `json:"prices"`
}

type SpotInstanceTypeSize struct {
	Size         string             `json:"size"`
	ValueColumns []*SpotValueColumn `json:"valueColumns"`
}

type SpotInstanceType struct {
	Type  string                  `json:"type"`
	Sizes []*SpotInstanceTypeSize `json:"sizes"`
}

type SpotRegion struct {
	Region        string              `json:"region"`
	InstanceTypes []*SpotInstanceType `json:"instanceTypes"`
}
type SpotPriceResponse struct {
	Rate    string        `json:"rate"`
	Regions []*SpotRegion `json:"regions"`
}

type SpotPriceResponseWrap struct {
	Config *SpotPriceResponse `json:"config"`
}

type SpotPriceCrawler struct {
	client          *http.Client
	Done            chan bool
	pricePerRegions map[string]map[string]*SpotPrice
}

func NewSpotPriceCrawler() *SpotPriceCrawler {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	s := SpotPriceCrawler{
		client:          client,
		Done:            make(chan bool),
		pricePerRegions: make(map[string]map[string]*SpotPrice),
	}

	return &s
}

func (s *SpotPriceCrawler) Fetch() error {
	t0 := time.Now()
	resp, err := s.client.Get(AWSSpotPriceUrl)
	if err != nil {
		log.Println("Error fetching spot price", err)
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Fail to read response body", err)
		return err
	}

	body = body[9 : len(body)-2]

	var priceWrap SpotPriceResponseWrap

	err = json.Unmarshal(body, &priceWrap)
	if err != nil {
		log.Println("Cannot parse json from spot request response", err)
		return err
	}

	price := priceWrap.Config

	for _, r := range price.Regions {
		s.pricePerRegions[r.Region] = make(map[string]*SpotPrice)

		for _, t := range r.InstanceTypes {
			total := 0
			for _, size := range t.Sizes {
				s.pricePerRegions[r.Region][size.Size] = &SpotPrice{}
				total += 1
				for _, vc := range size.ValueColumns {
					if vc.RawPrice.USD == "N/A*" {
						continue
					}

					if vc.Name == "linux" {
						if v, err := vc.RawPrice.Price(); err == nil {
							s.pricePerRegions[r.Region][size.Size].Linux = new(float64)
							*s.pricePerRegions[r.Region][size.Size].Linux = v
						}
					}

					if vc.Name == "mswin" {
						if v, err := vc.RawPrice.Price(); err == nil {
							s.pricePerRegions[r.Region][size.Size].MSWin = new(float64)
							*s.pricePerRegions[r.Region][size.Size].MSWin = v
						}
					}
				}

			}

			log.Printf("[spot price loader] found %d server price for region %s", total, r.Region)
		}
	}

	//log.Printf("Spot Price %+v", s.pricePerRegions)
	log.Println("Fetched spot price in", time.Now().Sub(t0), "at", time.Now())

	return nil
}

func (s *SpotPriceCrawler) SpotRegionName(region string) string {
	// The AWS API we're using has some funky names for some regions; e.g. eu-west-1 is referred to as eu-ireland
	// This function maps an "actual" region name to the one in this API call
	spotRegionMap := map[string]string{
		"us-east-1":      "us-east",
		"us-west-1":      "us-west",
		"eu-west-1":      "eu-ireland",
		"ap-southeast-1": "apac-sin",
		"ap-southeast-2": "apac-syd",
		"ap-northeast-1": "apac-tokyo",
	}
	spotRegionName, found := spotRegionMap[region]
	if found {
		return spotRegionName
	} else {
		return region
	}
}

func (s *SpotPriceCrawler) PriceForInstance(region string, instanceType string) (*SpotPrice, error) {
	m := s.pricePerRegions[region][instanceType]
	if m == nil {
		return nil, errors.New("Invalid instance type or region")
	}

	return m, nil
}

func (s *SpotPriceCrawler) Run() {
	ticker := time.NewTicker(150 * time.Second)

	go func() {
		for {
			select {
			case <-s.Done:
				return
			case <-ticker.C:
				s.Fetch()
			}
		}
	}()
	s.Fetch()
}
