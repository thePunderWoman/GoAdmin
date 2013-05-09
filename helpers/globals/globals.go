package globals

import (
	"flag"
	"os"
)

var (
	ListenAddr = flag.String("http", ":8080", "http listen address")
)

func SetGlobals() {
	flag.Parse()
	os.Setenv("TZ", "UTC")
}
