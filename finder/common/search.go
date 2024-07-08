package common

import (
	"regexp"
	"strconv"
	"strings"
)

type TermType int64

const (
	TextTermType TermType = 1
	ExprTermType TermType = 2
)

var (
	termRegex = regexp.MustCompile(`(\w+)([<>=]+)(\d+(\.\d+)?)`)
)

type AttbLookup interface {
	GetAttb(string) float64
}

type SearchFn func(AttbLookup) bool
type SearchTerm struct {
	Raw  string
	Type TermType

	SearchFn SearchFn
}

func (st *SearchTerm) IsText() bool {
	return st.Type == TextTermType
}

func (st *SearchTerm) IsExpr() bool {
	return st.Type == ExprTermType
}

func (st *SearchTerm) Text() string {
	return st.Raw
}

func NewSearchTerm(term string) *SearchTerm {
	st := &SearchTerm{
		Raw:  term,
		Type: TextTermType,
	}

	//if strings.HasPrefix(st.Raw, "mem") || strings.HasPrefix(st.Raw, "vcpu") ||
	//	strings.HasPrefix(st.Raw, "price") || strings.ContainsAny(st.Raw, "=><") {
	if strings.ContainsAny(st.Raw, "=><") {
		matches := termRegex.FindStringSubmatch(st.Raw)
		// This will return 4 element.
		// Example input: mem>=32
		// will return an array [mem>=32, mem, >=, 32 ]
		if len(matches) < 4 {
			// when it's malform, we consider simple search
			return st
		}

		st.Type = ExprTermType
		st.SearchFn = func(src AttbLookup) bool {
			target, err := strconv.ParseFloat(matches[3], 64)
			if err != nil {
				return false
			}

			lookup := src.GetAttb(matches[1])

			switch matches[2] {
			case ">":
				return lookup > target
			case "=":
				return lookup == target
			case ">=":
				return lookup >= target
			case "<=":
				return lookup <= target
			case "<":
				return lookup < target
			}

			return false
		}
	}

	return st
}

func ParseSearchTerm(q string) []*SearchTerm {
	var terms []*SearchTerm

	keywords := strings.Split(strings.Trim(q, " "), ",")
	for _, raw := range keywords {
		raw = strings.Trim(raw, " ")
		if len(raw) >= 1 {
			terms = append(terms, NewSearchTerm(raw))
		}
	}

	return terms
}
