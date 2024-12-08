package main

import (
	"context"
	"fmt"
	"strings"
	"time"

	"github.com/bxcodec/faker/v4"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
)

func generateFakeData(batchSize int) ([][]interface{}, error) {
	data := make([][]interface{}, batchSize)
	for i := 0; i < batchSize; i++ {
		username := strings.ReplaceAll(faker.Username(), ",", "")
		email := strings.ReplaceAll(faker.Email(), ",", "")
		data[i] = []interface{}{username, email}
	}
	return data, nil
}

func insertTestData(ctx context.Context, pool *pgxpool.Pool) error {
	const totalRecords = 100000
	start := time.Now()
	batchSize := 5000

	// Clean existing data
	_, err := pool.Exec(ctx, `TRUNCATE TABLE "Account", "Character", "Scores" CASCADE;`)
	if err != nil {
		return fmt.Errorf("failed to truncate tables: %v", err)
	}
	fmt.Println("Cleared existing data")

	// Get a dedicated connection for COPY
	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %v", err)
	}
	defer conn.Release()

	// Insert accounts in batches
	fmt.Println("Inserting accounts...")
	for processed := 0; processed < totalRecords; processed += batchSize {
		currentBatch := batchSize
		if processed+batchSize > totalRecords {
			currentBatch = totalRecords - processed
		}

		rows, err := generateFakeData(currentBatch)
		if err != nil {
			return fmt.Errorf("failed to generate fake data: %v", err)
		}

		copyCount, err := conn.CopyFrom(
			ctx,
			pgx.Identifier{"Account"},
			[]string{"username", "email"},
			pgx.CopyFromRows(rows),
		)
		if err != nil {
			return fmt.Errorf("failed to copy data: %v", err)
		}

		fmt.Printf("Inserted %d accounts... (%d%%)\n", 
			processed+int(copyCount), 
			(processed+int(copyCount))*100/totalRecords)
	}

	// Generate characters and scores
	fmt.Println("Generating characters and scores...")
	_, err = pool.Exec(ctx, `
		INSERT INTO "Character" (acc_id, class_id)
		SELECT 
			acc_id,
			floor(random() * 8 + 1)::int
		FROM "Account",
		generate_series(1, floor(random() * 8 + 1)::int);

		INSERT INTO "Scores" (char_id, reward_score)
		SELECT 
			char_id,
			floor(random() * 99001 + 1000)::int
		FROM "Character",
		generate_series(1, floor(random() * 5 + 1)::int);
	`)
	if err != nil {
		return fmt.Errorf("failed to generate characters and scores: %v", err)
	}

	// Get final counts
	var counts struct {
		accounts, characters, scores int
	}

	err = pool.QueryRow(ctx, `
		SELECT 
			(SELECT COUNT(*) FROM "Account"),
			(SELECT COUNT(*) FROM "Character"),
			(SELECT COUNT(*) FROM "Scores")
	`).Scan(&counts.accounts, &counts.characters, &counts.scores)
	if err != nil {
		return fmt.Errorf("failed to get counts: %v", err)
	}

	duration := time.Since(start)
	fmt.Printf("\nInsertion completed in %v\n", duration)
	fmt.Printf("Final counts:\n")
	fmt.Printf("- Accounts: %d\n", counts.accounts)
	fmt.Printf("- Characters: %d\n", counts.characters)
	fmt.Printf("- Scores: %d\n", counts.scores)

	return nil
}
