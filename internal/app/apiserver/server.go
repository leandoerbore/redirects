package apiserver

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"time"

	"github.com/google/uuid"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"github.com/leandoerbore/redirects/internal/app/model/redirect"
	"github.com/leandoerbore/redirects/internal/app/store"
	"github.com/sirupsen/logrus"
)

const (
	sessionName        = "leandoer"
	ctxKeyUser  ctxKey = iota
	ctxKeyRequestId
)

var (
	errIncorrectEmailOrPassword = errors.New("incorrect email or password")
	errNotAuthenticated         = errors.New("not authenticated")
)

type ctxKey int8

type server struct {
	router       *mux.Router
	logger       *logrus.Logger
	store        store.Store
	sessionStore sessions.Store
}

func newServer(store store.Store, sessionStore sessions.Store) *server {
	s := &server{
		router:       mux.NewRouter(),
		logger:       logrus.New(),
		store:        store,
		sessionStore: sessionStore,
	}

	s.configureRouter()

	return s
}

func (s *server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.router.ServeHTTP(w, r)
}

func (s *server) configureRouter() {
	s.router.Use(s.setRequestID)
	s.router.Use(s.logRequest)
	s.router.Use(handlers.CORS(handlers.AllowedOrigins([]string{"*"})))

	s.router.HandleFunc("/redirects", s.handleRedirectsGetAll()).Methods("GET")
	s.router.HandleFunc("/redirects", s.handleRedirectsCreate()).Methods("POST")
	s.router.HandleFunc("/redirects/{id}", s.handleRedirectsFind()).Methods("GET")
	s.router.HandleFunc("/redirects/{id}", s.handleRedirectsUpdate()).Methods("PATCH")
	s.router.HandleFunc("/redirects/{id}", s.handleRedirectsRemove()).Methods("DELETE")
}

func (s *server) setRequestID(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		id := uuid.New().String()
		w.Header().Set("X-Request-ID", id)
		next.ServeHTTP(w, r.WithContext(context.WithValue(r.Context(), ctxKeyRequestId, id)))
	})
}

func (s *server) logRequest(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		logger := s.logger.WithFields(logrus.Fields{
			"remote_addr": r.RemoteAddr,
			"request_id":  r.Context().Value(ctxKeyRequestId),
		})
		logger.Infof("started %s %s", r.Method, r.RequestURI)

		start := time.Now()
		rw := &responseWriter{w, http.StatusOK}
		next.ServeHTTP(rw, r)

		logger.Infof(
			"complited with %d %s in %v",
			rw.code,
			http.StatusText(rw.code),
			time.Now().Sub(start))
	})
}

func (s *server) handleRedirectsGetAll() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		arrRdir, err := s.store.Redirect().GetAll()
		if err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, arrRdir)
	}
}

func (s *server) handleRedirectsCreate() http.HandlerFunc {
	type request struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
		IsActive    bool   `json:"is_active"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		rdir := &redirect.Redirect{
			Source:      req.Source,
			Destination: req.Destination,
			IsActive:    req.IsActive,
		}

		if err := s.store.Redirect().Create(rdir); err != nil {
			s.error(w, r, http.StatusUnprocessableEntity, err)
			return
		}

		s.respond(w, r, http.StatusCreated, rdir)
	}
}

func (s *server) handleRedirectsFind() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		rdir, err := s.store.Redirect().Find(int(id))

		s.respond(w, r, http.StatusOK, rdir)
	}
}

func (s *server) handleRedirectsUpdate() http.HandlerFunc {
	type request struct {
		Source      string `json:"source"`
		Destination string `json:"destination"`
		IsActive    bool   `json:"is_active"`
	}

	return func(w http.ResponseWriter, r *http.Request) {
		req := &request{}
		if err := json.NewDecoder(r.Body).Decode(req); err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		rdir := &redirect.Redirect{
			ID:          id,
			Source:      req.Source,
			Destination: req.Destination,
			IsActive:    req.IsActive,
		}

		if err := s.store.Redirect().Update(rdir); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) handleRedirectsRemove() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			s.error(w, r, http.StatusBadRequest, err)
			return
		}

		if err := s.store.Redirect().Remove(id); err != nil {
			s.error(w, r, http.StatusInternalServerError, err)
			return
		}

		s.respond(w, r, http.StatusOK, nil)
	}
}

func (s *server) error(w http.ResponseWriter, r *http.Request, code int, err error) {
	s.respond(w, r, code, map[string]string{"error": err.Error()})
}

func (s *server) respond(w http.ResponseWriter, r *http.Request, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		json.NewEncoder(w).Encode(data)
	}
}
