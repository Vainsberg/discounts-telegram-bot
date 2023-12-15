package bottg

import (
	"database/sql"
	"log"
)

var db *sql.DB

func AddLincked(chatID string, text string) error {
	_, err := db.Exec(`
    INSERT INTO linked_accounts (name, query);
`, chatID, text)
	if err != nil {
		log.Fatal(err)
	}
	return nil

}
