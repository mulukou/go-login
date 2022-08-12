package main

import (
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type jwtCustomClaims struct {
	Admin bool `json:"admin"`
	jwt.StandardClaims
}

const keySecret = "windows_sucks"

func login(c echo.Context) error {
	username := c.FormValue("username")
	password := c.FormValue("password")

	if username != "test" || password != "test" {
		return echo.ErrUnauthorized
	}

	expiration := time.Now().Add(time.Hour * 72).Unix()

	claims := &jwtCustomClaims{
		true,
		jwt.StandardClaims{
			ExpiresAt: expiration,
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	t, err := token.SignedString([]byte(keySecret))
	if err != nil {
		return err
	}

	return c.JSON(http.StatusOK, echo.Map{
		"token":      t,
		"expiration": expiration,
	})
}

func main() {
	e := echo.New()

	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{http.MethodGet, http.MethodHead, http.MethodPut, http.MethodPatch, http.MethodPost, http.MethodDelete},
	}))
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())

	e.POST("/login", login)

	r := e.Group("/restricted")

	config := middleware.JWTConfig{
		Claims:     &jwtCustomClaims{},
		SigningKey: []byte(keySecret),
	}
	r.Use(middleware.JWTWithConfig(config))

	r.GET("", func(c echo.Context) error {
		return c.JSON(http.StatusOK, "You're Authenticated")
	})

	e.Logger.Fatal(e.Start(":1324"))
}
