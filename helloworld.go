package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime"
	"mime/multipart"
	"net/http"
	"net/mail"
	"strings"

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

	mediaType, params, err := mime.ParseMediaType(m.Header.Get("Content-Type"))
	if err != nil {
		gae_log.Errorf(ctx, "Error parsing media type: %v", err)
	}
	if strings.HasPrefix(mediaType, "multipart/") {
		// TODO: parse multipart recursively
		mr := multipart.NewReader(m.Body, params["boundary"])
		i := 1
		for {
			p, err := mr.NextPart()
			if err == io.EOF {
				break
			}
			if err != nil {
				gae_log.Errorf(ctx, "Error getting next part: %v", err)
			}
			slurp, err := ioutil.ReadAll(p)
			if err != nil {
				gae_log.Errorf(ctx, "Error reading part: %v", err)
			}
			// TODO: decode message
			gae_log.Infof(ctx, "Part%d %q", i, slurp)
			for key := range p.Header {
				gae_log.Infof(ctx, "%s: %s", key, p.Header.Get(key))
			}
			i++
		}
	} else {
		body, err := ioutil.ReadAll(m.Body)
		if err != nil {
			gae_log.Errorf(ctx, "Error reading message body: %v", err)
		}
		gae_log.Infof(ctx, "Body: %s", body)
	}
}
