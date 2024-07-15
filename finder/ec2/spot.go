package ec2

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"log"
	"net/http"
	"time"

	"github.com/yeo/ec2shop/finder/common"
)

const (
	AWSSpotPriceUrl = "https://website.spot.ec2.aws.a2z.com/spot.js"
)

type SpotPrice struct {
	Linux *float64
	MSWin *float64

	AdvisorLinux   *AdvisorInfo
	AdvisorWindows *AdvisorInfo
}

type SpotValueColumn struct {
	Name     string           `json:"name"`
	RawPrice *common.RawPrice `json:"prices"`
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

type SpotPriceFinder struct {
	client *http.Client
	Done   chan bool

	// nested map of
	// region.instance_type => SpotPrice
	pricePerRegions map[string]map[string]*SpotPrice
}

func NewSpotPriceFinder() *SpotPriceFinder {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	s := SpotPriceFinder{
		client:          client,
		Done:            make(chan bool),
		pricePerRegions: make(map[string]map[string]*SpotPrice),
	}

	return &s
}

func (s *SpotPriceFinder) Fetch() error {
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
		}
	}

	//log.Printf("Spot Price %+v", *(s.pricePerRegions["us-east-1"]["c5ad.2xlarge"].Linux))
	log.Println("Fetched spot price in", time.Now().Sub(t0), "at", time.Now())

	return nil
}

func (s *SpotPriceFinder) PriceForInstance(region string, instanceType string) (*SpotPrice, error) {
	m := s.pricePerRegions[region][instanceType]
	if m == nil {
		return nil, errors.New("Invalid instance type or region")
	}

	return m, nil
}

func (s *SpotPriceFinder) Run() {
	ticker := time.NewTicker(150 * time.Second)

	go func() {
		for {
			select {
			case <-s.Done:
				return
			case <-ticker.C:
				s.Fetch()
				s.FetchAdvisor()
			}
		}
	}()
	s.Fetch()
	s.FetchAdvisor()
}
