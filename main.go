package main

import (
	"flag"
	"fmt"
	"github.com/goji/httpauth"
	"github.com/ttlv/sms_admin/config"

	theplant_server "github.com/theplant/appkit/server"
	"github.com/ttlv/sms_admin/server"
	"net/http"
	"os"

	"github.com/ttlv/sms_admin/config/bindatafs"
)

func main() {
	cmdLine := flag.NewFlagSet(os.Args[0], flag.ContinueOnError)
	compileTemplate := cmdLine.Bool("compile-templates", false, "Compile Templates")
	cmdLine.Parse(os.Args[1:])
	c := config.MustGetConfig()
	mux, _ := server.NewServer()
	mw := theplant_server.Compose(httpAuthMiddleWare)
	if *compileTemplate {
		bindatafs.AssetFS.Compile()
	} else {
		fmt.Printf("Listening on: %v\n", c.ServerPort)
		if c.HTTPS {
			if err := http.ListenAndServeTLS(fmt.Sprintf(":%d", c.ServerPort), "config/local_certs/server.crt", "config/local_certs/server.key", mux); err != nil {
				panic(err)
			}
		} else {
			if err := http.ListenAndServe(fmt.Sprintf(":%d", c.ServerPort), mw(mux)); err != nil {
				panic(err)
			}
		}
	}
}

func httpAuthMiddleWare(handler http.Handler) http.Handler {
	cfg := config.MustGetConfig()
	fn := func(w http.ResponseWriter, r *http.Request) {
		h := httpauth.SimpleBasicAuth(cfg.HttpAuthName, cfg.HttpAuthPassword)(handler)
		h.ServeHTTP(w, r)
	}
	return http.HandlerFunc(fn)
}
