package activestandby

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
	Network      string
	Storage      string

	// On-demand price
	Cost         string
	MonthlyPrice string

	// Reserve price
	ActiveStandbyPrice string
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
			InstanceType: v.Attribute.InstanceType,
			Memory:       v.Attribute.Memory,
			VCPUS:        v.Attribute.VCPU,
			Network:      v.Attribute.NetworkPerformance,
			Storage:      v.Attribute.Storage,

			Cost:         common.ValueOrNA(v.Price),
			MonthlyPrice: common.ValueOrNA(common.MonthlyPrice(v.Price)),

			ActiveStandbyPrice: common.ValueOrNA(v.ActiveStandbyPrice),
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
		"vCPU",
		"Network",
		"Price",
		"Monthly",
		"Active Standby Price",
	)

	for _, price := range p {
		m := price.Attribute
		priceText += fmt.Sprintf(pattern,
			m.InstanceType,
			m.Memory,
			m.VCPU,
			m.NetworkPerformance,
			price.Price,
			common.MonthlyPrice(price.Price),
			common.ValueOrNA(price.ActiveStandbyPrice),
		)
	}

	return c.String(http.StatusOK, priceText)
}
