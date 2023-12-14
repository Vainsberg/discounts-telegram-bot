package bottg

import (
	"database/sql"
	"log"
)

var db *sql.DB

func AddLincked(user string, text string) error {
	_, err := db.Exec(`
    INSERT INTO linked_accounts (name, query, goods_id)
    VALUES (?, ?, (SELECT id FROM goods WHERE query = ?));
`, text, user, text)
	if err != nil {
		log.Fatal(err)
	}
	return nil

}
