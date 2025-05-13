package main

import (
	"log"

	verify "github.com/cursor/FMgo/scripts/verify/pkg"
)

func main() {
	if err := verify.VerifySupabase(); err != nil {
		log.Fatal(err)
	}
}
