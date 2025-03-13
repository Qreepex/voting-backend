package web

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/qreepex/voting-backend/internal/data"
	"github.com/qreepex/voting-backend/internal/redis"
)

var redisClient *redis.Redis
var database *data.Database

func Init(db *data.Database, redis *redis.Redis) {
	database = db
	redisClient = redis

	r := mux.NewRouter()

	v1Router := r.PathPrefix("/v1").Subrouter()

	r.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}).Methods(http.MethodGet)

	v1Router.HandleFunc("/candidates", getCandidatesHandler).Methods(http.MethodGet)
	v1Router.HandleFunc("/candidate/{id}", getCandidateHandler).Methods(http.MethodGet)
	v1Router.HandleFunc("/votes/{candidate}", getCandidateVotesHandler).Methods(http.MethodGet)
	v1Router.HandleFunc("/vote", postVoteHandler).Methods(http.MethodPost)

	http.Handle("/", r)

	http.ListenAndServe(":8080", nil)
}

func JSONRes(w http.ResponseWriter, err interface{}, code int) {
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	w.Header().Set("X-Content-Type-Options", "nosniff")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(err)
	log.Println("--> Response:", code)
	log.Println()
}
