package main

//docker:install graphviz

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net"
	"os"

	lightstep "github.com/lightstep/lightstep-tracer-go"
	opentracing "github.com/opentracing/opentracing-go"

	"sourcegraph.com/sourcegraph/sourcegraph/pkg/debugserver"
	"sourcegraph.com/sourcegraph/sourcegraph/xlang"
)

var (
	addr     = flag.String("addr", ":4388", "proxy server TCP listen address")
	profbind = flag.String("prof-http", ":6060", "net/http/pprof http bind address")
	trace    = flag.Bool("trace", false, "print traces of JSON-RPC 2.0 requests/responses")
)

func main() {
	flag.Parse()
	log.SetFlags(0)

	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func run() error {
	if t := os.Getenv("LIGHTSTEP_ACCESS_TOKEN"); t != "" {
		opentracing.InitGlobalTracer(lightstep.NewTracer(lightstep.Options{
			AccessToken: t,
		}))
	}

	if err := xlang.RegisterServersFromEnv(); err != nil {
		return err
	}

	lis, err := net.Listen("tcp", *addr)
	if err != nil {
		return err
	}
	fmt.Fprintln(os.Stderr, "lsp-proxy: listening on", lis.Addr())
	p := xlang.NewProxy()
	p.Trace = *trace
	if *profbind != "" {
		e := debugserver.Endpoint{
			Name:    "LSP-Proxy Connections",
			Path:    "/lsp-conns",
			Handler: &xlang.DebugHandler{Proxy: p},
		}
		go debugserver.Start(*profbind, e)
	}
	return p.Serve(context.Background(), lis)
}
