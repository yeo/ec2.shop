package ec2

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/yeo/ec2shop/finder/common"
)

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

	Reserved1yPrice            string
	Reserved3yPrice            string
	Reserved1yConveritblePrice string
	Reserved3yConveritblePrice string
}

type FriendlyPriceResponse struct {
	Prices []*FriendlyPrice
}

func (p SearchResult) RenderJSON(c echo.Context) error {
	formattedResp := &FriendlyPriceResponse{
		Prices: make([]*FriendlyPrice, len(p)),
	}

	for i, v := range p {
		formattedResp.Prices[i] = &FriendlyPrice{
			InstanceType:               v.Attribute.InstanceType,
			Memory:                     v.Attribute.Memory,
			VCPUS:                      v.Attribute.VCPU,
			Storage:                    v.Attribute.Storage,
			Network:                    v.Attribute.NetworkPerformance,
			Cost:                       v.Price,
			MonthlyPrice:               common.MonthlyPrice(v.Price),
			SpotPrice:                  v.SpotPriceHourly(),
			Reserved1yPrice:            common.ValueOrNA(v.Reserved1y),
			Reserved3yPrice:            common.ValueOrNA(v.Reserved3y),
			Reserved1yConveritblePrice: common.ValueOrNA(v.Reserved1yConveritble),
			Reserved3yConveritblePrice: common.ValueOrNA(v.Reserved3yConveritble),
		}

		if v.AdvisorSpotData != nil {
			formattedResp.Prices[i].SpotReclaimRate = v.AdvisorSpotData.FormatReclaim()
			formattedResp.Prices[i].SpotSavingRate = v.AdvisorSpotData.FormatSaving()
		}
	}

	return c.JSON(http.StatusOK, formattedResp)
}

func (p SearchResult) RenderText(c echo.Context) error {
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

	for _, price := range p {
		m := price.Attribute
		priceText += fmt.Sprintf(pattern,
			m.InstanceType,
			m.Memory,
			m.VCPU,
			m.Storage,
			m.NetworkPerformance,
			price.Price,
			common.MonthlyPrice(price.Price),
			price.SpotPriceHourly())
	}

	return c.String(http.StatusOK, priceText)
}
