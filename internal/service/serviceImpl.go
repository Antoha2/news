package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/Antoha2/news/internal/repository"
	"github.com/Antoha2/news/pkg/jwt"
	"github.com/Antoha2/news/pkg/logger/sl"
	"github.com/pkg/errors"
	"golang.org/x/crypto/bcrypt"
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
		Categories: validationCategories(news.Categories),
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
		Categories: validationCategories(news.Categories),
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

//remove duplicate values Categories
func validationCategories(n []int) []int {

	m := map[int]int{}
	s := make([]int, 0)

	for i := 0; i < len(n); i++ {
		m[n[i]]++
		if m[n[i]] < 2 {
			s = append(s, n[i])
		}
	}
	return s
}

//-------------------------------------------------------

func (s *authServImpl) RegisterNewUser(ctx context.Context, email string, pass string) (userID int64, err error) {
	const op = "serv.RegisterNewUser"
	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	log.Info("registring user")
	// Генерируем хэш и соль для пароля.
	passHash, err := bcrypt.GenerateFromPassword([]byte(pass), bcrypt.DefaultCost)
	if err != nil {
		log.Error("failed to generate password hash", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}
	// Сохраняем пользователя в БД
	id, err := s.rep.UserSaver(ctx, email, passHash)
	if err != nil {
		log.Error("failed to save user", sl.Err(err))
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	return id, nil
}

func (s *authServImpl) Login(ctx context.Context, email string, password string, appID int) (string, error) {
	const op = "serv.Login"

	log := s.log.With(
		slog.String("op", op),
		slog.String("email", email),
	)
	log.Info("attempting to login user")

	// Достаём пользователя из БД
	user, err := s.rep.UserProvider(ctx, email)
	if err != nil {
		if errors.Is(err, repository.ErrUserNotFound) {
			s.log.Warn("user not found", sl.Err(err))
			return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
		}
		s.log.Error("failed to get user", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}

	// Проверяем корректность полученного пароля
	if err := bcrypt.CompareHashAndPassword(user.PassHash, []byte(password)); err != nil {
		s.log.Info("invalid credentials", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, ErrInvalidCredentials)
	}

	log.Info("user logged in successfully")

	// Создаём токен авторизации
	token, err := jwt.NewToken(s.cfg, user)
	if err != nil {
		s.log.Error("failed to generate token", sl.Err(err))
		return "", fmt.Errorf("%s: %w", op, err)
	}
	return token, nil
}
