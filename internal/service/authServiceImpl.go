package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	"github.com/Antoha2/news/internal/repository"
	"github.com/Antoha2/news/pkg/jwt"
	"github.com/Antoha2/news/pkg/logger/sl"
	"golang.org/x/crypto/bcrypt"
)

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

// func (s *authServImpl) ParseToken(accesToken string) (int, error) {

// 	token, err :=  jwt.ParseWithClaims(accesToken, &tokenClaims{}, func(token *jwt.Token) (interface{}, error) {
// 		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
// 			return nil, errors.New("неверный метод подписи")
// 		}
// 		return []byte(signingKey), nil
// 	})
// 	if err != nil {
// 		return 0, nil
// 	}

// 	claims, ok := token.Claims.(*tokenClaims)
// 	if !ok {
// 		return 0, errors.New("token claims not of type *tokenClaims")
// 	}
// 	return claims.UserId, nil
// }
