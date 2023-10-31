package uuid

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/oklog/ulid/v2"
	"github.com/techrail/ground/logger"
	"math/rand"
	"time"
)

// GetNewUlid returns a new ULID value. If the ULID library fails, it creates a virtual ULID value instead
func GetNewUlid() ulid.ULID {
	entropy := ulid.Monotonic(rand.New(rand.NewSource(time.Now().UTC().UnixNano())), 0)
	newUlid, err := ulid.New(ulid.Timestamp(time.Now().UTC()), entropy)
	if err != nil {
		errMsg := fmt.Sprintf("E#1MRWPP - Could not create a new ULID. Error: %v", err)
		logger.Println(errMsg)
		return getVirtualUlid()
	}

	// NOTE: We can also use the following statement. But that one uses a 'MustNew' call underneath
	// 	and any failures will cause the program to panic. Also, the function does essentially the same
	// 	thing that we have done above!
	//newUlid = ulid.Make()

	return newUlid
}

// GetNewUlidAsUuidString returns the UUID string representation (with dashes) of a new ULID value
func GetNewUlidAsUuidString() string {
	return UlidToUuidString(GetNewUlid())
}

// GetNewUlidAsString returns the typical string representation of a new ULID value
func GetNewUlidAsString() string {
	return GetNewUlid().String()
}

// UlidToUuidString converts a ULID to a UUID string
func UlidToUuidString(ulidType ulid.ULID) string {
	return uuid.UUID(ulidType).String()
}

// getVirtualUlid function here tried to return a ULID type by filling the bits using two pieces of information
// The first 8 bytes are the current UTC Unix Timestamp with nanosecond precision. The last 8 bytes are random
func getVirtualUlid() ulid.ULID {
	ts := time.Now().UTC().UnixNano()
	r := rand.New(rand.NewSource(time.Now().UTC().UnixNano()))
	rnd := r.Uint64()
	return [16]byte{
		byte(0xff & (ts >> (8 * 7))),
		byte(0xff & (ts >> (8 * 6))),
		byte(0xff & (ts >> (8 * 5))),
		byte(0xff & (ts >> (8 * 4))),
		byte(0xff & (ts >> (8 * 3))),
		byte(0xff & (ts >> (8 * 2))),
		byte(0xff & (ts >> (8 * 1))),
		byte(0xff & (ts >> (8 * 0))),
		byte(0xff & (rnd >> (8 * 0))),
		byte(0xff & (rnd >> (8 * 1))),
		byte(0xff & (rnd >> (8 * 2))),
		byte(0xff & (rnd >> (8 * 3))),
		byte(0xff & (rnd >> (8 * 4))),
		byte(0xff & (rnd >> (8 * 5))),
		byte(0xff & (rnd >> (8 * 6))),
		byte(0xff & (rnd >> (8 * 7))),
	}
}
