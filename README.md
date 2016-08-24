Bxog is a simple and fast HTTP router for Go (HTTP request multiplexer).

[![API documentation](https://godoc.org/github.com/claygod/Bxog?status.svg)](https://godoc.org/github.com/claygod/Bxog)
[![Go Report Card](https://goreportcard.com/badge/github.com/claygod/Bxog)](https://goreportcard.com/report/github.com/claygod/Bxog)

# Usage

An example of using the multiplexer:
```go
package main

import (
	"github.com/claygod/Bxog"
	"github.com/claygod/Context"
	"io"
	"net/http"
)

// Handlers
func IHandler(w http.ResponseWriter, req *http.Request, c *Context.Context) {
	io.WriteString(w, "Welcome to Bxog!")
	if x := c.Get("answer"); x != nil {
		io.WriteString(w, "\n Context is used? "+x.(string)+"\n")
	}
}
func THandler(w http.ResponseWriter, req *http.Request, c *Context.Context) {
	io.WriteString(w, "Params:\n")
	if x := c.Get("par"); x != nil {
		io.WriteString(w, " 'par' -> "+x.(string)+"\n")
	}
	if x := c.Get("name"); x != nil {
		io.WriteString(w, " 'name' -> "+x.(string)+"\n")
	}
}
func PHandler(w http.ResponseWriter, req *http.Request, c *Context.Context) {
	// Getting parameters from URL
	io.WriteString(w, "Country:\n")
	io.WriteString(w, " 'name' -> "+c.Get("name").(string)+"\n")
	io.WriteString(w, " 'capital' -> "+c.Get("city").(string)+"\n")
	io.WriteString(w, " 'valuta' -> "+c.Get("money").(string)+"\n")

}

// Main
func main() {
	ctx := Context.New()
	ctx.Set("answer", "Yes!")

	m := bxog.New()
	m.Add("/", IHandler).
		Context(ctx.Fix()) // This context has access to the variable "answer"
	m.Add("/abc/:par", THandler)
	m.Add("/country/:name/capital/:city/valuta/:money", PHandler).
		Id("country"). // For ease indicate the short ID
		Method("GET")  // GET method do not need to write here, it is used by default (this is an example)
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

Bxog is the fastest router, showing the speed of query processing. Its speed is comparable to the speed of the popular multiplexers: Bone, Httprouter, Gorilla, Zeus.  In short (less time, the better):

- Bxog         330 ns/op
- HttpRouter   395 ns/op
- Zeus       23772 ns/op
- GorillaMux 30223 ns/op
- GorillaPat  1253 ns/op
- Bone       63656 ns/op

# API

Methods:
-  *New* - create a new multiplexer
-  *Add* - add a rule specifying the handler (the default method - GET, ID - as a string to this rule)
-  *Start* - start the server indicating the listening port
-  *Context* - for the route, which will be available from within the handler.
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
