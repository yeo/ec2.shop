package common

import (
	"slices"
	"strings"
)

func PriceFromRequest[T Inventory](priceData map[string]T, requestRegion string, keywords []*SearchTerm, sorters []*SortTerm) []T {
	prices := make([]T, 0)

	for _, price := range priceData {
		m := price.GetAttribute()
		// when search query is empty, match everything
		matched := len(keywords) == 0

		for _, kw := range keywords {
			if kw.IsText() {
				if strings.Contains(strings.ToLower(m.InstanceType), kw.Text()) ||
					strings.Contains(strings.ToLower(m.Storage), kw.Text()) ||
					strings.Contains(strings.ToLower(m.NetworkPerformance), kw.Text()) {
					matched = true
					// For text base, we do an OR, therefore we bait as soon as
					// we matched
					break
				}
			}
		}

		// For expression, we do `AND` we bail as soon as we failed to match
		for _, kw := range keywords {
			if kw.IsExpr() {
				if kw.SearchFn(price) {
					matched = true
				} else {
					matched = false
					break
				}
			}
		}

		if !matched {
			continue
		}

		prices = append(prices, price)
	}

	slices.SortFunc(prices, func(a, b T) int {
		for _, t := range sorters {
			switch t.Field {
			case "price":
				if a.GetAttribute().PriceFloat < b.GetAttribute().PriceFloat {
					return -t.Direction
				} else if a.GetAttribute().PriceFloat > b.GetAttribute().PriceFloat {
					return t.Direction
				}
			case "cpu":
				if a.GetAttribute().VCPUFloat < b.GetAttribute().VCPUFloat {
					return -t.Direction
				} else if a.GetAttribute().VCPUFloat > b.GetAttribute().VCPUFloat {
					return t.Direction
				}

			case "mem":
				if a.GetAttribute().MemoryGib < b.GetAttribute().MemoryGib {
					return -t.Direction
				} else if a.GetAttribute().MemoryGib > b.GetAttribute().MemoryGib {
					return t.Direction
				}
			}
		}

		return 0
	})

	return prices
}
