// @title       User Management API
// @version     1.0
// @BasePath    /
// @description REST API for user operations
package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"user-management/app"
	"user-management/cmd/migrations"
	"user-management/internal/user-management/domain/controller"
	"user-management/internal/user-management/domain/service"
	"user-management/internal/user-management/infrastructure/repository"

	"github.com/gorilla/mux"
	"github.com/rs/zerolog/log"

	_ "user-management/docs" // auto-generated Swagger docs

	httpSwagger "github.com/swaggo/http-swagger"
	"github.com/uptrace/bun/migrate"
	"github.com/urfave/cli/v2"
)

func main() {
	app := &cli.App{
		Name: "service",
		Flags: []cli.Flag{
			&cli.StringFlag{
				Name:  "env",
				Value: "dev",
				Usage: "environment",
			},
		},
		Commands: []*cli.Command{
			httpCommand,
			newDBCommand(migrations.Migrations),
		},
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal().Msg(err.Error())
	}
}

var httpCommand = &cli.Command{
	Name:  "http",
	Usage: "start grpc service",
	Flags: []cli.Flag{
		&cli.StringFlag{
			Name:  "addr",
			Value: ":8087",
			Usage: "serve address",
		},
	},
	Action: func(c *cli.Context) error {
		// read config from env
		_, app, err := app.Start(c.Context, "service", c.String("env"))
		if err != nil {
			return err
		}
		if app.Config().UserGeoApiToken == "" {
			return fmt.Errorf("User geo api token is missing")
		}
		defer app.Stop()

		repo := repository.NewUserRepository(app.DB())
		userService := service.NewUserService(repo, app.Config().UserGeoApiToken)
		controller := controller.NewController(userService)

		// create http router
		router := mux.NewRouter()
		router.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)
		router.HandleFunc("/users", controller.CreateUser).Methods("POST")
		router.HandleFunc("/users", controller.GetUsers).Methods("GET")
		router.HandleFunc("/users/{id:[0-9]+}", controller.GetUserByID).Methods("GET")
		router.HandleFunc("/users/{id:[0-9]+}", controller.UpdateUser).Methods("PUT")
		router.HandleFunc("/users/{id:[0-9]+}", controller.DeleteUser).Methods("DELETE")

		httpSrv := &http.Server{
			Addr:    c.String("addr"),
			Handler: router,
		}

		go func() {
			log.Info().Msgf("Starting HTTP server at %s", "8087")
			if err := httpSrv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
				log.Fatal().Msgf("HTTP server error: %v", err)
			}
		}()

		quit := make(chan os.Signal, 1)
		signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
		<-quit

		log.Info().Msg("Shutdown signal received, stopping services...")

		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := httpSrv.Shutdown(ctx); err != nil {
			log.Error().Msgf("HTTP shutdown error: %v", err)
		} else {
			log.Info().Msg("HTTP server stopped gracefully")
		}

		log.Info().Msg("All services stopped")

		return nil
	},
}

func newDBCommand(migrations *migrate.Migrations) *cli.Command {
	return &cli.Command{
		Name:  "db",
		Usage: "manage database migrations",
		Subcommands: []*cli.Command{
			{
				Name:  "init",
				Usage: "create migration tables",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					migrator := migrate.NewMigrator(app.DB(), migrations)
					return migrator.Init(ctx)
				},
			},
			{
				Name:  "migrate",
				Usage: "migrate database",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					migrator := migrate.NewMigrator(app.DB(), migrations)

					group, err := migrator.Migrate(ctx)
					if err != nil {
						return err
					}

					if group.ID == 0 {
						fmt.Printf("there are no new migrations to run\n")
						return nil
					}

					fmt.Printf("migrated to %s\n", group)
					return nil
				},
			},
			{
				Name:  "rollback",
				Usage: "rollback the last migration group",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					migrator := migrate.NewMigrator(app.DB(), migrations)

					group, err := migrator.Rollback(ctx)
					if err != nil {
						return err
					}

					if group.ID == 0 {
						fmt.Printf("there are no groups to roll back\n")
						return nil
					}

					fmt.Printf("rolled back %s\n", group)
					return nil
				},
			},
			{
				Name:  "lock",
				Usage: "lock migrations",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					migrator := migrate.NewMigrator(app.DB(), migrations)
					return migrator.Lock(ctx)
				},
			},
			{
				Name:  "unlock",
				Usage: "unlock migrations",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					migrator := migrate.NewMigrator(app.DB(), migrations)
					return migrator.Unlock(ctx)
				},
			},
			{
				Name:  "create_go",
				Usage: "create Go migration",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					migrator := migrate.NewMigrator(app.DB(), migrations)

					name := strings.Join(c.Args().Slice(), "_")
					mf, err := migrator.CreateGoMigration(ctx, name)
					if err != nil {
						return err
					}
					fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)

					return nil
				},
			},
			{
				Name:  "create_sql",
				Usage: "create up and down SQL migrations",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					migrator := migrate.NewMigrator(app.DB(), migrations)

					name := strings.Join(c.Args().Slice(), "_")
					files, err := migrator.CreateSQLMigrations(ctx, name)
					if err != nil {
						return err
					}

					for _, mf := range files {
						fmt.Printf("created migration %s (%s)\n", mf.Name, mf.Path)
					}

					return nil
				},
			},
			{
				Name:  "status",
				Usage: "print migrations status",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					migrator := migrate.NewMigrator(app.DB(), migrations)

					ms, err := migrator.MigrationsWithStatus(ctx)
					if err != nil {
						return err
					}
					fmt.Printf("migrations: %s\n", ms)
					fmt.Printf("unapplied migrations: %s\n", ms.Unapplied())
					fmt.Printf("last migration group: %s\n", ms.LastGroup())

					return nil
				},
			},
			{
				Name:  "mark_applied",
				Usage: "mark migrations as applied without actually running them",
				Action: func(c *cli.Context) error {
					ctx, app, err := app.StartCLI(c)
					if err != nil {
						return err
					}
					defer app.Stop()

					migrator := migrate.NewMigrator(app.DB(), migrations)

					group, err := migrator.Migrate(ctx, migrate.WithNopMigration())
					if err != nil {
						return err
					}

					if group.ID == 0 {
						fmt.Printf("there are no new migrations to mark as applied\n")
						return nil
					}

					fmt.Printf("marked as applied %s\n", group)
					return nil
				},
			},
		},
	}
}
