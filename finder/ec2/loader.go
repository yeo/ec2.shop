package ec2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"maps"
	"strconv"
	"strings"

	"github.com/yeo/ec2shop/finder/common"
)

var (
	gpuDetail map[string]*GPUInfo
)

func LoadPriceForType(r, generation string) map[string]*Price {
	filename := "./data/ec2/" + r + "-" + generation + ".json"
	priceList, err := common.LoadPriceJsonManifest(filename)
	if err != nil {
		return map[string]*Price{}
	}

	itemPrices := make(map[string]*Price)
	// return price data is a 2 nested map like this
	for _, regionalPriceItems := range priceList.Regions {
		for item, priceItem := range regionalPriceItems {
			priceItem.Build()

			serverTypeParts := strings.Split(item, " ")
			price := &Price{
				ID:        fmt.Sprintf("%s.%s", serverTypeParts[0], serverTypeParts[1]),
				Attribute: priceItem,
			}

			price.Price, _ = strconv.ParseFloat(priceItem.Price, 64)

			itemPrices[price.ID] = price
		}
	}

	return itemPrices
}

func Discover(r string) map[string]*Price {
	gpuDetail, _ = LoadGPUInfo("./data/gpu/gpu.json")
	regionalPrice := make(map[string]*Price)
	// build up a base array with server spec and on-demand price
	// this map hold all kind of servers including previous gen
	for _, generation := range []string{"ondemand", "previousgen-ondemand"} {
		onDemandPrice := LoadPriceForType(r, generation)
		maps.Copy(regionalPrice, onDemandPrice)
	}

	for id, reseveredPrice := range LoadPriceForType(r, "reservedinstance-1y") {
		if _, ok := regionalPrice[id]; ok == true {
			regionalPrice[id].Reserved1y = reseveredPrice.Price
		} else {
			fmt.Println("server has reserver data but not found in on-demand", id)
		}
	}

	for id, reseveredPrice := range LoadPriceForType(r, "reservedinstance-3y") {
		if _, ok := regionalPrice[id]; ok == true {
			regionalPrice[id].Reserved3y = reseveredPrice.Price
		} else {
			fmt.Println("server has reserver data but not found in on-demand", id)
		}
	}

	for id, reseveredPrice := range LoadPriceForType(r, "reservedinstance-convertible-1y") {
		if _, ok := regionalPrice[id]; ok == true {
			regionalPrice[id].Reserved1yConveritble = reseveredPrice.Price
		} else {
			fmt.Println("server has reserver data but not found in on-demand", id)
		}
	}

	for id, reseveredPrice := range LoadPriceForType(r, "reservedinstance-convertible-3y") {
		if _, ok := regionalPrice[id]; ok == true {
			regionalPrice[id].Reserved3yConveritble = reseveredPrice.Price
		} else {
			fmt.Println("server has reserver data but not found in on-demand", id)
		}
	}
	// TODO: Add other item such as reverse
	// go p.SpotPriceFinder.Run()

	return regionalPrice
}

type GPUInfo struct {
	Core    int    `json:"core"`
	Type    string `json:"type"`
	Mem     int    `json:"mem"`
	MemUnit string `json:"mem_unit"`
}

func LoadGPUInfo(filename string) (map[string]*GPUInfo, error) {
	gpuInfo := make(map[string]*GPUInfo)

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &gpuInfo); err != nil {
		return nil, err
	}

	return gpuInfo, nil
}
