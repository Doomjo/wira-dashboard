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

func generateUniqueUsernames(batchSize int) ([][]interface{}, error) {
	data := make([][]interface{}, batchSize)
	usernameMap := make(map[string]bool) // To track used usernames

	for i := 0; i < batchSize; i++ {
		var username string
		// Keep generating until we get a unique username
		for {
			// Get random number and handle error
			randNum, err := faker.RandomInt(1000, 9999)
			if err != nil {
				return nil, fmt.Errorf("failed to generate random number: %v", err)
			}

			// Combine random words with the number to make unique usernames
			username = fmt.Sprintf("%s%s%d",
				strings.ToLower(faker.Word()),
				strings.ToLower(faker.Word()),
				randNum[0]) // Use first number from the slice
			username = strings.ReplaceAll(username, ",", "")

			if !usernameMap[username] {
				usernameMap[username] = true
				break
			}
		}

		email := fmt.Sprintf("%s@example.com", username)
		data[i] = []interface{}{username, email}
	}
	return data, nil
}

func insertTestData(ctx context.Context, pool *pgxpool.Pool) error {
	const totalRecords = 100000
	start := time.Now()
	batchSize := 5000

	// Clean existing data
	_, err := pool.Exec(ctx, `TRUNCATE TABLE account, character, scores CASCADE;`)
	if err != nil {
		return fmt.Errorf("failed to truncate tables: %v", err)
	}
	fmt.Println("Cleared existing data")

	conn, err := pool.Acquire(ctx)
	if err != nil {
		return fmt.Errorf("failed to acquire connection: %v", err)
	}
	defer conn.Release()

	// Insert accounts in batches with unique usernames
	fmt.Println("Inserting accounts...")
	for processed := 0; processed < totalRecords; processed += batchSize {
		currentBatch := batchSize
		if processed+batchSize > totalRecords {
			currentBatch = totalRecords - processed
		}

		rows, err := generateUniqueUsernames(currentBatch)
		if err != nil {
			return fmt.Errorf("failed to generate unique usernames: %v", err)
		}

		copyCount, err := conn.CopyFrom(
			ctx,
			pgx.Identifier{"account"},
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
        -- Create temporary sequence for generating unique class assignments
        WITH RECURSIVE class_assignments AS (
            -- For each account, generate 1-8 rows with unique class_ids
            SELECT 
                a.acc_id,
                c.class_id
            FROM 
                account a
                CROSS JOIN LATERAL (
                    SELECT DISTINCT floor(random() * 8 + 1)::int as class_id
                    FROM generate_series(1, floor(random() * 8 + 1)::int)
                    WHERE floor(random() * 8 + 1)::int <= 8
                ) c
            GROUP BY a.acc_id, c.class_id  -- This ensures uniqueness
        )
        INSERT INTO character (acc_id, class_id)
        SELECT 
            acc_id,
            class_id
        FROM class_assignments;

        -- Generate scores as before
        INSERT INTO scores (char_id, reward_score)
        SELECT 
            char_id,
            floor(random() * 99001 + 1000)::int
        FROM character;
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
            (SELECT COUNT(*) FROM account),
            (SELECT COUNT(*) FROM character),
            (SELECT COUNT(*) FROM scores)
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
