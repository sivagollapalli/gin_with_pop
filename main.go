package main

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/markbates/pop"
	"github.com/sivagollapalli/gin_with_pop/models"
)

func main() {
	r := gin.Default()
	pop.Debug = true
	db, err := pop.Connect("development")

	if err != nil {
		log.Panic(err)
	}
	r.POST("/users", func(c *gin.Context) {

		user := models.User{
			Name:  c.Query("name"),
			Email: c.Query("email")}

		verrs, err := user.Validate(db)

		if err != nil {
			log.Panic(err)
		}

		if c.Query("password") != c.Query("password_confirmation") {
			verrs.Add("password", "Password should match")
		}

		if verrs.HasAny() {
			c.JSON(422, verrs)
			return
		}

		user.PasswordDigest, _ = user.HashPwd(c.Query("password"))

		_, _ = db.ValidateAndSave(&user)

		c.JSON(200, gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email})
	})

	r.GET("/users", func(c *gin.Context) {
		query := db.RawQuery("select id, name, email from users")
		users := []models.User{}
		err := query.All(&users)

		if err != nil {
			log.Println(err)
		}

		c.JSON(200, users)
	})

	r.GET("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		user := models.User{}
		if err := db.Find(&user, id); err != nil {
			log.Println(err)
			c.JSON(422, gin.H{"error": "user doesnt exists"})
			return
		}
		c.JSON(200, user)
	})
	r.PATCH("/users/:id", func(c *gin.Context) {
		id := c.Param("id")
		user := models.User{}
		if err := db.Find(&user, id); err != nil {
			log.Println(err)
			c.JSON(422, gin.H{"error": "user doesnt exists"})
			return
		}

		user.Name = c.Query("name")
		user.Email = c.Query("email")

		verrs, _ := db.ValidateAndSave(&user)

		if verrs.HasAny() {
			c.JSON(422, verrs)
			return
		}

		c.JSON(200, gin.H{
			"id":    user.ID,
			"name":  user.Name,
			"email": user.Email})

	})
	r.POST("/users/sign_in", func(c *gin.Context) {
		user := models.User{}
		email := c.Query("email")
		query := db.RawQuery("select * from users where email = ?", email)

		if err := query.First(&user); err != nil {
			log.Println(err)
			c.JSON(422, gin.H{"error": "User doesn't exists"})
			return
		}

		if user.Authenticate(c.Query("password")) {
			c.JSON(200, gin.H{"msg": "login success"})
		} else {
			c.JSON(422, gin.H{"msg": "login fail"})
		}
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}
