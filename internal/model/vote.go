package model

import "time"

type Vote struct {
	ID        string    `json:"_id,omitempty" bson:"_id,omitempty"`
	Candidate string    `json:"candidate"`
	Campaign  string    `json:"campaign"`
	IpAddress string    `json:"ipAddress"`
	Cookie    string    `json:"cookie"`
	Timestamp time.Time `json:"timestamp"`
}

type ApiVote struct {
	Candidate string `json:"candidate"`
	Campaign  string `json:"campaign"`
}

func NewVote(candidate, campaign, ipAddress, cookie string) *Vote {
	return &Vote{
		Candidate: candidate,
		Campaign:  campaign,
		IpAddress: ipAddress,
		Cookie:    cookie,
		Timestamp: time.Now(),
	}
}
