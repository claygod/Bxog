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
// Copyright © 2016-2018 Eduard Sesigin. Contacts: <claygod@yandex.ru>

import (
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestRoutingCore(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()
	muxx := New()
	muxx.Add("/", func(rw http.ResponseWriter, req *http.Request, r *Router) { req.Method = "CORE!" }).Method("GET")
	muxx.Test()
	muxx.ServeHTTP(res, req)

	if req.Method != "CORE!" {
		t.Error("Error in first url ('/')")
	}
}

func TestRouting(t *testing.T) {
	req, _ := http.NewRequest("GET", "/b/12345", nil)
	res := httptest.NewRecorder()
	muxx := New()
	muxx.Add("/a/:par", func(rw http.ResponseWriter, req *http.Request, r *Router) { req.Method = "ERR" }).Method("GET")
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
	muxx.Add("/a/:par", func(rw http.ResponseWriter, req *http.Request, r *Router) { rw.WriteHeader(777) }).Method("GET")
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
	muxx.Add("/a/:par", func(rw http.ResponseWriter, req *http.Request, r *Router) { rw.WriteHeader(777) }).Method("GET")
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
	muxx.Add("/a", func(rw http.ResponseWriter, req *http.Request, r *Router) { rw.WriteHeader(777) }).Method("GET")
	muxx.Add("/a/b", func(rw http.ResponseWriter, req *http.Request, r *Router) { rw.WriteHeader(778) }).Method("GET")
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
	muxx.Add("/a", func(rw http.ResponseWriter, req *http.Request, r *Router) { rw.WriteHeader(777) }).Method("GET")
	muxx.Add("/a/:par", func(rw http.ResponseWriter, req *http.Request, r *Router) { rw.WriteHeader(778) }).Method("GET")
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
	muxx.Add("/ab", func(rw http.ResponseWriter, req *http.Request, r *Router) { rw.WriteHeader(700) })
	muxx.Add("/abc", func(rw http.ResponseWriter, req *http.Request, r *Router) { rw.WriteHeader(701) })
	muxx.Add("/abcd", func(rw http.ResponseWriter, req *http.Request, r *Router) { rw.WriteHeader(702) })
	muxx.Test()
	muxx.ServeHTTP(res, req)

	if res.Code != 701 {
		t.Error("It does not work the method of default GET ", res.Code)
	}
}

func TestGetParam(t *testing.T) {
	req, _ := http.NewRequest("GET", "/abc/123", nil)
	res := httptest.NewRecorder()
	muxx := New()
	muxx.Add("/abc/:par", func(w http.ResponseWriter, req *http.Request, r *Router) {
		params := r.Params(req, "/abc/:par")
		req.Method = params["par"]

	})

	muxx.Test()
	muxx.ServeHTTP(res, req)

	if req.Method != "123" {
		t.Error("Error get param")
	}
}

func TestCreateUrl(t *testing.T) {
	req, _ := http.NewRequest("GET", "/abc/f", nil)
	res := httptest.NewRecorder()
	muxx := New()
	muxx.Add("/abc/:par", func(rw http.ResponseWriter, req *http.Request, r *Router) {}).Id("test")
	muxx.Test()
	muxx.ServeHTTP(res, req)

	if muxx.Create("test", map[string]string{"par": "456"}) != "/abc/456" {
		t.Error("Error creating URL")
	}
}

// Test default ID
func TestDefaultId(t *testing.T) {
	req, _ := http.NewRequest("GET", "/abc/f", nil)
	res := httptest.NewRecorder()
	muxx := New()
	muxx.Add("/abc/:par", func(rw http.ResponseWriter, req *http.Request, r *Router) {})
	muxx.Test()
	muxx.ServeHTTP(res, req)

	if muxx.Create("/abc/:par", map[string]string{"par": "456"}) != "/abc/456" {
		t.Error("Error default Id")
	}
}

// Test route "/"
func TestRouteSlash(t *testing.T) {
	req, _ := http.NewRequest("GET", "/", nil)
	res := httptest.NewRecorder()
	muxx := New()
	muxx.Add("/", func(rw http.ResponseWriter, req *http.Request, r *Router) { rw.WriteHeader(777) })
	muxx.Add("/abc", func(rw http.ResponseWriter, req *http.Request, r *Router) { rw.WriteHeader(700) })
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
	muxx.Add("/abc/:par1/:par2", func(rw http.ResponseWriter, req *http.Request, r *Router) {
		params := r.Params(req, "two")
		req.Method = params["par1"] + params["par2"]
	}).Id("two")
	muxx.Test()
	muxx.ServeHTTP(res, req)

	if req.Method != "p1p2" {
		t.Error("Error multiple routing variables", req.Method)
	}
}

func TestRoutingVariable(t *testing.T) {
	req, _ := http.NewRequest("GET", "/123", nil)
	res := httptest.NewRecorder()
	muxx := New()
	muxx.Add("/:abc", func(rw http.ResponseWriter, req *http.Request, r *Router) {
		params := r.Params(req, "/:abc")
		req.Method = params["abc"]
	})
	muxx.Test()
	muxx.ServeHTTP(res, req)
	//fmt.Println(req.Method)

	if req.Method != "123" {
		t.Error("Error routing variable")
	}
}

func TestSlashEnd(t *testing.T) {
	req, _ := http.NewRequest("GET", "/abc/", nil)
	res := httptest.NewRecorder()
	muxx := New()
	muxx.Add("/abc", func(rw http.ResponseWriter, req *http.Request, r *Router) { rw.WriteHeader(777) })
	muxx.Test()
	muxx.ServeHTTP(res, req)

	if res.Code == 777 {
		t.Error("Slash removing doesn't work !")
	}
}

func TestMoreRoutes(t *testing.T) {
	req, _ := http.NewRequest("GET", "/b/123/d", nil)
	res := httptest.NewRecorder()
	muxx := New()
	muxx.Add("/a/:par/d", func(rw http.ResponseWriter, req *http.Request, r *Router) { rw.WriteHeader(700) })
	muxx.Add("/b/:par/d", func(rw http.ResponseWriter, req *http.Request, r *Router) { rw.WriteHeader(701) })
	muxx.Add("/abc/def/:par", func(rw http.ResponseWriter, req *http.Request, r *Router) { rw.WriteHeader(702) })
	muxx.Test()
	muxx.ServeHTTP(res, req)

	if res.Code != 701 {
		t.Error("MORE ROUTES!! ", res.Code)
	}
}

func TestFool(t *testing.T) {
	req, _ := http.NewRequest("GET", "/a/xx/123", nil)
	res := httptest.NewRecorder()
	muxx := New()
	muxx.Add("/country/:par/money/:par", func(rw http.ResponseWriter, req *http.Request, r *Router) { rw.WriteHeader(699) })
	muxx.Add("/a/xx/:par", func(rw http.ResponseWriter, req *http.Request, r *Router) { rw.WriteHeader(700) })
	muxx.Add("/a/yy/:par", func(rw http.ResponseWriter, req *http.Request, r *Router) { rw.WriteHeader(701) })
	muxx.Add("/rtyrtyabc/def/:par", func(rw http.ResponseWriter, req *http.Request, r *Router) { rw.WriteHeader(702) })
	muxx.Test()
	muxx.ServeHTTP(res, req)

	if res.Code != 700 {
		t.Error("Enough is wrong, fool!", res.Code)
	}
}

func IHandler(w http.ResponseWriter, req *http.Request, r *Router) {
	io.WriteString(w, "Welcome to Bxog!")
}

func THandler(w http.ResponseWriter, req *http.Request, r *Router) {
	params := r.Params(req, "/abc/:par")
	io.WriteString(w, "Params:\n")
	io.WriteString(w, " 'par' -> "+params["par"]+"\n")
}
