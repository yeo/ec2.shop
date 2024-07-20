package ec2

import (
	"fmt"

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

func PriceFromRequest(priceData common.PriceByInstanceType[*Price], requestRegion string, keywords []*common.SearchTerm, sorters []*common.SortTerm) SearchResult {
	prices := common.PriceFromRequest(priceData, requestRegion, keywords, sorters)

	// Attempt to load spot price
	for i, price := range prices {
		m := price.Attribute
		if _spotPrice, err := spotPriceFinder.PriceForInstance(requestRegion, m.InstanceType); err == nil {
			if _spotPrice != nil && _spotPrice.Linux != nil {
				prices[i].SpotPrice = *_spotPrice.Linux
				prices[i].AdvisorSpotData = _spotPrice.AdvisorLinux
			}
		}
	}

	return prices
}

var (
	spotPriceFinder *SpotPriceFinder
)

func MonitorSpot() {
	spotPriceFinder = NewSpotPriceFinder()
	spotPriceFinder.Run()
}
