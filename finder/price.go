package finder

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"maps"
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

	Reserved1y            float64 `json:"-"`
	Reserved3y            float64 `json:"-"`
	Reserved1yConveritble float64 `json:"-"`
	Reserved3yConveritble float64 `json:"-"`

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

type PriceByInstanceType = map[string]*Price
type PriceByRegion = map[string]PriceByInstanceType
type PriceFinder struct {
	regions PriceByRegion

	SpotPriceFinder *SpotPriceCrawler
}

// Load price from db for all regions
func (p *PriceFinder) Discover() {
	p.regions = make(map[string]map[string]*Price)

	for _, r := range AvailableRegions {
		regionalPrice := make(map[string]*Price)
		// build up a base array with server spec and on-demand price
		// this map hold all kind of servers including previous gen
		for _, generation := range []string{"ondemand", "previousgen-ondemand"} {
			onDemandPrice := p.loadRegion(r, generation)
			maps.Copy(regionalPrice, onDemandPrice)
		}

		for id, reseveredPrice := range p.loadRegion(r, "reservedinstance-1y") {
			if _, ok := regionalPrice[id]; ok == true {
				regionalPrice[id].Reserved1y = reseveredPrice.Price
			} else {
				fmt.Println("server has reserver data but not found in on-demand", id)
			}
		}

		for id, reseveredPrice := range p.loadRegion(r, "reservedinstance-3y") {
			if _, ok := regionalPrice[id]; ok == true {
				regionalPrice[id].Reserved3y = reseveredPrice.Price
			} else {
				fmt.Println("server has reserver data but not found in on-demand", id)
			}
		}

		for id, reseveredPrice := range p.loadRegion(r, "reservedinstance-convertible-1y") {
			if _, ok := regionalPrice[id]; ok == true {
				regionalPrice[id].Reserved1yConveritble = reseveredPrice.Price
			} else {
				fmt.Println("server has reserver data but not found in on-demand", id)
			}
		}

		for id, reseveredPrice := range p.loadRegion(r, "reservedinstance-convertible-3y") {
			if _, ok := regionalPrice[id]; ok == true {
				regionalPrice[id].Reserved3yConveritble = reseveredPrice.Price
			} else {
				fmt.Println("server has reserver data but not found in on-demand", id)
			}
		}
		p.regions[r] = regionalPrice
	}

	// TODO: Add other item such as reverse
	go p.SpotPriceFinder.Run()

}

func (p *PriceFinder) loadRegion(r string, generation string) map[string]*Price {
	fmt.Printf("load price for region %s generation %s", r, generation)
	var priceList PriceManifest

	filename := "./data/" + r + "-" + generation + ".json"
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Printf("error %s %+v\n", filename, err)
		return map[string]*Price{}
	}

	err = json.Unmarshal(content, &priceList)
	if err != nil {
		fmt.Printf("error process %s %v\n", filename, err)
		return map[string]*Price{}
	}

	itemPrices := make(map[string]*Price)
	// price is a 2 nested map like this
	for _, regionalPriceItems := range priceList.Regions {
		for item, priceItem := range regionalPriceItems {
			serverTypeParts := strings.Split(item, " ")
			price := &Price{
				ID:        fmt.Sprintf("%s.%s", serverTypeParts[0], serverTypeParts[1]),
				Attribute: &priceItem,
			}

			price.Price, _ = strconv.ParseFloat(priceItem.Price, 64)
			price.Attribute.VCPU, _ = strconv.ParseInt(priceItem.RawVCPU, 10, 64)

			itemPrices[price.ID] = price
		}
	}

	return itemPrices
}

func (p *PriceFinder) PriceListByRegion(region string) map[string]*Price {
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
