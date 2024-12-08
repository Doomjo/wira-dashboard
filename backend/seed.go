package seed

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"

	_ "github.com/lib/pq"
)

var db *sql.DB

func InsertTestData() {
	// Insert into Account
	var accID int
	accountQuery := `INSERT INTO Account (username, email) VALUES ($1, $2) RETURNING acc_id`
	err := db.QueryRow(accountQuery, "testuser", "testuser@example.com").Scan(&accID)
	if err != nil {
		log.Fatalf("Failed to insert into Account: %v", err)
	}
	fmt.Printf("Inserted into Account with acc_id: %d\n", accID)

	// Insert into Character
	for i := 1; i <= 8; i++ {
		var charID int
		characterQuery := `INSERT INTO Character (acc_id, class_id) VALUES ($1, $2) RETURNING char_id`
		err := db.QueryRow(characterQuery, accID, i).Scan(&charID)
		if err != nil {
			log.Printf("Failed to insert into Character for class_id %d: %v", i, err)
			continue
		}
		fmt.Printf("Inserted into Character with char_id: %d for class_id: %d\n", charID, i)

		// Insert into Scores
		scoreQuery := `INSERT INTO Scores (char_id, reward_score) VALUES ($1, $2)`
		_, err = db.Exec(scoreQuery, charID, rand.Intn(1000))
		if err != nil {
			log.Printf("Failed to insert into Scores for char_id %d: %v", charID, err)
		} else {
			fmt.Printf("Inserted score for char_id: %d\n", charID)
		}
	}
}
