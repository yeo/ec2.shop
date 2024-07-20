package common

import (
	"fmt"
	"slices"
	"strings"
)

func PriceFromRequest[T Inventory](priceData map[string]T, requestRegion string, keywords []*SearchTerm, sorters []*SortTerm) []T {
	prices := make([]T, 0)

	for _, price := range priceData {
		m := price.GetAttribute()
		// when search query is empty, match everything
		if len(keywords) == 0 {
			prices = append(prices, price)
			continue
		}

		// We start, default to a not match, and looking for item that has a
		// match
		matched := false

		for _, kw := range keywords {
			if kw.IsText() {
				switch kw.TextOp {

				case ExcludeOpType:
					fmt.Println("evaluate exclude")
					if strings.Contains(strings.ToLower(m.InstanceType), kw.Text()) {
						matched = false
						break
					}

					// when it's long enough, we look into storage, to exclude
					// thing ssd/ebs
					if len(kw.Text()) >= 3 {
						if strings.Contains(strings.ToLower(m.Storage), kw.Text()) {
							matched = false
							break
						}
					}

					// if we reach here, the check satitsifed
					matched = true

				case IncludeOpType:
					fmt.Println("evaluate include", strings.Contains(strings.ToLower(m.InstanceType), kw.Text()), m.InstanceType, kw.Text())
					if len(kw.Text()) < 3 {
						if strings.Contains(strings.ToLower(m.Family), kw.Text()) {
							matched = true
							break
						}
					}

					if strings.Contains(strings.ToLower(m.InstanceType), kw.Text()) {
						matched = true
						break
					}

					if strings.Contains(strings.ToLower(m.Storage), kw.Text()) {
						matched = true
						// For text base, we do an OR, therefore we bait as soon as
						// we matched
						break
					}
				}
			}
		}

		if !matched {
			// bail early if the keyword isn't a match
			continue
		}

		// now , narrow down the result with expression
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
