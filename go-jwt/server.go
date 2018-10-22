package main

import (
	"time"

	"github.com/codegangsta/martini"
	"github.com/codegangsta/martini-contrib/binding"
	"github.com/codegangsta/martini-contrib/render"
	"github.com/dgrijalva/jwt-go"
)

// User model
type User struct {
	UserId   string `form:"userid" json:"userid" binding:"required"`
	Password string `form:"password" json:"password" binding:"required"`
}

const (
	ValidUser = "user"
	ValidPass = "passw0rd"
	SecretKey = "1234567890"
)

func main() {

	m := martini.Classic()

	m.Use(martini.Static("static"))
	m.Use(render.Renderer())

	// Authenticate user
	m.Post("/sample/auth", binding.Bind(User{}), func(user User, r render.Render) {

		if user.UserId == ValidUser && user.Password == ValidPass {

			token := jwt.New(jwt.GetSigningMethod("HS256"))
			token.Claims = jwt.MapClaims{
				"exp":    time.Now().Add(time.Minute * 5).Unix(),
			}
			tokenString, err := token.SignedString([]byte(SecretKey))
			if err != nil {
				r.HTML(201, "error", nil)
				return
			}
			data := map[string]string{
				"token": tokenString,
			}
			r.JSON(201, data)
		} else {
			r.JSON(403, nil)
		}

	})

	// Check Key is ok
	m.Get("/sample/:token", func(params martini.Params, r render.Render) {
		if authentication(params["token"]) {
			r.JSON(200, "OK")
		} else {
			r.JSON(403, nil)
		}
	})

	m.Run()
}

func authentication(requestToken string) bool {
	token, err := jwt.Parse(requestToken, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err == nil && token.Valid {
		return true
	} else {
		return false
	}
}
