package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"strconv"
	"strings"
)

type RawPrice struct {
	USD string `json:"USD"`
}

type PriceManifest struct {
	Regions map[string]PriceMap `json:"regions"`
}

type PriceAttribute struct {
	Price          string `json:"price"`
	InstanceFamily string `json:"Instance Family"`

	RawVCPU string `json:"vCPU"`
	VCPU    int64  `json:"-"`

	PriceFloat float64 `json:"-"`

	MemoryGib float64 `json:"-"`
	VCPUFloat float64 `json:"-"`

	InstanceType       string `json:"Instance Type"`
	Memory             string `json:"Memory"`
	Storage            string `json:"Storage"`
	NetworkPerformance string `json:"Network Performance"`

	// Reverse Instance price
	LeaseContractLength   string  `json:"LeaseContractLength"`
	RiUpfront             string  `json:"riupfront:PricePerUnit"`
	RiEffectiveHourlyRate float64 `json:"-"`
	RiUpfrontFloat        float64 `json:"-"`
	PurchaseOption        string  `json:"PurchaseOption"`
}
type PriceMap = map[string]*PriceAttribute

func (r *RawPrice) Price() (float64, error) {
	return strconv.ParseFloat(r.USD, 64)
}

// Build internal data structure for price to make it searchable. Such as
// convert string to float
func (a *PriceAttribute) Build() {
	gib := strings.Split(a.Memory, " ")
	if len(gib) >= 2 {
		a.MemoryGib, _ = strconv.ParseFloat(gib[0], 64)
	}
	a.VCPU, _ = strconv.ParseInt(a.RawVCPU, 10, 64)
	a.VCPUFloat = float64(a.VCPU)

	a.PriceFloat, _ = strconv.ParseFloat(a.Price, 64)

	if a.RiUpfront != "" {
		a.RiUpfrontFloat, _ = strconv.ParseFloat(a.RiUpfront, 64)
		if a.LeaseContractLength == "1yr" {
			a.RiEffectiveHourlyRate = a.PriceFloat + (a.RiUpfrontFloat / 365 / 24)
		} else if a.LeaseContractLength == "3yr" {
			a.RiEffectiveHourlyRate = a.PriceFloat + (a.RiUpfrontFloat / 365 / 24 / 3)
		}

		a.RiEffectiveHourlyRate = math.Floor(a.RiEffectiveHourlyRate*1000) / 1000
	}
}

// LoadPriceJsonManifest parses the price data on json file
func LoadPriceJsonManifest(filename string) (*PriceManifest, error) {
	content, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	var _priceList PriceManifest
	err = json.Unmarshal(content, &_priceList)

	if err != nil {
		return nil, err
	}

	var priceList PriceManifest
	priceList.Regions = make(map[string]map[string]*PriceAttribute)

	for region, value := range _priceList.Regions {
		priceList.Regions[RegionMaps[region].ID] = value
	}

	return &priceList, err
}

func ValueOrNA(v float64) string {
	if v > 0 {
		return fmt.Sprintf("%.3f", v)
	}

	return "NA"
}

func MonthlyPrice(p float64) float64 {
	// Assume 730 hours per month, similar to aws calculator https://aws.amazon.com/calculator/calculator-assumptions/
	value := p * 730

	// workaround to round a float64 to 4 decimals
	return math.Round(value*1000) / 1000
}
