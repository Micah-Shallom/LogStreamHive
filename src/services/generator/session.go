package main

import (
	"math/rand"
	"time"
)

type UserSession struct {
	State      string
	LastUpdate time.Time
}

func (lg *LogGenerator) updateUserSession(user_id string) {
	states := []string{
		"login", "browse", "search", "view_item", "add_to_cart", "checkout", "purchase", "logout",
	}

	if _, exist := lg.userSessions[user_id]; !exist {
		// setup a new user session
		lg.userSessions[user_id] = UserSession{
			State:      "login",
			LastUpdate: time.Now(),
		}
		return
	}

	user_state := lg.userSessions[user_id].State
	var currentIdx int
	for i := range states {
		if states[i] == user_state {
			currentIdx = i
			break
		}
	}

	next_idx := (currentIdx + 1) % len(states)
	lg.userSessions[user_id] = UserSession{
		State:      states[next_idx],
		LastUpdate: time.Now(),
	}

	//cleanup if more than 100 sessions
	if len(lg.userSessions) > 100 {
		current_time := time.Now()

		for uid, session := range lg.userSessions {
			if current_time.Sub(session.LastUpdate) > 5*time.Minute {
				delete(lg.userSessions, uid)
			}
		}
	}
}

func (lg *LogGenerator) createMessageFromPattern(log_type, user_id string) string {
	var state string
	if _, exists := lg.userSessions[user_id]; exists {
		state = lg.userSessions[user_id].State
		sessionMessages := map[string]string{
			"login":       "User logged in successfully",
			"browse":      "User browsing product catalog",
			"search":      "User performed search query",
			"view_item":   "User viewing product details",
			"add_to_cart": "User added item to cart",
			"checkout":    "User initiated checkout process",
			"purchase":    "User completed purchase",
			"logout":      "User logged out",
		}

		if msg, ok := sessionMessages[state]; ok {
			return msg
		}
	}

	messages := map[string][]string{
		"INFO": {
			"User logged in successfully",
			"Page loaded in 0.2 seconds",
			"Database connection established",
			"Cache refreshed successfully",
			"API request completed",
		},
		"WARNING": {
			"High memory usage detected",
			"API response time exceeding threshold",
			"Database connection pool running low",
			"Retry attempt for failed operation",
			"Cache miss rate increasing",
		},
		"ERROR": {
			"Failed to connect to database",
			"API request timeout",
			"Invalid user credentials",
			"Processing error in data pipeline",
			"Out of memory error",
		},
		"DEBUG": {
			"Function X called with parameters Y",
			"SQL query execution details",
			"Cache lookup performed",
			"Request headers processed",
			"Internal state transition",
		},
	}

	if msgs, ok := messages[log_type]; ok {
		return msgs[rand.Intn(len(msgs))]
	}

	return "Sample log message for " + log_type

}
