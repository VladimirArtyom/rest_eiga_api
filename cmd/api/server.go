package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func (app *application) serve() error {
	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", app.cfg.port),
		Handler:      app.routes(),
		ErrorLog:     log.New(app.logger, "", 0), // 新しい実証 ( New implementation  )
		IdleTimeout:  time.Minute,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 30 * time.Second,
	}
	shutdownErrorChannel := make(chan error)

	go func() {
		quit := make(chan os.Signal, 1)
		syscalls := []os.Signal{syscall.SIGINT, syscall.SIGTERM}

		signal.Notify(quit, syscalls...)

		s := <-quit

		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()

		app.logger.PrintInfo("Shutting down the server...", map[string]string{
			"Signal": s.String(),
		})

		err := server.Shutdown(ctx)
		if err != nil {
			shutdownErrorChannel <- err
		}

		app.logger.PrintInfo("Completing background tasks...", map[string]string{
			"addr": server.Addr,
		})	

		app.wg.Wait()
		shutdownErrorChannel <- nil
	}()

	app.logger.PrintInfo("starting server", map[string]string{
		"addr": server.Addr,
		"env":  app.cfg.env,
	})

	err := server.ListenAndServe()
	
	if !errors.Is(err, http.ErrServerClosed) {
		return err 
	}

	err = <-shutdownErrorChannel
	if err != nil {
		return err
	}
	
	app.logger.PrintInfo("Stopped Server", map[string]string{
		"addr": server.Addr,
	})
	
	return nil
}
