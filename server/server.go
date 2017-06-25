package server

import (
	"github.com/Goryudyuma/kstmisucon1/middleware"
	"github.com/Goryudyuma/kstmisucon1/router"

	"github.com/gin-gonic/contrib/sessions"
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

func Setup(db *gorm.DB) *gin.Engine {
	r := gin.Default()
	store := sessions.NewCookieStore([]byte("kstmisucon1"))
	r.Use(sessions.Sessions("kstmsession", store))
	r.Use(middleware.SetDBtoContext(db))
	router.Initialize(r)
	return r
}
