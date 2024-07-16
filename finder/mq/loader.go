package mq

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
		priceItem.Build()

		nameParts := strings.Split(name, " ")
		var id string

		if family == "RabbitMQ" {
			if strings.Contains(name, "RabbitMQ Single Instance") {
				id = nameParts[len(nameParts)-1]
			} else {
				continue
			}
		}

		if family == "ActiveMQ" {
			if strings.Contains(name, "RabbitMQ") {
				continue
			}

			if strings.Contains(name, "Single Instance") || strings.Contains(name, "Active Standby") {
				// Active Standby mq m5.4xlarge
				// Single Instance mq t3.micro
				// "Active Standby mqCRDR m5.xlarge ActiveMQ CRDR
				id = nameParts[3]
			} else {
				continue
			}
		}

		price := &activestandby.Price{
			Attribute: priceItem,
		}
		price.ID = fmt.Sprintf("mq.%s", id)

		a := common.InstanceToAttb[id]
		if a == nil {
			// Just left its empty, we still have its id and price
			a = &common.PriceAttribute{}
		}
		price.Attribute = &common.PriceAttribute{
			InstanceType: price.ID,
			VCPU:         a.VCPU,
			MemoryGib:    a.MemoryGib,
			Memory:       a.Memory,
		}

		if _, ok := itemPrices[price.ID]; !ok {
			itemPrices[price.ID] = price
		}

		if family == "RabbitMQ" {
			itemPrices[price.ID].Price = priceItem.PriceFloat
		} else {
			if strings.Contains(name, "Active Standby") {
				itemPrices[price.ID].ActiveStandbyPrice = priceItem.PriceFloat
			} else {
				itemPrices[price.ID].Price = priceItem.PriceFloat
			}
		}
	}

	return itemPrices
}

func Discover(family, region string) map[string]*activestandby.Price {
	regionalPrice := make(map[string]*activestandby.Price)

	// build up a base array with server spec and on-demand price
	// this map hold all kind of servers including previous gen
	onDemandPrice := LoadPriceForType(
		"./data/mq/mq.json",
		region,
		family)

	maps.Copy(regionalPrice, onDemandPrice)

	return regionalPrice
}
