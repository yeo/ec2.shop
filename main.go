package main

import (
	"errors"
	"html/template"
	"io"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"

	"github.com/yeo/ec2shop/finder"
	"github.com/yeo/ec2shop/finder/common"
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

var (
	logger *slog.Logger
	e      *echo.Echo

	bindTo string
	debug  bool

	priceFinder *finder.PriceFinder
)

func loadEnv() {
	logger = slog.New(slog.NewTextHandler(os.Stderr, nil))
	debug = os.Getenv("DEBUG") == "1"

	bindTo = os.Getenv("BIND_TO")
	if bindTo == "" {
		bindTo = "127.0.0.1:6001"
	}

	// Echo instance
	e = echo.New()

	if err := common.LoadRegions(); err != nil {
		panic(err)
	}

	priceFinder = finder.New()
	priceFinder.Discover()

	t := &Template{
		templates: template.Must(template.ParseGlob("views/*.html")),
	}
	e.Renderer = t

}

func setupRouter() {
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Static("/static", "static")

	// Routes
	e.GET("/", GetPriceHandler(debug, priceFinder))
	e.GET("/:svc", GetPriceHandler(debug, priceFinder))
	e.GET("/:svc/", GetPriceHandler(debug, priceFinder))
}

func main() {
	loadEnv()
	setupRouter()

	e.Logger.Fatal(e.Start(bindTo))
}

func GetPriceHandler(debug bool, p *finder.PriceFinder) func(echo.Context) error {
	ts := time.Now()

	return func(c echo.Context) error {
		if debug {
			ts = time.Now()

			e.Renderer = &Template{
				templates: template.Must(template.ParseGlob("views/*.html")),
			}
		}

		awsSvc := c.Param("svc")
		if awsSvc == "" {
			awsSvc = "ec2"
		}

		// If user not select, default to us-east1
		currentRegion := "us-east-1"
		if region := c.QueryParam("region"); region != "" {
			currentRegion = region
		}

		if !debug {
			c.Response().Header().Set("Cache-Control", "public, max-age=300, stale-while-revalidate=60, stale-if-error=10800")
		}

		prices := p.SearchPriceFromRequest(c)

		if prices == nil {
			return errors.New("Invalid region")
		}

		if IsJson(c) {
			return prices.RenderJSON(c)
		}

		if IsText(c) {
			return prices.RenderText(c)
		}

		return c.Render(http.StatusOK, "index.html", map[string]interface{}{
			"ts":                ts,
			"priceData":         prices,
			"currentRegion":     currentRegion,
			"regions":           common.AvailableRegions,
			"regionIDToNames":   common.RegionIDToNames,
			"svc":               awsSvc,
			"availableServices": finder.AvailableServices,
		})
	}
}
