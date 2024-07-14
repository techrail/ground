package main

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/oklog/ulid/v2"
)

func main() {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	fmt.Println(ulid.New(ms, entropy))
	fmt.Println(ulid.New(ms, entropy))
}
