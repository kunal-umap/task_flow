package utils

import (
	"net/http"

	"taskflow/internal/types"

	"github.com/google/uuid"
)

// IMPORTANT: exported (capital G)
func GetUserID(r *http.Request) (uuid.UUID, bool) {
	claims, ok := r.Context().Value(types.UserContextKey).(*JWTClaims)
	if !ok {
		return uuid.Nil, false
	}

	id, err := uuid.Parse(claims.UserID)
	if err != nil {
		return uuid.Nil, false
	}

	return id, true
}
