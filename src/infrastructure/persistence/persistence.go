package persistence

import (
	"database/sql"
	"fmt"
	"log"

	config "github.com/bncunha/erp-api/src/main"
	_ "github.com/lib/pq"
)

type Persistence struct{
	cfg *config.Config
}

func NewPersistence(cfg *config.Config) *Persistence {
	return &Persistence{
		cfg: cfg,
	}
}

func (p *Persistence) ConnectDb() (*sql.DB, error) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", p.cfg.DB_USER, p.cfg.DB_PASS, p.cfg.DB_HOST, p.cfg.DB_NAME) 

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Print("Erro ao conectar:", err)
		return db, err
	}

	err = db.Ping()
	if err != nil {
		log.Print("Conexão falhou:", err)
		return db, err
	}

	log.Println("Conectado com sucesso!")
	return db, nil
}

func (p *Persistence) CloseConnection(db *sql.DB) {
	db.Close()
}