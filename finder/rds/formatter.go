package rds

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
	MultiAZ      string
	MultiAZ2     string

	// Reserve price
	Reserved1yPrice   string
	Reserved1yPartial string
	Reserved3yPrice   string

	ReservedMultiAZ1y        string
	ReservedMultiAZ1yPartial string
	ReservedMultiAZ3y        string
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
			MultiAZ:      common.ValueOrNA(v.MultiAZ),
			MultiAZ2:     common.ValueOrNA(v.MultiAZ2),

			Reserved1yPrice:   common.ValueOrNA(v.Reserved1y),
			Reserved1yPartial: common.ValueOrNA(v.Reserved1yPartial),
			Reserved3yPrice:   common.ValueOrNA(v.Reserved3y),

			ReservedMultiAZ1y:        common.ValueOrNA(v.ReservedMultiAZ1y),
			ReservedMultiAZ1yPartial: common.ValueOrNA(v.ReservedMultiAZ1yPartial),
			ReservedMultiAZ3y:        common.ValueOrNA(v.ReservedMultiAZ3y),
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
		"MultiAZ",
		"MultiAZ(2 standby)",
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
			price.MultiAZ,
			price.MultiAZ2)
	}

	return c.String(http.StatusOK, priceText)
}
