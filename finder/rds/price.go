package rds

import (
	"fmt"
	"maps"
	"slices"
	"strconv"
	"strings"

	"github.com/yeo/ec2shop/finder/common"
)

type PriceByInstanceType = map[string]*Price

// Price structure for a given ec2 instance
type Price struct {
	ID string `json:"id"`

	RawPrice *common.RawPrice `json:"price"`

	Price float64 `json:"-"`

	Reserved1y            float64 `json:"-"`
	Reserved3y            float64 `json:"-"`
	Reserved1yConveritble float64 `json:"-"`
	Reserved3yConveritble float64 `json:"-"`

	Attribute *common.PriceAttribute `json:"attributes"`
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
	}

	return lookup
}

func LoadPriceForType(r, generation string) map[string]*Price {
	filename := "./data/rds/" + generation + ".json"
	priceList, err := common.LoadPriceJsonManifest(filename)
	if err != nil {
		panic(err)
	}

	itemPrices := make(map[string]*Price)
	fmt.Println("[rds] loaded price", r, priceList.Regions[r])

	for _, priceItem := range priceList.Regions[r] {
		priceItem.Build()

		price := &Price{
			ID:        priceItem.InstanceType,
			Attribute: priceItem,
		}

		price.Price, _ = strconv.ParseFloat(priceItem.Price, 64)

		itemPrices[price.ID] = price
	}

	return itemPrices
}

func Discover(r string) map[string]*Price {
	regionalPrice := make(map[string]*Price)
	// build up a base array with server spec and on-demand price
	// this map hold all kind of servers including previous gen
	for _, generation := range []string{"postgresql-ondemand"} {
		onDemandPrice := LoadPriceForType(r, generation)
		maps.Copy(regionalPrice, onDemandPrice)
	}

	fmt.Printf("[rds]found %d rds price for region %s\n", len(regionalPrice), r)

	return regionalPrice
}

type SearchResult []*Price

func PriceFromRequest(priceData PriceByInstanceType, requestRegion string, keywords []*common.SearchTerm, sorters []*common.SortTerm) SearchResult {
	prices := make([]*Price, 0)

	for _, price := range priceData {
		m := price.Attribute
		// when search query is empty, match everything
		matched := len(keywords) == 0

		for _, kw := range keywords {
			if kw.IsText() {
				if strings.Contains(m.InstanceType, kw.Text()) ||
					strings.Contains(m.Storage, kw.Text()) ||
					strings.Contains(m.NetworkPerformance, kw.Text()) {
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