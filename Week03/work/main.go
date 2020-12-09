package main

import (
	"context"
	"github.com/pkg/errors"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/sync/errgroup"
)

func main() {
	exit := make(chan struct{}, 1)
	g, _ := errgroup.WithContext(context.Background())
	g.Go(func() error {
		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGHUP, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGINT)
		switch <-quit {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT, syscall.SIGHUP:
			log.Println("got quit signal, exit ...")
			close(quit)
			close(exit)
		}
		return nil
	})

	g.Go(func() error {
		s := http.Server{Addr: ":8000"}
		mux := http.NewServeMux()
		mux.HandleFunc("/", helloworld)
		s.Handler = mux
		errServer := make(chan error, 1)
		go func() {
			errServer <- s.ListenAndServe()
		}()

		select {
		case <-exit:
			s.Shutdown(context.Background())
			return nil
		case err := <-errServer:
			close(exit)
			return errors.Wrap(err, "listen and serve")
		}
	})

	if err := g.Wait(); err != nil {
		log.Println(err)
	}
}

func helloworld(w http.ResponseWriter, r *http.Request) {
	// 模拟耗时处理
	time.Sleep(5 * time.Second)
	w.Write([]byte("hello world"))
}
