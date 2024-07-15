package ec2

import (
	"fmt"
	"maps"
	"strconv"
	"strings"

	"github.com/yeo/ec2shop/finder/common"
)

// Price structure for a given ec2 instance
type Price struct {
	ID string `json:"id"`

	// RawPrice can be a float or a string or a NA
	RawPrice *common.RawPrice `json:"price"`

	Price     float64 `json:"-"`
	SpotPrice float64 `json:"-"`

	AdvisorSpotData *AdvisorInfo `json:"-"`

	Reserved1y            float64 `json:"-"`
	Reserved3y            float64 `json:"-"`
	Reserved1yConveritble float64 `json:"-"`
	Reserved3yConveritble float64 `json:"-"`

	Attribute *common.PriceAttribute `json:"attributes"`
}
type SearchResult []*Price

func (p *Price) GetAttribute() *common.PriceAttribute {
	return p.Attribute
}

func (p *Price) GetAttb(key string) float64 {
	lookup := float64(0)
	switch key {
	case "mem":
		lookup = p.Attribute.MemoryGib
	case "cpu", "vcpu", "core":
		lookup = p.Attribute.VCPUFloat
	case "price":
		lookup = p.Price
	case "spot":
		lookup = p.SpotPrice
	}

	return lookup
}

func (p *Price) SpotPriceHourly() string {
	txtSpotPrice := "NA"

	if p.SpotPrice > 0 {
		txtSpotPrice = fmt.Sprintf("%.4f", p.SpotPrice)
	}

	return txtSpotPrice
}

func LoadPriceForType(r, generation string) map[string]*Price {
	filename := "./data/ec2/" + r + "-" + generation + ".json"
	priceList, err := common.LoadPriceJsonManifest(filename)
	if err != nil {
		return map[string]*Price{}
	}

	itemPrices := make(map[string]*Price)
	// return price data is a 2 nested map like this
	for _, regionalPriceItems := range priceList.Regions {
		for item, priceItem := range regionalPriceItems {
			priceItem.Build()

			serverTypeParts := strings.Split(item, " ")
			price := &Price{
				ID:        fmt.Sprintf("%s.%s", serverTypeParts[0], serverTypeParts[1]),
				Attribute: priceItem,
			}

			price.Price, _ = strconv.ParseFloat(priceItem.Price, 64)

			itemPrices[price.ID] = price
		}
	}

	return itemPrices
}

func Discover(r string) map[string]*Price {
	regionalPrice := make(map[string]*Price)
	// build up a base array with server spec and on-demand price
	// this map hold all kind of servers including previous gen
	for _, generation := range []string{"ondemand", "previousgen-ondemand"} {
		onDemandPrice := LoadPriceForType(r, generation)
		maps.Copy(regionalPrice, onDemandPrice)
	}

	for id, reseveredPrice := range LoadPriceForType(r, "reservedinstance-1y") {
		if _, ok := regionalPrice[id]; ok == true {
			regionalPrice[id].Reserved1y = reseveredPrice.Price
		} else {
			fmt.Println("server has reserver data but not found in on-demand", id)
		}
	}

	for id, reseveredPrice := range LoadPriceForType(r, "reservedinstance-3y") {
		if _, ok := regionalPrice[id]; ok == true {
			regionalPrice[id].Reserved3y = reseveredPrice.Price
		} else {
			fmt.Println("server has reserver data but not found in on-demand", id)
		}
	}

	for id, reseveredPrice := range LoadPriceForType(r, "reservedinstance-convertible-1y") {
		if _, ok := regionalPrice[id]; ok == true {
			regionalPrice[id].Reserved1yConveritble = reseveredPrice.Price
		} else {
			fmt.Println("server has reserver data but not found in on-demand", id)
		}
	}

	for id, reseveredPrice := range LoadPriceForType(r, "reservedinstance-convertible-3y") {
		if _, ok := regionalPrice[id]; ok == true {
			regionalPrice[id].Reserved3yConveritble = reseveredPrice.Price
		} else {
			fmt.Println("server has reserver data but not found in on-demand", id)
		}
	}
	// TODO: Add other item such as reverse
	// go p.SpotPriceFinder.Run()

	return regionalPrice
}

var (
	spotPriceFinder *SpotPriceFinder
)

func MonitorSpot() {
	spotPriceFinder = NewSpotPriceFinder()
	spotPriceFinder.Run()
}
