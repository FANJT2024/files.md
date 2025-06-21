package server

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"log/slog"
	"net/http"
	"strconv"
	"sync"
	"time"

	"github.com/spf13/afero"

	"zakirullin/stuffbot/config"
	"zakirullin/stuffbot/internal/fs"
)

const (
	TokenLength        = 32
	OneTimeTokenExpiry = 10 * time.Minute
)

var (
	oneTimeTokens = make(map[string]oneTimeToken)
	mu            sync.RWMutex
)

type oneTimeToken struct {
	userID    int64
	expiresAt time.Time
}

func GenerateOneTimeToken(userID int64) string {
	token := generateToken()

	mu.Lock()
	oneTimeTokens[token] = oneTimeToken{
		userID:    userID,
		expiresAt: time.Now().Add(OneTimeTokenExpiry),
	}
	mu.Unlock()

	return token
}

func UserID(token string) (int64, bool) {
	tokens, err := fs.NewFS(config.BotCfg.TokensDir, afero.NewOsFs())
	if err != nil {
		slog.Error("Failed to create file system for tokens", "error", err)
		return 0, false
	}

	data, err := tokens.Read(fs.DirRoot, token)
	if err != nil {
		return 0, false
	}

	userID, err := strconv.ParseInt(data, 10, 64)
	if err != nil {
		return 0, false
	}

	return userID, true
}

func IssueToken(w http.ResponseWriter, r *http.Request) {
	defer func() {
		if r := recover(); r != nil {
			log.Printf("PANIC in IssueToken: %v", r)
			http.Error(w, "Internal server error", 500)
		}
	}()

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		OneTimeToken string `json:"oneTimeToken"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	permanentToken, ok := issueNewToken(req.OneTimeToken)
	if !ok {
		http.Error(w, "Invalid or expired token", http.StatusUnauthorized)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	err := json.NewEncoder(w).Encode(map[string]string{"token": permanentToken})
	if err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
		return
	}
}

// TODO CHECK that user id belongs to oneTimeToken ID, or get user id by oneTimeToken
func authMiddleware(next http.HandlerFunc, tokensDir string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		token := r.Header.Get("Authorization")
		if token != AuthToken {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		next(w, r)
	}
}

func issueNewToken(oneTimeToken string) (string, bool) {
	mu.Lock()
	data, exists := oneTimeTokens[oneTimeToken]
	if !exists || time.Now().After(data.expiresAt) {
		mu.Unlock()
		return "", false
	}
	delete(oneTimeTokens, oneTimeToken)
	mu.Unlock()

	token := generateToken()
	tokens, err := fs.NewFS(config.BotCfg.TokensDir, afero.NewOsFs())
	if err != nil {
		slog.Error("Failed to create file system for tokens", "error", err)
		return "", false
	}
	err = tokens.Write(fs.DirRoot, token, strconv.FormatInt(data.userID, 10))
	if err != nil {
		return "", false
	}

	return token, true
}

func generateToken() string {
	bytes := make([]byte, TokenLength)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
