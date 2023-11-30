package db

import (
	"context"
	"fmt"

	pgx "github.com/jackc/pgx/v5"
	"github.com/pgvector/pgvector-go"
	"github.com/rcleveng/assistant/server/env"
)

type EmbeddingsDB interface {
	Add(author int64, text string, embeddings []float32) (int64, error)
	Close()
}

type Embeddings struct {
	ctx  context.Context
	conn *pgx.Conn
}

// returns chunk id
func (emb *Embeddings) Add(author int64, text string, embeddings []float32) (int64, error) {
	sql := `
INSERT INTO embeddings(
	content, tokens, author, created, embedding
) VALUES(
	$1, $2, $3, NOW(), $4
) RETURNING id;`
	var id int64
	if err := emb.conn.QueryRow(emb.ctx, sql, text, 0, author, pgvector.NewVector(embeddings)).Scan(&id); err != nil {
		return 0, err
	}
	return id, nil
}

func (emb *Embeddings) Close() {
	emb.conn.Close(emb.ctx)
}

func NewEmbeddings(env *env.Environment) (*Embeddings, error) {
	// urlExample := "postgres://username:password@localhost:5432/database_name"
	dbname := env.DatabaseDatabase
	if len(dbname) == 0 {
		dbname = "assistant"
	}
	dbport := 5432

	url := fmt.Sprintf("postgres://%s:%s@%s:%d/%s", env.DatabaseUserName, env.DatabasePassword, env.DatabaseHostname, dbport, dbname)
	conn, err := pgx.Connect(context.Background(), url)
	if err != nil {
		return nil, err
	}
	return &Embeddings{
		ctx:  context.Background(),
		conn: conn,
	}, err

}
