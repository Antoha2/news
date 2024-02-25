package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

//add user
func (r *RepImpl) AddNews(ctx context.Context, news *RepNews) (*RepNews, error) {

	rNews := &RepNews{}

	query := "INSERT INTO news (title, content) VALUES ($1, $2) RETURNING id, title, content"
	row := r.DB.QueryRowContext(ctx, query, news.Title, news.Content)
	if err := row.Scan(&rNews.Id, &rNews.Title, &rNews.Content); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("sql add News failed, query: %s", query))
	}

	rNews.Categories = news.Categories
	queryConstrain, args := buildAddQueryConstrain(rNews)

	query = fmt.Sprintf("INSERT INTO news_categories(news_Id, categories_id) VALUES%s", queryConstrain)
	_, err := r.DB.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("sql add News failed, query: %s", query))
	}
	rNews.Categories = news.Categories
	return rNews, nil
}

//edit News
func (r *RepImpl) EditNews(ctx context.Context, news *RepNews) (*RepNews, error) {

	count := 0

	if err := r.DB.QueryRowContext(ctx, "SELECT COUNT(id) FROM news WHERE id = $1", news.Id).Scan(&count); err != nil {
		return nil, errors.New("sql edit News failed")
	}
	if count == 0 {
		return nil, errors.New("sql edit News failed, no such ID exists")
	}

	tx, err := r.DB.Begin()
	if err != nil {
		return nil, errors.Wrap(err, "sql edit News failed ")
	}

	rNews := &RepNews{
		Id:         news.Id,
		Categories: news.Categories,
	}

	if news.Title != "" || news.Content != "" {
		queryConstrain, args := buildEditQueryConstrain(news)
		query := fmt.Sprintf("UPDATE news SET%s WHERE id=%d RETURNING id, title, content", queryConstrain, news.Id)
		row := tx.QueryRowContext(ctx, query, args...)
		if err := row.Scan(&rNews.Id, &rNews.Title, &rNews.Content); err != nil {
			tx.Rollback()
			return nil, errors.Wrap(err, fmt.Sprintf("sql edit News failed, query: %s", query))
		}
		// } else {
		// 	query := "SELECT title, content FROM news WHERE id = $1"
		// 	row := tx.QueryRowContext(ctx, query, news.Id)
		// 	if err := row.Scan(&rNews.Title, &rNews.Content); err != nil {
		// 		tx.Rollback()
		// 		return nil, errors.Wrap(err, fmt.Sprintf("sql edit News failed, query: %s", query))

		// 	}
	}

	if len(rNews.Categories) != 0 {

		query := "DELETE FROM news_categories WHERE news_id = $1"
		_, err = tx.ExecContext(ctx, query, news.Id)
		if err != nil {
			tx.Rollback()
			return nil, errors.Wrap(err, fmt.Sprintf("sql edit News failed, query: %s", query))
		}

		queryConstrain, args := buildAddQueryConstrain(rNews)
		query = fmt.Sprintf("INSERT INTO news_categories(news_Id, categories_id) VALUES%s", queryConstrain)
		_, err = r.DB.ExecContext(ctx, query, args...)
		if err != nil {
			tx.Rollback()
			return nil, errors.Wrap(err, fmt.Sprintf("sql edit News failed, query: %s", query))
		}
	} else {

		query := "SELECT categories_id FROM news_categories WHERE news_id=$1"
		rows, err := r.DB.QueryContext(ctx, query, rNews.Id)
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
		rNews.Categories = categories
	}

	err = tx.Commit()
	if err != nil {
		tx.Rollback()
		return nil, errors.Wrap(err, "sql edit News failed")
	}

	return rNews, nil

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
func buildAddQueryConstrain(rNews *RepNews) (string, []any) {

	constrains := make([]string, 0, len(rNews.Categories))
	args := make([]any, 0, len(rNews.Categories))

	y := 1
	for i := 0; i < len(rNews.Categories); i++ {
		constrains = append(constrains, fmt.Sprintf("($%d, $%d)", y, y+1))

		args = append(args, rNews.Id)
		args = append(args, rNews.Categories[i])

		y += 2
	}

	queryConstrain := fmt.Sprintf(" %s", strings.Join(constrains, ","))
	return queryConstrain, args
}

//build edit query string
func buildEditQueryConstrain(rNews *RepNews) (string, []any) {

	i := 1
	constrains := make([]string, 0, 3)
	args := make([]any, 0, 3)

	if rNews.Title != "" {
		s := fmt.Sprintf("title=$%d", i)
		i++

		constrains = append(constrains, s)
		args = append(args, rNews.Title)
	}
	if rNews.Content != "" {
		s := fmt.Sprintf("content=$%d", i)
		i++

		constrains = append(constrains, s)
		args = append(args, rNews.Content)
	}

	queryConstrain := fmt.Sprintf(" %s", strings.Join(constrains, ", "))

	return queryConstrain, args
}
