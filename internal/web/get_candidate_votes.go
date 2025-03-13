package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

func getCandidateVotesHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	candidateId := vars["candidate"]

	if candidateId == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	votes, err := database.GetCandidateVotes(candidateId)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	if votes == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	timestamps := make([]int64, 0)
	for _, vote := range votes {
		timestamps = append(timestamps, vote.Timestamp.Unix())
	}

	JSONRes(w, timestamps, http.StatusOK)
}
