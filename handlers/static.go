package handlers

import "github.com/labstack/echo"

type Info struct {
	Streams []string
}

func (h *Handlers) ControlPage(c echo.Context) error {
	return nil
}
