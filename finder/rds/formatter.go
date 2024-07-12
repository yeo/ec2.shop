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
	Cost         float64
	// This is weird because spot instance sometime have price list as NA so we use this to make it as not available
	MonthlyPrice float64

	Reserved1yPrice            float64
	Reserved3yPrice            float64
	Reserved1yConveritblePrice float64
	Reserved3yConveritblePrice float64
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
			Network:                    v.Attribute.NetworkPerformance,
			Cost:                       v.Price,
			MonthlyPrice:               common.MonthlyPrice(v.Price),
			Reserved1yPrice:            v.Reserved1y,
			Reserved3yPrice:            v.Reserved3y,
			Reserved1yConveritblePrice: v.Reserved1yConveritble,
			Reserved3yConveritblePrice: v.Reserved3yConveritble,
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
	)

	for _, price := range p {
		m := price.Attribute
		priceText += fmt.Sprintf(pattern,
			m.InstanceType,
			m.Memory,
			m.VCPU,
			m.NetworkPerformance,
			price.Price,
			common.MonthlyPrice(price.Price))
	}

	return c.String(http.StatusOK, priceText)
}
