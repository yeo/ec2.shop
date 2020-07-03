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

func main() {
	debug := os.Getenv("DEBUG") == "1"

	s := NewSpotPriceCrawler()
	s.Run()

	p := &PriceFinder{
		SpotPriceFinder: s,
	}
	p.Load()

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
	e.GET("/", GetPriceHandler(debug, p))

	// Start server

	listen_on := os.Getenv("BIND_TO")
	if listen_on == "" {
		listen_on = "127.0.0.1:6001"
	}

	e.Logger.Fatal(e.Start(listen_on))
}

func GetPriceHandler(debug bool, p *PriceFinder) func(echo.Context) error {
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
					InstanceType: v.Attribute.EC2InstanceType,
					Memory:       v.Attribute.EC2Memory,
					VCPUS:        v.Attribute.EC2VCPU,
					Storage:      v.Attribute.EC2Storage,
					Network:      v.Attribute.EC2NetworkPerformance,
					Cost:         v.Price,
					MonthlyPrice: v.MonthlyPrice(),
					SpotPrice:    v.FormatSpotPrice(),
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
					m.EC2InstanceType,
					m.EC2Memory,
					m.EC2VCPU,
					m.EC2Storage,
					m.EC2NetworkPerformance,
					price.Price,
					price.MonthlyPrice(),
					price.FormatSpotPrice())
			}

			//priceText += "└──────────────────────────────────────────────────────────────────────────────────────────────────────┘\n"

			return c.String(http.StatusOK, priceText)
		}

		currentRegion := "us-east-1"
		if region := c.QueryParam("region"); region != "" {
			currentRegion = region
		}

		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"ts":            ts,
			"priceData":     prices,
			"currentRegion": currentRegion,
		})
	}
}
