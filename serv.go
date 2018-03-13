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
	cert := flag.String("cert", "", "certificate")
	cred := flag.String("cred", "", "path to file with base64 encoded credentials (user:password)")
	file := flag.String("file", "", "path to the file to serve")
	key := flag.String("key", "", "key")
	port := flag.Int("port", 4554, "port to listen on (1024 - 65535")
	flag.Parse()

	// Make sure we have a valid credentials file
	if *cred == "" {
		errorf("missing the -cred flag")
	}

	// Read the contents of the file
	data, err := ioutil.ReadFile(*cred)
	if err != nil {
		errorf("can't open credentials file, %v", err)
	}

	// Decode the input
	credentials, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		errorf("can't decode the credentials, %v", err)
	}

	// Separate the user and password
	parts := strings.Split(string(credentials), ":")
	if len(parts) != 2 {
		errorf("unexpected format for credentials, expected foo:bar")
	}
	user := parts[0]
	password := parts[1]

	// We need a certificate and a key
	if *cert == "" {
		errorf("missing the -cert flag")
	}
	if *key == "" {
		errorf("missing the -key flag")
	}

	// We also need a file to serve on a successful request
	if *file == "" {
		errorf("missing the -file flag")
	}

	// And finally, we need to listen on a valid port
	if *port <= 1024 || *port >= 65535 {
		errorf("%v is not within the valid port range", *port)
	}

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
	if err := http.ListenAndServeTLS(fmt.Sprintf(":%d", *port), *cert, *key, nil); err != nil {
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
