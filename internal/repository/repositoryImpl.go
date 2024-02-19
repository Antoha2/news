package repository

import (
	"context"
)

//add user
func (r *RepImpl) AddNews(ctx context.Context, news *RepNews) (*RepNews, error) {

	//r.DB.Save(news)
	// dbNews := &DBNews{
	// 	Title:   news.Title,
	// 	Content: news.Content,
	// }
	// if err := r.DB.Save(dbNews); err != nil {
	// 	log.Fatal(err)
	// }
	//log.Println("rep !!!!!!!!!!!!!!!!!!!!!!!! ", news)
	//repNews := RepNews{}

	// query := "INSERT INTO news (title, content) VALUES ($1, $2) RETURNING id, title, content"
	// row := r.DB.QueryRowContext(ctx, query, news.Title, news.Content)
	// if err := row.Scan(&repNews.Id, &repNews.Title, &repNews.Content, &repNews.Categories); err != nil {
	// 	return nil, errors.Wrap(err, fmt.Sprintf("sql add news failed, query: %s", query))
	// }

	// query = "INSERT INTO NewsCategories (newsId, CategoryId) VALUES ($1, $2)"
	// for i := 0; i < len(news.Categories); i++ {
	// 	row,err := r.DB.ExecContext(ctx, query, news.Id, news.Categories[i])
	// }

	return news, nil
}

//edit News
func (r *RepImpl) EditNews(ctx context.Context, id int, news *RepNews) (*RepNews, error) {
	return news, nil
}

//get News
func (r *RepImpl) GetNews(ctx context.Context) ([]*RepNews, error) {
	return nil, nil
}
