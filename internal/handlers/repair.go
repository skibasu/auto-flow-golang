package handlers

import (
	"encoding/json"
	"net/http"
)

type Repair struct {
	Id string `json:"id"`
}

func (h *Handler) GetRepairs(w http.ResponseWriter, r *http.Request) {
	repairs := []Repair{
		{Id: "1"},
		{Id: "2"},
		{Id: "3"},
	}
	// 📦 response
	json.NewEncoder(w).Encode(repairs)
}
