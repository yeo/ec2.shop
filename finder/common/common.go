package common

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"slices"
	"strings"
	"time"
)

var (
	RegionMaps       = make(map[string]*Region)
	AvailableRegions = []string{}
)

type Region struct {
	Name      string `json:"name"`
	ID        string `json:"code"`
	Type      string `json:"type"`
	Label     string `json:"label"`
	Continent string `json:"continent"`
}

// LoadRegions populate our region <-> name mapping map
func LoadRegions() error {
	client := &http.Client{
		Timeout: 30 * time.Second,
	}
	resp, err := client.Get(fmt.Sprintf("https://b0.p.awsstatic.com/locations/1.0/aws/current/locations.json?timestamp=%d", time.Now()))
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var regionData map[string]*Region
	err = json.Unmarshal(body, &regionData)
	if err != nil {
		return err
	}

	for k, v := range regionData {
		AvailableRegions = append(AvailableRegions, v.ID)
		RegionMaps[k] = &Region{
			Name:      v.Name,
			ID:        v.ID,
			Type:      v.Type,
			Label:     v.Label,
			Continent: v.Continent,
		}
	}

	slices.SortFunc(AvailableRegions, func(r1, r2 string) int {
		if strings.HasPrefix(r1, "us-") && strings.HasPrefix(r2, "us-") {
			if r1 < r2 {
				return -1
			} else if r1 == r2 {
				return 0
			} else {
				return 1
			}
		}
		if strings.HasPrefix(r1, "us-") {
			return -1
		}

		if strings.HasPrefix(r2, "us-") {
			return 1
		}

		if r1 < r2 {
			return -1
		} else if r1 == r2 {
			return 0
		} else {
			return 1
		}

	})
	return nil
}
