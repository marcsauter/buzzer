package main

import (
	"log"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"syscall"

	"github.com/marcsauter/buzzer/pkg/pitch"
	"github.com/marcsauter/buzzer/pkg/ticker"
)

func main() {

	device := os.Getenv("TICKER_DEVICE")
	if _, err := os.Stat(device); os.IsNotExist(err) {
		log.Fatal("TICKER_DEVICE missing or not valid")
	}
	url, err := url.Parse(os.Getenv("TICKER_PITCH_URL"))
	if err != nil || !url.IsAbs() {
		log.Fatal("TICKER_PITCH_URL missing or not valid")
	}
	interval, err := strconv.Atoi(os.Getenv("TICKER_PITCH_CHECK_INTERVAL"))
	if err != nil {
		log.Fatal("TICKER_PITCH_CHECK_INTERVAL missing or not valid")
	}

	t, err := ticker.NewTicker(device)
	if err != nil {
		log.Fatal(err)
	}

	p := pitch.NewPitch(url)
	p.StartCheckNext(interval, t)

	//
	cancel := make(chan os.Signal, 1)
	signal.Notify(cancel, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
	for {
		select {
		case <-cancel:
			log.Fatalln("signal received - exiting")
		}
	}
}
