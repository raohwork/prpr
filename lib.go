package main

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	mario "github.com/njasm/marionette_client"
)

const getContentF = `return document.querySelector('html').outerHTML;`

type resp struct {
	Value string `json:"value"`
}

type Handler struct {
	client  *mario.Client
	windows chan string
	token   chan byte
	Secret  string
}

func New(server string, port, max int) (ret *Handler) {
	c := mario.NewClient()
	if err := c.Connect(server, port); err != nil {
		log.Fatalf("cannot connect to firefox: %s", err)
	}
	if _, err := c.NewSession("", nil); err != nil {
		log.Fatalf("cannot create new session: %s", err)
	}

	// fetch current opening windows
	ws, err := c.WindowHandles()
	if err != nil {
		log.Fatalf("cannot fetch info of currently opened windows: %s", err)
	}

	// create windows
	for l := len(ws); l < max; l++ {
		_, err = c.ExecuteScript(`window.open('about:blank')`, []interface{}{1}, 1000, false)
		if err != nil {
			log.Fatalf("cannot open new window: %s", err)
		}
	}
	// close windows
	for l := len(ws); l > max; l-- {
		if err = c.SwitchToWindow(ws[l-1]); err != nil {
			log.Fatalf("cannot switch to window %s: %s", ws[l-1], err)
		}
		if _, err = c.CloseWindow(); err != nil {
			log.Fatalf("cannot close windows %s: %s", ws[l-1], err)
		}
	}

	// fetch again
	if ws, err = c.WindowHandles(); err != nil {
		log.Fatalf("cannot fetch info of currently opened windows: %s", err)
	}

	ret = &Handler{
		client:  c,
		windows: make(chan string, max),
		token:   make(chan byte, 1),
	}
	for _, w := range ws {
		ret.windows <- w
	}
	ret.token <- 0

	return
}

func (h *Handler) allocate(w string) (err error) {
	<-h.token
	return h.client.SwitchToWindow(w)
}

func (h *Handler) release(w string) {
	h.token <- 0
}

func (h *Handler) Grab(uri, wait string) (ret string, err error) {
	// allocate a window
	w := <-h.windows
	defer func() {
		h.windows <- w
		h.token <- 0
	}()

	// get token
	if err = h.allocate(w); err != nil {
		return
	}

	if _, err = h.client.Navigate(uri); err != nil {
		return
	}

	// release token and wait 1 second for page loaing
	h.release(w)
	time.Sleep(1 * time.Second)

	// check 10 times, wait 1 second between each check
	ok := false
	for x := 0; x < 10; x++ {
		if err = h.allocate(w); err != nil {
			return
		}
		we, e := h.client.FindElement(
			mario.CSS_SELECTOR,
			wait,
		)
		if e != nil {
			err = e
			h.client.Navigate("about:blank")
			return
		}
		h.release(w)
		if we != nil {
			ok = true
			break
		}
		time.Sleep(1 * time.Second)
	}

	if err = h.allocate(w); err != nil {
		return
	}
	if !ok {
		err = errors.New("timed out")
		h.client.Navigate("about:blank")
		return
	}

	r, err := h.client.ExecuteScript(getContentF, []interface{}{1}, 1000, false)
	if err != nil {
		h.client.Navigate("about:blank")
		return
	}

	var data resp
	if err = json.Unmarshal([]byte(r.Value), &data); err != nil {
		h.client.Navigate("about:blank")
		return
	}

	ret = data.Value
	h.client.Navigate("about:blank")
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
