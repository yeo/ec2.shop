package ec2

import "fmt"

type FriendlyPrice struct {
	InstanceType string
	Memory       string
	VCPUS        int64
	Storage      string
	Network      string
	Cost         float64
	// This is weird because spot instance sometime have price list as NA so we use this to make it as not available
	MonthlyPrice float64

	SpotPrice       string
	SpotReclaimRate string
	SpotSavingRate  string

	Reserved1yPrice            float64
	Reserved3yPrice            float64
	Reserved1yConveritblePrice float64
	Reserved3yConveritblePrice float64
}

type FriendlyPriceResponse struct {
	Prices []*FriendlyPrice
}

func (p *PriceFinder) RenderJSON(prices []*Price) *FriendlyPriceResponse {
	formattedResp := &FriendlyPriceResponse{
		Prices: make([]*FriendlyPrice, len(prices)),
	}

	for i, v := range prices {
		formattedResp.Prices[i] = &FriendlyPrice{
			InstanceType:               v.Attribute.InstanceType,
			Memory:                     v.Attribute.Memory,
			VCPUS:                      v.Attribute.VCPU,
			Storage:                    v.Attribute.Storage,
			Network:                    v.Attribute.NetworkPerformance,
			Cost:                       v.Price,
			MonthlyPrice:               v.MonthlyPrice(),
			SpotPrice:                  v.SpotPriceHourly(),
			Reserved1yPrice:            v.Reserved1y,
			Reserved3yPrice:            v.Reserved3y,
			Reserved1yConveritblePrice: v.Reserved1yConveritble,
			Reserved3yConveritblePrice: v.Reserved3yConveritble,
		}

		if v.AdvisorSpotData != nil {
			formattedResp.Prices[i].SpotReclaimRate = v.AdvisorSpotData.FormatReclaim()
			formattedResp.Prices[i].SpotSavingRate = v.AdvisorSpotData.FormatSaving()
		}
	}

	return formattedResp
}

func (p *PriceFinder) RenderText(prices []*Price) string {
	header := "%-15s  %-12s  %4s vCPUs  %-20s  %-18s  %-10s  %-10s  %-10s\n"
	pattern := "%-15s  %-12s  %4d vCPUs  %-20s  %-18s  %-10.4f  %-10.3f  %-10s\n"

	priceText := ""
	priceText += fmt.Sprintf(header,
		"Instance Type",
		"Memory",
		"",
		"Storage",
		"Network",
		"Price",
		"Monthly",
		"Spot Price")

	for _, price := range prices {
		m := price.Attribute
		priceText += fmt.Sprintf(pattern,
			m.InstanceType,
			m.Memory,
			m.VCPU,
			m.Storage,
			m.NetworkPerformance,
			price.Price,
			price.MonthlyPrice(),
			price.SpotPriceHourly())
	}

	return priceText
}
