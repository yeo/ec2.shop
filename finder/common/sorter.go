package common

import "strings"

const (
	AscSort  int = 1
	DescSort int = -1
)

type SortTerm struct {
	Direction int
	Field     string
}

func ParseSortTerm(q string) []*SortTerm {
	var terms []*SortTerm

	keywords := strings.Split(strings.Trim(q, " "), ",")
	for _, raw := range keywords {
		raw = strings.Trim(raw, " ")
		if len(raw) <= 1 {
			continue
		}

		t := SortTerm{
			Field:     raw,
			Direction: AscSort,
		}

		if raw[0:1] == "-" {
			t.Field = raw[1:]
			t.Direction = DescSort
		}

		terms = append(terms, &t)
	}

	if len(terms) == 0 {
		terms = append(terms, &SortTerm{
			Field:     "price",
			Direction: AscSort,
		})
	}

	return terms
}
