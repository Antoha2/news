package transport

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"strconv"

	"github.com/Antoha2/news/internal/service"
	"github.com/Antoha2/news/pkg/logger/sl"
	"github.com/gofiber/fiber"
	"github.com/pkg/errors"
)

func (a *apiImpl) StartHTTP() error {
	app := fiber.New()

	app.Get("/api/list", a.getNewsHandler)
	app.Post("/api/:id", a.editNewsHandler)
	app.Post("/add", a.addNewsHandler)

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
func (a *apiImpl) getNewsHandler(c *fiber.Ctx) {

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
		return
	}

	log.Info("run get News")

	news, err := a.service.GetNews(c.Context(), pNews)
	if err != nil {
		a.log.Error("occurred error for GetNews", sl.Err(err))
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"error":   err,
		})
		return
	}
	c.JSON(&fiber.Map{
		"success": true,
		"news":    news,
	})
}

//edit News
func (a *apiImpl) editNewsHandler(c *fiber.Ctx) {

	const op = "edit News"
	log := a.log.With(slog.String("op", op))

	id, err := strconv.Atoi(c.Params("id"))
	if err != nil {
		log.Error("cant convert param", sl.Err(err))
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"error":   err,
		})
		return
	}

	news := &service.News{}
	if err := c.BodyParser(&news); err != nil {
		log.Error("cant unmarshall", sl.Err(err))
		c.Status(http.StatusBadRequest).JSON(&fiber.Map{
			"success": false,
			"error":   err,
		})
		return
	}

	log.Info("run edit News", sl.Atr("News", news))

	news, err = a.service.EditNews(c.Context(), id, news)
	if err != nil {
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"error":   err,
		})
		return
	}
	c.JSON(&fiber.Map{
		"success": true,
		"news":    news,
	})
}

//add News
func (a *apiImpl) addNewsHandler(c *fiber.Ctx) {
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
		return
	}

	log.Info("run add News", sl.Atr("News", news))

	rNews, err := a.service.AddNews(c.Context(), news)
	if err != nil {
		a.log.Error("occurred error for addNews", sl.Err(err))
		c.Status(http.StatusInternalServerError).JSON(&fiber.Map{
			"success": false,
			"error":   err,
		})
		return
	}

	c.JSON(&fiber.Map{
		"success": true,
		"news_id": rNews.Id,
	})
}
