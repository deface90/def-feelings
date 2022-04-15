package rest

import (
	"context"
	"fmt"
	"github.com/deface90/def-feelings/storage"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	"net/http"
	"sync"
	"time"

	"github.com/go-chi/chi/v5"
	log "github.com/sirupsen/logrus"
)

type Rest struct {
	Auth   Auth
	engine storage.Engine
	config storage.Config

	httpServer *http.Server
	lock       sync.Mutex

	logger *log.Logger
}

func NewRestService(e storage.Engine, c storage.Config, l *log.Logger) *Rest {
	return &Rest{
		Auth:   Auth{Engine: e},
		engine: e,
		config: c,
		logger: l,
	}
}

func (s *Rest) Run() {
	s.logger.Infof("Starting HTTP rest server on port %s", s.config.Port)

	s.lock.Lock()
	s.httpServer = s.makeHTTPServer(s.config.Port, s.routes())
	s.lock.Unlock()

	err := s.httpServer.ListenAndServe()
	s.logger.Infof("HTTP rest server terminated, %s", err)
}

func (s *Rest) routes() chi.Router {
	router := chi.NewRouter()
	router.Use(middleware.Heartbeat("/ping"))

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders: []string{"*"},
		ExposedHeaders: []string{"Link", "X-Pagination-Total-Count", "X-Pagination-Page", "X-Pagination-Per-Page",
			"X-List-Id", "X-List-Title"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Route("/api/v1", func(routerFront chi.Router) {
		routerFront.Group(func(routerPrivate chi.Router) {
			routerPrivate.Use(s.Auth.Middleware)

			routerPrivate.Post("/user/edit/{id}", s.editUserHandler)
			routerPrivate.Get("/user/delete/{id}", s.deleteUserHandler)
			routerPrivate.Get("/user/get/{id}", s.getUserHandler)
			routerPrivate.Post("/user/list", s.listUsersHandler)
			routerPrivate.Post("/user/subscribe/{id}", s.subscribeUserHandler)

			routerPrivate.Post("/feeling/list", s.listFeelingsHandler)
			routerPrivate.Post("/feeling/frequency", s.feelingsFrequencyHandler)

			routerPrivate.Post("/status/create", s.createStatusHandler)
			routerPrivate.Get("/status/get/{id}", s.getStatusHandler)
			routerPrivate.Post("/status/list", s.listStatusesHandler)

			routerPrivate.Get("/logout", s.logoutHandler)
		})

		routerFront.Post("/user/create", s.createUserHandler)
		routerFront.Post("/auth/login", s.loginUserCtrl)
		routerFront.Post("/auth/session", s.checkSessionCtrl)
	})

	return router
}

// Shutdown rest http server
func (s *Rest) Shutdown() {
	s.logger.Infof("Shutdown HTTP rest server")
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	s.lock.Lock()
	if s.httpServer != nil {
		if err := s.httpServer.Shutdown(ctx); err != nil {
			s.logger.WithError(err).Warnf("HTTP rest shutdown error")
		}
		s.logger.Infof("Shutdown HTTP rest server completed")
	}
	s.lock.Unlock()
	err := s.engine.Shutdown()
	if err != nil {
		s.logger.WithError(err).Warnf("Failed to gracefull shutdown engine")
	}
}

func (s *Rest) makeHTTPServer(port string, router http.Handler) *http.Server {
	return &http.Server{
		Addr:              fmt.Sprintf(":%s", port),
		Handler:           router,
		ReadHeaderTimeout: 5 * time.Second,
		IdleTimeout:       30 * time.Second,
	}
}
