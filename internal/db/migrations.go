package db

import "log"

// CreateSchema creates initial schema if no table created before.
func (db *SQLDB) CreateSchema() error {
	if db.instance == nil {
		err := db.Connect()
		if err != nil {
			log.Println("not able to connect")
			return err
		}
	}

	_, err := db.instance.Exec(`
		CREATE TABLE IF NOT EXISTS links (
			id BIGSERIAL NOT NULL,
			short_id VARCHAR(255) NOT NULL,
			original_url VARCHAR(255) NOT NULL,
			PRIMARY KEY (id)
		)
	`)

	if err != nil {
		log.Println("not able to create `links` table")
		return err
	}

	_, err = db.instance.Exec(`
		CREATE TABLE IF NOT EXISTS user_links (
			id BIGSERIAL NOT NULL,
			user_id VARCHAR(64) NOT NULL,
			link_id INT,
			PRIMARY KEY (id),
			CONSTRAINT fk_link
			  	FOREIGN KEY(link_id)
					REFERENCES links(id)
		)
	`)

	if err != nil {
		log.Println("not able to create `user_links` table")
		return err
	}

	_, err = db.instance.Exec(`
		ALTER TABLE "links" DROP CONSTRAINT IF EXISTS "unique_original_url"
	`)

	if err != nil {
		log.Println("not able to drop constraint created before")
		return err
	}

	_, err = db.instance.Exec(`
		ALTER TABLE "links" ADD CONSTRAINT "unique_original_url" UNIQUE (original_url)
	`)

	if err != nil {
		log.Println("not able to create constraint")
		return err
	}

	_, err = db.instance.Exec(`
		ALTER TABLE "links" ADD COLUMN IF NOT EXISTS "is_deleted" BOOLEAN NULL DEFAULT FALSE
	`)

	log.Println("app schema was successfully restored")
	return nil
}
