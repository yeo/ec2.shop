package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"

	"github.com/jedib0t/go-pretty/table"
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
	EC2VCPU                   string `json:"aws:ec2:vcpu"`

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
	Price float64 `json:"-"`

	Attribute *Attribute `json:"attributes"`
}

type MetaPrice struct {
	Prices []Price `json:"prices"`
}

type PriceFinder struct {
	regions map[string][]Price
}

func (p *PriceFinder) Load() {
	p.regions = make(map[string][]Price)

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
		var priceList MetaPrice

		p.regions[r] = make([]Price, 0)

		content, err := ioutil.ReadFile("./data/" + r + "-ondemand.json")
		if err != nil {
			fmt.Println("error %v", err)
			continue
		}

		err = json.Unmarshal(content, &priceList)
		if err != nil {
			fmt.Println("error %v", err)
			continue
		}
		p.regions[r] = priceList.Prices

		for i, price := range p.regions[r] {
			p.regions[r][i].Price, err = strconv.ParseFloat(price.RawPrice.USD, 64)
			if err != nil {
				fmt.Println("Error when converting price", err)
			}
		}
	}

}

func (p *PriceFinder) PriceListByRegion(region string) []Price {
	return p.regions[region]
}

func main() {
	p := &PriceFinder{}
	p.Load()

	// Echo instance
	e := echo.New()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	// Routes
	e.GET("/", GetPriceHandler(p))

	// Start server

	listen_on := os.Getenv("BIND_TO")
	if listen_on == "" {
		listen_on = "127.0.0.1:6000"
	}
	e.Logger.Fatal(e.Start(listen_on))
}

// Handler
func GetPriceHandler(p *PriceFinder) func(echo.Context) error {
	header := "│ %s%%-15s │ %s%%-12s │ %s%%4s vCPUs │ %s%%-20s │ %s%%-18s │ %s%%-10s │\n"
	colorizeHeader := fmt.Sprintf(header, Green, White, White, White, White, Red)

	pattern := "│ %s%%-15s │ %s%%-12s │ %s%%4s vCPUs │ %s%%-20s │ %s%%-18s │ %s%%-10.4f │\n"
	colorizePattern := fmt.Sprintf(pattern, Green, White, White, White, Yellow, Red)

	return func(c echo.Context) error {
		requestRegion := c.QueryParam("region")
		if requestRegion == "" {
			requestRegion = c.QueryParam("r")
		}

		if requestRegion == "" {
			requestRegion = "us-east-1"
		}

		prices := p.PriceListByRegion(requestRegion)
		if prices == nil {
			return errors.New("Invalid region")
		}

		tw := table.NewWriter()
		tw.AppendHeader(table.Row{"#", "First Name", "Last Name", "Salary"})

		priceText := "┌──────────────────────────────────────────────────────────────────────────────────────────────────────┐\n"
		priceText += fmt.Sprintf(colorizeHeader,
			"Instance Type",
			"Memory",
			"",
			"Storage",
			"Network",
			"Price")

		for _, price := range prices {
			priceText += "├──────────────────────────────────────────────────────────────────────────────────────────────────────┤\n"
			m := price.Attribute
			priceText += fmt.Sprintf(colorizePattern,
				m.EC2InstanceType,
				m.EC2Memory,
				m.EC2VCPU,
				m.EC2Storage,
				m.EC2NetworkPerformance,
				price.Price)

		}

		priceText += "└──────────────────────────────────────────────────────────────────────────────────────────────────────┘\n"

		c.Response().Header().Set("Cache-Control", "public, max-age=300, stale-while-revalidate=60, stale-if-error=10800")
		return c.String(http.StatusOK, priceText)
	}
}
