Bxog is a simple and fast HTTP router for Go (HTTP request multiplexer).

[![API documentation](https://godoc.org/github.com/claygod/Bxog?status.svg)](https://godoc.org/github.com/claygod/Bxog)
[![Go Report Card](https://goreportcard.com/badge/github.com/claygod/Bxog)](https://goreportcard.com/report/github.com/claygod/Bxog)

# Usage

An example of using the multiplexer:
```go
package main

import (
	"io"
	"net/http"
	"github.com/claygod/bxog"
)

// Handlers
func IHandler(w http.ResponseWriter, req *http.Request, r *bxog.Router) {
	io.WriteString(w, "Welcome to Bxog!")
}
func THandler(w http.ResponseWriter, req *http.Request, r *bxog.Router) {
	params := r.Params(req, "/abc/:par")
	io.WriteString(w, "Params:\n")
	io.WriteString(w, " 'par' -> "+params["par"]+"\n")
}
func PHandler(w http.ResponseWriter, req *http.Request, r *bxog.Router) {
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
	m := bxog.New()
	m.Add("/", IHandler)
	m.Add("/abc/:par", THandler)
	m.Add("/country/:name/capital/:city/valuta/:money", PHandler).
		Id("country"). // For ease indicate the short ID
		Method("GET") // GET method do not need to write here, it is used by default (this is an example)
	m.Start(":80")
}
```

Click URLs:
- http://localhost
- http://localhost/abc/123
- http://localhost/country/USA/capital/Washington/valuta/dollar

# Settings

Necessary changes in the configuration of the multiplexer can be made in the configuration file [config.go](https://github.com/claygod/Bxog/blob/master/config.go)

# Perfomance

Bxog is the fastest router, showing the speed of query processing. Its speed is comparable to the speed of the popular multiplexers: Bone, Httprouter, Gorilla, Zeus.  The test is done on a computer with a i7-6700T processor and 8 GB RAM. Detailed benchmark [here](https://github.com/claygod/bxogtest). In short (less time, the better):

- Bxog-8		103 ns/op
- HttpRouter-8 		201 ns/op
- Zeus-8 		13420 ns/op
- Gorilla-8 		17350 ns/op
- GorillaPat-8 		693 ns/op
- Bone-8 		49633 ns/op

# API

Methods:
-  *New* - create a new multiplexer
-  *Add* - add a rule specifying the handler (the default method - GET, ID - as a string to this rule)
-  *Start* - start the server indicating the listening port
-  *Params* - extract parameters from URL
-  *Create* - generate URL of the available options
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
