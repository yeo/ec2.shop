package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strconv"
	"strings"

	"github.com/labstack/echo/v4"
)

type Attribute struct {
	EC2CapacityStatus         string `json:"aws:ec2:capacitystatus"`
	EC2ClockSpeed             string `json:"aws:ec2:clockSpeed"`
	EC2CurrentGeneration      string `json:"aws:ec2:currentGeneration"`
	EC2DedicatedEbsThroughput string `json:"aws:ec2:dedicatedEbsThroughput"`
	EC2ECU                    string `json:"aws:ec2:ecu"`
	EC2EnhancedNetworking     string `json:"aws:ec2:enhancedNetworkingSupported"`
	EC2InstanceFamily         string `json:"aws:ec2:instanceFamily"`
	EC2InstanceType           string `json:"aws:ec2:instanceType"`
	EC2LicenseModel           string `json:"aws:ec2:licenseModel"`
	EC2Memory                 string `json:"aws:ec2:memory"`
	EC2NetworkPerformance     string `json:"aws:ec2:networkPerformance"`
	EC2OperatingSystem        string `json:"aws:ec2:operatingSystem"`
	EC2PhysicalProcessor      string `json:"aws:ec2:physicalProcessor"`
	EC2ProcessorArchitecture  string `json:"aws:ec2:processorArchitecture"`
	EC2ProcessorFeatures      string `json:"aws:ec2:processorFeatures"`
	EC2Storage                string `json:"aws:ec2:storage"`
	EC2Tenancy                string `json:"aws:ec2:tenancy"`
	EC2Term                   string `json:"aws:ec2:term"`
	EC2UsageType              string `json:"aws"ec2:usagetype"`
	RawEC2VCPU                string `json:"aws:ec2:vcpu"`
	EC2VCPU                   int64  `json:"-"`

	ProductFamily string `json:"aws:productFamily"`
	Service       string `json:"aws:service"`
	SKU           string `json:"aws:sku"`
}

type RawPrice struct {
	USD string `json:"USD"`
}

func (r *RawPrice) Price() (float64, error) {
	return strconv.ParseFloat(r.USD, 64)
}

type Price struct {
	ID        string    `json:"id"`
	Unit      string    `json:"unit"`
	RawPrice  *RawPrice `json:"price"`
	Price     float64   `json:"-"`
	SpotPrice float64   `json:"-"`

	Attribute *Attribute `json:"attributes"`
}

func (p *Price) MonthlyPrice() float64 {
	// Assume 730 hours per month, similar to aws calculator https://aws.amazon.com/calculator/calculator-assumptions/
	return p.Price * 730
}

func (p *Price) FormatSpotPrice() string {
	txtSpotPrice := "NA"

	if p.SpotPrice > 0 {
		txtSpotPrice = fmt.Sprintf("%.4f", p.SpotPrice)
	}

	return txtSpotPrice
}

type FriendlyPrice struct {
	InstanceType string
	Memory       string
	VCPUS        int64
	Storage      string
	Network      string
	Cost         float64
	// This is weird because spot instance sometime have price list as NA so we use this to make it as not available
	MonthlyPrice float64
	SpotPrice    string
}

type FriendlyPriceResponse struct {
	Prices []*FriendlyPrice
}

type MetaPrice struct {
	Prices []*Price `json:"prices"`
}

type SpotPriceFinder interface {
	PriceForInstance(region string, instanceType string) (*SpotPrice, error)
}

type PriceFinder struct {
	regions         map[string][]*Price
	SpotPriceFinder SpotPriceFinder
}

func (p *PriceFinder) Load() {
	p.regions = make(map[string][]*Price)

	regions := []string{
		"af-south-1",
		"ap-south-1",
		"eu-north-1",
		"eu-west-3",
		"eu-south-1",
		"eu-west-2",
		"eu-west-1",
		"ap-northeast-3",
		"ap-northeast-2",
		"us-gov-east-1",
		"ap-northeast-1",
		"us-west-2-lax-1",
		"me-south-1",
		"ca-central-1",
		"sa-east-1",
		"ap-east-1",
		"us-gov-west-1",
		"ap-southeast-1",
		"ap-southeast-2",
		"eu-central-1",
		"us-east-1",
		"us-east-2",
		"us-west-1",
		"us-west-2",
	}

	for _, r := range regions {
		p.regions[r] = make([]*Price, 0)
		for _, generation := range []string{"ondemand", "ondemand-previous-generation"} {
			var priceList MetaPrice

			filename := "./data/" + r + "-" + generation + ".json"
			content, err := ioutil.ReadFile(filename)
			if err != nil {
				fmt.Printf("error %s %+v\n", filename, err)
				continue
			}

			err = json.Unmarshal(content, &priceList)
			if err != nil {
				fmt.Printf("error process %s %v\n", filename, err)
				continue
			}
			p.regions[r] = append(p.regions[r], priceList.Prices...)

			for i, price := range p.regions[r] {
				p.regions[r][i].Price, err = price.RawPrice.Price()

				if err != nil {
					fmt.Printf("Error when converting price %+v\n", err)
				}
				p.regions[r][i].Attribute.EC2VCPU, err = strconv.ParseInt(price.Attribute.RawEC2VCPU, 10, 64)
			}
		}
	}

}

func (p *PriceFinder) PriceListByRegion(region string) []*Price {
	return p.regions[region]
}

func (p *PriceFinder) PriceListFromRequest(c echo.Context) []*Price {
	requestRegion := c.QueryParam("region")
	if requestRegion == "" {
		requestRegion = c.QueryParam("r")
	}

	if requestRegion == "" {
		requestRegion = "us-east-1"
	}

	prices := make([]*Price, 0)

	filter := c.QueryParam("filter")
	keywords := strings.Split(filter, ",")

	for _, price := range p.PriceListByRegion(requestRegion) {
		m := price.Attribute
		matched := false
		for _, kw := range keywords {
			if strings.Contains(m.EC2InstanceType, kw) ||
				strings.Contains(m.EC2Storage, kw) ||
				strings.Contains(m.EC2NetworkPerformance, kw) {
				matched = true
			}
		}
		if !matched {
			continue
		}

		// Attempt to load spot price
		if _spotPrice, err := p.SpotPriceFinder.PriceForInstance(requestRegion, m.EC2InstanceType); err == nil {
			if _spotPrice != nil && _spotPrice.Linux != nil {
				price.SpotPrice = *_spotPrice.Linux
			}
		}

		prices = append(prices, price)
	}

	return prices
}
