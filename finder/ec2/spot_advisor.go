package ec2

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"time"
)

const (
	AWSSpotAdvisorDataUrl = "https://spot-bid-advisor.s3.amazonaws.com/spot-advisor-data.json"
)

var (
	// Load the saving ranges from spot advisor data
	savingRanges = []string{
		"<5%",
		"5-10%",
		"10-15%",
		"15-20%",
		">20%",
	}
)

type AdvisorInfo struct {
	Saving  *int64 `json:"s"`
	Reclaim *int64 `json:"r"`
}

func (a *AdvisorInfo) FormatReclaim() string {
	if a.Reclaim == nil {
		return "NA"
	}

	return savingRanges[*a.Reclaim]
}

func (a *AdvisorInfo) FormatSaving() string {
	if a.Saving == nil {
		return "NA"
	}
	return fmt.Sprintf("%v%%", *a.Saving)
}

type SavingRange struct {
	Index int    `json:"index"`
	Label string `json:"label"`
	Dots  int    `json:"dots"`
	Max   int    `json:"max"`
}

type AdvisorRegionData struct {
	Linux   map[string]AdvisorInfo `json:"Linux"`
	Windows map[string]AdvisorInfo `json:"Windows"`
}
type SpotAdvisorDataResp struct {
	SpotAdvisor map[string]AdvisorRegionData `json:"spot_advisor"`
}

func (s *SpotPriceFinder) FetchAdvisor() error {
	t0 := time.Now()
	resp, err := s.client.Get(AWSSpotAdvisorDataUrl)
	if err != nil {
		log.Println("Error fetching spot price", err)
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Println("Fail to read response body", err)
		return err
	}

	var respWrap SpotAdvisorDataResp

	err = json.Unmarshal(body, &respWrap)
	if err != nil {
		log.Println("Cannot parse json from spot request response", err)
		return err
	}

	for regionName, spotAdvisorRegionalData := range respWrap.SpotAdvisor {
		for instanceType, advisorData := range spotAdvisorRegionalData.Linux {
			//log.Printf("[advisor data loader] found %d server price for region %s", total, r.Region)
			if spotPriceData, ok := s.pricePerRegions[regionName][instanceType]; ok {
				spotPriceData.AdvisorLinux = &advisorData
				s.pricePerRegions[regionName][instanceType] = spotPriceData
			}
		}
		// TODO: Integrate with mswindow data
		//for instanceType, advisorData := range spotAdvisorRegionalData.Windows {
		//	//log.Printf("[advisor data loader] found %d server price for region %s", total, r.Region)
		//	s.pricePerRegions[regionName][instanceType].AdvisorWindows = &advisorData
		//}
	}

	log.Println("Fetched advisor spot info in", time.Now().Sub(t0), "at", time.Now())

	return nil
}
