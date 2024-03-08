package main

import (
	"github.com/labstack/echo/v4"
	"strings"
)

func IsJson(c echo.Context) bool {
	contentType := c.Request().Header.Get("Content-Type")
	accept := c.Request().Header.Get("Accept")
	qa := c.QueryString()

	return strings.Contains(contentType, "json") || strings.Contains(accept, "json") || strings.Contains(qa, "json")
}

func IsText(c echo.Context) bool {
	ua := c.Request().Header.Get("User-Agent")
	accept := c.Request().Header.Get("Accept")

	if strings.Contains(accept, "html") {
		return false
	}

	if strings.Contains(ua, "Chrome") ||
		strings.Contains(ua, "Safari") ||
		strings.Contains(ua, "Gecko") ||
		strings.Contains(ua, "HTML") ||
		strings.Contains(ua, "iPhone") ||
		strings.Contains(ua, "iPad") ||
		strings.Contains(ua, "Mozilla") {
		return false
	}

	return true
}
