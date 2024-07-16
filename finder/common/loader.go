package common

import (
	"encoding/json"
	"io/ioutil"
)

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
		if region == "Any" {
			continue
		}

		if _, ok := RegionMaps[region]; ok {
			priceList.Regions[RegionMaps[region].ID] = value
		}
	}

	return &priceList, err
}
