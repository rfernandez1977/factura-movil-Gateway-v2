package main

import (
	"log"

	verify "github.com/fmgo/scripts/verify/pkg"
)

func main() {
	if err := verify.VerifySupabase(); err != nil {
		log.Fatal(err)
	}
}
