package globals

import (
	"flag"
)

var (
	ListenAddr     = flag.String("http", ":8080", "http listen address")
	StandardLayout = []string{"templates/shared/head.html", "templates/shared/header.html", "templates/shared/navigation.html", "templates/shared/footer.html"}
)

func SetGlobals() {
	flag.Parse()
}
