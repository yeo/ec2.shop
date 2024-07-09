package main

import (
	"errors"
	"html/template"
	"io"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/yeo/ec2shop/finder/common"
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
			return c.JSON(http.StatusOK, p.RenderJSON(prices))
		}

		if IsText(c) {
			// When loading by shell we can pass these param
			priceText := p.RenderText(prices)

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
			"regions":       common.AvailableRegions,
		})
	}
}
