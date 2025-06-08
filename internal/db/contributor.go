package db

// CheckContributorExists checks if a contributor with the given ID exists in the database
func CheckContributorExists(id int) (bool, error) {
	var exists bool
	query := `SELECT EXISTS(SELECT 1 FROM contributors WHERE id = $1)`
	err := db.QueryRow(query, id).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}
