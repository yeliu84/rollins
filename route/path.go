package route

import (
	"errors"
	"fmt"
	"regexp"
	"strings"
)

const (
	namedSub         = `(/([^/]+))`
	namedEagerSub    = `(/([\w/]+))`
	namedOptionalSub = `(/([^/]+))?`
)

var (
	ErrDupParam = errors.New("duplicate path parameter error")
	pathParamRe = regexp.MustCompile(`/:([A-Za-z]\w*)([?*])?`)
)

type pathParam struct {
	match string
	label string
	flag  string
}

func (pp *pathParam) getSub() string {
	switch pp.flag {
	case "*":
		return namedEagerSub
	case "?":
		return namedOptionalSub
	default:
		return namedSub
	}
}

type pathParams []pathParam

func (params pathParams) containsLabel(label string) bool {
	for _, p := range params {
		if p.label == label {
			return true
		}
	}
	return false
}

func (params pathParams) extractValues(pathRe *regexp.Regexp, path string) (map[string]string, []string) {
	matches := pathRe.FindStringSubmatch(path)
	if matches != nil {
		values := make(map[string]string)
		for i, value := range matches[1:] {
			if i%2 != 0 {
				p := params[(i-1)/2]
				values[p.label] = value
			}
		}
		return values, matches
	}
	return nil, nil
}

func compilePath(path string) (*regexp.Regexp, pathParams, error) {
	path = fmt.Sprintf("^%s$", path)
	m := pathParamRe.FindAllStringSubmatch(path, -1)
	if m == nil {
		re, err := regexp.Compile(path)
		if err != nil {
			return nil, nil, err
		}
		return re, pathParams{}, nil
	}
	// build path parameters list
	params := make(pathParams, 0, len(m))
	for _, param := range m {
		label := param[1]
		if params.containsLabel(label) {
			return nil, nil, ErrDupParam
		}
		params = append(params, pathParam{
			match: param[0],
			label: label,
			flag:  param[2],
		})
	}
	// transform path to regular expression
	for _, param := range params {
		path = strings.Replace(path, param.match, param.getSub(), -1)
	}
	re, err := regexp.Compile(path)
	return re, params, err
}
