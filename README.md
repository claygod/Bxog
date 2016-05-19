Bxog is a simple and fast HTTP router for Go (HTTP request multiplexer).

# Usage

An example of using the multiplexer:
`
package main

import (
	"io"
	"net/http"
	"github.com/claygod/bxog"
)

// Handlers
func IHandler(w http.ResponseWriter, req *http.Request, r *bxog.Router) { //, ro *bxog.Router
	io.WriteString(w, "Welcome to Bxog!")
}
func THandler(w http.ResponseWriter, req *http.Request, r *bxog.Router) {
	params := r.Params(req, "/abc/:para")
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
	io.WriteString(w, "Creating a URL from route:\n")
	io.WriteString(w, r.Create("country", map[string]string{"name": "Russia", "city": "1867", "money": "rouble"}))
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
`

Click URLs:
- http://localhost
- http://localhost/abc/123
- http://localhost/country/USA/capital/Washington/valuta/dollar

# Settings

Necessary changes in the configuration of the multiplexer can be made in the configuration file [config.go](https://github.com/claygod/bxog/config.go)

# Perfomance

Bxog is the fastest router, showing the speed of query processing. Its speed is comparable to the speed of the popular multiplexers: Bone, Httprouter, Gorilla, Zeus.  Detailed benchmark [here](https://github.com/claygod/bxogtest). In short (less time, the better):
`
BenchmarkBxogMux-4      	 5000000	       330 ns/op
BenchmarkHttpRouterMux-4	 3000000	       395 ns/op
BenchmarkZeusMux-4      	  100000	     23772 ns/op
BenchmarkGorillaMux-4   	   50000	     30223 ns/op
BenchmarkGorillaPatMux-4	 1000000	      1253 ns/op
BenchmarkBoneMux2-4     	   20000	     63656 ns/op
`
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

Arguments in the rules designated Roth colon. Example route: */abc/:param* , where *abc* is a static section and *:param* - the dynamic section(argument).

# Static files

The directory path to the file and its nickname as part of URL specified in the configuration file. This constants *FILE_PREF* and *FILE_PATH*
