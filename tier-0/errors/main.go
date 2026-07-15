package main

import (
	"errors"
	"fmt"
	exceptions "tier-0/errors/domain"
)

var ErrPostgresConnection = errors.New("sql: database pg_hsm is down")

func SearchKey() error {
	errSql := ErrPostgresConnection

	if errSql != nil {
		return exceptions.NewDomainError(
			exceptions.ErrCodeInternal,
			"error connecting to the database",
			errSql,
		)
	}

	return nil
}
func main() {
	err := SearchKey()

	if err != nil {
		if de, ok := err.(*exceptions.DomainError); ok {
			de.WithDetails("retry_count", 3).WithDetails("table", "hsm_leys")
		}

		code := exceptions.GetErrorCode(err)
		fmt.Println("Error code:", code)

		var DomainErr *exceptions.DomainError

		if errors.As(err, &DomainErr) {
			fmt.Println("Domain error message:", DomainErr.Message)
			fmt.Println("Domain error details:", DomainErr.Details)
		}

		if errors.Is(err, ErrPostgresConnection) {
			fmt.Println("The error is related to the Postgres connection.")
		} else {
			fmt.Println("The error is not related to the Postgres connection.")
		}
	}
}
