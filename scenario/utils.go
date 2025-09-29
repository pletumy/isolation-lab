package scenario

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
)

func QueryAccounts(tx pgx.Tx, label string) error {
	sql_cmd := "SELECT * FROM accounts"
	rows, err := tx.Query(context.Background(), sql_cmd)
	if err != nil {
		return err
	}
	defer rows.Close()

	for rows.Next() {
		var id, balance int
		var customer_name string

		if err := rows.Scan(&id, &customer_name, &balance); err != nil {
			return err
		}
		log.Printf("[%s] ID: %d, Name: %s, Balance: %d", label, id, customer_name, balance)
	}
	return rows.Err()
}

func PrintStep(session, msg string) {
	fmt.Printf("[%s] %s\n", session, msg)
}
