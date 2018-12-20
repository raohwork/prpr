package main

import (
	"encoding/json"
	"log"
	"net/http"
	"sync"
	"time"

	mario "github.com/njasm/marionette_client"
)

const getContentF = `return document.querySelector('html').outerHTML;`

type resp struct {
	Value string `json:"value"`
}

type Handler struct {
	client *mario.Client
	lock   sync.Mutex
	Secret string
}

func New(server string, port int) (ret *Handler) {
	c := mario.NewClient()
	if err := c.Connect(server, port); err != nil {
		log.Fatalf("cannot connect to firefox: %s", err)
	}
	if _, err := c.NewSession("", nil); err != nil {
		log.Fatalf("cannot create new session: %s", err)
	}

	return &Handler{client: c}
}

func (h *Handler) Grab(uri, wait string) (ret string, err error) {
	h.lock.Lock()
	defer h.lock.Unlock()

	defer h.client.Navigate("about:blank")

	if _, err = h.client.Navigate(uri); err != nil {
		return
	}

	cond := mario.ElementIsPresent(
		mario.By(mario.CSS_SELECTOR),
		wait,
	)
	_, _, err = mario.Wait(h.client).For(10 * time.Second).Until(cond)
	if err != nil {
		return
	}

	r, err := h.client.ExecuteScript(getContentF, []interface{}{1}, 1000, false)
	if err != nil {
		return
	}

	var data resp
	if err = json.Unmarshal([]byte(r.Value), &data); err != nil {
		return
	}

	ret = data.Value
	return
}

func (h *Handler) Accept(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	uri := r.FormValue("uri")
	wait := r.FormValue("wait")
	secret := r.FormValue("secret")
	if h.Secret != "" && secret != h.Secret {
		w.WriteHeader(400)
		log.Printf("incorrect secret: %s", r.RemoteAddr)
		return
	}

	log.Printf("Grabing: %s (%s)", uri, wait)

	data, err := h.Grab(uri, wait)
	if err != nil {
		w.WriteHeader(500)
		log.Printf("cannot grab: %s", err)
		return
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(data))
}
