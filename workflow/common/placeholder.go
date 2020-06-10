package common

import (
	"fmt"
	"strconv"
	"strings"
)

const placeholderPrefix = "$placeholder-"

// placeholderGenerator is to generate dynamically-generated placeholder strings.
type placeholderGenerator struct {
	index int
}

// NewPlaceholderGenerator returns a placeholderGenerator.
func NewPlaceholderGenerator() *placeholderGenerator {
	return &placeholderGenerator{}
}

// NextPlaceholder returns an arbitrary string to perform mock substitution of variables
func (p *placeholderGenerator) NextPlaceholder() string {
	s := fmt.Sprintf("%s%d", placeholderPrefix, p.index)
	p.index = p.index + 1
	return s
}

func (p *placeholderGenerator) IsPlaceholder(s string) bool {
	if !strings.HasPrefix(s, placeholderPrefix) {
		return false
	}
	index, err := strconv.Atoi(s[len(placeholderPrefix):])
	if err != nil {
		return false
	}
	if index < 0 || index >= p.index {
		return false
	}
	return true
}
