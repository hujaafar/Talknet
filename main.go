package main

import (
	"context"
	"flag"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"talknet/internal/forum"
	"time"
)

func main() {
	seed := flag.Bool("seed-demo", false, "add fictional sample discussions to an empty database")
	health := flag.Bool("healthcheck", false, "check the running server and exit")
	flag.Parse()
	path := os.Getenv("TALKNET_DB")
	if path == "" {
		path = "data/talknet.db"
	}
	addr := os.Getenv("TALKNET_ADDR")
	if addr == "" {
		addr = ":8080"
	}
	if *health {
		_, port, err := net.SplitHostPort(addr)
		if err != nil {
			log.Fatal(err)
		}
		client := http.Client{Timeout: 3 * time.Second}
		response, err := client.Get("http://127.0.0.1:" + port + "/healthz")
		if err != nil {
			log.Fatal(err)
		}
		response.Body.Close()
		if response.StatusCode != 200 {
			os.Exit(1)
		}
		return
	}
	db, err := forum.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	if *seed {
		if err := forum.SeedDemo(db); err != nil {
			log.Fatal(err)
		}
		log.Print("Sample-data setup complete (existing communities are unchanged)")
		return
	}
	app := forum.New(db, os.Getenv("TALKNET_SECURE_COOKIES") == "true")
	srv := &http.Server{Addr: addr, Handler: app, ReadHeaderTimeout: 5 * time.Second, ReadTimeout: 15 * time.Second, WriteTimeout: 30 * time.Second, IdleTimeout: 60 * time.Second}
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()
	stopped := make(chan struct{})
	go func() {
		defer close(stopped)
		<-ctx.Done()
		shutdown, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := srv.Shutdown(shutdown); err != nil {
			log.Printf("shutdown: %v", err)
		}
	}()
	log.Printf("Talknet listening on %s", addr)
	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		log.Fatal(err)
	}
	<-stopped
}
