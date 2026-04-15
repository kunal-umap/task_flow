package types

// contextKey is a private type to avoid collisions in context
type contextKey string

// UserContextKey is used to store JWT claims in request context
const UserContextKey = contextKey("user")
