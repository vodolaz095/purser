package postgresql

import (
	"context"
	"embed"
	"time"

	"github.com/jackc/pgx/v5"
	_ "github.com/jackc/pgx/v5/stdlib" // https://stackoverflow.com/questions/76865674/how-to-use-goose-migrations-with-pgx
	"github.com/pressly/goose/v3"

	"github.com/vodolaz095/purser/model"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

type Repository struct {
	DatabaseConnectionString string
	conn                     *pgx.Conn
}

func (r *Repository) Ping(ctx context.Context) error {
	return r.conn.Ping(ctx)
}

func (r *Repository) Init(ctx context.Context) error {
	conn, err := pgx.Connect(ctx, r.DatabaseConnectionString)
	if err != nil {
		return err
	}
	r.conn = conn
	db, err := goose.OpenDBWithDriver("pgx", r.DatabaseConnectionString)
	if err != nil {
		return err
	}
	goose.SetBaseFS(embedMigrations)
	err = goose.SetDialect("postgres")
	if err != nil {
		return err
	}
	err = goose.Up(db, "migrations")
	return err
}

func (r *Repository) Close(ctx context.Context) error {
	return r.conn.Close(ctx)
}

func (r *Repository) Create(ctx context.Context, body string, meta map[string]string) (model.Secret, error) {
	return model.Secret{}, nil
}

func (r *Repository) FindByID(ctx context.Context, id string) (model.Secret, error) {
	//	row := r.conn.QueryRow(ctx, "SELECT * FROM secret WHERE id = $1", id)
	return model.Secret{}, nil
}

func (r *Repository) DeleteByID(ctx context.Context, id string) error {
	_, err := r.conn.Exec(ctx, "DELETE FROM secret WHERE id = $1", id)
	return err
}

func (r *Repository) Prune(ctx context.Context) error {
	_, err := r.conn.Exec(ctx, "DELETE FROM secret WHERE created_at < $1", time.Now().Add(-model.TTL))
	return err
}
