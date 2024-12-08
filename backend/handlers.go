// handlers.go
package main

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
)

type Handler struct {
	db *sql.DB
}

type Player struct {
	CharID      int     `json:"char_id"`
	Username    string  `json:"username"`
	ClassID     int     `json:"class_id"`
	RewardScore float64 `json:"reward_score"`
}

func NewHandler(db *sql.DB) *Handler {
	return &Handler{
		db: db,
	}
}

func (h *Handler) GetPlayers(w http.ResponseWriter, r *http.Request) {
	players := []Player{}

	query := `
        SELECT 
            c.char_id,
            a.username,
            c.class_id,
            COALESCE(s.reward_score, 0) as reward_score
        FROM 
            character c
            JOIN account a ON c.acc_id = a.acc_id
            LEFT JOIN scores s ON c.char_id = s.char_id
        ORDER BY 
            s.reward_score DESC NULLS LAST
    `

	rows, err := h.db.Query(query)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	for rows.Next() {
		var p Player
		err := rows.Scan(&p.CharID, &p.Username, &p.ClassID, &p.RewardScore)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		players = append(players, p)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(players)
}
