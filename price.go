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
	Price          string `json:"price"`
	InstanceFamily string `json:"Instance Family"`

	RawVCPU string `json:"vCPU"`
	VCPU    int64  `json:"-"`

	InstanceType       string `json:"Instance Type"`
	Memory             string `json:"Memory"`
	Storage            string `json:"Storage"`
	NetworkPerformance string `json:"Network Performance"`

	plcOperatingSystem string `json:"plc:OperatingSystem"`
	plcInstanceFamily  string `json:"plc:InstanceFamily"`
}

type PriceMap = map[string]Attribute

type PriceManifest struct {
	Regions map[string]PriceMap `json:"regions"`
}

type RawPrice struct {
	USD string `json:"USD"`
}

func (r *RawPrice) Price() (float64, error) {
	return strconv.ParseFloat(r.USD, 64)
}

// Our own data
type Price struct {
	ID string `json:"id"`

	// RawPrice can be a float or a string or a NA
	RawPrice *RawPrice `json:"price"`

	Price     float64 `json:"-"`
	SpotPrice float64 `json:"-"`

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

type SpotPriceFinder interface {
	PriceForInstance(region string, instanceType string) (*SpotPrice, error)
}

type PriceFinder struct {
	regions         map[string][]*Price
	SpotPriceFinder SpotPriceFinder
}

// Load price from db for all regions
func (p *PriceFinder) Load() {
	p.regions = make(map[string][]*Price)

	for _, r := range availableRegions {
		p.regions[r] = make([]*Price, 0)

		for _, generation := range []string{"ondemand", "previousgen-ondemand"} {
			p.loadRegion(r, generation)
		}
	}
}

func (p *PriceFinder) loadRegion(r string, generation string) {
	fmt.Printf("load price for region %s generation %s", r, generation)
	var priceList PriceManifest

	filename := "./data/" + r + "-" + generation + ".json"
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("error %s %+v\n", filename, err)
		return
	}

	err = json.Unmarshal(content, &priceList)
	if err != nil {
		fmt.Printf("error process %s %v\n", filename, err)
		return
	}

	// price is a 2 nested map like this
	for _, regionalPriceItems := range priceList.Regions {
		for item, priceItem := range regionalPriceItems {
			fmt.Printf("%s server: %s\n", r, item)
			price := &Price{
				Attribute: &priceItem,
			}
			price.Price, _ = strconv.ParseFloat(priceItem.Price, 64)
			price.Attribute.VCPU, _ = strconv.ParseInt(priceItem.RawVCPU, 10, 64)

			p.regions[r] = append(p.regions[r], price)
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
			if strings.Contains(m.InstanceType, kw) ||
				strings.Contains(m.Storage, kw) ||
				strings.Contains(m.NetworkPerformance, kw) {
				matched = true
			}
		}
		if !matched {
			continue
		}

		// Attempt to load spot price
		if _spotPrice, err := p.SpotPriceFinder.PriceForInstance(requestRegion, m.InstanceType); err == nil {
			if _spotPrice != nil && _spotPrice.Linux != nil {
				price.SpotPrice = *_spotPrice.Linux
			}
		}

		prices = append(prices, price)
	}

	return prices
}

var (
	availableRegions = []string{
		"af-south-1",
		"ap-east-1",
		"ap-northeast-1",
		"ap-northeast-2",
		"ap-northeast-3",
		"ap-south-1",
		"ap-south-2",
		"ap-southeast-1",
		"ap-southeast-2",
		"ap-southeast-3",
		"ap-southeast-4",
		"ca-central-1",
		"ca-west-1",
		"eu-central-1",
		"eu-central-2",
		"eu-north-1",
		"eu-south-1",
		"eu-south-2",
		"eu-west-1",
		"eu-west-2",
		"eu-west-3",
		"il-central-1",
		"me-central-1",
		"me-south-1",
		"sa-east-1",
		"us-east-1",
		"us-east-2",
		"us-east-2-mci-1",
		"us-gov-east-1",
		"us-gov-west-1",
		"us-west-1",
		"us-west-2",
		"ap-northeast-1-wl1-kix1",
		"ap-northeast-1-wl1-nrt1",
		"ap-northeast-2-wl1-cjj1",
		"ap-northeast-2-wl1-sel1",
		"ca-central-1-wl1-yto1",
		"eu-central-1-wl1-ber1",
		"eu-central-1-wl1-dtm1",
		"eu-central-1-wl1-muc1",
		"eu-west-2-wl1-lon1",
		"eu-west-2-wl1-man1",
		"eu-west-2-wl2-man1",
		"us-east-1-wl1",
		"us-east-1-wl1-atl1",
		"us-east-1-wl1-bna1",
		"us-east-1-wl1-chi1",
		"us-east-1-wl1-clt1",
		"us-east-1-wl1-dfw1",
		"us-east-1-wl1-dtw1",
		"us-east-1-wl1-iah1",
		"us-east-1-wl1-mia1",
		"us-east-1-wl1-msp1",
		"us-east-1-wl1-nyc1",
		"us-east-1-wl1-tpa1",
		"us-east-1-wl1-was1",
		"us-west-2-wl1",
		"us-west-2-wl1-den1",
		"us-west-2-wl1-las1",
		"us-west-2-wl1-lax1",
		"us-west-2-wl1-phx1",
		"us-west-2-wl1-sea1",
		"af-south-1-los-1",
		"ap-northeast-1-tpe-1",
		"ap-south-1-ccu-1",
		"ap-south-1-del-1",
		"ap-southeast-1-bkk-1",
		"ap-southeast-1-mnl-1",
		"ap-southeast-2-akl-1",
		"ap-southeast-2-per-1",
		"eu-central-1-ham-1",
		"eu-central-1-waw-1",
		"eu-north-1-cph-1",
		"eu-north-1-hel-1",
		"me-south-1-mct-1",
		"us-east-1-atl-1",
		"us-east-1-bos-1",
		"us-east-1-bue-1",
		"us-east-1-chi-1",
		"us-east-1-dfw-1",
		"us-east-1-iah-1",
		"us-east-1-lim-1",
		"us-east-1-mci-1",
		"us-east-1-mia-1",
		"us-east-1-msp-1",
		"us-east-1-nyc-1",
		"us-east-1-phl-1",
		"us-east-1-qro-1",
		"us-east-1-scl-1",
		"us-west-2-den-1",
		"us-west-2-las-1",
		"us-west-2-lax-1",
		"us-west-2-pdx-1",
		"us-west-2-phx-1",
		"us-west-2-sea-1",
	}
)
