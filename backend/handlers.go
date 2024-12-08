// handlers.go
package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/patrickmn/go-cache"
)

type Handler struct {
	db    *sql.DB
	cache *cache.Cache
}

type Player struct {
	CharID      int     `json:"char_id"`
	Username    string  `json:"username"`
	ClassID     int     `json:"class_id"`
	RewardScore float64 `json:"reward_score"`
}

type PaginatedResponse struct {
	Players     []Player `json:"players"`
	TotalCount  int      `json:"total_count"`
	CurrentPage int      `json:"current_page"`
	TotalPages  int      `json:"total_pages"`
}

func NewHandler(db *sql.DB) *Handler {
	if db == nil {
		log.Fatal("Database connection is nil")
	}
	// Initialize cache with 5 minute expiration and 10 minute cleanup interval
	c := cache.New(5*time.Minute, 10*time.Minute)
	return &Handler{
		db:    db,
		cache: c,
	}
}

func (h *Handler) GetPlayers(w http.ResponseWriter, r *http.Request) {
	// Get query parameters
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))
	search := r.URL.Query().Get("search")

	if page < 1 {
		page = 1
	}
	if limit < 1 || limit > 100 {
		limit = 10 // Default limit
	}

	offset := (page - 1) * limit

	// Create cache key based on parameters
	cacheKey := getCacheKey(page, limit, search)

	// Try to get from cache first
	if cached, found := h.cache.Get(cacheKey); found {
		w.Header().Set("Content-Type", "application/json")
		w.Header().Set("X-Cache", "HIT")
		json.NewEncoder(w).Encode(cached)
		return
	}

	// Build query with search and pagination
	var query string
	var args []interface{}
	var countQuery string

	if search != "" {
		query = `
            SELECT 
                c.char_id,
                a.username,
                c.class_id,
                COALESCE(s.reward_score, 0) as reward_score
            FROM 
                character c
                JOIN account a ON c.acc_id = a.acc_id
                LEFT JOIN scores s ON c.char_id = s.char_id
            WHERE 
                LOWER(a.username) LIKE LOWER($1)
            ORDER BY 
                s.reward_score DESC NULLS LAST
            LIMIT $2 OFFSET $3
        `
		countQuery = `
            SELECT COUNT(*)
            FROM character c
            JOIN account a ON c.acc_id = a.acc_id
            WHERE LOWER(a.username) LIKE LOWER($1)
        `
		args = []interface{}{"%" + search + "%", limit, offset}
	} else {
		query = `
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
            LIMIT $1 OFFSET $2
        `
		countQuery = `
            SELECT COUNT(*)
            FROM character c
        `
		args = []interface{}{limit, offset}
	}

	// Get total count
	var totalCount int
	if search != "" {
		err := h.db.QueryRow(countQuery, "%"+search+"%").Scan(&totalCount)
		if err != nil {
			log.Printf("Error getting total count: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	} else {
		err := h.db.QueryRow(countQuery).Scan(&totalCount)
		if err != nil {
			log.Printf("Error getting total count: %v", err)
			http.Error(w, "Internal server error", http.StatusInternalServerError)
			return
		}
	}

	// Execute main query
	rows, err := h.db.Query(query, args...)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	players := []Player{}
	for rows.Next() {
		var p Player
		err := rows.Scan(&p.CharID, &p.Username, &p.ClassID, &p.RewardScore)
		if err != nil {
			log.Printf("Error scanning row: %v", err)
			continue
		}
		players = append(players, p)
	}

	totalPages := (totalCount + limit - 1) / limit

	response := PaginatedResponse{
		Players:     players,
		TotalCount:  totalCount,
		CurrentPage: page,
		TotalPages:  totalPages,
	}

	// Store in cache for 5 minutes
	h.cache.Set(cacheKey, response, cache.DefaultExpiration)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Cache", "MISS")
	json.NewEncoder(w).Encode(response)
}

func getCacheKey(page, limit int, search string) string {
	return fmt.Sprintf("players:page=%d:limit=%d:search=%s", page, limit, search)
}
