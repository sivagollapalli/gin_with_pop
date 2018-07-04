package main

import (
	"gin_with_pop/actions"
	"gin_with_pop/models"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/pop"
	jwt "gopkg.in/appleboy/gin-jwt.v2"
)

func init() {
	if os.Getenv("GO_ENV") == "" {
		os.Setenv("GO_ENV", "development")
	}
	_, err := pop.Connect(os.Getenv("GO_ENV"))

	if err != nil {
		log.Fatal("Unable to connect to db server")
	}
}

func main() {
	r := gin.New()

	r.Use(gin.Logger())
	r.Use(gin.Recovery())

	pop.Debug = true

	log.Println(os.Getenv("GO_ENV"))
	db, _ := pop.Connect(os.Getenv("GO_ENV"))

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
			user := models.User{}

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

		r.POST("/users/upload", actions.ImageUpload)
		auth.POST("/users", actions.CreateUser)
		auth.GET("/users", actions.ListUsers)
		auth.GET("/users/:id", actions.ShowUser)
		auth.PATCH("/users/:id", actions.UpdateUser)
	}

	r.Run(":3001") // listen and serve on 0.0.0.0:8080
}
