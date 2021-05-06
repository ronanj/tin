package tin

import (
	"net/http/httptest"
	"testing"
)

func TestRoutingNoParams(t *testing.T) {

	path, params := extractPath("/aaaa/bbbb/cccc")

	if path != "/aaaa/bbbb/cccc" {
		t.Fatalf("Expect path to be unchanged")
	}

	if params != nil {
		t.Fatalf("Expect to have no params")
	}
}

func TestRoutingWith1Params(t *testing.T) {

	path, params := extractPath("/aaaa/bbbb/:cccc")

	if path != "/aaaa/bbbb/" {
		t.Fatalf("Expect path to be /aaaa/bbbb/ but got " + path)
	}

	if params == nil || len(params) != 1 {
		t.Fatalf("Expect to have one params")
	}

	if params["cccc"] != 3 {
		t.Fatalf("Expect to have one params")
	}
}

func TestRoutingWith2Params(t *testing.T) {

	path, params := extractPath("/aaaa/:bbbb/:cccc")

	if path != "/aaaa/" {
		t.Fatalf("Expect path to be /aaaa/ but got " + path)
	}

	if params == nil || len(params) != 2 {
		t.Fatalf("Expect to have one params")
	}

	if params["cccc"] != 3 {
		t.Fatalf("Expect to have param cccc")
	}

	if params["bbbb"] != 2 {
		t.Fatalf("Expect to have param bbbb")
	}
}

func TestRoutingGetParam(t *testing.T) {

	_, params := extractPath("/aaaa/:bbbb/:cccc")

	req := httptest.NewRequest("GET", "/p1/p2/p3", nil)
	tin := &Tin{}
	ctx := tin.newContext(nil, req, params)

	val := ctx.Param("cccc")
	if val != "p3" {
		t.Fatalf("Expect cccc to be p3 but got " + val)
	}

}
