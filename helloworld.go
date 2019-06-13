package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/mail"

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
	m, err := mail.ReadMessage(&b)
	if err != nil {
		gae_log.Errorf(ctx, "Error reading message: %v", err)
		return
	}
	gae_log.Infof(ctx, "From: %s", m.Header.Get("From"))
	gae_log.Infof(ctx, "Subject: %s", m.Header.Get("Subject"))
	body, err := ioutil.ReadAll(m.Body)
	if err != nil {
		gae_log.Errorf(ctx, "Error reading message body: %v", err)
	}
	gae_log.Infof(ctx, "Body: %s", body)
	// TODO: decode MIME multipart
}
