package repository

import (
	"context"
	"fmt"
	"strings"

	"github.com/pkg/errors"
)

//add user
func (r *RepImpl) AddNews(ctx context.Context, news *RepNews) (*RepNews, error) {

	repNews := &RepNews{}
	repNews.Categories = news.Categories

	query := "INSERT INTO news (title, content) VALUES ($1, $2) RETURNING id, title, content"
	row := r.DB.QueryRowContext(ctx, query, news.Title, news.Content)
	if err := row.Scan(&repNews.Id, &repNews.Title, &repNews.Content); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("sql add News failed, query: %s", query))
	}

	queryConstrain, args := buildAddQueryConstrain(repNews)

	_, err := r.DB.ExecContext(ctx, queryConstrain, args...)
	if err != nil {
		return news, errors.Wrap(err, fmt.Sprintf("sql add News failed, query: %s", query))
	}
	repNews.Categories = news.Categories
	return repNews, nil
}

//edit News
func (r *RepImpl) EditNews(ctx context.Context, id int, news *RepNews) (*RepNews, error) {

	repNews := &RepNews{}

	queryConstrain, args := buildEditQueryConstrain(news, id)

	query := fmt.Sprintf("UPDATE news SET%s RETURNING id, title, content", queryConstrain)
	row := r.DB.QueryRowContext(ctx, query, args...)
	if err := row.Scan(&repNews.Id, &repNews.Title, &repNews.Content); err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("sql edit News failed, query: %s", query))
	}
	return repNews, nil

}

//get News
func (r *RepImpl) GetNews(ctx context.Context, pNews *SearchTerms) ([]*RepNews, error) {

	rNews := make([]*RepNews, 0)

	query := "SELECT news.id, news.title, news.content, news_categories.categories_id FROM news, news_categories WHERE news_categories.news_id = news.id ORDER BY id asc LIMIT $1 OFFSET $2"

	rows, err := r.DB.QueryContext(ctx, query, pNews.Limit, pNews.Offset)
	if err != nil {
		return nil, errors.Wrap(err, fmt.Sprintf("sql select Users failed, query: %s", query))
	}

	for rows.Next() {
		news := &RepNews{}
		err := rows.Scan(&news.Id, &news.Title, &news.Content, &news.Categories)
		if err != nil {
			return nil, errors.Wrap(err, "sql scan Users failed")
		}
		rNews = append(rNews, news)

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

	queryConstrain := fmt.Sprintf("INSERT INTO news_categories(news_Id, categories_id) VALUES %s", strings.Join(constrains, ","))
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
