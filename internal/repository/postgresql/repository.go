package postgresql

import (
	"context"
	"embed"
	"errors"
	"time"

	"github.com/exaring/otelpgx"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
	_ "github.com/jackc/pgx/v5/stdlib" // https://stackoverflow.com/questions/76865674/how-to-use-goose-migrations-with-pgx
	"github.com/pressly/goose/v3"
	"github.com/vodolaz095/purser/model"
)

//go:embed migrations/*.sql
var embedMigrations embed.FS

// Repository реализует интерфейс SecretRepo с базой данных postgresql внутри
type Repository struct {
	DatabaseConnectionString string
	conn                     *pgx.Conn
}

// Ping проверяет соединение с базой данных
func (r *Repository) Ping(ctx context.Context) error {
	return r.conn.Ping(ctx)
}

// Init настраивает соединение с базой данных
func (r *Repository) Init(ctx context.Context) error {
	opts, err := pgx.ParseConfig(r.DatabaseConnectionString)
	if err != nil {
		return err
	}
	opts.Tracer = otelpgx.NewTracer()
	conn, err := pgx.ConnectConfig(ctx, opts)
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
	if err != nil {
		if errors.Is(err, goose.ErrNoNextVersion) {
			return nil
		}
		return err
	}
	return nil
}

// Close закрывает соединение с базой данных
func (r *Repository) Close(ctx context.Context) error {
	return r.conn.Close(ctx)
}

// Create создаёт новый model.Secret
func (r *Repository) Create(ctx context.Context, body string, meta map[string]string) (model.Secret, error) {
	var secret model.Secret
	now := time.Now()
	secret.Body = body
	secret.Meta = meta
	secret.CreatedAt = now
	secret.ExpireAt = now.Add(model.TTL)

	dbMeta := make(pgtype.Hstore, 0)
	for k, v := range meta {
		dbMeta[k] = &v
	}
	row := r.conn.QueryRow(ctx,
		`INSERT INTO secret (body, meta, created_at) VALUES ($1,$2::hstore,$3) RETURNING id;`,
		body, dbMeta, now,
	)
	err := row.Scan(&secret.ID)
	if err != nil {
		return model.Secret{}, err
	}
	return secret, nil
}

// FindByID ищет model.Secret по идентификатору
func (r *Repository) FindByID(ctx context.Context, id string) (model.Secret, error) {
	var secret model.Secret
	dbMeta := make(pgtype.Hstore, 0)
	row := r.conn.QueryRow(ctx, "SELECT body,meta,created_at FROM secret WHERE id = $1::uuid", id)
	err := row.Scan(&secret.Body, &dbMeta, &secret.CreatedAt)
	if err != nil {
		if err == pgx.ErrNoRows {
			return model.Secret{}, model.ErrSecretNotFound
		}
		return model.Secret{}, err
	}
	secret.Meta = make(map[string]string, len(dbMeta))
	for k := range dbMeta {
		secret.Meta[k] = *dbMeta[k]
	}
	secret.ID = id
	secret.ExpireAt = secret.CreatedAt.Add(model.TTL)
	return secret, nil
}

// DeleteByID удаляет секрет по идентификатору
func (r *Repository) DeleteByID(ctx context.Context, id string) error {
	_, err := r.conn.Exec(ctx, "DELETE FROM secret WHERE id = $1::uuid", id)
	return err
}

// Prune удаляет старые секреты
func (r *Repository) Prune(ctx context.Context) error {
	_, err := r.conn.Exec(ctx, "DELETE FROM secret WHERE created_at < $1", time.Now().Add(-model.TTL))
	return err
}
