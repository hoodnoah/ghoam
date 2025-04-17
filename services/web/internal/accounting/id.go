package accounting

import (
	// std
	"math/rand"
	"time"

	// external
	ulid "github.com/oklog/ulid/v2"
)

// generates a new unique ID using ULID;
// the ID is lexicographically sortable and based on time,
// allowing chronology to be indicated by the ID itself.
func NewID() string {
	t := time.Now().UTC()

	// Generate a ULID with a timestamp
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	return ulid.MustNew(ulid.Timestamp(t), entropy).String()
}
