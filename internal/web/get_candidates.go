package web

import "net/http"

func getCandidatesHandler(w http.ResponseWriter, r *http.Request) {
	candidates, err := database.GetCandidates()
	if err != nil {
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	if candidates == nil {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}

	JSONRes(w, candidates, http.StatusOK)
}
