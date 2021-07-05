package main

import "regexp"

type BranchPattern struct {
	exp *regexp.Regexp
}

func (p *BranchPattern) Match(str string) bool {
	return p.exp.MatchString(str)
}

func NewBranchPattern(pattern string) (*BranchPattern, error) {
	exp, err := regexp.Compile(pattern)
	if err != nil {
		return nil, err
	}

	return &BranchPattern{
		exp: exp,
	}, nil
}
