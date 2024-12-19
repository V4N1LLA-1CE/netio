package main

import (
	"fmt"
	"github.com/V4N1LLA-1CE/netio"
	"regexp"
)

func main() {
	v := netio.NewValidator()

	email := "invalid-email"
	age := 15
	role := "superuser"
	interests := []string{"coding", "coding"}

	v.Check(regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`).MatchString(email), "email", "Invalid email")
	v.Check(age >= 18, "age", "Must be 18 or older")
	v.Check(netio.IsIn(role, "admin", "user", "moderator"), "role", "Invalid role")
	v.Check(!netio.HasDuplicates(interests), "interests", "Duplicate interests found")

	if !v.Valid() {
		fmt.Println("Errors:", v.Errors)
	}
}
