package tin

import (
	"testing"
)

func TestRouting_WrongURLPrefix(t *testing.T) {

	router := newTinRouter(New())
	router.add(extractPath("/aaa/:bbb"), "GET", nil)

	if route, _, err := router.findRoute("/v2/aaa/p1", "GET"); err != nil {
		t.Fatalf("Expect no error to be raised")

	} else if route != nil {
		t.Fatalf("Expect route not to be found, but got " + route.path.url)

	}

}

func TestRouting_WrongURLSuffix(t *testing.T) {

	router := newTinRouter(New())
	router.add(extractPath("/aaa/:bbb"), "GET", nil)

	if route, _, err := router.findRoute("/aaa/p1/alpha", "GET"); err != nil {
		t.Fatalf("Expect no error to be raised")

	} else if route != nil {
		t.Fatalf("Expect route not to be found, but got " + route.path.url)

	}

}

func TestRouting_TrailingSlash(t *testing.T) {

	router := newTinRouter(New())
	router.add(extractPath("/aaa/:bbb"), "GET", nil)

	if route, _, err := router.findRoute("/aaa/p1/", "GET"); err != nil {
		t.Fatalf("Expect no error to be raised")

	} else if route != nil {
		t.Fatalf("Expect route not to be found, but got " + route.path.url)

	}

}

func TestRouting_HeadSlash(t *testing.T) {

	router := newTinRouter(New())
	router.add(extractPath("/aaa/:bbb"), "GET", nil)

	if route, _, err := router.findRoute("//aaa/p1/", "GET"); err != nil {
		t.Fatalf("Expect no error to be raised")

	} else if route != nil {
		t.Fatalf("Expect route not to be found, but got " + route.path.url)

	}

}

func TestRouting_WrongMethod(t *testing.T) {

	router := newTinRouter(New())
	router.add(extractPath("/aaa/:bbb"), "GET", nil)

	if _, valid, err := router.findRoute("/aaa/p1", "DELETE"); err != nil || valid {
		t.Fatalf("Expect error to be raised")

	}

}

func TestRouting_CorrectURL_WithParams(t *testing.T) {

	router := newTinRouter(New())
	router.add(extractPath("/aaa/:bbb"), "GET", nil)

	if route, _, err := router.findRoute("/aaa/p1", "GET"); err != nil {
		t.Fatalf("Expect no error to be raised")

	} else if route == nil {
		t.Fatalf("Expect route to be found")

	} else if route.path.param("/aaa/p1", "bbb") != "p1" {
		t.Fatalf("Expect param to be p1")

	}

}

func TestRouting_CorrectURL_WithComplexParams(t *testing.T) {

	router := newTinRouter(New())

	route1 := extractPath("/aaa/bbb/:server/ccc")
	router.add(route1, "GET", nil)

	route2 := extractPath("/aaa/bbb/:server/:name")
	router.add(route2, "GET", nil)

	route3 := extractPath("/aaa/bbb/:server")
	router.add(route3, "GET", nil)

	if route, _, _ := router.findRoute("/aaa/bbb/alpha/ccc", "GET"); route.path != route1 {
		t.Fatalf("Expect to get route1")
	}

	if route, _, _ := router.findRoute("/aaa/bbb/alpha/other", "GET"); route.path != route2 {
		t.Fatalf("Expect to get route2")
	}

	if route, _, _ := router.findRoute("/aaa/bbb/beta", "GET"); route.path != route3 {
		t.Fatalf("Expect to get route3")
	}
}
