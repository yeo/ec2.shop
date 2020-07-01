package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var (
	Reset  = "\033[0m"
	Red    = "\033[31m"
	Green  = "\033[32m"
	Yellow = "\033[33m"
	Blue   = "\033[34m"
	Purple = "\033[35m"
	Cyan   = "\033[36m"
	Gray   = "\033[37m"
	White  = "\033[97m"
)

type Attribute struct {
	EC2CapacityStatus         string `json:"aws:ec2:capacitystatus"`
	EC2ClockSpeed             string `json:"aws:ec2:clockSpeed"`
	EC2CurrentGeneration      string `json:"aws:ec2:currentGeneration"`
	EC2DedicatedEbsThroughput string `json:"aws:ec2:dedicatedEbsThroughput"`
	EC2ECU                    string `json:"aws:ec2:ecu"`
	EC2EnhancedNetworking     string `json:"aws:ec2:enhancedNetworkingSupported"`
	EC2InstanceFamily         string `json:"aws:ec2:instanceFamily"`
	EC2InstanceType           string `json:"aws:ec2:instanceType"`
	EC2LicenseModel           string `json:"aws:ec2:licenseModel"`
	EC2Memory                 string `json:"aws:ec2:memory"`
	EC2NetworkPerformance     string `json:"aws:ec2:networkPerformance"`
	EC2OperatingSystem        string `json:"aws:ec2:operatingSystem"`
	EC2PhysicalProcessor      string `json:"aws:ec2:physicalProcessor"`
	EC2ProcessorArchitecture  string `json:"aws:ec2:processorArchitecture"`
	EC2ProcessorFeatures      string `json:"aws:ec2:processorFeatures"`
	EC2Storage                string `json:"aws:ec2:storage"`
	EC2Tenancy                string `json:"aws:ec2:tenancy"`
	EC2Term                   string `json:"aws:ec2:term"`
	EC2UsageType              string `json:"aws"ec2:usagetype"`
	RawEC2VCPU                string `json:"aws:ec2:vcpu"`
	EC2VCPU                   int64  `json:"-"`

	ProductFamily string `json:"aws:productFamily"`
	Service       string `json:"aws:service"`
	SKU           string `json:"aws:sku"`
}

type Price struct {
	ID       string `json:"id"`
	Unit     string `json:"unit"`
	RawPrice struct {
		USD string `json:"USD"`
	} `json:"price"`
	Price        float64 `json:"-"`
	MonthlyPrice float64 `json:"-"`

	Attribute *Attribute `json:"attributes"`
}

type FriendlyPrice struct {
	InstanceType string
	Memory       string
	VCPUS        int64
	Storage      string
	Network      string
	Cost         float64
}

type FriendlyPriceResponse struct {
	Prices []*FriendlyPrice
}

type MetaPrice struct {
	Prices []*Price `json:"prices"`
}

type PriceFinder struct {
	regions map[string][]*Price
}

func (p *PriceFinder) Load() {
	p.regions = make(map[string][]*Price)

	regions := []string{
		"af-south-1",
		"ap-south-1",
		"eu-north-1",
		"eu-west-3",
		"eu-south-1",
		"eu-west-2",
		"eu-west-1",
		"ap-northeast-3",
		"ap-northeast-2",
		"us-gov-east-1",
		"ap-northeast-1",
		"us-west-2-lax-1",
		"me-south-1",
		"ca-central-1",
		"sa-east-1",
		"ap-east-1",
		"us-gov-west-1",
		"ap-southeast-1",
		"ap-southeast-2",
		"eu-central-1",
		"us-east-1",
		"us-east-2",
		"us-west-1",
		"us-west-2",
	}

	for _, r := range regions {
		p.regions[r] = make([]*Price, 0)
		for _, generation := range []string{"ondemand", "ondemand-previous-generation"} {
			var priceList MetaPrice

			filename := "./data/" + r + "-" + generation + ".json"
			content, err := ioutil.ReadFile(filename)
			if err != nil {
				fmt.Printf("error %s %+v\n", filename, err)
				continue
			}

			err = json.Unmarshal(content, &priceList)
			if err != nil {
				fmt.Printf("error process %s %v\n", filename, err)
				continue
			}
			p.regions[r] = append(p.regions[r], priceList.Prices...)

			for i, price := range p.regions[r] {
				p.regions[r][i].Price, err = strconv.ParseFloat(price.RawPrice.USD, 64)
				// Assume 730 hours per month, similar to aws calculator https://aws.amazon.com/calculator/calculator-assumptions/
				p.regions[r][i].MonthlyPrice = p.regions[r][i].Price * 730

				if err != nil {
					fmt.Printf("Error when converting price %+v\n", err)
				}
				p.regions[r][i].Attribute.EC2VCPU, err = strconv.ParseInt(price.Attribute.RawEC2VCPU, 10, 64)
			}
		}
	}

}

func (p *PriceFinder) PriceListByRegion(region string) []*Price {
	return p.regions[region]
}

func (p *PriceFinder) PriceListFromRequest(c echo.Context) []*Price {
	requestRegion := c.QueryParam("region")
	if requestRegion == "" {
		requestRegion = c.QueryParam("r")
	}

	if requestRegion == "" {
		requestRegion = "us-east-1"
	}

	prices := make([]*Price, 0)

	filter := c.QueryParam("filter")
	keywords := strings.Split(filter, ",")

	for _, price := range p.PriceListByRegion(requestRegion) {
		m := price.Attribute
		matched := false
		for _, kw := range keywords {
			if strings.Contains(m.EC2InstanceType, kw) ||
				strings.Contains(m.EC2Storage, kw) ||
				strings.Contains(m.EC2NetworkPerformance, kw) {
				matched = true
			}
		}
		if !matched {
			continue
		}

		prices = append(prices, price)
	}

	return prices
}

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

	p := &PriceFinder{}
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

func IsJson(c echo.Context) bool {
	contentType := c.Request().Header.Get("Content-Type")
	accept := c.Request().Header.Get("Accept")
	qa := c.QueryString()

	return strings.Contains(contentType, "json") || strings.Contains(accept, "json") || strings.Contains(qa, "json")
}

func IsText(c echo.Context) bool {
	ua := c.Request().Header.Get("User-Agent")
	accept := c.Request().Header.Get("Accept")

	if strings.Contains(c.QueryString(), "txt") {
		return true
	}

	if strings.Contains(accept, "html") {
		return false
	}

	if strings.Contains(ua, "Chrome") || strings.Contains(ua, "Safari") || strings.Contains(ua, "Mozilla") {
		return false
	}

	return true
}

func GetPriceHandler(debug bool, p *PriceFinder) func(echo.Context) error {
	header := "%-15s  %-12s  %4s vCPUs  %-20s  %-18s  %-10s  %-10s\n"

	pattern := "%-15s  %-12s  %4d vCPUs  %-20s  %-18s  %-10.4f  %-8.3f\n"

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
				"Monthly")

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
					price.MonthlyPrice)

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
