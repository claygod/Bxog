package bxog

// Test
// Check the operation of the different modes of the multiplexer Bxog
//
//   ░░░░██▄
//   ░░░██▀    ▐
//   ▌░███▄    ▐
//   ▌▐███░▀▄███▄▄▄██▄▄
//   ▌█████▌░░▌░░░░░░▌
//   ▌▀▀▀▌▐█░░▌░░░░░░▌
//   ▌▀▀▀▌▐█░░▌░░░░░░▌
//   ▌░░░▌░█▄ ▌░░░░░░▌
//
// 2016 Eduard Sesigin. Contacts: <claygod@yandex.ru>

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRouting(t *testing.T) {
	req, _ := http.NewRequest("GET", "/b/12345", nil)
	res := httptest.NewRecorder()
	muxx := New()
	muxx.Add("/a/:par", func(rw http.ResponseWriter, req *http.Request) { req.Method = "ERR" }).Method("GET")
	muxx.Test()
	muxx.ServeHTTP(res, req)

	if req.Method == "ERR" {
		t.Error("handler should not be called")
	}
}

func TestError404(t *testing.T) {
	req, _ := http.NewRequest("GET", "/b/12345", nil)
	res := httptest.NewRecorder()
	muxx := New()
	muxx.Add("/a/:par", func(rw http.ResponseWriter, req *http.Request) { rw.WriteHeader(777) }).Method("GET")
	muxx.Test()
	muxx.ServeHTTP(res, req)

	if res.Code == 777 {
		t.Errorf("expecting error code 404, got %v", res.Code)
	}
}

func TestRoutingMethod(t *testing.T) {
	req, _ := http.NewRequest("POST", "/a/12345", nil)
	res := httptest.NewRecorder()
	muxx := New()
	muxx.Add("/a/:par", func(rw http.ResponseWriter, req *http.Request) { rw.WriteHeader(777) }).Method("GET")
	muxx.Test()
	muxx.ServeHTTP(res, req)

	if res.Code == 777 {
		t.Error("response to a wrong method")
	}
}

// Test if the mux don't handle by prefix (static)
func TestRoutingPathStatic(t *testing.T) {
	req, _ := http.NewRequest("POST", "/a/b", nil)
	res := httptest.NewRecorder()
	muxx := New()
	muxx.Add("/a", func(rw http.ResponseWriter, req *http.Request) { rw.WriteHeader(777) }).Method("GET")
	muxx.Add("/a/b", func(rw http.ResponseWriter, req *http.Request) { rw.WriteHeader(778) }).Method("GET")
	muxx.Test()
	muxx.ServeHTTP(res, req)

	if res.Code == 777 {
		t.Error("response with the wrong path")
	}
}

// Test if the mux don't handle by prefix (dinamic)
func TestRoutingPathDinamic(t *testing.T) {
	req, _ := http.NewRequest("POST", "/a/b", nil)
	res := httptest.NewRecorder()
	muxx := New()
	muxx.Add("/a", func(rw http.ResponseWriter, req *http.Request) { rw.WriteHeader(777) }).Method("GET")
	muxx.Add("/a/:par", func(rw http.ResponseWriter, req *http.Request) { rw.WriteHeader(778) }).Method("GET")
	muxx.Test()
	muxx.ServeHTTP(res, req)

	if res.Code == 777 {
		t.Error("response with the wrong path")
	}
}

func TestDefaultMethodGet(t *testing.T) {
	req, _ := http.NewRequest("GET", "/abc", nil)
	res := httptest.NewRecorder()
	muxx := New()
	muxx.Add("/ab", func(rw http.ResponseWriter, req *http.Request) { rw.WriteHeader(700) })
	muxx.Add("/abc", func(rw http.ResponseWriter, req *http.Request) { rw.WriteHeader(701) })
	muxx.Add("/abcd", func(rw http.ResponseWriter, req *http.Request) { rw.WriteHeader(702) })
	muxx.Test()
	muxx.ServeHTTP(res, req)

	if res.Code != 701 {
		t.Error("It does not work the method of default GET")
	}
}

func TestGetParam(t *testing.T) {
	req, _ := http.NewRequest("GET", "/abc/123", nil)
	res := httptest.NewRecorder()
	muxx := New()
	muxx.Add("/abc/:par", func(w http.ResponseWriter, req *http.Request) {})

	muxx.Test()
	muxx.ServeHTTP(res, req)
	if req.Header.Get("par") != "123" {
		t.Error("Error get param")
	}
}

// Test route "/"
func TestRouteSlash(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()
	muxx := New()
	muxx.Add("/", func(rw http.ResponseWriter, req *http.Request) { rw.WriteHeader(777) })
	muxx.Add("/abc", func(rw http.ResponseWriter, req *http.Request) { rw.WriteHeader(700) })
	muxx.Test()
	muxx.ServeHTTP(res, req)
	if res.Code != 777 {
		t.Error("Error route '/'")
	}
}

func TestMultipleRoutingVariables(t *testing.T) {
	req, _ := http.NewRequest("GET", "/abc/p1/p2", nil)
	res := httptest.NewRecorder()
	muxx := New()
	muxx.Add("/abc/:par1/:par2", func(rw http.ResponseWriter, req *http.Request) {}).Id("two")
	muxx.Test()
	muxx.ServeHTTP(res, req)
	if req.Header.Get("par1") != "p1" || req.Header.Get("par2") != "p2" {
		t.Error("Error multiple routing variables")
	}
}

func TestRoutingVariable(t *testing.T) {
	req, _ := http.NewRequest("GET", "/123", nil)
	res := httptest.NewRecorder()
	muxx := New()
	muxx.Add("/:abc", func(rw http.ResponseWriter, req *http.Request) {})
	muxx.Test()
	muxx.ServeHTTP(res, req)
	if req.Header.Get("abc") != "123" {
		t.Error("Error routing variable")
	}
}

func TestSlashEnd(t *testing.T) {
	req, _ := http.NewRequest("GET", "/abc/", nil)
	res := httptest.NewRecorder()
	muxx := New()
	muxx.Add("/abc", func(rw http.ResponseWriter, req *http.Request) { rw.WriteHeader(777) })
	muxx.Test()
	muxx.ServeHTTP(res, req)
	if res.Code == 777 {
		t.Error("Slash removing doesn't work !")
	}
}

func TestLongUrl(t *testing.T) {
	req, _ := http.NewRequest("GET", "/country/USA/capital/Washington/valuta/dollar", nil)
	res := httptest.NewRecorder()
	muxx := New()
	muxx.Add("/country/:name/capital/:city/valuta/:money", func(rw http.ResponseWriter, req *http.Request) {
	})
	muxx.Test()
	muxx.ServeHTTP(res, req)
	if req.Header.Get("name") != "USA" {
		t.Error("Error long url: ", req.Header.Get("name"))
	}
}
