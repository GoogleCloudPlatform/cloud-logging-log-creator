// Copyright 2021 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"cloud.google.com/go/compute/metadata"
	"cloud.google.com/go/logging"
)

var s *SeverityLogger

func main() {

	fmt.Printf("Spinning up chatty logs\n")
	mclient := metadata.NewClient(&http.Client{Transport: userAgentTransport{
		userAgent: "chattylogs",
		base:      http.DefaultTransport,
	}})

	fmt.Printf("Getting project id.\n")
	projectID, err := mclient.ProjectID()
	if err != nil {
		fmt.Printf("cannot get id from metadata, defaulting to env (%s), \n", err)
		projectID = os.Getenv("PROJECT_ID")
	}

	fmt.Printf("Project ID: %s\n", projectID)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	s, err = NewSeverityLogger("chattylogs", projectID)
	if err != nil {
		log.Fatalf("Failed to create logger: %v", err)
	}

	fillLogs(s)

	http.HandleFunc("/healthz", handleHealth)
	http.HandleFunc("/", handleHealth)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatal(err)
	}
}

func handleHealth(w http.ResponseWriter, r *http.Request) {
	s.loghttp(r)
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, "ok")
}

func fillLogs(s *SeverityLogger) {

	go func() {
		for {
			s.log(logging.Info, "an informational log entry")
			time.Sleep(5 * time.Second)
		}
	}()

	go func() {
		for {
			time.Sleep(20 * time.Second)
			s.log(logging.Warning, "be prepared")

		}
	}()

	go func() {
		for {
			time.Sleep(15 * time.Second)
			s.log(logging.Error, "error: %s", fmt.Errorf("testing error issues"))

		}
	}()

	go func() {
		for {
			time.Sleep(2 * time.Second)
			s.log(logging.Debug, "small detail of interest to developer")

		}
	}()

	go func() {
		for {
			time.Sleep(60 * time.Second)
			s.log(logging.Critical, "super important")

		}
	}()

	go func() {
		for {
			time.Sleep(6 * time.Second)
			fmt.Printf("default content outputted to stdout\n")

		}
	}()
}

// userAgentTransport sets the User-Agent header before calling base.
type userAgentTransport struct {
	userAgent string
	base      http.RoundTripper
}

// RoundTrip implements the http.RoundTripper interface.
func (t userAgentTransport) RoundTrip(req *http.Request) (*http.Response, error) {
	req.Header.Set("User-Agent", t.userAgent)
	return t.base.RoundTrip(req)
}

// NewSeverityLogger creates a severitylogger for the logging of the things.
func NewSeverityLogger(name, project string) (*SeverityLogger, error) {

	client, err := logging.NewClient(context.Background(), project)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %v", err)
	}
	defer client.Close()

	return &SeverityLogger{name, client}, nil
}

// SeverityLogger is a struct for aggregating making calls to log a little
// easier.
type SeverityLogger struct {
	name   string
	client *logging.Client
}

func (s *SeverityLogger) log(severity logging.Severity, msg string, arg ...interface{}) {
	s.client.Logger(s.name).StandardLogger(severity).Printf(msg, arg...)
}

func (s *SeverityLogger) loghttp(r *http.Request) {
	d := time.Now().Format("[02/Jan/2006:15:04:06 -0700]")
	fmt.Printf("%s %s %s \"%s %s %s\" %d %d\n", r.Host, r.UserAgent(), d, r.Method, r.URL.Path, r.Proto, http.StatusOK, r.ContentLength)
}
