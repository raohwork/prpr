package main // import "git.ronmi.tw/raoh/prpr"

import (
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"
)

func main() {
	secret := os.Getenv("SECRET")
	firefox := os.Getenv("FIREFOX")
	if firefox == "" {
		firefox = "firefox"
	}
	bind := os.Getenv("BIND")
	if bind == "" {
		bind = ":9801"
	}

	ch := make(chan error)

	go runFX(firefox, ch)

	// wait few second for firefox to start
	time.Sleep(10 * time.Second)
	go runWeb(bind, secret, ch)

	log.Fatal(<-ch)
}

func runFX(firefox string, ch chan error) {
	// starting firefox
	cmd := exec.Command(
		firefox,
		"--marionette",
		"--headless",
		"--safe-mode",
	)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	ch <- cmd.Run()
}

func runWeb(bind, secret string, ch chan error) {
	h := New("", 0)
	if secret != "" {
		h.Secret = secret
	}
	http.HandleFunc("/grab", h.Accept)
	log.Print("Starting prpr...")
	ch <- http.ListenAndServe(bind, nil)
}
