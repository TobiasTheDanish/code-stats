package server

import (
	"net/http"
	"tobiasthedanish/code-stats/internal/session"
	view "tobiasthedanish/code-stats/internal/view"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

func (s *Server) RegisterRoutes() http.Handler {
	e := echo.New()
	e.Static("/assets", "assets")

	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.GET("/", s.IndexHandler)

	e.GET("/health", s.healthHandler)

	e.GET("/daily", s.dailyHandler)
	e.GET("/weekly", s.weeklyHandler)

	return e
}

func (s *Server) IndexHandler(c echo.Context) error {
	daily, err := s.sessionStore.ForPeriod(session.Day)
	if err != nil {
		view.Error(err).Render(c.Request().Context(), c.Response().Writer)
	}

	c.Response().Header().Set(echo.HeaderContentType, echo.MIMETextHTML)

	return view.Index(daily.ToViewModel()).Render(c.Request().Context(), c.Response().Writer)
}

func (s *Server) healthHandler(c echo.Context) error {
	return c.JSON(http.StatusOK, s.db.Health())
}

func (s *Server) dailyHandler(c echo.Context) error {
	sessions, err := s.sessionStore.ForPeriod(session.Day)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, sessions)
}

func (s *Server) weeklyHandler(c echo.Context) error {
	sessions, err := s.sessionStore.ForPeriod(session.Week)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, map[string]string{
			"error": err.Error(),
		})
	}
	return c.JSON(http.StatusOK, sessions)
}
