package mydatabase

import "database/sql"

func storeFileMetadata(db *sql.DB, id, Filename, key string, deleteAt string) error {
	query := `INSERT INTO files (id,filename,encryption_key,delete_at) VALUES($1,$2,$3,$4)`
	_, err := db.Exec(query, id, Filename, key, deleteAt)
	return err
}
