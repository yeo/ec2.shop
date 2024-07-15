package finder

import (
	"sync"

	"github.com/labstack/echo/v4"

	"github.com/yeo/ec2shop/finder/common"
	"github.com/yeo/ec2shop/finder/common/simpleri"

	"github.com/yeo/ec2shop/finder/ec2"
	"github.com/yeo/ec2shop/finder/elasticache"
	"github.com/yeo/ec2shop/finder/es"
	"github.com/yeo/ec2shop/finder/rds"
)

type PriceByService struct {
	EC2 common.PriceByInstanceType[*ec2.Price]

	//RDS        rds.PriceByInstanceType
	RDS        common.PriceByInstanceType[*rds.Price]
	RDSMariaDB common.PriceByInstanceType[*rds.Price]
	RDSMySQL   common.PriceByInstanceType[*rds.Price]

	// Elasticache
	Elasticache common.PriceByInstanceType[*simpleri.Price]

	Opensearch common.PriceByInstanceType[*simpleri.Price]
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

			p.Regions[loadedRegion].Elasticache = elasticache.Discover("Redis", loadedRegion)
			p.Regions[loadedRegion].Opensearch = es.Discover(loadedRegion)
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
		return rds.SearchResult(common.PriceFromRequest[*rds.Price](p.Regions[requestRegion].RDS, requestRegion, keywords, sorters))
	case "rds-mariadb":
		return rds.SearchResult(common.PriceFromRequest[*rds.Price](p.Regions[requestRegion].RDSMariaDB, requestRegion, keywords, sorters))
	case "rds-mysql":
		return rds.SearchResult(common.PriceFromRequest[*rds.Price](p.Regions[requestRegion].RDSMySQL, requestRegion, keywords, sorters))

	case "elasticache":
		return simpleri.SearchResult(common.PriceFromRequest[*simpleri.Price](p.Regions[requestRegion].Elasticache, requestRegion, keywords, sorters))

	case "opensearch":
		return simpleri.SearchResult(common.PriceFromRequest[*simpleri.Price](p.Regions[requestRegion].Opensearch, requestRegion, keywords, sorters))
	}

	return ec2.PriceFromRequest(p.Regions[requestRegion].EC2, requestRegion, keywords, sorters)
}
