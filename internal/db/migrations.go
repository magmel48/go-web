package db

import "log"

// CreateSchema creates initial schema if no table created before.
func CreateSchema() error {
	if DB == nil {
		err := Connect()
		if err != nil {
			log.Println("not able to connect")
			return err
		}
	}

	_, err := DB.Exec(`
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

	_, err = DB.Exec(`
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

	log.Println("app schema was successfully restored")
	return nil
}
