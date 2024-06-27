package finder

import (
	"strings"
)

type TermType int64

const (
	TextTermType TermType = 1
	ExprTermType TermType = 2
)

type SearchFn func(p *Price) bool
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
		st.Type = ExprTermType
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
