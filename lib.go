package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/raohwork/marionette-go/mnclient"
	"github.com/raohwork/marionette-go/mnsender"
	"github.com/raohwork/marionette-go/tabmgr"
)

const getContentF = `return document.querySelector('html').outerHTML;`

type resp struct {
	Value string `json:"value"`
}

type Handler struct {
	client *tabmgr.TabManager
	tabs   chan string
	Secret string
}

func New(server string, port, max int) (ret *Handler) {
	if max < 1 {
		max = 1
	}
	if server == "" {
		server = "127.0.0.1"
	}
	if port < 1 {
		port = 2828
	}

	c, err := mnsender.NewTCPSender(fmt.Sprintf("%s:%d", server, port), 0)
	if err != nil {
		log.Fatalf("cannot connect to marionette server: %s", err)
	}
	if err = c.Start(); err != nil {
		log.Fatalf("cannot start marionette: %s", err)
	}
	cl := &mnclient.Commander{Sender: c}
	if _, _, err = cl.NewSession(); err != nil {
		log.Fatalf("cannot create new session: %s", err)
	}

	tabNames := make([]string, max)
	ch := make(chan string, max)
	for x := 1; x <= max; x++ {
		name := "tab" + strconv.Itoa(x)
		tabNames[x-1] = name
		ch <- name
	}

	tm, err := tabmgr.New(c, tabNames)
	if err != nil {
		log.Fatalf("cannot init TabManager: %s", err)
	}

	ret = &Handler{
		client: tm,
		tabs:   ch,
	}

	return
}

func (h *Handler) allocate() (ret *tabmgr.Tab) {
	str := <-h.tabs
	return h.client.GetTab(str)
}

func (h *Handler) release(tab *tabmgr.Tab) {
	h.tabs <- tab.GetName()
}

func (h *Handler) Grab(uri, wait string) (ret string, err error) {
	// allocate a window
	tab := h.allocate()
	defer h.release(tab)

	ch := tab.NavigateAsync(uri)
	if err = <-ch; err != nil {
		return
	}
	defer tab.Navigate("about:blank")

	if _, err = tab.WaitFor(wait, 10); err != nil {
		return
	}

	err = tab.ExecuteScript(string(getContentF), &ret)
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
