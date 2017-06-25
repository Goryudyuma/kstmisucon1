package controllers

import (
	"encoding/json"
	"strconv"

	dbpkg "github.com/Goryudyuma/kstmisucon1/db"
	"github.com/Goryudyuma/kstmisucon1/helper"
	"github.com/Goryudyuma/kstmisucon1/models"
	"github.com/Goryudyuma/kstmisucon1/sessions"
	"github.com/Goryudyuma/kstmisucon1/version"

	"golang.org/x/crypto/bcrypt"

	"github.com/gin-gonic/gin"
)

func GetUserById(c *gin.Context) {
	db := dbpkg.DBInstance(c)

	id := c.Params.ByName("id")

	ret := models.User{}
	db.Raw(`SELECT * FROM users WHERE id = ? limit 1`, id).Scan(&ret)

	ret.Password = "hidden"

	c.JSON(200, ret)
}

func GetUsers(c *gin.Context) {
	ver, err := version.New(c)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.User{})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db, err = parameter.Paginate(db)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db = parameter.SetPreloads(db)
	db = parameter.SortRecords(db)
	db = parameter.FilterFields(db)
	users := []models.User{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.User{}, fields)

	if err := db.Select(queryFields).Find(&users).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	index := 0

	if len(users) > 0 {
		index = int(users[len(users)-1].ID)
	}

	if err := parameter.SetHeaderLink(c, index); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if version.Range("1.0.0", "<=", ver) && version.Range(ver, "<", "2.0.0") {
		// conditional branch by version.
		// 1.0.0 <= this version < 2.0.0 !!
	}

	if _, ok := c.GetQuery("stream"); ok {
		enc := json.NewEncoder(c.Writer)
		c.Status(200)

		for _, user := range users {
			fieldMap, err := helper.FieldToMap(user, fields)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			if err := enc.Encode(fieldMap); err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}
		}
	} else {
		fieldMaps := []map[string]interface{}{}

		for _, user := range users {
			fieldMap, err := helper.FieldToMap(user, fields)
			if err != nil {
				c.JSON(400, gin.H{"error": err.Error()})
				return
			}

			fieldMaps = append(fieldMaps, fieldMap)
		}

		if _, ok := c.GetQuery("pretty"); ok {
			c.IndentedJSON(200, fieldMaps)
		} else {
			c.JSON(200, fieldMaps)
		}
	}
}

func makeUser(c *gin.Context, user *models.User) {
	db := dbpkg.DBInstance(c)

	if len(user.Password) < 8 {
		c.JSON(400, gin.H{"error": "パスワードは8文字以上で設定してください"})
		return
	}

	if 64 < len(user.Password) {
		c.JSON(400, gin.H{"error": "パスワードは64文字以下で設定してください"})
		return
	}

	password, err := bcrypt.GenerateFromPassword([]byte(user.Password), 0)
	if err != nil {
		c.JSON(400, gin.H{"error": "パスワードに使えない文字が含まれています:" + err.Error()})
		return
	}
	user.Password = string(password)

	if len(user.UserName) < 5 {
		c.JSON(400, gin.H{"error": "ユーザー名は5文字以上で設定してください"})
		return
	}

	if 16 < len(user.UserName) {
		c.JSON(400, gin.H{"error": "ユーザー名は16文字以下で設定してください"})
		return
	}

	users := []models.User{}
	db.Raw("SELECT * FROM users WHERE user_name = ?", user.UserName).Scan(&users)
	if len(users) != 0 {
		c.JSON(400, gin.H{"error": "すでにそのユーザー名は登録されています"})
		return
	}

	if len(user.ScreenName) < 5 {
		c.JSON(400, gin.H{"error": "表示名は5文字以上で設定してください"})
		return
	}
	if 64 < len(user.ScreenName) {
		c.JSON(400, gin.H{"error": "表示名は64文字以下で設定してください"})
		return
	}
}

func CreateUser(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	user := models.User{}

	if err := c.Bind(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	makeUser(c, &user)

	if err := db.Create(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	sessions.Login(c)

	c.JSON(201, nil)
}

func UpdateUser(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	id, err := sessions.LoginID(c)
	if err != nil {
		return
	}
	user := models.User{}

	if db.First(&user, id).Error != nil {
		content := gin.H{"error": "user with id#" + strconv.Itoa(int(id)) + " not found"}
		c.JSON(404, content)
		return
	}

	if err := c.Bind(&user); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	makeUser(c, &user)

	if err := db.Save(&user).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db.Exec("UPDATE comments SET writer_name = ? WHERE writer_id = ?", user.ScreenName, id)

	c.JSON(200, nil)
}
