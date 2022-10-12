// Copyright The OpenTelemetry Authors
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Example using OTLP exporters + collector + third-party backends. For
// information about using the exporter, see:
// https://pkg.go.dev/go.opentelemetry.io/otel/exporters/otlp?tab=doc#example-package-Insecure
package main

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"

	log "github.com/sirupsen/logrus"
)

type HelloHandler struct {
	ctx      context.Context
	response string
}

func (h *HelloHandler) helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, h.response)
}

type CounterHandler struct {
	ctx     context.Context
	counter int
}

func (ct *CounterHandler) counterHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Println(ct.counter)
	ct.counter++
	msg := fmt.Sprintf("Counter: %d", ct.counter)
	fmt.Fprintln(w, msg)
}

type NotFoundHandler struct {
	ctx context.Context
}

func (nf *NotFoundHandler) notFoundHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		w.WriteHeader(404)
		w.Write([]byte("404 - not found\n"))
		msg := "404 - not found"
	    fmt.Fprintln(w, msg)
		return
	}
	msg := "This page does nothing, add a '/count' or a '/hello'"
	fmt.Fprintln(w, msg)
}

func listenAndServe(ctx context.Context, port string, handler http.Handler) {
	log.Infof("serving on %s", port)
	_, err := os.Stat("/etc/tls-config/tls.crt")
	if err == nil {
		if err := http.ListenAndServeTLS(":"+port, "/etc/tls-config/tls.crt", "/etc/tls-config/tls.key", handler); err != nil {
			msg := fmt.Sprintf("ListenAndServe: " + err.Error())
			log.Panic(msg)
		}
		return
	}
	if errors.Is(err, os.ErrNotExist) {
		if err := http.ListenAndServe(":"+port, handler); err != nil {
			msg := fmt.Sprintf("ListenAndServe: " + err.Error())
			log.Panic(msg)
		}
		return
	}
	log.Fatalf("failed to start serving: %w", err)
}

func init() {
	log.SetFormatter(&log.TextFormatter{
		DisableColors: true,
		FullTimestamp: true,
	})

	log.SetOutput(os.Stdout)
	log.SetLevel(log.TraceLevel)
}

func main() {
	log.Infof("Starting server")
	ctx := context.Background()

	helloResponse := os.Getenv("RESPONSE")
	if len(helloResponse) == 0 {
		helloResponse = "Hello OpenTelemetry!"
	}
	hello := &HelloHandler{
		response: helloResponse,
		ctx:      ctx,
	}

	count := &CounterHandler{
		ctx:     ctx,
		counter: 0,
	}

	notFound := &NotFoundHandler{
		ctx: ctx,
	}

	port := os.Getenv("PORT")
	if len(port) == 0 {
		port = "8080"
	}
	mux := http.NewServeMux()
	mux.HandleFunc("/hello", hello.helloHandler)
	mux.HandleFunc("/count", count.counterHandler)
	mux.HandleFunc("/", notFound.notFoundHandler)
	go listenAndServe(ctx, port, mux)

	select {}
}
