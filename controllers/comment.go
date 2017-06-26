package controllers

import (
	"encoding/json"

	dbpkg "github.com/Goryudyuma/kstmisucon1/db"
	"github.com/Goryudyuma/kstmisucon1/helper"
	"github.com/Goryudyuma/kstmisucon1/models"
	"github.com/Goryudyuma/kstmisucon1/sessions"

	"github.com/gin-gonic/gin"
)

func GetComments(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.Comment{})
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
	comments := []models.Comment{}
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.Comment{}, fields)

	if err := db.Select(queryFields).Find(&comments).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	index := 0

	if len(comments) > 0 {
		index = int(comments[len(comments)-1].ID)
	}

	if err := parameter.SetHeaderLink(c, index); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if _, ok := c.GetQuery("stream"); ok {
		enc := json.NewEncoder(c.Writer)
		c.Status(200)

		for _, comment := range comments {
			fieldMap, err := helper.FieldToMap(comment, fields)
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

		for _, comment := range comments {
			fieldMap, err := helper.FieldToMap(comment, fields)
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

func GetComment(c *gin.Context) {
	db := dbpkg.DBInstance(c)
	parameter, err := dbpkg.NewParameter(c, models.Comment{})
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	db = parameter.SetPreloads(db)
	comment := models.Comment{}
	id := c.Params.ByName("id")
	fields := helper.ParseFields(c.DefaultQuery("fields", "*"))
	queryFields := helper.QueryFields(models.Comment{}, fields)

	if err := db.Select(queryFields).First(&comment, id).Error; err != nil {
		content := gin.H{"error": "comment with id#" + id + " not found"}
		c.JSON(404, content)
		return
	}

	fieldMap, err := helper.FieldToMap(comment, fields)
	if err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if _, ok := c.GetQuery("pretty"); ok {
		c.IndentedJSON(200, fieldMap)
	} else {
		c.JSON(200, fieldMap)
	}
}

func CreateComment(c *gin.Context) {
	userID, err := sessions.LoginID(c)
	if err != nil {
		return
	}

	db := dbpkg.DBInstance(c)
	comment := models.Comment{}

	if err := c.Bind(&comment); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	comment.WriterID = userID
	if rows, err := db.Raw("SELECT user_name FROM users WHERE id = ?", userID).Rows(); err == nil {
		rows.Next()
		rows.Scan(&comment.WriterName)
		rows.Close()
	} else {
		c.JSON(400, gin.H{"error": err.Error()})
		rows.Close()
		return
	}
	if rows, err := db.Raw("SELECT comment FROM comments WHERE id = ?", comment.ParentID).Rows(); err == nil {
		rows.Next()
		rows.Scan(&comment.ParentComment)
		rows.Close()
	} else {
		c.JSON(400, gin.H{"error": err.Error()})
		rows.Close()
		return
	}

	if err := db.Create(&comment).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	c.JSON(201, comment)
}

func UpdateComment(c *gin.Context) {
	userID, err := sessions.LoginID(c)
	if err != nil {
		return
	}
	db := dbpkg.DBInstance(c)
	id := c.Params.ByName("id")
	comment := models.Comment{}

	if db.First(&comment, id).Error != nil {
		content := gin.H{"error": "comment with id#" + id + " not found"}
		c.JSON(404, content)
		return
	}

	if comment.WriterID != userID {
		content := gin.H{"error": "他人のコメントは編集できません"}
		c.JSON(504, content)
		return
	}

	if err := c.Bind(&comment); err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	if err := db.Save(&comment).Error; err != nil {
		c.JSON(400, gin.H{"error": err.Error()})
		return
	}

	childrenComments := []models.Comment{}
	db.Raw("SELECT * FROM comments WHERE parent_id = ?", id).Scan(&childrenComments)
	for _, one := range childrenComments {
		one.ParentComment = comment.Comment
		if err := db.Save(&one).Error; err != nil {
			c.JSON(400, gin.H{"error": err.Error()})
			return
		}
	}

	c.JSON(200, comment)
}
