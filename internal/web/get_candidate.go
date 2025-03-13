package web

import (
	"net/http"

	"github.com/gorilla/mux"
)

func getCandidateHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	candidateId := vars["id"]

	if candidateId == "" {
		http.Error(w, "Bad request", http.StatusBadRequest)
		return
	}

	candidate, err := database.GetCandidate(candidateId)
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	if candidate == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	JSONRes(w, candidate, http.StatusOK)
}
