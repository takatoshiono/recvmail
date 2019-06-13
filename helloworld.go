package main

import (
	"bytes"
	"fmt"
	"net/http"

	"google.golang.org/appengine"
	gae_log "google.golang.org/appengine/log"
)

func main() {
	http.HandleFunc("/", indexHandler)
	http.HandleFunc("/_ah/mail/", incomingMail)
	appengine.Main()
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	fmt.Fprint(w, "Hello, World!")
}

func incomingMail(w http.ResponseWriter, r *http.Request) {
	ctx := appengine.NewContext(r)
	defer r.Body.Close()
	var b bytes.Buffer
	if _, err := b.ReadFrom(r.Body); err != nil {
		gae_log.Errorf(ctx, "Error reading body: %v", err)
		return
	}
	gae_log.Infof(ctx, "Received mail: %v", b)
}
