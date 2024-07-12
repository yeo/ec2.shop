package common

import "github.com/labstack/echo/v4"

type SearchResult interface {
	RenderJSON(echo.Context) error
	RenderText(echo.Context) error
}
