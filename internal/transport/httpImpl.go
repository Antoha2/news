package transport

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"
	"strings"

	"github.com/Antoha2/news/internal/service"
	"github.com/Antoha2/news/pkg/logger/sl"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/pkg/errors"
)

func (a *apiImpl) StartHTTP() error {
	app := fiber.New()

	app.Post("/auth/login", a.loginHandler)
	app.Post("/auth/register", a.registerHandler)

	app.Use("/api", adaptor.HTTPMiddleware(a.userIdentify))

	app.Get("/api/list", a.getNewsHandler)
	app.Post("/api/add", a.addNewsHandler)
	app.Post("/api/:id", a.editNewsHandler)

	err := app.Listen(fmt.Sprintf(":%s", a.cfg.HTTP.HostPort))
	if err != nil {
		return errors.Wrap(err, "ocurred error StartHTTP")
	}
	return nil
}

func (a *apiImpl) Stop() {
	if err := a.server.Shutdown(context.TODO()); err != nil {
		panic(errors.Wrap(err, "ocurred error Stop"))
	}
}

//get News
func (a *apiImpl) getNewsHandler(c *fiber.Ctx) error {

	const op = "getNews"
	log := a.log.With(slog.String("op", op))

	pNews := &service.SearchTerms{
		Limit:  a.cfg.DBConfig.DefaultPropertyLimit,
		Offset: a.cfg.DBConfig.DefaultPropertyOffset,
	}

	if err := c.BodyParser(&pNews); err != nil {
		log.Error("cant unmarshall", sl.Err(err))
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"error":   err,
		})
		return err
	}

	log.Info("run get News")

	news, err := a.service.GetNews(c.Context(), pNews)
	if err != nil {
		a.log.Error("occurred error for GetNews", sl.Err(err))
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"error":   err,
		})
		return err
	}
	c.JSON(&fiber.Map{
		"success": true,
		"news":    news,
	})
	return nil
}

//edit News
func (a *apiImpl) editNewsHandler(c *fiber.Ctx) error {

	const op = "edit News"
	log := a.log.With(slog.String("op", op))

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Error("wrong format id", sl.Err(err))
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"error":   err,
		})
		return err
	}

	news := &service.News{}
	if err := c.BodyParser(&news); err != nil {
		log.Error("cant unmarshall", sl.Err(err))
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"error":   err,
		})
		return err
	}

	if err = requestValidation(news); err != nil {
		log.Error("occurred error for edit News, ", sl.Err(err))
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"error":   err,
		})
		return err
	}

	log.Info("run edit News", sl.Atr("News", news))

	news, err = a.service.EditNews(c.Context(), id, news)
	if err != nil {
		log.Error("occurred error for edit News, ", sl.Err(err))
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"error":   err,
		})
		return err
	}
	c.JSON(&fiber.Map{
		"success": true,
		"news":    news,
	})
	return nil
}

//add News
func (a *apiImpl) addNewsHandler(c *fiber.Ctx) error {
	const op = "add News"
	log := a.log.With(slog.String("op", op))

	log.Info("run add News")

	news := &service.News{}
	if err := c.BodyParser(&news); err != nil {
		log.Error("cant unmarshall", sl.Err(err))
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"error":   err,
		})
		return err
	}

	log.Info("run add News", sl.Atr("News", news))

	rNews, err := a.service.AddNews(c.Context(), news)
	if err != nil {
		a.log.Error("occurred error for addNews", sl.Err(err))
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"error":   err,
		})
		return err
	}

	c.JSON(&fiber.Map{
		"success": true,
		"news_id": rNews.Id,
	})
	return nil
}

//request Validation
func requestValidation(n *service.News) error {

	if n.Id == 0 && n.Title == "" && n.Content == "" && len(n.Categories) == 0 {
		return errors.New("empty request")
	}

	vErr := make([]string, 0, 4)
	if len(n.Title) >= 255 {
		vErr = append(vErr, "long length Title")
	}
	if len(n.Content) > 1000 {
		vErr = append(vErr, "long length Content")
	}
	if len(n.Categories) > 10 {
		vErr = append(vErr, "too mush Categories")
	}

	if len(vErr) != 0 {
		return errors.New(strings.Join(vErr, ", "))
	}

	return nil
}

//-----------------------------------------------------------------------------------

func (a *apiImpl) loginHandler(c *fiber.Ctx) error {
	const op = "login user"
	log := a.log.With(slog.String("op", op))
	log.Info("run login user")

	req := &service.LoginRequest{}

	if err := c.BodyParser(&req); err != nil {
		log.Error("cant unmarshall", sl.Err(err))
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"error":   err,
		})
		return err
	}

	token, err := a.authService.Login(c.Context(), req.Username, req.Password, 0)
	if err != nil {
		a.log.Error("occurred error for Login user", sl.Err(err))
		c.Status(http.StatusInternalServerError).JSON(&fiber.ErrInternalServerError)
	}
	c.JSON(&fiber.Map{
		"success": true,
		"token":   token,
	})
	return nil

}

func (a *apiImpl) registerHandler(c *fiber.Ctx) error {
	const op = "register user"
	log := a.log.With(slog.String("op", op))

	req := &service.RegisterRequest{}

	if err := c.BodyParser(&req); err != nil {
		log.Error("cant unmarshall", sl.Err(err))
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"error":   err,
		})
		return err
	}

	fmt.Println(req)

	log.Info("run register user")
	userId, err := a.authService.RegisterNewUser(c.Context(), req.Email, req.Password)
	if err != nil {
		a.log.Error("occurred error for register user", sl.Err(err))
		c.Status(http.StatusInternalServerError).JSON(&fiber.ErrInternalServerError)
	}
	c.JSON(&fiber.Map{
		"success": true,
		"userId":  userId,
	})
	return nil
}
