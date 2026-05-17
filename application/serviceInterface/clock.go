package serviceInterface

import "time"

// Clock provides the current time, injectable for testing.
type Clock interface {
	Now() time.Time
}
