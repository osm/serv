package main

import (
	"crypto/subtle"
	"encoding/base64"
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

func main() {
	// Parameters from either the environment or the command line
	cert := envFlag("cert", "path to the certificate")
	cred := envFlag("cred", "path or string with base64 encoded credentials (user:password)")
	file := envFlag("file", "path to the file to serve")
	key := envFlag("key", "path to the key")
	port := envFlag("port", "port to listen on (defaults to 4554)")
	flag.Parse()

	// Make sure we our required parameters
	if *cred == "" {
		errorf("missing the -cred flag or CRED environment variable")
	}
	if *cert == "" {
		errorf("missing the -cert flag or CERT environment variable")
	}
	if *key == "" {
		errorf("missing the -key flag or KEY environment variable")
	}
	if *file == "" {
		errorf("missing the -file flag or FILE environment variable")
	}
	if *port == "" {
		*port = "4554"
	}

	// Get basic auth credentials
	user, password := getBasicAuthCredentials(*cred)

	// Setup the one and only handler
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Make sure we are authenticated
		u, p, ok := r.BasicAuth()
		if !ok || subtle.ConstantTimeCompare([]byte(u), []byte(user)) != 1 || subtle.ConstantTimeCompare([]byte(p), []byte(password)) != 1 {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted area"`)
			w.WriteHeader(401)
			w.Write([]byte("unauthorized\n"))
			return
		}

		// Read the file that we want to serve
		d, err := ioutil.ReadFile(*file)
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("can't find the file to serve"))
		}

		// Serve the contents
		w.Write(d)
	})

	// Start the server
	if err := http.ListenAndServeTLS(fmt.Sprintf(":%s", *port), *cert, *key, nil); err != nil {
		errorf("can't start server, %v", err)
	}
}

// errorf prints to stderr and exits.
func errorf(format string, args ...interface{}) {
	m := fmt.Sprintf(format, args...)
	if !strings.HasSuffix(m, "\n") {
		m += "\n"
	}
	fmt.Fprintf(os.Stderr, m)
	os.Exit(1)
}

// envFlag reads the value from the environment or command line parameter
func envFlag(key, description string) *string {
	if env := os.Getenv(strings.ToUpper(key)); env != "" {
		return &env
	}
	return flag.String(key, "", description)
}

// getBasicAuthCredentials gets the credentials from the src string.
// src can be either a base64 encoded string or a file path.
// If src is a path we expect the contents to be base64 encoded.
func getBasicAuthCredentials(src string) (string, string) {
	var encoded string

	if _, err := os.Stat(src); err == nil {
		tmp, _ := ioutil.ReadFile(src)
		encoded = string(tmp)
	} else {
		encoded = src
	}

	decoded, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		errorf("doesn't look like we got a base64 encoded string, %v", err)
	}

	parts := strings.Split(string(decoded), ":")
	if len(parts) != 2 {
		errorf("unexpected format for credentials, expected foo:bar")
	}

	return parts[0], parts[1]
}
