package app

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"test-backend-1-kuprinvv/internal/di"
	"test-backend-1-kuprinvv/internal/middleware"
	"time"

	"github.com/go-chi/chi/v5"
	chiMiddleware "github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"
)

type App struct {
	container *di.Container
	router    chi.Router
	server    *http.Server
}

func NewApp() (*App, error) {
	container, err := di.NewContainer()
	if err != nil {
		return nil, err
	}

	a := &App{container: container}
	a.setupRouter()

	return a, nil
}

func (a *App) Run() error {
	log.Printf("server started on %s", a.server.Addr)

	go func() {
		if err := a.server.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			log.Printf("server error: %v", err)
		}
	}()

	return a.gracefulShutdown()
}

func (a *App) setupRouter() {
	ctx := context.Background()
	r := chi.NewRouter()

	r.Use(chiMiddleware.Logger)
	r.Use(chiMiddleware.Heartbeat("/_info"))
	r.Use(chiMiddleware.Recoverer)

	r.Get("/swagger/*", httpSwagger.WrapHandler)

	r.Route("/", func(r chi.Router) {
		r.Post("/dummyLogin", a.container.AuthHandler(ctx).DummyLogin)
		r.Post("/register", a.container.AuthHandler(ctx).Register)
		r.Post("/login", a.container.AuthHandler(ctx).Login)
	})

	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware(a.container.JWTConfig()))

		r.Route("/rooms", func(r chi.Router) {
			r.Get("/list", a.container.RoomHandler(ctx).ListRooms)
			r.Get("/{roomId}/slots/list", a.container.SlotHandler(ctx).GetSlots)

			r.Route("/{roomId}", func(r chi.Router) {
				r.Use(middleware.RequireRole("admin"))
				r.Post("/schedule/create", a.container.ScheduleHandler(ctx).CreateSchedule)
			})

			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRole("admin"))
				r.Post("/create", a.container.RoomHandler(ctx).CreateRoom)
			})
		})

		r.Route("/bookings", func(r chi.Router) {
			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRole("admin"))
				r.Get("/list", a.container.BookingHandler(ctx).ListBookings)
			})

			r.Group(func(r chi.Router) {
				r.Use(middleware.RequireRole("user"))
				r.Post("/create", a.container.BookingHandler(ctx).CreateBooking)
				r.Get("/my", a.container.BookingHandler(ctx).GetMyBookings)
				r.Post("/{bookingId}/cancel", a.container.BookingHandler(ctx).CancelBooking)
			})
		})
	})

	a.router = r
	a.server = &http.Server{
		Addr:              ":8080",
		Handler:           r,
		ReadHeaderTimeout: 5 * time.Second,
	}
}

func (a *App) gracefulShutdown() error {
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	log.Println("shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := a.server.Shutdown(ctx); err != nil {
		return fmt.Errorf("server shutdown failed: %w", err)
	}

	log.Println("server stopped")
	return nil
}
