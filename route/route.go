package route

import (
	"errors"
	"net/http"
	"regexp"
)

var (
	ErrDupRoute           = errors.New("duplicate route error")
	ErrRouteNotFound      = errors.New("route not found error")
	ErrTooManyRoutesFound = errors.New("too many routes found error")

	routes = make(map[string]*Route)
)

type HttpMethod string

const (
	HttpGet     HttpMethod = "GET"
	HttpPut                = "PUT"
	HttpPost               = "POST"
	HttpDelete             = "DELETE"
	HttpOptions            = "OPTIONS"
	HttpAny                = "*"
)

type Params map[string]string

type Response map[string]interface{}

type Error struct {
	Status  int
	Message string
}

type RouteHandler interface {
	GetAllowedMethods() []HttpMethod
	Handle(*http.Request, Params) (Response, *Error)
}

type Route struct {
	path    string
	handler RouteHandler
	pathRe  *regexp.Regexp
	params  pathParams
}

func (r Route) Exec(req *http.Request, params Params) (Response, *Error) {
	methods := r.handler.GetAllowedMethods()
	ok := false
	if len(methods) > 0 {
		if methods[0] == HttpAny {
			ok = true
		} else {
			for _, m := range methods {
				if HttpMethod(req.Method) == m {
					ok = true
					break
				}
			}
		}
	}
	if ok {
		return r.handler.Handle(req, params)
	}
	return nil, &Error{
		Status:  http.StatusMethodNotAllowed,
		Message: "method not allowed",
	}
}

type candidate struct {
	r      *Route
	match  string
	params Params
}

func (c candidate) countNonEmptyParams() int {
	count := 0
	for _, v := range c.params {
		if v != "" {
			count++
		}
	}
	return count
}

func Add(path string, handler RouteHandler) error {
	if _, exists := routes[path]; exists {
		return ErrDupRoute
	}
	re, params, err := compilePath(path)
	if err != nil {
		return err
	}
	routes[path] = &Route{
		path:    path,
		handler: handler,
		pathRe:  re,
		params:  params,
	}
	return nil
}

func Find(path string) (*Route, Params, error) {
	// check for exact match
	route, ok := routes[path]
	if ok {
		return route, nil, nil
	}
	// build candidate list
	candidates := make([]candidate, 0)
	for _, route := range routes {
		values, m := route.params.extractValues(route.pathRe, path)
		if m != nil {
			candidates = append(candidates, candidate{
				r:      route,
				match:  m[0],
				params: values,
			})
		}
	}
	if len(candidates) < 1 {
		return nil, nil, ErrRouteNotFound
	}
	// find the best candidate
	best := candidates[0]
	for _, c := range candidates[1:] {
		bLen := len(best.match)
		cLen := len(c.match)
		if cLen > bLen {
			best = c
		} else if cLen == bLen {
			bNonEmpty := best.countNonEmptyParams()
			cNonEmpty := c.countNonEmptyParams()
			if cNonEmpty > bNonEmpty {
				best = c
			} else if cNonEmpty == bNonEmpty {
				return nil, nil, ErrTooManyRoutesFound
			}
		}
	}
	return best.r, best.params, nil
}
