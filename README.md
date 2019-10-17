Bxog is a simple and fast HTTP router for Go (HTTP request multiplexer).

[![API documentation](https://godoc.org/github.com/claygod/Bxog?status.svg)](https://godoc.org/github.com/claygod/Bxog)
[![Go Report Card](https://goreportcard.com/badge/github.com/claygod/Bxog)](https://goreportcard.com/report/github.com/claygod/Bxog)
[![Mentioned in Awesome Go](https://awesome.re/mentioned-badge.svg)](https://github.com/avelino/awesome-go)

# Usage

An example of using the multiplexer:
```go
package main

import (
	"io"
	"net/http"

	bx "github.com/claygod/Bxog"
)

// Handlers
func IHandler(w http.ResponseWriter, req *http.Request, r *bx.Router) {
	io.WriteString(w, "Welcome to Bxog!")
}
func THandler(w http.ResponseWriter, req *http.Request, r *bx.Router) {
	params := r.Params(req, "/abc/:par")
	io.WriteString(w, "Params:\n")
	io.WriteString(w, " 'par' -> "+params["par"]+"\n")
}
func PHandler(w http.ResponseWriter, req *http.Request, r *bx.Router) {
	// Getting parameters from URL
	params := r.Params(req, "country")
	io.WriteString(w, "Country:\n")
	io.WriteString(w, " 'name' -> "+params["name"]+"\n")
	io.WriteString(w, " 'capital' -> "+params["city"]+"\n")
	io.WriteString(w, " 'valuta' -> "+params["money"]+"\n")
	// Creating a URL string
	io.WriteString(w, "Creating a URL from route (This is an example of creating another URL):\n")
	io.WriteString(w, r.Create("country", map[string]string{"name": "Russia", "capital": "Moscow", "money": "rouble"}))
}

// Main
func main() {
	m := bx.New()
	m.Add("/", IHandler)
	m.Add("/abc/:par", THandler)
	m.Add("/country/:name/capital/:city/valuta/:money", PHandler).
		Id("country"). // For ease indicate the short ID
		Method("GET")  // GET method do not need to write here, it is used by default (this is an example)
	m.Test()
	m.Start(":9999")
}

```

Click URLs:
- http://localhost:9999
- http://localhost:9999/abc/123
- http://localhost:9999/country/USA/capital/Washington/valuta/dollar

# Settings

Necessary changes in the configuration of the multiplexer can be made in the configuration file [config.go](https://github.com/claygod/Bxog/blob/master/config.go)

# Perfomance

Bxog is the fastest router, showing the speed of query processing. Its speed is comparable to the speed of the popular multiplexers: Bone, Httprouter, Gorilla, Zeus. The test is done on a computer with a i3-6320 3.7GHz processor and 8 GB RAM. In short (less time, the better):

- Bxog 163 ns/op
- HttpRouter 183 ns/op
- Zeus 12302 ns/op
- GorillaMux 14928 ns/op
- GorillaPat 618 ns/op
- Bone 47333 ns/op

Detailed benchmark [here](https://github.com/claygod/BxogTest)

# API

Methods:
-  *New* - create a new multiplexer
-  *Add* - add a rule specifying the handler (the default method - GET, ID - as a string to this rule)
-  *Start* - start the server indicating the listening port
-  *Params* - extract parameters from URL
-  *Create* - generate URL of the available options
-  *Shutdown* - graceful stop the server
-  *Stop* - aggressive stop the server
-  *Test* - Start analogue (for testing only)

Example:
`
	m := bxog.New()
	m.Add("/", IHandler)
`

# Named parameters

Arguments in the rules designated route colon. Example route: */abc/:param* , where *abc* is a static section and *:param* - the dynamic section(argument).

# Static files

The directory path to the file and its nickname as part of URL specified in the configuration file. This constants *FILE_PREF* and *FILE_PATH*
