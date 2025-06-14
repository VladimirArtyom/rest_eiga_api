package main

import (
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	"golang.org/x/time/rate"
)

func (app *application) recoverPanic(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			err := recover()
			if err != nil {
				w.Header().Set("Connection", "close")
				app.serverErrorResponse(w, r, fmt.Errorf("%s", err))
			}
		}()

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimitGlobal(next http.Handler) http.Handler {

	var limiter *rate.Limiter = rate.NewLimiter(2, 4)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !limiter.Allow() {
			app.rateLimitExceededResponse(w, r)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func (app *application) rateLimit(next http.Handler) http.Handler {
	type client struct {
		limiter  *rate.Limiter
		lastSeen time.Time
	}

	var (
		clients = make(map[string]*client)
		mutex   sync.Mutex
	)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if app.cfg.limiter.enabled {

			ip, _, err := net.SplitHostPort(r.RemoteAddr)
			if err != nil {
				app.serverErrorResponse(w, r, err)
				return
			}

			go func() {

				time.Sleep(1 * time.Minute)

				for {
					mutex.Lock()
					for ip, client := range clients {
						if time.Since(client.lastSeen) > 3*time.Minute {
							delete(clients, ip)
						}
					}

					mutex.Unlock()
				}
			}()

			mutex.Lock()

			_, found := clients[ip]
			if !found {
				clients[ip] = &client{
					limiter: rate.NewLimiter(rate.Limit(app.cfg.limiter.rps),
						app.cfg.limiter.burst),
				}
			}

			clients[ip].lastSeen = time.Now()
			if !clients[ip].limiter.Allow() {
				mutex.Unlock()
				app.rateLimitExceededResponse(w, r)
				return
			}

			mutex.Unlock()

		}
		next.ServeHTTP(w, r)
	})
}
