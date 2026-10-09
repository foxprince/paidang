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
			authed.GET("/coach/me", h.GetMe)
			authed.PUT("/coach/profile", h.UpdateProfile)
			authed.POST("/coach/publish", h.Publish)
			authed.GET("/coach/today", h.Today)

			authed.GET("/bookings", h.ListBookings)
			authed.PATCH("/bookings/:id/confirm", h.ConfirmBooking)
			authed.PATCH("/bookings/:id/reject", h.RejectBooking)
			authed.PATCH("/bookings/:id/complete", h.CompleteBooking)
			authed.PATCH("/bookings/:id/cancel", h.CancelBooking)
			authed.PATCH("/bookings/:id/mark-paid", h.MarkPaid)

			authed.GET("/students", h.ListStudents)
			authed.POST("/students", h.CreateStudent)
			authed.GET("/students/:id", h.GetStudent)
			authed.POST("/students/:id/notes", h.AddNote)

			authed.GET("/income/summary", h.IncomeSummary)
			authed.POST("/complaints", h.CreateComplaint)
		}

		// 学员端公开接口
		v1.GET("/public/coach/:id", h.PublicCoach)
		v1.GET("/public/coach/:id/slots", h.PublicSlots)
		v1.POST("/public/bookings", h.PublicBooking)
		v1.POST("/public/sms-code", h.SmsCode)
		v1.POST("/public/reviews", h.PublicReview)
	}

	log.Fatal(r.Run(":" + cfg.Port))
}
