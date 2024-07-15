package simpleri

import (
	"fmt"
	"maps"

	"github.com/yeo/ec2shop/finder/common"
)

type FilterFunc func(string, *common.PriceAttribute) bool

type DiscoverRequest struct {
	OndemandFile string
	Region       string
	Family       string
	RiPrefixPath string
	FilterFunc   FilterFunc

	NodeTypes []string
}

func LoadPriceForType(filename string, r string, resourceClassFamily string, filter FilterFunc) map[string]*Price {
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

		if !filter(name, priceItem) {
			continue
		}

		if _, ok := itemPrices[price.ID]; !ok {
			itemPrices[price.ID] = price
		}

		itemPrices[price.ID].Price = priceItem.PriceFloat
	}

	return itemPrices
}

func Discover(d *DiscoverRequest) map[string]*Price {
	regionalPrice := make(map[string]*Price)

	if d.FilterFunc == nil {
		d.FilterFunc = func(_ string, _ *common.PriceAttribute) bool {
			return true
		}
	}
	// build up a base array with server spec and on-demand price
	// this map hold all kind of servers including previous gen
	onDemandPrice := LoadPriceForType(
		d.OndemandFile,
		d.Region,
		d.Family,
		d.FilterFunc)

	maps.Copy(regionalPrice, onDemandPrice)

	for _, generation := range []string{
		"1%20year-No%20Upfront",
		"1%20year-Partial%20Upfront",
		"1%20year-All%20Upfront",

		"3%20year-No%20Upfront",
		"3%20year-Partial%20Upfront",
		"3%20year-All%20Upfront",
	} {
		for _, instanceClass := range d.NodeTypes {
			filename := d.RiPrefixPath + generation + "-" + d.Region + "-" + instanceClass + ".json"
			riPriceList, err := common.LoadPriceJsonManifest(filename)
			if err != nil {
				//panic(err)
				continue
			}

			for _, priceItem := range riPriceList.Regions[d.Region] {
				priceItem.Build()

				if priceItem.InstanceType == "" {
					continue
				}

				switch generation {
				case "1%20year-No%20Upfront":
					regionalPrice[priceItem.InstanceType].Reserved1y = priceItem.PriceFloat
				case "1%20year-Partial%20Upfront":
					regionalPrice[priceItem.InstanceType].Reserved1yPartial = priceItem.RiEffectiveHourlyRate
				case "1%20year-All%20Upfront":
					regionalPrice[priceItem.InstanceType].Reserved1yAll = priceItem.RiEffectiveHourlyRate
				case "3%20year-No%20Upfront":
					regionalPrice[priceItem.InstanceType].Reserved3y = priceItem.PriceFloat
				case "3%20year-Partial%20Upfront":
					regionalPrice[priceItem.InstanceType].Reserved3yPartial = priceItem.RiEffectiveHourlyRate
				case "3%20year-All%20Upfront":
					regionalPrice[priceItem.InstanceType].Reserved3yAll = priceItem.RiEffectiveHourlyRate
				}
			}
		}

	}

	return regionalPrice
}
