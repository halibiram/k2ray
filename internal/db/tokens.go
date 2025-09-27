package db

import (
	"database/sql"
	"time"
)

// BlocklistToken adds a token's JTI to the revoked_tokens table.
func BlocklistToken(jti string, expiresAt time.Time) error {
	insertSQL := `INSERT INTO revoked_tokens (jti, expires_at) VALUES (?, ?)`
	_, err := DB.Exec(insertSQL, jti, expiresAt.Unix())
	return err
}

// IsTokenBlocklisted checks if a token's JTI exists in the revoked_tokens table.
func IsTokenBlocklisted(jti string) (bool, error) {
	var foundJti string
	querySQL := `SELECT jti FROM revoked_tokens WHERE jti = ?`
	err := DB.QueryRow(querySQL, jti).Scan(&foundJti)
	if err != nil {
		if err == sql.ErrNoRows {
			// The token JTI was not found in the blocklist, so it's not blocklisted.
			return false, nil
		}
		// Some other database error occurred.
		return false, err
	}
	// A row was found, meaning the token is blocklisted.
	return true, nil
}

// CleanupExpiredTokens removes tokens from the blocklist that have expired.
func CleanupExpiredTokens() (int64, error) {
	deleteSQL := `DELETE FROM revoked_tokens WHERE expires_at < ?`
	now := time.Now().Unix()
	result, err := DB.Exec(deleteSQL, now)
	if err != nil {
		return 0, err
	}
	return result.RowsAffected()
}
