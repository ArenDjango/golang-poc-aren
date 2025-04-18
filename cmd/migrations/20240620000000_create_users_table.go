package migrations

import (
	"context"
	"user-management/internal/user-management/infrastructure/model"

	"github.com/uptrace/bun"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewCreateTable().
			Model((*model.User)(nil)).
			Exec(ctx)
		if err != nil {
			panic(err)
		}

		return nil
	}, func(ctx context.Context, db *bun.DB) error {
		_, err := db.NewDropTable().Model((*model.User)(nil)).IfExists().Exec(ctx)
		if err != nil {
			panic(err)
		}
		return nil
	})
}
