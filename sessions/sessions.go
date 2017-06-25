package sessions

import (
	"errors"
	"net/http"

	dbpkg "github.com/Goryudyuma/kstmisucon1/db"
	"github.com/Goryudyuma/kstmisucon1/models"

	"golang.org/x/crypto/bcrypt"

	sessionpkg "github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
)

func Login(c *gin.Context) {
	db := dbpkg.DBInstance(c)

	c.Request.ParseForm()
	id := c.Request.Form["userid"][0]
	password := c.Request.Form["password"][0]

	user := models.User{}
	db.Raw(`SELECT * FROM users WHERE id = ? limit 1`, id).Scan(&user)

	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))

	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	session := sessionpkg.Default(c)
	session.Set("userID", id)
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
