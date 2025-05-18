package main

import (
	"log"

	verify "FMgo/scripts/verify/pkg"
)

func main() {
	if err := verify.VerifySupabase(); err != nil {
		log.Fatal(err)
	}
}
