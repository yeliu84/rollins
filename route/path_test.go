package route

import (
	"fmt"
	"testing"
)

func TestCompilePath(t *testing.T) {
	path := "/simple/path"
	re, params, err := compilePath(path)
	if re.String() != fmt.Sprintf("^%s$", path) || len(params) != 0 || err != nil {
		t.Errorf("testing %s failed\n", path)
	}

	path = "/wrong/(path"
	re, params, err = compilePath(path)
	if err == nil {
		t.Errorf("testing %s failed\n", path)
	}

	path = "/two/identical/:label/:label/path"
	re, params, err = compilePath(path)
	if err != ErrDupParam {
		t.Errorf("testing %s failed\n", path)
	}

	path = "/some/path/:with/two/:params"
	re, params, err = compilePath(path)
	if re == nil || len(params) != 2 || err != nil {
		t.Errorf("testing %s failed\n", path)
	}
	if re.String() != "^/some/path(/([^/]+))/two(/([^/]+))$" {
		t.Errorf("testing %s failed\n", path)
	}

	path = "/some/path/:with/:eager*/param"
	re, params, err = compilePath(path)
	if re == nil || len(params) != 2 || err != nil {
		t.Errorf("testing %s failed\n", path)
	}
	if re.String() != "^/some/path(/([^/]+))(/([\\w/]+))/param$" {
		t.Errorf("testing %s failed\n", path)
	}

	path = "/some/path/:with/:optional?/param"
	re, params, err = compilePath(path)
	if re == nil || len(params) != 2 || err != nil {
		t.Errorf("testing %s failed\n", path)
	}
	if re.String() != "^/some/path(/([^/]+))(/([^/]+))?/param$" {
		t.Errorf("testing %s failed\n", path)
	}
}
