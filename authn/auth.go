package authn

import (
	"strings"
	"time"

	"github.com/google/uuid"
)

// Auth is an auth record
type Auth struct {
	Auth      AuthKey
	CreatedAt *time.Time
}

// AuthKey is an API auth keys
type AuthKey string

// NilAuthKey is a nil auth key
var NilAuthKey = AuthKey("")

// NewAuthKey generates a new auth key
func NewAuthKey() AuthKey {
	s := strings.ReplaceAll(uuid.NewString(), "-", "")
	return AuthKey(s)
}
