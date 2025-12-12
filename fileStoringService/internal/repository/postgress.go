package repository

import (
	"context"
	"fmt"
	"log"

	"fileStoringService/internal/domain/entities"
	"fileStoringService/internal/domain/repository"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
)

type PostgresDB struct {
	Pool *pgxpool.Pool
}

func NewPostgresDB(PostgresDBUrl string) repository.Repo {
	if PostgresDBUrl == "" {
		log.Fatal("DATABASE_URL не установлено")
	}

	pool, err := pgxpool.New(context.Background(), PostgresDBUrl)
	if err != nil {
		log.Fatalf("не получилось подключиться к DB %v", err)
	}

	return &PostgresDB{Pool: pool}
}

func (p *PostgresDB) CreateFile(ctx context.Context, file entities.File) error {
	if file.ID == uuid.Nil || file.FileName == "" || file.ContentType == "" {
		return fmt.Errorf("некорректные данные файла")
	}

	_, err := p.Pool.Exec(ctx,
		"INSERT INTO files (id, file_name, content_type) VALUES ($1, $2, $3)",
		file.ID, file.FileName, file.ContentType)
	if err != nil {
		return fmt.Errorf("ошибка при создании файла %w", err)
	}
	return nil
}

func (p *PostgresDB) CreateWork(ctx context.Context, work entities.Work) error {
	if work.ID == uuid.Nil || work.File.ID == uuid.Nil || work.UserName == "" || work.TypeWork == "" {
		return fmt.Errorf("некорректные данные работы")
	}

	_, err := p.Pool.Exec(ctx,
		"INSERT INTO works (id, user_name, created_at, type_work, file_id) VALUES ($1, $2, $3, $4, $5)",
		work.ID, work.UserName, work.CreatedAt, work.TypeWork, work.File.ID)
	if err != nil {
		return fmt.Errorf("ошибка при создании работы %w", err)
	}
	return nil
}

func (p *PostgresDB) GetWork(ctx context.Context, id string) (entities.Work, error) {
	var w entities.Work
	var f entities.File

	err := p.Pool.QueryRow(ctx,
		`SELECT w.id, w.user_name, w.created_at, w.type_work,
		        f.id, f.file_name, f.content_type
		   FROM works w
		   JOIN files f ON w.file_id = f.id
		  WHERE w.id=$1`, id).
		Scan(&w.ID, &w.UserName, &w.CreatedAt, &w.TypeWork,
			&f.ID, &f.FileName, &f.ContentType)
	if err != nil {
		return entities.Work{}, fmt.Errorf("ошибка при получении работы %w", err)
	}

	w.File = f
	return w, nil
}

func (p *PostgresDB) GetWorksByType(ctx context.Context, typeWork string) ([]entities.Work, error) {
	rows, err := p.Pool.Query(ctx,
		`SELECT w.id, w.user_name, w.created_at, w.type_work,
		        f.id, f.file_name, f.content_type
		   FROM works w
		   JOIN files f ON w.file_id = f.id
		  WHERE w.type_work=$1`, typeWork)
	if err != nil {
		return nil, fmt.Errorf("ошибка при получении работ по типу %w", err)
	}
	defer rows.Close()

	var works []entities.Work
	for rows.Next() {
		var w entities.Work
		var f entities.File
		if err := rows.Scan(&w.ID, &w.UserName, &w.CreatedAt, &w.TypeWork,
			&f.ID, &f.FileName, &f.ContentType); err != nil {
			return nil, fmt.Errorf("ошибка при сканировании строки %w", err)
		}
		w.File = f
		works = append(works, w)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("ошибка при итерации по строкам %w", err)
	}

	return works, nil
}

func (p *PostgresDB) Close() {
	if p.Pool != nil {
		p.Pool.Close()
	}
}
