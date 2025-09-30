package scenario

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

// desc:
// -	T1
// -	T2
// -	T3

func RunReadCommitted() {
	pool, err := InitDBPool()
	if err != nil {
		log.Fatalf("Failed to innit DB Pool: %v", err)
	}

	// T1
	go func() {
		// init conn
		conn, err := pool.Acquire(context.Background())
		if err != nil {
			log.Fatalf("Session A: acquire conn: %v", err)
		}
		defer conn.Release()

		// begin tx
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{
			IsoLevel: pgx.ReadCommitted,
		})
		if err != nil {
			log.Fatalf("Session A: begin tx: %v", err)
		}
		defer safeRollBack("A", tx)

		// exec
		fmt.Println("[Session A] Start transaction, insert row (not committed yet)")

		sql_cmd := "INSERT INTO accounts (customer_name, balance) VALUES ($1, $2)"
		_, err = tx.Exec(context.Background(),
			sql_cmd, "user", 999)
		if err != nil {
			log.Fatalf("Session A: insert err: %v", err)
		}

		fmt.Println("[Session A] Insert done but not committed yet. Sleeping 10s ...")
		time.Sleep(10 * time.Second)

		// commit
		if err := tx.Commit(context.Background()); err != nil {
			log.Fatalf("Session A: commit: %v", err)
		}

		fmt.Println("[Session A] Commit done")
	}()

	time.Sleep(2 * time.Second)

	go func() {
		conn, err := pool.Acquire(context.Background())
		if err != nil {
			log.Fatalf("Session B: acquire conn: %v", err)
		}

		defer conn.Release()
		fmt.Println("[Session B] Start queries with Read Committed")

		for i := 0; i < 3; i++ {
			var count int
			sql_cmd := "SELECT COUNT(*) FROM accounts WHERE balance = 999"
			err = conn.QueryRow(context.Background(), sql_cmd).Scan(&count)
			if err != nil {
				log.Fatalf("Session B: query err: %v", err)
			}
			fmt.Printf("[Session B] Query %d: count rows with balance=999 → %d\n", i, count)
			time.Sleep(5 * time.Second)
		}
	}()

	time.Sleep(20 * time.Second)

	// T2
}
