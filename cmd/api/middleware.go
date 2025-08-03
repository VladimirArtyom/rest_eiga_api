package main

import (
	"errors"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time" 
	"github.com/VladimirArtyom/rest_eiga_api/internal/data"
	"github.com/VladimirArtyom/rest_eiga_api/internal/validator"
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


func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r* http.Request) {

		w.Header().Add("Vary", "Authorization")

		var authorizationHeader string = r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = app.contextSetUser(r, data.AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}


		var headerParts []string = strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		var token string = headerParts[1]

		var v *validator.Validator = validator.New()

		data.ValidateToken(v, token)

		if !v.Valid() {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		user, err := app.models.Users.GetForToken(data.ScopeAuthentication, token)
		if err != nil {
			switch {
			case errors.Is(err, data.ErrRecordNotFound):
				app.invalidAuthenticationTokenResponse(w, r)
			default:
				app.serverErrorResponse(w, r, err)
			}
			return 
		}
		r = app.contextSetUser(r, user)
		next.ServeHTTP(w, r)
	})
}
 
func (app *application) enableCORS(next http.Handler) http.Handler{
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Vary is to give additional cache key information
		// It is to distinguish every response that was made to the server.
		w.Header().Add("Vary", "Origin")
		
		w.Header().Add("Vary", "Access-Control-Request-Method")

		origin := r.Header.Get("origin")
		if (app.cfg.cors.origins[origin]) {
			w.Header().Set("Access-Control-Allow-Origin", origin)
			w.Header().Set("Access-Control-Allow-Credentials", "true")
			
			// Si la request est Preflight
			if r.Method == http.MethodOptions && r.Header.Get("Access-Control-Request-Method") != "" {
				w.Header().Set("Access-Control-Allow-Methods", "OPTIONS, PUT, PATCH, DELETE" )
				w.Header().Set("Access-Control-Allow-Headers", "Authorization, Content-Type")
				
				w.Header().Set("Access-Control-Max-Age", "300")
				w.WriteHeader(http.StatusOK)
				return 
			}
		}
		next.ServeHTTP(w, r)
	})
} 

func (app *application) requirePermission(code string, next http.HandlerFunc) http.HandlerFunc {
		fn := func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		fmt.Println(user.ID)
		permissions, err := app.models.Permissions.GetAllForUser(user.ID)
		fmt.Println(permissions)
		if err != nil {
			app.notPermittedResponse(w,r)
			return
		}
		if !permissions.Include(code) {
			app.notPermittedResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	} 

	return app.requireActivatedUser(fn)
}

func (app* application) requireActivatedUser(next http.HandlerFunc) http.HandlerFunc {
	fn :=  func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		if !user.Activated{
			app.inactiveAccountResponse(w, r)
			return
		}
		next.ServeHTTP(w,r)
	}

	return app.requireAuthenticatedUser(fn)
}

func (app *application) requireAuthenticatedUser(next http.HandlerFunc) http.HandlerFunc {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request ) {
		user := app.contextGetUser(r)
		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	})
	
}

