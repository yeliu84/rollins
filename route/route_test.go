package route

import (
	"net/http"
	"testing"
)

type Noop struct{}

func (n Noop) GetAllowedMethods() []HttpMethod {
	return []HttpMethod{HttpAny}
}

func (n Noop) Handle(_ *http.Request, _ Params) (Response, *Error) {
	return nil, nil
}

var noop = Noop{}

func clearRoutes() {
	for k := range routes {
		delete(routes, k)
	}
}

func TestAdd(t *testing.T) {
	clearRoutes()

	path := "/some/:path/with/:param"
	err := Add(path, noop)
	if err != nil {
		t.Fatalf("testing add %s failed", path)
	}
	if _, ok := routes[path]; !ok {
		t.Fatalf("testing add %s failed", path)
	}

	err = Add(path, noop)
	if err != ErrDupRoute {
		t.Fatalf("testing add duplicate route failed")
	}

	path = "/new/:path"
	err = Add(path, noop)
	if err != nil {
		t.Fatalf("testing add %s failed", path)
	}
	if _, ok := routes[path]; !ok {
		t.Fatalf("testing add %s failed", path)
	}
	if len(routes) != 2 {
		t.Fatalf("testing add %s failed", path)
	}
}

func TestFind(t *testing.T) {
	clearRoutes()

	paths := []string{
		"/some/simple/path",
		"/some/:path/with/:param",
		"/some/:path/with/:param/:optional?",
		"/some/:longpath*/with/:param",
	}
	for _, path := range paths {
		err := Add(path, noop)
		if err != nil {
			t.Error("adding paths failed for testing Find()")
		}
	}

	path := "/some/simple/path"
	_, _, err := Find(path)
	if err != nil {
		t.Fatalf("testing Find(\"%s\") failed\n", path)
	}

	path = "/some/value1/with/value2/value3"
	_, params, err := Find(path)
	if err != nil {
		t.Fatalf("testing Find(\"%s\") failed\n", path)
	}
	if params["path"] != "value1" || params["param"] != "value2" || params["optional"] != "value3" {
		t.Fatalf("testing Find(\"%s\") failed", path)
	}

	path = "/some/value1/abc/def/with/value2"
	_, params, err = Find(path)
	if err != nil {
		t.Fatalf("testing Find(\"%s\") failed\n", path)
	}
	if params["longpath"] != "value1/abc/def" || params["param"] != "value2" {
		t.Fatalf("testing Find(\"%s\") failed", path)
	}

	path = "/some/value1/with/value2"
	_, _, err = Find(path)
	if err != ErrTooManyRoutesFound {
		t.Fatalf("testing Find(\"%s\") failed\n", path)
	}

	path = "/prefix/some/value1/with/value2"
	_, _, err = Find(path)
	if err != ErrRouteNotFound {
		t.Fatalf("testing Find(\"%s\") failed\n", path)
	}

	path = "/another/path"
	_, _, err = Find(path)
	if err != ErrRouteNotFound {
		t.Fatalf("testing Find(\"%s\") failed\n", path)
	}
}
