package main

import (
	"errors"
	"fmt"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/yeo/ec2shop/finder/ec2"
)

type Template struct {
	templates *template.Template
}

func (t *Template) Render(w io.Writer, name string, data interface{}, c echo.Context) error {

	// Add global methods if data is a map
	if viewContext, isMap := data.(map[string]interface{}); isMap {
		viewContext["reverse"] = c.Echo().Reverse
	}

	return t.templates.ExecuteTemplate(w, name, data)
}

type FriendlyPrice struct {
	InstanceType string
	Memory       string
	VCPUS        int64
	Storage      string
	Network      string
	Cost         float64
	// This is weird because spot instance sometime have price list as NA so we use this to make it as not available
	MonthlyPrice               float64

	SpotPrice                  string
	SpotReclaimRate string
	SpotSavingRate string

	Reserved1yPrice            float64
	Reserved3yPrice            float64
	Reserved1yConveritblePrice float64
	Reserved3yConveritblePrice float64
}
type FriendlyPriceResponse struct {
	Prices []*FriendlyPrice
}

func main() {
	debug := os.Getenv("DEBUG") == "1"

	priceFinder := ec2.New()
	priceFinder.Discover()

	// Echo instance
	e := echo.New()

	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Renderer = t

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/static", "static")

	// Routes
	e.GET("/", GetPriceHandler(debug, priceFinder))

	// Start server

	listen_on := os.Getenv("BIND_TO")
	if listen_on == "" {
		listen_on = "127.0.0.1:6001"
	}

	e.Logger.Fatal(e.Start(listen_on))
}

func GetPriceHandler(debug bool, p *ec2.PriceFinder) func(echo.Context) error {
	header := "%-15s  %-12s  %4s vCPUs  %-20s  %-18s  %-10s  %-10s  %-10s\n"

	pattern := "%-15s  %-12s  %4d vCPUs  %-20s  %-18s  %-10.4f  %-10.3f  %-10s\n"

	ts := time.Now()

	return func(c echo.Context) error {
		c.Response().Header().Set("Cache-Control", "public, max-age=300, stale-while-revalidate=60, stale-if-error=10800")
		if debug {
			ts = time.Now()
		}

		prices := p.PriceListFromRequest(c)

		if prices == nil {
			return errors.New("Invalid region")
		}

		if IsJson(c) {
			friendlyPrices := &FriendlyPriceResponse{
				Prices: make([]*FriendlyPrice, len(prices)),
			}

			for i, v := range prices {
				friendlyPrices.Prices[i] = &FriendlyPrice{
					InstanceType:               v.Attribute.InstanceType,
					Memory:                     v.Attribute.Memory,
					VCPUS:                      v.Attribute.VCPU,
					Storage:                    v.Attribute.Storage,
					Network:                    v.Attribute.NetworkPerformance,
					Cost:                       v.Price,
					MonthlyPrice:               v.MonthlyPrice(),
					SpotPrice:                  v.FormatSpotPrice(),
					Reserved1yPrice:            v.Reserved1y,
					Reserved3yPrice:            v.Reserved3y,
					Reserved1yConveritblePrice: v.Reserved1yConveritble,
					Reserved3yConveritblePrice: v.Reserved3yConveritble,
				}

				if v.AdvisorSpotData != nil {
					friendlyPrices.Prices[i].SpotReclaimRate = v.AdvisorSpotData.FormatReclaim()
					friendlyPrices.Prices[i].SpotSavingRate = v.AdvisorSpotData.FormatSaving()
				}
			}

			return c.JSON(http.StatusOK, friendlyPrices)
		}

		if IsText(c) {
			// When loading by shell we can pass these param
			priceText := ""
			//priceText += "┌──────────────────────────────────────────────────────────────────────────────────────────────────────┐\n"
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

				//priceText += "├──────────────────────────────────────────────────────────────────────────────────────────────────────┤\n"
				priceText += fmt.Sprintf(pattern,
					m.InstanceType,
					m.Memory,
					m.VCPU,
					m.Storage,
					m.NetworkPerformance,
					price.Price,
					price.MonthlyPrice(),
					price.FormatSpotPrice())
			}

			//priceText += "└──────────────────────────────────────────────────────────────────────────────────────────────────────┘\n"

			return c.String(http.StatusOK, priceText)
		}

		// If user not select, default to us-east1
		currentRegion := "us-east-1"
		if region := c.QueryParam("region"); region != "" {
			currentRegion = region
		}

		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"ts":            ts,
			"priceData":     prices,
			"currentRegion": currentRegion,
			"regions":       ec2.AvailableRegions,
		})
	}
}
