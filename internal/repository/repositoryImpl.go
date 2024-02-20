package repository

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/Antoha2/news/internal/lib/models"
	"github.com/pkg/errors"
)

//add user
func (r *RepImpl) AddNews(ctx context.Context, news *RepNews) (*RepNews, error) {

	repNews := &RepNews{}

	query := "INSERT INTO news (title, content) VALUES ($1, $2) RETURNING id, title, content"
	row := r.DB.QueryRowContext(ctx, query, news.Title, news.Content)
	if err := row.Scan(&repNews.Id, &repNews.Title, &repNews.Content); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("sql add News failed, query: %s", query))
	}

	queryConstrain, args := buildAddQueryConstrain(repNews)

	query = fmt.Sprintf("INSERT INTO news_categories(news_Id, categories_id) VALUES%s", queryConstrain)
	_, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("sql add News failed, query: %s", query))
	}
	repNews.Categories = news.Categories
	return repNews, nil
}

//edit News
func (r *RepImpl) EditNews(ctx context.Context, id int, news *RepNews) (*RepNews, error) {

	count := 0

	if err := r.DB.QueryRowContext(ctx, "SELECT COUNT(id) FROM news WHERE id = $1", id).Scan(&count); err != nil {
		return nil, errors.Wrap(err, "sql edit News failed ")
	}

	if count == 0 {
		return nil, errors.New("sql edit News failed, no such ID exists")
	}

	tx, err := r.DB.Begin()

	if err != nil {
		return nil, errors.Wrap(err, "sql edit News failed ")
	}

	query := "DELETE FROM news_categories WHERE news_id = $1"
	_, err = tx.ExecContext(ctx, query, id)
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, fmt.Sprintf("sql edit News failed, query: %s", query))
	}

	repNews := &RepNews{}
	if news.Id != 0 || news.Title != "" || news.Content != "" {
		queryConstrain, args := buildEditQueryConstrain(news, id)
		query = fmt.Sprintf("UPDATE news SET%s RETURNING id, title, content", queryConstrain)
		row := tx.QueryRowContext(ctx, query, args...)
		if err := row.Scan(&repNews.Id, &repNews.Title, &repNews.Content); err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("sql edit News failed, query: %s", query))
		}
	}
	if news.Id == 0 {
		news.Id = id
	}

	if len(news.Categories) != 0 {
		queryConstrain, args := buildAddQueryConstrain(news)
		query = fmt.Sprintf("INSERT INTO news_categories(news_Id, categories_id) VALUES%s", queryConstrain)
		_, err = r.DB.ExecContext(ctx, query, args...)
		if err != nil {
			tx.Rollback()
			return nil, errors.Wrap(err, fmt.Sprintf("sql edit News failed, query: %s", query))
		}
	}
	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "sql edit News failed")
	}
	return repNews, nil

}

//get News
func (r *RepImpl) GetNews(ctx context.Context, pNews *SearchTerms) ([]*RepNews, error) {

	rNews := make([]*RepNews, 0)

	query := "SELECT id, title, content FROM news ORDER BY id asc LIMIT $1 OFFSET $2"

	rows, err := r.DB.QueryContext(ctx, query, pNews.Limit, pNews.Offset)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("sql get News failed, query: %s", query))
	}

	for rows.Next() {
		news := &RepNews{}
		err := rows.Scan(&news.Id, &news.Title, &news.Content)
		if err != nil {
			return nil, errors.Wrap(err, "sql get News failed")
		}
		rNews = append(rNews, news)
	}

	for i := 0; i < len(rNews); i++ {
		query = "SELECT categories_id FROM news_categories WHERE news_id=$1"
		rows, err := r.DB.QueryContext(ctx, query, rNews[i].Id)
		if err != nil {
			return nil, errors.Wrap(err, fmt.Sprintf("sql get News failed, query: %s", query))
		}
		categories := make([]int, 0)
		for rows.Next() {
			category := 0

			err := rows.Scan(&category)
			if err != nil {
				return nil, errors.Wrap(err, "sql get News failed")
			}
			categories = append(categories, category)
		}
		rNews[i].Categories = categories
	}

	return rNews, nil
}

//build add query string
func buildAddQueryConstrain(repNews *RepNews) (string, []any) {

	constrains := make([]string, 0, len(repNews.Categories))
	args := make([]any, 0, len(repNews.Categories))

	y := 1
	for i := 0; i < len(repNews.Categories); i++ {
		constrains = append(constrains, fmt.Sprintf("($%d, $%d)", y, y+1))

		args = append(args, repNews.Id)
		args = append(args, repNews.Categories[i])

		y += 2
	}

	queryConstrain := fmt.Sprintf(" %s", strings.Join(constrains, ","))
	return queryConstrain, args
}

//build edit query string
func buildEditQueryConstrain(repNews *RepNews, id int) (string, []any) {

	i := 1
	constrains := make([]string, 0, 3)
	args := make([]any, 0, 3)

	if repNews.Id != 0 {
		s := fmt.Sprintf("id=$%d", i)
		i++

		constrains = append(constrains, s)
		args = append(args, repNews.Id)
	}
	if repNews.Title != "" {
		s := fmt.Sprintf("title=$%d", i)
		i++

		constrains = append(constrains, s)
		args = append(args, repNews.Title)
	}
	if repNews.Content != "" {
		s := fmt.Sprintf("content=$%d", i)
		i++

		constrains = append(constrains, s)
		args = append(args, repNews.Content)
	}

	queryConstrain := fmt.Sprintf(" %s WHERE id=$%d", strings.Join(constrains, ", "), i)
	args = append(args, id)

	return queryConstrain, args
}

//------------------------------------------------------------

// SaveUser saves user to db.
func (r *RepImpl) UserSaver(ctx context.Context, email string, passHash []byte) (int64, error) {
	const op = "repository.postgres.SaveUser"

	//checking user uniqueness
	var count int
	stmt := r.DB.QueryRowContext(ctx, "SELECT COUNT(id) FROM users WHERE email = $1", email)

	if err := stmt.Scan(&count); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	if count > 0 {
		return 0, fmt.Errorf("%s: %w", op, ErrUserExists)
	}

	//adding a user to the database
	var id int64
	stmt = r.DB.QueryRowContext(ctx, "INSERT INTO users (email, pass_hash) VALUES ($1, $2) RETURNING id", email, passHash)
	if err := stmt.Scan(&id); err != nil {
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

// User returns user by email.
func (r *RepImpl) UserProvider(ctx context.Context, email string) (models.User, error) {
	const op = "repository.postgres.User"

	user := new(models.User)
	stmt := r.DB.QueryRowContext(ctx, "SELECT id, email, pass_hash FROM users WHERE email = $1", email)

	err := stmt.Scan(&user.ID, &user.Email, &user.PassHash)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.User{}, fmt.Errorf("%s: %w", op, ErrUserNotFound)
		}

		return models.User{}, fmt.Errorf("%s: %w", op, err)
	}

	return *user, nil
}
