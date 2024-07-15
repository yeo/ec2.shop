package simpleri

import (
	"maps"

	"github.com/yeo/ec2shop/finder/common"
)

type DiscoverRequest struct {
	OndemandFile string
	Region       string
	Family       string
	RiPrefixPath string

	NodeTypes []string
}

func Discover(d *DiscoverRequest) map[string]*Price {
	regionalPrice := make(map[string]*Price)

	// build up a base array with server spec and on-demand price
	// this map hold all kind of servers including previous gen
	onDemandPrice := LoadPriceForType(
		d.OndemandFile,
		d.Region,
		d.Family,
		func(name string) bool {
			return true
		})

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
