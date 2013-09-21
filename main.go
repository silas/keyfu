package main

import (
	"flag"
)

var httpInterface = flag.String("interface", ":8000", "http server interface")
var staticPath = flag.String("static", "./static", "static path")
var templatesPath = flag.String("templates", "./templates", "templates path")

func main() {
	flag.Parse()
	run(*httpInterface, *staticPath, *templatesPath)
}
