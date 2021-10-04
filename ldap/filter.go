package ldap

import (
	ber "github.com/go-asn1-ber/asn1-ber"
	goldap "github.com/go-ldap/ldap/v3"
	"strings"
)

type FilterBuilder struct {
	packet *ber.Packet
}

func NewFilterBuilder(filter string) (*FilterBuilder, error) {
	f := normalizeFilter(filter)
	if len(strings.TrimSpace(f)) == 0 {
		return &FilterBuilder{}, nil
	}
	p, err := goldap.CompileFilter(f)
	if err != nil {
		return &FilterBuilder{}, ErrInvalidFilter
	}
	return &FilterBuilder{packet: p}, nil
}

func (f *FilterBuilder) And(filterB *FilterBuilder) *FilterBuilder {
	if f.packet == nil {
		return filterB
	}
	if filterB.packet == nil {
		return f
	}
	p := ber.Encode(ber.ClassContext, ber.TypeConstructed, goldap.FilterAnd, nil, goldap.FilterMap[goldap.FilterAnd])
	p.AppendChild(f.packet)
	p.AppendChild(filterB.packet)
	return &FilterBuilder{packet: p}
}

func (f *FilterBuilder) String() (string, error) {
	if f.packet == nil {
		return "", nil
	}
	return goldap.DecompileFilter(f.packet)
}

func normalizeFilter(filter string) string {
	norFilter := strings.TrimSpace(filter)
	if len(norFilter) == 0 {
		return norFilter
	}
	if strings.HasPrefix(norFilter, "(") && strings.HasSuffix(norFilter, ")") {
		return norFilter
	}
	return "(" + norFilter + ")"
}
