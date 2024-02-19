package service

import (
	"context"
	"log"

	"github.com/Antoha2/news/internal/repository"
	"github.com/pkg/errors"
)

//get userS
func (s *ServImpl) GetNews(ctx context.Context) ([]*News, error) {
	return nil, nil
	// readFilter := &repository.RepQueryFilter{
	// 	Id:          filter.Id,
	// 	Name:        filter.Name,
	// 	SurName:     filter.SurName,
	// 	Patronymic:  filter.Patronymic,
	// 	Age:         filter.Age,
	// 	Gender:      filter.Gender,
	// 	Nationality: filter.Nationality,
	// 	Limit:       filter.Limit,
	// 	Offset:      filter.Offset,
	// }
	// repUsers, err := s.rep.GetUsers(ctx, readFilter)
	// if err != nil {
	// 	return nil, errors.Wrap(err, "occurred error GetUsers")
	// }

	// users := make([]*User, len(repUsers))
	// for index, user := range repUsers {
	// 	t := &User{
	// 		Id:          user.Id,
	// 		Name:        user.Name,
	// 		SurName:     user.SurName,
	// 		Patronymic:  user.Patronymic,
	// 		Age:         user.Age,
	// 		Gender:      user.Gender,
	// 		Nationality: user.Nationality,
	// 	}
	// 	users[index] = t
	// }
	// return users, nil
}

//add user
func (s *ServImpl) AddNews(ctx context.Context, news *News) (*News, error) {

	repNews := &repository.RepNews{
		Title:      news.Title,
		Content:    news.Content,
		Categories: news.Categories,
	}
	log.Println("serv !!!!!!!!!!!!!!!!!!!!! ", repNews)
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
