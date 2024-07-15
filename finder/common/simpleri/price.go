package simpleri

import (
	"fmt"

	"github.com/yeo/ec2shop/finder/common"
)

// Price structure for a given ec2 instance
type Price struct {
	ID string `json:"id"`

	RawPrice *common.RawPrice `json:"price"`

	Price float64 `json:"-"`

	Reserved1y        float64 `json:"-"`
	Reserved1yPartial float64 `json:"-"`
	Reserved1yAll     float64 `json:"-"`

	Reserved3y        float64 `json:"-"`
	Reserved3yPartial float64 `json:"-"`
	Reserved3yAll     float64 `json:"-"`

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

func LoadPriceForType(filename string, r string, resourceClassFamily string, filter func(string) bool) map[string]*Price {
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

		if !filter(name) {
			continue
		}

		if _, ok := itemPrices[price.ID]; !ok {
			itemPrices[price.ID] = price
		}

		itemPrices[price.ID].Price = priceItem.PriceFloat
	}

	return itemPrices
}
