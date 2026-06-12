// Package ids generates UUIDv7 identifiers: time-ordered and index-friendly,
// so the id doubles as the pagination cursor (docs/DATA_MODEL.md).
package ids

import "github.com/google/uuid"

func New() uuid.UUID {
	id, err := uuid.NewV7()
	if err != nil {
		// Only possible on entropy-source failure — not a recoverable
		// application error.
		panic(err)
	}
	return id
}
