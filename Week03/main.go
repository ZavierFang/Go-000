package main

import (
	"context"
	"errors"
	"fmt"
	"golang.org/x/sync/errgroup"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func main() {
	g, ctx := errgroup.WithContext(context.Background())

	g.Go(func() (err error) {
		defer func() {
			if recoverErr := recover(); recoverErr != nil {
				log.Printf("panic by %+v", recoverErr)
				err = errors.New("panic in listen system signal")
			}
		}()

		osSig := make(chan os.Signal, 1)
		signal.Notify(osSig, syscall.SIGINT, syscall.SIGTERM, syscall.SIGKILL)
		select {
		case sig := <-osSig:
			log.Printf("exist by os signal(%d)", sig)
			return errors.New("exist by os signal")
		case <-ctx.Done():
			log.Printf("exist cause stop")
			return ctx.Err()
		}
	})

	g.Go(func() (err error) {
		defer func() {
			if recoverErr := recover(); recoverErr != nil {
				log.Printf("panic by %+v", recoverErr)
				err = errors.New("panic in run http server")
			}
		}()

		err = runHttpServer(ctx)
		if err != nil {
			log.Printf("run http server failed, err: %+v", err)
		}
		return err
	})

	if err := g.Wait(); err != nil {
		log.Printf("call errgroup wait failed, err: %+v", err)
	}

	//等待正在处理的请求结束
	time.Sleep(15 * time.Second)

	return
}

func runHttpServer(ctx context.Context) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/", func(writer http.ResponseWriter, request *http.Request) {
		fmt.Fprintln(writer, "Hello week03")
	})
	s := http.Server{
		Addr:    "127.0.0.1:8080",
		Handler: mux,
	}

	go func() {
		<-ctx.Done()
		log.Printf("shutdown http server")
		_ = s.Shutdown(ctx)
	}()

	return s.ListenAndServe()
}
