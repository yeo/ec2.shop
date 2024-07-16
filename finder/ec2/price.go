package ec2

import (
	"fmt"
	"slices"
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

func PriceFromRequest(priceData common.PriceByInstanceType[*Price], requestRegion string, keywords []*common.SearchTerm, sorters []*common.SortTerm) SearchResult {
	prices := make([]*Price, 0)

	for _, price := range priceData {
		m := price.Attribute
		// when search query is empty, match everything
		matched := len(keywords) == 0

		for _, kw := range keywords {
			if kw.IsText() {
				if strings.Contains(strings.ToLower(m.InstanceType), kw.Text()) ||
					strings.Contains(strings.ToLower(m.Storage), kw.Text()) ||
					strings.Contains(strings.ToLower(m.NetworkPerformance), kw.Text()) {
					matched = true
					// For text base, we do an OR, therefore we bait as soon as
					// we matched
					break
				}
			}
		}

		// For expression, we do `AND` we bail as soon as we failed to match
		for _, kw := range keywords {
			if kw.IsExpr() {
				if kw.SearchFn(price) {
					matched = true
				} else {
					matched = false
					break
				}
			}
		}

		if !matched {
			continue
		}

		// Attempt to load spot price
		if _spotPrice, err := spotPriceFinder.PriceForInstance(requestRegion, m.InstanceType); err == nil {
			if _spotPrice != nil && _spotPrice.Linux != nil {
				price.SpotPrice = *_spotPrice.Linux
				price.AdvisorSpotData = _spotPrice.AdvisorLinux
			}
		}

		prices = append(prices, price)
	}

	slices.SortFunc(prices, func(a, b *Price) int {
		for _, t := range sorters {
			switch t.Field {
			case "price":
				if a.Price < b.Price {
					return -t.Direction
				} else if a.Price > b.Price {
					return t.Direction
				}
			case "cpu":
				if a.Attribute.VCPUFloat < b.Attribute.VCPUFloat {
					return -t.Direction
				} else if a.Attribute.VCPUFloat > b.Attribute.VCPUFloat {
					return t.Direction
				}

			case "mem":
				if a.Attribute.MemoryGib < b.Attribute.MemoryGib {
					return -t.Direction
				} else if a.Attribute.MemoryGib > b.Attribute.MemoryGib {
					return t.Direction
				}
			}
		}

		return 0
	})

	return prices
}

var (
	spotPriceFinder *SpotPriceFinder
)

func MonitorSpot() {
	spotPriceFinder = NewSpotPriceFinder()
	spotPriceFinder.Run()
}
