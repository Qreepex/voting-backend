package web

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"time"

	"github.com/qreepex/voting-backend/internal/config"
	"github.com/qreepex/voting-backend/internal/model"
	"github.com/qreepex/voting-backend/internal/service"
	"github.com/redis/go-redis/v9"
	"go.mongodb.org/mongo-driver/mongo"
)

var cookieName = config.EnvMustGet("COOKIE_NAME")
var ipAddressHeader = config.EnvMustGet("IP_ADDRESS_HEADER")

// checks cookies and if cookie is limited returns seconds until eligible
// return cookie, ip, pending, error
// pending = 0 -> not pending / no cookie
// pending = -1 -> error
func CheckCookie(w http.ResponseWriter, r *http.Request) (string, string, int64, error) {
	cookie, err := r.Cookie(cookieName)
	if err != nil && !errors.Is(err, http.ErrNoCookie) {
		return "", "", -1, nil
	}

	if cookie != nil && cookie.Value != "" {
		log.Printf("Checking cookie: %v", cookie.Value)

		redisCookie, err := redisClient.CheckCookie(cookie.Value)
		if err != nil && err != redis.Nil {
			log.Printf("Failed to check cookie: %v", err)
			http.Error(w, "Internal error", http.StatusInternalServerError)
			return cookie.Value, "", -1, err
		}

		if redisCookie == nil {
			return cookie.Value, "", 0, nil
		}

		return cookie.Value, redisCookie.Ip, redisCookie.Pending, nil
	}

	return "", "", 0, nil
}

// checks IP and if IP is limited returns seconds until eligible
// return ip, cookie, pending, error
// pending = 0 -> not pending
// pending = -1 -> error
func CheckIP(w http.ResponseWriter, r *http.Request) (string, string, int64, error) {
	ipAddress := r.Header.Get(ipAddressHeader)
	if ipAddress == "" {
		log.Printf("Failed to get IP address from header: %v", ipAddressHeader)

		ipAddress = "127.0.0.99"
	}

	hashedIpAdress := service.HashIp(ipAddress)

	log.Printf("Checking IP: %v", hashedIpAdress)

	redisIp, err := redisClient.CheckIP(hashedIpAdress)
	if err != nil && err != redis.Nil {
		log.Printf("Failed to check IP: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return hashedIpAdress, "", -1, err
	}

	if redisIp == nil {
		return hashedIpAdress, "", 0, nil
	}

	return hashedIpAdress, redisIp.Cookie, redisIp.Pending, nil
}

func postVoteHandler(w http.ResponseWriter, r *http.Request) {
	cookieValue, cookieIp, cookiePending, err := CheckCookie(w, r)
	if err != nil || cookiePending == -1 {
		log.Printf("Failed to check cookie: %v", err)
		JSONRes(w, map[string]string{"error": "Internal error"}, http.StatusInternalServerError)
		return
	}

	ipValue, ipCookie, ipPending, err := CheckIP(w, r)
	if err != nil || ipPending == -1 {
		log.Printf("Failed to check IP: %v", err)
		JSONRes(w, map[string]string{"error": "Internal error"}, http.StatusInternalServerError)
		return
	}

	if cookiePending > 0 || ipPending > 0 {
		if cookieValue != ipCookie || cookiePending != ipPending {
			log.Printf("Cookie and IP mismatch: %v %v %v %v %v %v", cookieValue, ipValue, cookieIp, ipCookie, cookiePending, ipPending)
			if ipPending > cookiePending {
				newCookieValeue := cookieValue
				if cookieValue == "" {
					if ipCookie != "" {
						newCookieValeue = ipCookie
					} else {
						newCookieValeue = service.GenerateUniqueCookie()
						err := redisClient.SetIP(ipValue, time.Now().Add(time.Duration(cookiePending)*time.Second), newCookieValeue)
						if err != nil {
							log.Printf("Failed to set IP: %v", err)
							JSONRes(w, map[string]string{"error": "Internal error"}, http.StatusInternalServerError)
							return
						}
					}
				}

				err := redisClient.SetCookie(newCookieValeue, time.Now().Add(time.Duration(ipPending)*time.Second), ipValue)
				newCookie := http.Cookie{
					Name:     cookieName,
					Value:    newCookieValeue,
					Path:     "/",
					Expires:  time.Now().Add(time.Duration(ipPending) * time.Second),
					HttpOnly: false,
					SameSite: http.SameSiteLaxMode,
					Secure:   false,
				}
				http.SetCookie(w, &newCookie)

				if err != nil {
					log.Printf("Failed to set cookie: %v", err)
					JSONRes(w, map[string]string{"error": "Internal error"}, http.StatusInternalServerError)
					return
				}
			} else {

				err := redisClient.SetIP(ipValue, time.Now().Add(time.Duration(cookiePending)*time.Second), cookieValue)
				if err != nil {
					log.Printf("Failed to set IP: %v", err)
					JSONRes(w, map[string]string{"error": "Internal error"}, http.StatusInternalServerError)
					return
				}
			}

		}

		JSONRes(w, map[string]string{"error": "You have already voted"}, http.StatusTooManyRequests)
		return
	}

	var vote model.ApiVote
	err = json.NewDecoder(r.Body).Decode(&vote)
	if err != nil {
		log.Printf("Failed to decode vote: %v", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	log.Printf("Got vote for %v in campaign %v", vote.Candidate, vote.Campaign)

	candidate, err := database.GetCandidate(vote.Candidate)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Printf("Failed to get candidate: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	if candidate == nil || err == mongo.ErrNoDocuments {
		JSONRes(w, map[string]string{"error": "Candidate not found"}, http.StatusNotFound)
		return
	}

	campaign, err := database.GetCampaign(vote.Campaign)
	if err != nil && err != mongo.ErrNoDocuments {
		log.Printf("Failed to get campaign: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	if campaign == nil || err == mongo.ErrNoDocuments {
		JSONRes(w, map[string]string{"error": "Campaign not found"}, http.StatusNotFound)
		return
	}

	voteModel := model.NewVote(vote.Candidate, vote.Campaign, ipValue, cookieValue)

	_, err = database.CreateVote(*voteModel)
	if err != nil {
		log.Printf("Failed to create vote: %v", err)
		http.Error(w, "Internal error", http.StatusInternalServerError)
		return
	}

	newCookieValue := cookieValue
	if cookieValue == "" {
		newCookieValue = service.GenerateUniqueCookie()
	}

	redisClient.SetVote(ipValue, newCookieValue)

	newCookie := http.Cookie{
		Name:     cookieName,
		Value:    newCookieValue,
		Path:     "/",
		Expires:  time.Now().Add(time.Hour * 24),
		HttpOnly: false,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
	}

	http.SetCookie(w, &newCookie)
	JSONRes(w, map[string]string{"message": "Vote registered"}, http.StatusOK)
}
