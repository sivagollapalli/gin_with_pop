package main

import (
	"gin_with_pop/actions"
	"gin_with_pop/models"
	"log"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/pop"
	"gopkg.in/appleboy/gin-jwt.v2"
)

func main() {
	r := gin.New()
	r.Use(gin.Logger())
	r.Use(gin.Recovery())
	pop.Debug = true
	db, _ := pop.Connect("development")

	authMiddleware := &jwt.GinJWTMiddleware{
		Realm:            "test zone",
		SigningAlgorithm: "HS512",
		Key:              []byte("secret key"),
		Timeout:          time.Hour * 24,
		MaxRefresh:       time.Hour * 24,
		PayloadFunc: func(userID string) map[string]interface{} {
			claims := make(map[string]interface{})
			claims["email"] = userID
			return claims
		},
		Authenticator: func(userId string, password string, c *gin.Context) (string, bool) {
			log.Println(userId)
			log.Println(password)
			user := models.User{}
			log.Println(userId)
			log.Println(password)

			query := db.RawQuery("select * from users where email = ?", userId)

			if err := query.First(&user); err != nil {
				log.Println(err)
				return "", false
			}
			if user.Authenticate(password) {
				return user.Email, true
			}

			return "", false
		},
		Authorizator: func(userId string, c *gin.Context) bool {
			if userId != "" {
				return true
			}

			return false
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.JSON(code, gin.H{
				"code":    code,
				"message": message,
			})
		},
		// TokenLookup is a string in the form of "<source>:<name>" that is used
		// to extract token from the request.
		// Optional. Default value "header:Authorization".
		// Possible values:
		// - "header:<name>"
		// - "query:<name>"
		// - "cookie:<name>"
		TokenLookup: "header:Authorization",
		// TokenLookup: "query:token",
		// TokenLookup: "cookie:token",

		// TokenHeadName is a string in the header. Default value is "Bearer"
		TokenHeadName: "Bearer",

		// TimeFunc provides the current time. You can override it to use another time value. This is useful for testing or if your server uses a different time zone than your tokens.
		TimeFunc: time.Now,
	}

	r.POST("/login", authMiddleware.LoginHandler)

	auth := r.Group("/")

	auth.Use(authMiddleware.MiddlewareFunc())
	{
		auth.POST("/users", actions.CreateUser)
		auth.GET("/users", actions.ListUsers)
		auth.GET("/users/:id", actions.ShowUser)
		auth.PATCH("/users/:id", actions.UpdateUser)
	}

	/*r.POST("/users/sign_in", func(c *gin.Context) {
		user := models.User{}
		email := c.Query("email")
		query := db.RawQuery("select * from users where email = ?", email)

		if err := query.First(&user); err != nil {
			log.Println(err)
			c.JSON(422, gin.H{"error": "User doesn't exists"})
			return
		}

		if user.Authenticate(c.Query("password")) {
			log.Println(user)
			token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
				"identifier": user.ID,
				"email":      user.Email})
			tokenString, error := token.SignedString([]byte("secret"))

			if error != nil {
				fmt.Println(error)
			}
			c.JSON(200, gin.H{
				"msg":   "login success",
				"token": tokenString})
		} else {
			c.JSON(422, gin.H{"msg": "login fail"})
		}
	})

	r.GET("/verify", func(c *gin.Context) {
		tokenParam := c.Query("token")

		token, _ := jwt.Parse(tokenParam, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("There was an error")
			}
			return []byte("secret"), nil
		})

		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			var user models.User
			mapstructure.Decode(claims, &user)
			log.Println(user)
			c.JSON(200, gin.H{
				"id":    user.Identifier,
				"email": user.Email})
		} else {
			c.JSON(422, gin.H{"msg": "Invalid authorization token"})
		}
	})*/

	r.Run(":8080") // listen and serve on 0.0.0.0:8080
}
