package main

import (
	"context"
	smtp "ktn-go/smtp"
	web "ktn-go/web"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	log "github.com/sirupsen/logrus"
)

func main() {
	log.Println("Starting KTN Server...")

	smtpSrv := smtp.SmtpServer()
	go func() {
		log.Info("Starting SMTP Server on port ", smtpSrv.Addr)
		if err := smtpSrv.ListenAndServe(); err != nil {
			log.Fatalf("SMTP Server failure: %s\n", err)
		}
	}()

	httpSrv := web.HttpServer()
	go func() {
		// service connections
		log.Info("Starting HTTP Server on port ", httpSrv.Addr)
		if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("HTTP Server failure: %s\n", err)
		}
	}()

	// Wait for interrupt signal to gracefully shutdown the server with
	// a timeout of 5 seconds.
	quit := make(chan os.Signal, 1)
	// kill (no param) default send syscanll.SIGTERM
	// kill -2 is syscall.SIGINT
	// kill -9 is syscall. SIGKILL but can"t be catch, so don't need add it
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	sign := <-quit
	log.Info("KTN Shutdown: Received signal ", sign)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpSrv.Shutdown(ctx); err != nil {
		log.Fatal("KTN Shutdown: HTTP Shutdown failed:", err)
	} else {
		log.Info("KTN Shutdown: HTTP closed")

	}
	if err := smtpSrv.Close(); err != nil {
		log.Fatal("KTN Shutdown: SMTP Shutdown failed:", err)
	} else {
		log.Info("KTN Shutdown: SMTP closed")

	}
	// catching ctx.Done(). timeout of 5 seconds.
	select {
	case <-ctx.Done():
		log.Error("KTN Shutdown: Timeout of 5 seconds.")
	default:
		log.Info("KTN Shutdown: Server exiting...")
	}
}
