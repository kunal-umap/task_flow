package utils

import (
	"net/http"

	"github.com/google/uuid"
)

// 🔥 shared context key
type contextKey string

const UserContextKey = contextKey("user")

// JWTClaims must already exist in utils/jwt.go

func GetUserID(r *http.Request) (uuid.UUID, bool) {

	val := r.Context().Value(UserContextKey)
	if val == nil {
		return uuid.Nil, false
	}

	claims, ok := val.(*JWTClaims)
	if !ok {
		return uuid.Nil, false
	}

	id, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil, false
	}

	return id, true
}
