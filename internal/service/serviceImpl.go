package service

import (
	"context"

	"github.com/Antoha2/news/internal/repository"
	"github.com/pkg/errors"
)

//get userS
func (s *ServImpl) GetNews(ctx context.Context, pNews *SearchTerms) ([]*News, error) {

	rNews := &repository.SearchTerms{
		Limit:  pNews.Limit,
		Offset: pNews.Offset,
	}

	repNews, err := s.rep.GetNews(ctx, rNews)
	if err != nil {
		return nil, errors.Wrap(err, "occurred error GetNews")
	}

	newsS := make([]*News, len(repNews))
	for index, news := range repNews {
		t := &News{
			Id:         news.Id,
			Title:      news.Title,
			Content:    news.Content,
			Categories: news.Categories,
		}
		newsS[index] = t
	}
	return newsS, nil
}

//add user
func (s *ServImpl) AddNews(ctx context.Context, news *News) (*News, error) {

	repNews := &repository.RepNews{
		Title:      news.Title,
		Content:    news.Content,
		Categories: news.Categories,
	}
	repNews, err := s.rep.AddNews(ctx, repNews)
	if err != nil {
		return nil, errors.Wrap(err, "occurred error AddNews")
	}

	respUser := &News{
		Id:         repNews.Id,
		Title:      repNews.Title,
		Content:    repNews.Content,
		Categories: repNews.Categories,
	}
	return respUser, nil
}

//edit News
func (s *ServImpl) EditNews(ctx context.Context, id int, news *News) (*News, error) {

	reposNews := &repository.RepNews{
		Id:         news.Id,
		Title:      news.Title,
		Content:    news.Content,
		Categories: news.Categories,
	}

	reposNews, err := s.rep.EditNews(ctx, id, reposNews)
	if err != nil {
		return nil, errors.Wrap(err, "occurred error edit News")
	}
	respNews := &News{
		Id:         reposNews.Id,
		Title:      reposNews.Title,
		Content:    reposNews.Content,
		Categories: reposNews.Categories,
	}

	return respNews, nil
}
