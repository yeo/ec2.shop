package elasticache

import (
	"maps"

	"github.com/yeo/ec2shop/finder/common"
	"github.com/yeo/ec2shop/finder/common/simpleri"
)

func Discover(elasticacheFamily, r string) map[string]*simpleri.Price {
	regionalPrice := make(map[string]*simpleri.Price)

	// build up a base array with server spec and on-demand price
	// this map hold all kind of servers including previous gen
	onDemandPrice := simpleri.LoadPriceForType(
		"./data/elasticache/elasticache.json",
		r,
		elasticacheFamily,
		func(name string) bool {
			// Elasticache use the same price for memcache+redis so no filtering
			// needed
			return true
			//return strings.Contains(name, elasticacheFamily)
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
		for _, instanceClass := range []string{
			"Standard",
			"Network%20optimized",
			"Memory%20optimized",
		} {
			filename := "./data/elasticache/elasticache-reservedinstance-" + generation + "-" + r + "-" + instanceClass + ".json"
			riPriceList, err := common.LoadPriceJsonManifest(filename)
			if err != nil {
				continue
			}

			for name, priceItem := range riPriceList.Regions[r] {
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
