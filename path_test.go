package tin

import (
	"net/http/httptest"
	"testing"
)

func TestPathNoParams(t *testing.T) {

	path := extractPath("/aaaa/bbbb/cccc")

	if path.url != "/aaaa/bbbb/cccc" {
		t.Fatalf("Expect path to be unchanged, but got " + path.url)
	}

	if path.params != nil {
		t.Fatalf("Expect to have no params")
	}
}

func TestPathWith1Params(t *testing.T) {

	path := extractPath("/aaaa/bbbb/:cccc")

	if path.url != "/aaaa/bbbb/[^/]*" {
		t.Fatalf("Expect path to be /aaaa/bbbb/ but got " + path.url)
	}

	if path.params == nil || len(path.params) != 1 {
		t.Fatalf("Expect to have one params")
	}

	if path.params["cccc"] != 3 {
		t.Fatalf("Expect to have one params")
	}
}

func TestPathWith2Params(t *testing.T) {

	path := extractPath("/aaaa/:bbbb/:cccc")

	if path.url != "/aaaa/[^/]*/[^/]*" {
		t.Fatalf("Expect path to be /aaaa/ but got " + path.url)
	}

	if path.params == nil || len(path.params) != 2 {
		t.Fatalf("Expect to have one params")
	}

	if path.params["cccc"] != 3 {
		t.Fatalf("Expect to have param cccc")
	}

	if path.params["bbbb"] != 2 {
		t.Fatalf("Expect to have param bbbb")
	}
}

func TestPathGetParam(t *testing.T) {

	params := extractPath("/aaaa/:bbbb/:cccc")

	req := httptest.NewRequest("GET", "/p1/p2/p3", nil)
	ctx := newContext(nil, req, params)

	val := ctx.Param("cccc")
	if val != "p3" {
		t.Fatalf("Expect cccc to be p3 but got " + val)
	}

}
