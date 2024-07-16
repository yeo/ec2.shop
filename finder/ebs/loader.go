package ebs

import (
	"fmt"
	"maps"

	"github.com/yeo/ec2shop/finder/common"
	"github.com/yeo/ec2shop/finder/common/activestandby"
)

func LoadPriceForType(filename string, r string, family string) map[string]*activestandby.Price {
	priceList, err := common.LoadPriceJsonManifest(filename)
	if err != nil {
		panic(fmt.Errorf("error load json manifest: %w", err))
	}

	itemPrices := make(map[string]*activestandby.Price)

	for name, priceItem := range priceList.Regions[r] {
		priceItem.Build()

		priceItem.InstanceType = name
		price := &activestandby.Price{
			ID:        name,
			Attribute: priceItem,
		}

		itemPrices[price.ID] = price
		if _, ok := itemPrices[price.ID]; ok {
			itemPrices[price.ID].Price = priceItem.PriceFloat
		} else {
			fmt.Println("missing id", price.ID)
		}
	}

	return itemPrices
}

func Discover(family, region string) map[string]*activestandby.Price {
	regionalPrice := make(map[string]*activestandby.Price)

	// build up a base array with server spec and on-demand price
	// this map hold all kind of servers including previous gen
	onDemandPrice := LoadPriceForType(
		"./data/ebs/ebs.json",
		region,
		family)

	maps.Copy(regionalPrice, onDemandPrice)

	return regionalPrice
}
