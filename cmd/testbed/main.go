package main

import (
	"fmt"
	"github.com/oklog/ulid/v2"
	"math/rand"
	"time"
)

func main() {
	entropy := rand.New(rand.NewSource(time.Now().UnixNano()))
	ms := ulid.Timestamp(time.Now())
	fmt.Println(ulid.New(ms, entropy))
	fmt.Println(ulid.New(ms, entropy))
}
