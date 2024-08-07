// packageri provides a common primitive to load price for resource that just
// have a price, no az or ative/standby setup and also suppoer reserver instance
package simpleri

import (
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
