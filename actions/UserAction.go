package actions

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/pop"
	"github.com/sivagollapalli/gin_with_pop/models"
)

func CreateUser(c *gin.Context) {
	pop.Debug = true
	db, err := pop.Connect("development")

	if err != nil {
		log.Panic(err)
	}
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
}

func ListUsers(c *gin.Context) {

	db, _ := pop.Connect("development")
	query := db.RawQuery("select id, name, email from users")
	users := []models.User{}
	err := query.All(&users)

	if err != nil {
		log.Println(err)
	}

	c.JSON(200, users)
}

func ShowUser(c *gin.Context) {
	db, _ := pop.Connect("development")
	id := c.Param("id")
	user := models.User{}
	if err := db.Find(&user, id); err != nil {
		log.Println(err)
		c.JSON(422, gin.H{"error": "user doesnt exists"})
		return
	}
	c.JSON(200, user)
}

func UpdateUser(c *gin.Context) {
	db, _ := pop.Connect("development")
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

}
