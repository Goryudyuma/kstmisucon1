package sessions

import (
	"errors"
	"net/http"

	dbpkg "github.com/Goryudyuma/kstmisucon1/db"
	"github.com/Goryudyuma/kstmisucon1/models"

	"golang.org/x/crypto/bcrypt"

	"github.com/davecgh/go-spew/spew"
	sessionpkg "github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	db := dbpkg.DBInstance(c)

	c.Request.ParseForm()
	if _, ok := c.Request.Form["username"]; !ok {
		c.JSON(500, gin.H{"error": "usernameが含まれていません"})
		return
	}
	username := c.Request.Form["username"][0]

	if _, ok := c.Request.Form["password"]; !ok {
		c.JSON(500, gin.H{"error": "passwordが含まれていません"})
		return
	}
	password := c.Request.Form["password"][0]

	user := models.User{}
	if err := db.Raw(`SELECT * FROM users WHERE user_name = ? limit 1`, username).Scan(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	spew.Dump(user.Password)
	spew.Dump(password)
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	session := sessionpkg.Default(c)
	session.Set("userID", user.ID)
	session.Save()

	return
}

func Logout(c *gin.Context) {
	session := sessionpkg.Default(c)
	session.Clear()
	session.Save()

	c.JSON(201, nil)
}

func LoginID(c *gin.Context) (uint, error) {
	session := sessionpkg.Default(c)

	id := session.Get("userID")
	if id == nil {
		c.Redirect(http.StatusMovedPermanently, "/login")
		return 0, errors.New("loginしてください")
	}
	return id.(uint), nil
}
