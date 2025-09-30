package scenario

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

func RunRepeatableRead() {
	pool, err := InitDBPool()
	if err != nil {
		log.Fatalf("Init DB pool failed: %v", err)
	}

	go func() {
		conn, err := pool.Acquire(context.Background())
		if err != nil {
			log.Fatalf("Session A: acquire conn: %v", err)
		}
		defer conn.Release()

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{
			IsoLevel: pgx.RepeatableRead,
		})
		if err != nil {
			log.Fatalf("Session A: begin tx: %v", err)
		}
		defer safeRollBack("A", tx)

		fmt.Println("[Session A] Begin with Repeatable Read")
		var count int
		sql_cmd := "SELECT COUNT(*) FROM accounts WHERE balance=777"
		if err := tx.QueryRow(context.Background(), sql_cmd).Scan(&count); err != nil {
			log.Fatalf("Session A: second select: %v", err)
		}
		fmt.Printf("[Session A] First count balance=777 → %d\n", count)

		fmt.Println("[Session A] Sleeping 12s before re-query...")
		time.Sleep(12 * time.Second)

		// requery
		if err := tx.QueryRow(context.Background(), sql_cmd).Scan(&count); err != nil {
			log.Fatalf("Session A: second select: %v", err)
		}
		fmt.Printf("[Session A] Second count balance=777 → %d\n", count)

		if err := tx.Commit(context.Background()); err != nil {
			log.Printf("[Session A] Commit error: %v\n", err)
		} else {
			fmt.Println("[Session A] Commit done")
		}
	}()

	time.Sleep(2 * time.Second)

	// B
	go func() {
		conn, err := pool.Acquire(context.Background())
		if err != nil {
			log.Fatalf("Session B: acquire conn: %v", err)
		}
		defer conn.Release()

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{
			IsoLevel: pgx.RepeatableRead,
		})
		if err != nil {
			log.Fatalf("Session B: begin tx: %v", err)
		}
		defer safeRollBack("B", tx)

		sql_cmd := "INSERT INTO accounts (customer_name, balance) VALUES ($1, $2)"
		fmt.Println("[Session B] Inserting new row balance=777")
		_, err = tx.Exec(context.Background(), sql_cmd, "temp", 777)
		if err != nil {
			log.Fatalf("Session B: insert: %v", err)
		}

		time.Sleep(5 * time.Second)
		if err := tx.Commit(context.Background()); err != nil {
			log.Fatalf("Session B: commit: %v", err)
		}
		fmt.Println("[Session B] Commit done")
	}()
	time.Sleep(20 * time.Second)
}
