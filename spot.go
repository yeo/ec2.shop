package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	AWSSpotPriceUrl       = "https://website.spot.ec2.aws.a2z.com/spot.js"
	AWSSpotAdvisorDataUrl = "https://spot-bid-advisor.s3.amazonaws.com/spot-advisor-data.json"
)

type SpotPrice struct {
	Linux            *float64
	LinuxSavings     *int
	LinuxReclaimRate *int
	MSWin            *float64
	MSWinSavings     *int
	MSWinReclaimRate *int
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

type SpotInstanceTypeDetails struct {
	Savings     int `json:"s"`
	ReclaimRate int `json:"r"`
}

type SpotAdvisorRegion struct {
	Windows map[string]*SpotInstanceTypeDetails `json:"Windows"`
	Linux   map[string]*SpotInstanceTypeDetails `json:"Linux"`
}

type SpotAdvisorResponseWrap struct {
	Regions map[string]*SpotAdvisorRegion `json:"spot_advisor"`
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
		fmt.Println("Error fetching spot price", err)
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Fail to read response body", err)
		return err
	}

	body = body[9 : len(body)-2]

	var priceWrap SpotPriceResponseWrap

	err = json.Unmarshal(body, &priceWrap)
	if err != nil {
		fmt.Println("Cannot parse json from spot request response", err)
		return err
	}

	price := priceWrap.Config

	for _, r := range price.Regions {
		region := s.SpotRegionName(r.Region)
		s.pricePerRegions[region] = make(map[string]*SpotPrice)

		for _, t := range r.InstanceTypes {
			for _, size := range t.Sizes {
				s.pricePerRegions[region][size.Size] = &SpotPrice{}
				for _, vc := range size.ValueColumns {
					if vc.RawPrice.USD == "N/A*" {
						continue
					}

					if vc.Name == "linux" {
						if v, err := vc.RawPrice.Price(); err == nil {
							s.pricePerRegions[region][size.Size].Linux = new(float64)
							*s.pricePerRegions[region][size.Size].Linux = v
						}
					}

					if vc.Name == "mswin" {
						if v, err := vc.RawPrice.Price(); err == nil {
							s.pricePerRegions[region][size.Size].MSWin = new(float64)
							*s.pricePerRegions[region][size.Size].MSWin = v
						}
					}
				}

			}
		}
	}

	//fmt.Printf("Spot Price %+v", s.pricePerRegions)
	fmt.Println("Fetched spot price in", time.Now().Sub(t0), "at", time.Now())

	return nil
}

func (s *SpotPriceCrawler) FetchSpotAdvisor() error {
	t0 := time.Now()
	resp, err := s.client.Get(AWSSpotAdvisorDataUrl)
	if err != nil {
		fmt.Println("Error fetching spot advisor data", err)
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Fail to read response body", err)
		return err
	}

	var spotAdvisorWrap SpotAdvisorResponseWrap

	err = json.Unmarshal(body, &spotAdvisorWrap)
	if err != nil {
		fmt.Println("Cannot parse json from spot request response", err)
		return err
	}

	regions := spotAdvisorWrap.Regions

	for region, instanceTypes := range regions {
		fmt.Println("At region", region)
		for instanceType, instanceDetails := range instanceTypes.Windows {
			s.pricePerRegions[region][instanceType].MSWinSavings = new(int)
			*s.pricePerRegions[region][instanceType].MSWinSavings = instanceDetails.Savings
			s.pricePerRegions[region][instanceType].MSWinReclaimRate = new(int)
			*s.pricePerRegions[region][instanceType].MSWinReclaimRate = instanceDetails.ReclaimRate + 1
		}
		for instanceType, instanceDetails := range instanceTypes.Linux {
			s.pricePerRegions[region][instanceType].LinuxSavings = new(int)
			*s.pricePerRegions[region][instanceType].LinuxSavings = instanceDetails.Savings
			s.pricePerRegions[region][instanceType].LinuxReclaimRate = new(int)
			*s.pricePerRegions[region][instanceType].LinuxReclaimRate = instanceDetails.ReclaimRate + 1
		}
	}

	//fmt.Printf("Spot Price %+v", s.pricePerRegions)
	fmt.Println("Fetched spot advisor data in", time.Now().Sub(t0), "at", time.Now())

	return nil
}

func (s *SpotPriceCrawler) SpotRegionName(region string) string {
	// The AWS API we're using has some funky names for some regions; e.g. eu-west-1 is referred to as eu-ireland
	// This function maps region name from this API call to actual region names
	spotRegionMap := map[string]string{
		"us-east":    "us-east-1",
		"us-west":    "us-west-1",
		"eu-ireland": "eu-west-1",
		"apac-sin":   "ap-southeast-1",
		"apac-syd":   "ap-southeast-2",
		"apac-tokyo": "ap-northeast-1",
	}
	fixedRegionName, found := spotRegionMap[region]
	if found {
		return fixedRegionName
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
	s.FetchSpotAdvisor()
}
