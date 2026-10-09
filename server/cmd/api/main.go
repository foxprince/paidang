package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"

	"github.com/foxprince/paidang/server/internal/config"
	"github.com/foxprince/paidang/server/internal/handler"
	"github.com/foxprince/paidang/server/internal/middleware"
	"github.com/foxprince/paidang/server/internal/repository"
	"github.com/foxprince/paidang/server/internal/service"
)

func main() {
	cfg := config.Load()

	db, err := gorm.Open(postgres.Open(cfg.DatabaseURL), &gorm.Config{})
	if err != nil {
		log.Fatalf("db: %v", err)
	}

	store := repository.New(db)
	svc := service.New(store, cfg)
	h := handler.New(svc, cfg)

	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(middleware.RateLimit(120))

	r.GET("/health", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"ok": true}) })

	v1 := r.Group("/api/v1")
	{
		v1.POST("/auth/wechat-login", h.WechatLogin)

		authed := v1.Group("", middleware.JWT(cfg.JWTSecret))
		{
			authed.GET("/coach/me", h.NotImpl)
			authed.PUT("/coach/profile", h.NotImpl)
			authed.POST("/coach/publish", h.NotImpl)
			authed.GET("/coach/today", h.NotImpl)

			authed.GET("/bookings", h.NotImpl)
			authed.PATCH("/bookings/:id/confirm", h.NotImpl)
			authed.PATCH("/bookings/:id/reject", h.NotImpl)
			authed.PATCH("/bookings/:id/complete", h.NotImpl)
			authed.PATCH("/bookings/:id/cancel", h.NotImpl)
			authed.PATCH("/bookings/:id/mark-paid", h.NotImpl)

			authed.GET("/students", h.NotImpl)
			authed.POST("/students", h.NotImpl)
			authed.GET("/students/:id", h.NotImpl)
			authed.POST("/students/:id/notes", h.NotImpl)

			authed.GET("/income/summary", h.NotImpl)
			authed.POST("/complaints", h.NotImpl)
		}

		// 学员端公开接口
		v1.GET("/public/coach/:id", h.NotImpl)
		v1.GET("/public/coach/:id/slots", h.NotImpl)
		v1.POST("/public/bookings", h.NotImpl)
		v1.POST("/public/sms-code", h.NotImpl)
	}

	log.Fatal(r.Run(":" + cfg.Port))
}
