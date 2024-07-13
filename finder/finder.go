package finder

import (
	"sync"

	"github.com/labstack/echo/v4"

	"github.com/yeo/ec2shop/finder/common"
	"github.com/yeo/ec2shop/finder/ec2"
	"github.com/yeo/ec2shop/finder/rds"
)

type PriceByService struct {
	EC2 ec2.PriceByInstanceType

	RDS        rds.PriceByInstanceType
	RDSMariaDB rds.PriceByInstanceType
	RDSMySQL   rds.PriceByInstanceType
}

type PriceFinder struct {
	Regions map[string]*PriceByService
}

func New() *PriceFinder {
	p := &PriceFinder{
		Regions: make(map[string]*PriceByService),
	}

	for _, r := range common.AvailableRegions {
		p.Regions[r] = &PriceByService{
			//EC2: make(map[string]*ec2.Price),
			EC2: make(map[string]*ec2.Price),
		}
	}

	return p
}

// Load price from db for all regions
func (p *PriceFinder) Discover() {
	// Load price for all supported service

	var wg sync.WaitGroup
	for r, _ := range p.Regions {
		wg.Add(1)
		go func(loadedRegion string) {
			defer wg.Done()
			p.Regions[loadedRegion].EC2 = ec2.Discover(loadedRegion)
			p.Regions[loadedRegion].RDS = rds.Discover("rds-postgresql", loadedRegion)
			p.Regions[loadedRegion].RDSMariaDB = rds.Discover("rds-mariadb", loadedRegion)
			p.Regions[loadedRegion].RDSMySQL = rds.Discover("rds-mysql", loadedRegion)
		}(r)
	}
	wg.Wait()

	go ec2.MonitorSpot()
}

func (p *PriceFinder) SearchPriceFromRequest(c echo.Context) common.SearchResult {
	requestRegion := c.QueryParam("region")
	if requestRegion == "" {
		requestRegion = c.QueryParam("r")
	}

	if requestRegion == "" {
		requestRegion = "us-east-1"
	}

	awsSvc := c.Param("svc")
	if awsSvc == "" {
		awsSvc = "ec2"
	}

	keywords := common.ParseSearchTerm(c.QueryParam("filter"))
	sorters := common.ParseSortTerm(c.QueryParam("sort"))

	switch awsSvc {
	case "rds":
		return rds.PriceFromRequest(p.Regions[requestRegion].RDS, requestRegion, keywords, sorters)
	case "rds-mariadb":
		return rds.PriceFromRequest(p.Regions[requestRegion].RDSMariaDB, requestRegion, keywords, sorters)
	case "rds-mysql":
		return rds.PriceFromRequest(p.Regions[requestRegion].RDSMySQL, requestRegion, keywords, sorters)

	case "ec2":
		return ec2.PriceFromRequest(p.Regions[requestRegion].EC2, requestRegion, keywords, sorters)
	}

	return ec2.PriceFromRequest(p.Regions[requestRegion].EC2, requestRegion, keywords, sorters)
}
