package postgres

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/jmoiron/sqlx" // Импортируем sqlx
	_ "github.com/lib/pq"     // Импортируем драйвер PostgreSQL
	"log"
	"os"
	"pet/internal/domain/models"
	"pet/internal/storage"
)

func MustGetEnv(key string) string {
	val := os.Getenv(key)
	if val == "" {
		log.Fatalf("FATAL: Environment variable %s not set", key)
	}
	return val
}

var ConnString = MustGetEnv("DATABASE_URL")

type Storage struct {
	log *log.Logger
	db  *sqlx.DB
}

// Конструктор Storage
func New() (*Storage, error) {
	const op = "storage.postgres.New"

	// Указываем строку подключения к БД
	db := sqlx.MustConnect("postgres", ConnString)

	return &Storage{db: db}, nil
}

func (s *Storage) SaveGoods(ctx context.Context, brand string, placeSave int64, storeHouse int64, worker string) (int64, error) {
	const op = "storage.postgres.SaveGoods"

	res, err := s.db.ExecContext(ctx, "INSERT INTO goods (brand, place_save, store_house, worker) VALUES ($1, $2, $3, $4)", brand, placeSave, storeHouse, worker)
	if err != nil {
		s.log.Errorf("%s: Error inserting goods: %v", op, err)
		return 0, fmt.Errorf("%s: %w", op, err)
	}

	id, err := res.LastInsertId()
	if err != nil {
		s.log.Errorf("%s: Unable to retrieve last insert id: %v", op, err)
		return 0, fmt.Errorf("%s: unable to retrieve last insert id: %w", op, err)
	}

	s.log.Infof("%s: Successfully saved goods with id=%d", op, id)
	return id, nil
}

//func (s *Storage) SaveGoods(ctx context.Context, brand string, placeSave int64, storeHouse int64, worker string) (int64, error) {
//	const op = "storage.postgres.SaveGoods"
//
//	// Выполняем запрос к БД
//	//stmt, err := s.db.Prepare("INSERT INTO goods (brand, place_save, store_house, worker) VALUES ($1, $2, $3, $4)")
//	//if err != nil {
//	//	return 0, fmt.Errorf("%s: %w", op, err)
//	//}
//	//res, err := stmt.ExecContext(ctx, brand, placeSave, storeHouse, worker)
//	//if err != nil {
//	//	if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" { // Код ошибки для уникального ограничения
//	//		return 0, fmt.Errorf("%s: %w", op, storage.ErrGoodsExists)
//	//	}
//	//	return 0, fmt.Errorf("%s: %w", op, err)
//	//}
//
//	res, err := s.db.ExecContext(ctx, "INSERT INTO goods (brand, place_save, store_house, worker) VALUES ($1, $2, $3, $4)", brand, placeSave, storeHouse, worker)
//	if err != nil {
//		if pgErr, ok := err.(*pq.Error); ok && pgErr.Code == "23505" { // Код ошибки для уникального ограничения
//			return 0, fmt.Errorf("%s: %w", op, storage.ErrGoodsExists)
//		}
//		return 0, fmt.Errorf("%s: %w", op, err)
//	}
//
//	id, err := res.LastInsertId() // Получаем последний вставленный ID
//	if err != nil {
//		return 0, fmt.Errorf("%s: unable to retrieve last insert id: %w", op, err)
//	}
//
//	return id, nil
//}

func (s *Storage) Goods(ctx context.Context, id int64) (models.Goods, error) {
	const op = "storage.postgres.Goods"

	var goods models.Goods
	err := s.db.GetContext(ctx, &goods, "SELECT id, brand, place_save, store_house, worker FROM goods WHERE id = $1", id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return models.Goods{}, fmt.Errorf("%s: %w", op, storage.ErrGoodsNotFound)
		}
		return models.Goods{}, fmt.Errorf("%s: %w", op, err)
	}
	return goods, nil
}

func (s *Storage) DeleteGoods(ctx context.Context, id int64) (bool, error) {
	const op = "storage.postgres.DeleteGoods"

	// Выполняем запрос на удаление
	res, err := s.db.ExecContext(ctx, "DELETE FROM goods WHERE id = $1", id)
	if err != nil {
		return false, fmt.Errorf("%s: %w", op, err)
	}

	// Проверяем количество затронутых строк
	rowsAffected, err := res.RowsAffected()
	if err != nil {
		return false, fmt.Errorf("%s: unable to retrieve affected rows: %w", op, err)
	}

	// Если ни одна строка не была затронута, значит товар не найден
	if rowsAffected == 0 {
		return false, fmt.Errorf("%s: %w", op, storage.ErrGoodsNotFound)
	}

	return true, nil // Удаление прошло успешно
}