package rds

import (
	"fmt"
	"maps"
	"strings"

	"github.com/yeo/ec2shop/finder/common"
)

// Price structure for a given ec2 instance
type Price struct {
	ID string `json:"id"`

	RawPrice *common.RawPrice `json:"price"`

	Price    float64 `json:"-"`
	MultiAZ  float64 `json:"-"`
	MultiAZ2 float64 `json:"-"`

	Reserved1y        float64 `json:"-"`
	Reserved1yPartial float64 `json:"-"`
	Reserved3y        float64 `json:"-"`

	ReservedMultiAZ1y        float64 `json:"-"`
	ReservedMultiAZ1yPartial float64 `json:"-"`
	ReservedMultiAZ3y        float64 `json:"-"`

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
	}

	return lookup
}

func LoadPriceForType(r, generation string) map[string]*Price {
	filename := "./data/rds/" + generation + ".json"
	priceList, err := common.LoadPriceJsonManifest(filename)
	if err != nil {
		panic(fmt.Errorf("error load json manifest: %w", err))
	}

	itemPrices := make(map[string]*Price)

	for name, priceItem := range priceList.Regions[r] {
		priceItem.Build()

		price := &Price{
			ID:        priceItem.InstanceType,
			Attribute: priceItem,
		}

		if price.ID == "" {
			continue
		}

		if _, ok := itemPrices[price.ID]; !ok {
			itemPrices[price.ID] = price
		}

		if strings.Contains(name, "Multi-AZ readable") {
			itemPrices[price.ID].MultiAZ2 = priceItem.PriceFloat
		} else if strings.Contains(name, "Single-AZ") {
			itemPrices[price.ID].Price = priceItem.PriceFloat
		} else {
			itemPrices[price.ID].MultiAZ = priceItem.PriceFloat
		}

	}

	return itemPrices
}

func Discover(rdsType, r string) map[string]*Price {
	regionalPrice := make(map[string]*Price)
	// build up a base array with server spec and on-demand price
	// this map hold all kind of servers including previous gen
	for _, generation := range []string{
		rdsType + "-ondemand",
	} {
		onDemandPrice := LoadPriceForType(r, generation)
		maps.Copy(regionalPrice, onDemandPrice)
	}

	for _, generation := range []string{
		rdsType + "-reservedinstance-multi-az-1y",
		rdsType + "-reservedinstance-multi-az-1y-partial",
		rdsType + "-reservedinstance-multi-az-3y",

		rdsType + "-reservedinstance-single-az-1y",
		rdsType + "-reservedinstance-single-az-1y-partial",
		rdsType + "-reservedinstance-single-az-3y",
	} {
		filename := "./data/rds/" + r + "-" + generation + ".json"
		riPriceList, err := common.LoadPriceJsonManifest(filename)
		if err != nil {
			continue
		}

		for _, priceItem := range riPriceList.Regions[r] {
			priceItem.Build()

			switch generation {
			case rdsType + "-reservedinstance-multi-az-1y":
				regionalPrice[priceItem.InstanceType].ReservedMultiAZ1y = priceItem.PriceFloat
			case rdsType + "-reservedinstance-multi-az-1y-partial":
				regionalPrice[priceItem.InstanceType].ReservedMultiAZ1yPartial = priceItem.RiEffectiveHourlyRate
			case rdsType + "-reservedinstance-multi-az-3y":
				regionalPrice[priceItem.InstanceType].ReservedMultiAZ3y = priceItem.RiEffectiveHourlyRate

			case rdsType + "-reservedinstance-single-az-1y":
				regionalPrice[priceItem.InstanceType].Reserved1y = priceItem.PriceFloat
			case rdsType + "-reservedinstance-single-az-1y-partial":
				regionalPrice[priceItem.InstanceType].Reserved1yPartial = priceItem.RiEffectiveHourlyRate
			case rdsType + "-reservedinstance-single-az-3y":
				regionalPrice[priceItem.InstanceType].Reserved3y = priceItem.RiEffectiveHourlyRate
			}
		}

	}

	return regionalPrice
}
