package msk

import (
	"fmt"
	"maps"
	"strings"

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
		if !strings.Contains(name, "RunBroker") {
			// Active Standby mq m5.4xlarge
			// Single Instance mq t3.micro
			// "Active Standby mqCRDR m5.xlarge ActiveMQ CRDR
			continue
		}
		priceItem.Build()

		nameParts := strings.Split(name, " ")
		id := nameParts[len(nameParts)-1]
		if strings.Contains(id, "-") {
			continue
		}

		fmt.Println("msk name", name, id)

		price := &activestandby.Price{
			ID:        fmt.Sprintf("kafka.%s", id),
			Attribute: priceItem,
		}

		price.Attribute.InstanceType = price.ID
		if _, ok := itemPrices[price.ID]; !ok {
			itemPrices[price.ID] = price
		}
		itemPrices[price.ID].Price = priceItem.PriceFloat
	}

	return itemPrices
}

func Discover(family, region string) map[string]*activestandby.Price {
	regionalPrice := make(map[string]*activestandby.Price)

	// build up a base array with server spec and on-demand price
	// this map hold all kind of servers including previous gen
	onDemandPrice := LoadPriceForType(
		"./data/msk/msk.json",
		region,
		family)

	maps.Copy(regionalPrice, onDemandPrice)

	return regionalPrice
}
