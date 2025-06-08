package db

import "gutondeka/internal/models"

// GetArticlesCount returns the count of articles grouped by month and year
func GetArticlesCount() ([]models.ArticleCount, error) {
	query := `
		SELECT 
			strftime('%m-%Y', created_at) as month_year,
			COUNT(*) as count
		FROM articles
		GROUP BY month_year
		ORDER BY created_at`

	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var counts []models.ArticleCount
	for rows.Next() {
		var count models.ArticleCount
		err := rows.Scan(&count.MonthYear, &count.Count)
		if err != nil {
			return nil, err
		}
		counts = append(counts, count)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return counts, nil
}

func InsertArticle(article *models.Article) error {
	query := `
	INSERT INTO articles (contributor_id, title, author, source, date, category, file_path)
	VALUES (?, ?, ?, ?, ?, ?, ?)`

	_, err := db.Exec(query,
		article.ContributorID,
		article.Title,
		article.Author,
		article.Source,
		article.Date,
		article.Category,
		article.FilePath)

	return err
}

func Close() error {
	if db != nil {
		return db.Close()
	}
	return nil
}
