package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"
	"github.com/sirupsen/logrus"

	"go-clean-architecture/pkg/config"
	"go-clean-architecture/pkg/logger"
	pkgmongodb "go-clean-architecture/pkg/mongodb"
	pkgvalidator "go-clean-architecture/pkg/validator"
	todohttpdelivery "go-clean-architecture/todo/delivery/http"
	todorepository "go-clean-architecture/todo/repository"
	todoservice "go-clean-architecture/todo/service"
	responseutil "go-clean-architecture/utils/response"
)

func Routes() *chi.Mux {
	router := chi.NewRouter()
	router.Use(
		render.SetContentType(render.ContentTypeJSON), // Set content-Type headers as application/json
		middleware.Logger, // Log API request calls
		// middleware.DefaultCompress, // Compress results, mostly gzipping assets and json
		middleware.RedirectSlashes, // Redirect slashes to no slash URL versions
		middleware.Recoverer,       // Recover from panics without crashing server
	)

	return router
}

// PrintAllRoutes - printing all routes
func PrintAllRoutes(router *chi.Mux) {
	walkFunc := func(method string, route string, handler http.Handler, middlewares ...func(http.Handler) http.Handler) error {
		logrus.Printf("%s %s\n", method, route) // Walk and print out all routes
		return nil
	}
	if err := chi.Walk(router, walkFunc); err != nil {
		logger.Error(err)
	}
}

func main() {
	pkgvalidator.New()

	// Load environment variables
	err := config.LoadConfig()
	if err != nil {
		logger.Error(err)
	}

	// Init MongoDB
	_, cancel, client := pkgmongodb.InitMongoDB()
	defer cancel()

	router := Routes()

	router.Get("/", func(w http.ResponseWriter, r *http.Request) {
		render.JSON(w, r, responseutil.H{
			"success": "true",
			"code":    200,
			"message": "Services run properly",
		})
	})

	// Repository
	todoRepo := todorepository.New(client)

	// Service
	todoService := todoservice.New(todoRepo)

	// Handler
	todoHandler := todohttpdelivery.New(router, todoService)
	todoHandler.RegisterRoutes(router)

	// Print
	PrintAllRoutes(router)

	addr := fmt.Sprintf("%s%s", ":", os.Getenv("PORT"))
	err = http.ListenAndServe(addr, router)
	if err != nil {
		logger.Error(err)
	}
}
