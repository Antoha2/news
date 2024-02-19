package transport

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"

	"github.com/Antoha2/news/internal/service"
	"github.com/Antoha2/news/pkg/logger/sl"
	"github.com/gofiber/fiber"
	"github.com/pkg/errors"
)

func (a *apiImpl) StartHTTP() error {
	app := fiber.New()

	app.Get("/api/getNews", a.getNewsHandler)
	app.Post("/", a.editNewsHandler)
	app.Post("/api/addNews", a.addNewsHandler)

	err := app.Listen(fmt.Sprintf(":%s", a.cfg.HTTP.HostPort))
	if err != nil {
		return errors.Wrap(err, "ocurred error StartHTTP")
	}
	a.log.Info(fmt.Sprintf("Запуск HTTP-сервера на http://127.0.0.1:%s", a.cfg.HTTP.HostPort))
	return nil
}

func (a *apiImpl) Stop() {
	if err := a.server.Shutdown(context.TODO()); err != nil {
		panic(errors.Wrap(err, "ocurred error Stop"))
	}
}

func (a *apiImpl) getNewsHandler(c *fiber.Ctx) {

	const op = "getNews"
	log := a.log.With(slog.String("op", op))

	log.Info("run get News")

	news, err := a.service.GetNews(c.Context())
	if err != nil {
		a.log.Error("occurred error for GetNews", sl.Err(err))
		c.Status(http.StatusInternalServerError).JSON(err.Error())
		return
	}

	if err != nil || len(news) == 0 {
		c.Status(404).JSON(&fiber.Map{
			"success": false,
			"error":   "There are no posts!",
		})
		return
	}
	c.JSON(&fiber.Map{
		"success": true,
		"posts":   news,
	})
}

func (a *apiImpl) editNewsHandler(c *fiber.Ctx) {

	// const op = "getUser"
	// log := a.log.With(slog.String("op", op))

	fmt.Println("!!!!!!!!!!!!!!!!!!!!!!!!!")
}

//add News
func (a *apiImpl) addNewsHandler(c *fiber.Ctx) {
	const op = "addNews"
	log := a.log.With(slog.String("op", op))

	log.Info("run get News")

	news := &service.News{}
	if err := c.BodyParser(&news); err != nil {
		log.Error("cant unmarshall", sl.Err(err))
		c.Status(http.StatusBadRequest).JSON(err.Error())
		return
	}
	fmt.Print(news)

	log.Info("run add News", sl.Atr("News", news))

	news, err := a.service.AddNews(c.Context(), news)
	if err != nil {
		a.log.Error("occurred error for GetNews", sl.Err(err))
		c.Status(http.StatusInternalServerError).JSON(err.Error())
		return
	}

	// if err != nil || len(news) == 0 {
	// 	c.Status(404).JSON(&fiber.Map{
	// 		"success": false,
	// 		"error":   "There are no posts!",
	// 	})
	// 	return
	// }
	// c.JSON(&fiber.Map{
	// 	"success": true,
	// 	"posts":   news,
	// })
}
