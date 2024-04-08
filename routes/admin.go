package routes

import (
	"fmt"
	"net/http"
	"os"
	"search_engine/db"
	"search_engine/utils"
	"search_engine/views"
	"strconv"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

func LoginHandler(c *fiber.Ctx) error {
	return render(c, views.Login())
}

type loginForm struct {
	Email    string `form:"email"`
	Password string `form:"password"`
}

func LoginPostHandler(c *fiber.Ctx) error {
	input := loginForm{}
	if err := c.BodyParser(&input); err != nil {
		c.Status(http.StatusInternalServerError)
		return c.SendString("<h2>Error: Something went wrong !</h2>")
	}

	fmt.Println("input: ", input)

	user := &db.User{}
	user, err := user.LoginAsAdmin(input.Email, input.Password)
	fmt.Println("err: ", err)
	if err != nil {
		c.Status(http.StatusUnauthorized)
		return c.SendString("<h2>Error: Unauthorised</h2>")
	}

	signedToken, err := utils.CreateNewAuthToken(user.ID, user.Email, user.IsAdmin)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.SendString("<h2>Error: Something went wrong logging in!</h2>")
	}

	cookie := &fiber.Cookie{
		Name:     "admin",
		Value:    signedToken,
		Expires:  time.Now().Add(24 * time.Hour),
		HTTPOnly: true,
	}

	c.Cookie(cookie)
	c.Append("HX-Redirect", "/")

	return c.SendStatus(http.StatusOK)
}

func LogoutHandler(c *fiber.Ctx) error {
	c.ClearCookie("admin")
	c.Set("HX-Redirect", "/login")

	return c.SendStatus(http.StatusOK)
}

type AdminClaims struct {
	User                 string `json:"user"`
	Id                   uint   `json:"id"`
	jwt.RegisteredClaims `json:"claims"`
}

func AuthMiddleWare(c *fiber.Ctx) error {
	cookie := c.Cookies("admin")
	if cookie == "" {
		return c.Redirect("/login", http.StatusMovedPermanently)
	}

	token, err := jwt.ParseWithClaims(cookie, &AdminClaims{}, func(t *jwt.Token) (interface{}, error) {
		return []byte(os.Getenv("SECRET_KEY")), nil
	})
	if err != nil {
		return c.Redirect("/login", http.StatusMovedPermanently)
	}

	_, ok := token.Claims.(*AdminClaims)
	if ok && token.Valid {
		return c.Next()
	}

	return c.Redirect("/login", http.StatusMovedPermanently)
}

func DashboardHandler(c *fiber.Ctx) error {
	settings := db.SearchSettings{}
	err := settings.Get()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.SendString("<h2>Error: Cannot get settings !</h2>")
	}

	amount := strconv.FormatUint(uint64(settings.Amount), 10)
	return render(c, views.Home(amount, settings.SearchOn, settings.AddNew))
}

type SettingsForm struct {
	Amount   int    `form:"amount"`
	SearchOn string `form:"searchOn"`
	AddNew   string `form:"addNew"`
}

func DashboardPostHandler(c *fiber.Ctx) error {
	input := SettingsForm{}
	if err := c.BodyParser(&input); err != nil {
		c.Status(http.StatusInternalServerError)
		return c.SendString("<h2>Error: Something went wrong !</h2>")
	}

	addNew := false
	if input.AddNew == "on" {
		addNew = true
	}

	searchOn := false
	if input.SearchOn == "on" {
		searchOn = true
	}

	settings := &db.SearchSettings{}
	settings.Amount = uint(input.Amount)
	settings.AddNew = addNew
	settings.SearchOn = searchOn
	settings.ID = 1

	err := settings.Update()
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return c.SendString("<h2>Error: Cannot update settings !</h2>")
	}

	c.Append("HX-Refresh", "true")

	return c.SendStatus(http.StatusOK)
}
