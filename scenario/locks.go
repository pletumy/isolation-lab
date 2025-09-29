package scenario

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/jackc/pgx/v5"
)

// RunLocks demo FOR UPDATE vs FOR UPDATE NOWAIT
func RunLocks() {
	pool, err := InitDBPool()
	if err != nil {
		log.Fatalf("Initiate pool connection failed <3")
	}

	// goroutine A: use FOR UPDATE
	go func() {
		// get conn from pool
		conn, err := pool.Acquire(context.Background())
		if err != nil {
			log.Fatalf("Session A: acquire conn: %v", err)
		}
		defer conn.Release()

		// start tx
		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		if err != nil {
			log.Fatalf("Session A: begin tx: %v", err)
		}
		defer safeRollBack("A", tx)

		fmt.Println("[Session A] Locking row id=1 with FOR UPDATE...")

		sql_cmd := "SELECT * FROM accounts WHERE id=1 FOR UPDATE"
		_, err = tx.Exec(context.Background(), sql_cmd)
		if err != nil {
			log.Fatalf("Session A: query: %v", err)
		}

		fmt.Println("[Session A] Row locked. Sleeping 10s...")
		time.Sleep(10 * time.Second)

		if err := tx.Commit(context.Background()); err != nil {
			log.Fatalf("Session A: commit: %v", err)
		}
		fmt.Println("[Session A] Commit done")
	}()

	time.Sleep(2 * time.Second)

	// goroutine B: use FOR UPDATE NOWAIT
	// go func() {
	// 	conn, err := pool.Acquire(context.Background())
	// 	if err != nil {
	// 		log.Fatalf("Session B: acquire conn: %v", err)
	// 	}
	// 	defer conn.Release()

	// 	tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
	// 	if err != nil {
	// 		log.Fatalf("Session B: begin tx: %v", err)
	// 	}
	// 	defer safeRollBack("B", tx)

	// 	fmt.Println("[Session B] Trying to lock row id=1 with FOR UPDATE NOWAIT...")
	// 	sql_cmd := "SELECT * FROM accounts WHERE id=1 FOR UPDATE NOWAIT"

	// 	_, err = tx.Exec(context.Background(), sql_cmd)
	// 	if err != nil {
	// 		log.Fatalf("Session B: query: %v", err)
	// 		return
	// 	}

	// 	if err := tx.Commit(context.Background()); err != nil {
	// 		log.Fatalf("Session B: commit: %v", err)
	// 	}
	// 	fmt.Println("[Session B] Commit done")
	// }()

	// goroutine C: use FOR UPDATE
	go func() {
		conn, err := pool.Acquire(context.Background())
		if err != nil {
			log.Fatalf("Session C: acquire conn: %v", err)
		}
		defer conn.Release()

		tx, err := conn.BeginTx(context.Background(), pgx.TxOptions{})
		if err != nil {
			log.Fatalf("Session C: begin tx: %v", err)
		}
		defer safeRollBack("C", tx)

		fmt.Println("[Session C] Locking row id=1 with FOR UPDATE...")

		sql_cmd := "SELECT * from accounts where id=1 FOR UPDATE"
		_, err = tx.Exec(context.Background(), sql_cmd)
		if err != nil {
			log.Fatalf("Session C: query: %v", err)
		}

		fmt.Println("[Session C] Row locked. Sleeping 10s...")
		time.Sleep(10 * time.Second)

		if err := tx.Commit(context.Background()); err != nil {
			log.Fatalf("Session C: commit: %v", err)
		}
		fmt.Println("[Session C] Commit done")
	}()

	time.Sleep(15 * time.Second)
}

func safeRollBack(session string, tx pgx.Tx) {
	if err := tx.Rollback(context.Background()); err != nil && err != pgx.ErrTxClosed {
		fmt.Printf("[Session %s] rollback err: %w\n", session, err)
	}
}
