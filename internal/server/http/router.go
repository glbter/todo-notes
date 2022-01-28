package http

import (
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"os"
	"strconv"
	"time"
	"todoNote/internal/repo"
	auth2 "todoNote/internal/server/http/auth"
	"todoNote/internal/server/http/log"

	"todoNote/internal/server/http/handler"
	md "todoNote/internal/server/http/middleware"
	"todoNote/internal/usecase"
)

const(
	JwtLifetimeMillisEnv = "JWT_LIFETIME_MILLIS"
	publicKeyEnv = "PUBLIC_KEY"
	privateKeyEnv = "PRIVATE_KEY"
)

func NewRouter(repo Repositories) (chi.Router, error) {
	r := chi.NewRouter()

	l, err := strconv.Atoi(os.Getenv(JwtLifetimeMillisEnv))
	if err != nil {
		return nil, err
	}

	auth, err := auth2.NewJwtAuth(
		time.Duration(l) * time.Millisecond,
		os.Getenv(privateKeyEnv),
		os.Getenv(publicKeyEnv))

	if err != nil {
		return nil, err
	}

	usecaseUser := usecase.NewUserUsecase(repo.User)
	usecaseNote := usecase.NewNoteUsecase(repo.Note)

	logger := log.MyLogger{}

	ah := handler.NewAuthHandler(usecaseUser, auth, logger)
	uh := handler.NewUserHandler(usecaseUser, logger)
	nh := handler.NewNoteHandler(usecaseNote, usecaseUser, logger)
	md := md.New(auth)

	r.Group(func(r chi.Router) {
		r.Use(middleware.Timeout(30 * time.Second))
		r.Use(middleware.RequestID)
		r.Use(middleware.Recoverer)
		r.Use(middleware.Logger)
		r.Use(md.ApmMiddleware)
		//r.Use(apmchi.Middleware())
		r.Route("/api/v1", func(r chi.Router) {

			r.Route("/notes", func(r chi.Router) {
				r.Use(md.AuthMiddleware)

				r.Post("/", nh.CreateNote)
				r.Get("/", nh.GetNotes)

				r.Route("/{noteId}", func(r chi.Router) {
					r.Get("/", nh.GetNote)
					r.Patch("/", nh.PartialUpdateNote)
					r.Delete("/", nh.DeleteNote)
				})
			})

			r.Route("/users", func(r chi.Router) {
				r.Post("/", uh.CreateUser)

				r.Group(func(r chi.Router) {
					r.Use(md.AuthMiddleware)

					r.Patch("/", uh.PartialUpdateUser)
					r.Delete("/", uh.DeleteUser)
				})
			})

			r.Post("/login", ah.Login)
		})
	})

	r.Handle("/metrics", promhttp.Handler())

	return r, nil
}

type Repositories struct {
	User repo.IRepoUser
	Note repo.IRepoNote
}
