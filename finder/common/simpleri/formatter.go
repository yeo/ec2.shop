package simpleri

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

	// On-demand price
	Cost         string
	MonthlyPrice string

	// Reserve price
	Reserved1y        string
	Reserved1yPartial string
	Reserved1yAll     string

	Reserved3y        string
	Reserved3yPartial string
	Reserved3yAll     string
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

			Cost:         common.ValueOrNA(v.Price),
			MonthlyPrice: common.ValueOrNA(common.MonthlyPrice(v.Price)),

			Reserved1y:        common.ValueOrNA(v.Reserved1y),
			Reserved1yPartial: common.ValueOrNA(v.Reserved1yPartial),
			Reserved1yAll:     common.ValueOrNA(v.Reserved1yAll),

			Reserved3y:        common.ValueOrNA(v.Reserved3y),
			Reserved3yPartial: common.ValueOrNA(v.Reserved3yPartial),
			Reserved3yAll:     common.ValueOrNA(v.Reserved3yAll),
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
		"Reserved1y",
		"Reserved1y Partial",
		"Reserved1y All",
		"Reserved3y",
		"Reserved3y Partial",
		"Reserved3y All",
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
			price.Reserved1y,
			price.Reserved1yPartial,
			price.Reserved1yPartial,
		)
	}

	return c.String(http.StatusOK, priceText)
}
