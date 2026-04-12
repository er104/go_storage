package postgres

import (
	"database/sql"                        // интерфейс для работы с SQL
	"fmt"                                 
	"log"                                 
	"go_storage/internal/models"
	"time"                                
	"github.com/jackc/pgx/v5"            // драйвер PostgreSQL
	"github.com/jackc/pgx/v5/stdlib"     // адаптер для использования pgx с database/sql
	"go_storage/internal/storage"         
)

// PGStorage - структура, реализующая интерфейс ProductStorage
// Хранит подключение к PostgreSQL базе данных
type PGStorage struct {
	db *sql.DB // пул соединений с базой данных
}

// Config - конфигурация подключения к PostgreSQL
type Config struct {
	Host     string // хост сервера БД 
	Port     string // порт 
	User     string // имя пользователя
	Password string // пароль
	DBName   string // имя базы данных
}

// NewPGStorage - конструктор PostgreSQL хранилища
// Параметры:
// - cfg: конфигурация подключения
// Возвращает:
// - ProductStorage: реализацию интерфейса
// - error: ошибку при подключении или создании таблицы
func NewPGStorage(cfg Config) (storage.ProductStorage, error) {
	// Data Source Name - строка, содержащая все параметры для подключения
	dsn := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host, cfg.Port, cfg.User, cfg.Password, cfg.DBName)

	// Парсинг DSN строки в структуру ConnConfig
	// Позволяет валидировать параметры и использовать pgx
	connConfig, err := pgx.ParseConfig(dsn)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	// Открытие соединения с БД через интерфейс database/sql
	db := stdlib.OpenDB(*connConfig)

	// Проверка связи с БД
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("ping error: %w", err)
	}

	// Создание таблицы, если она не существует
	if err := createTable(db); err != nil {
		return nil, fmt.Errorf("create table error: %w", err)
	}

	return &PGStorage{db: db}, nil
}

// createTable - создание таблицы products в БД
// SQL запрос с IF NOT EXISTS гарантирует, что повторный вызов не вызовет ошибку
func createTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS products (
	id SERIAL PRIMARY KEY,           -- автоинкрементируемый первичный ключ
	name TEXT NOT NULL,              -- название товара (обязательное поле)
	sbin TEXT NOT NULL,              -- серийный номер (обязательное поле)
	expiry_date DATE NOT NULL        -- срок годности (тип DATE)
	);
	`
	_, err := db.Exec(query) // выполнение запроса
	return err
}

// Add - добавление товара в базу данных
// Реализует метод интерфейса ProductStorage
func (storage *PGStorage) Add(product *models.Product) {
	query := `INSERT INTO products(name, sbin, expiry_date) VALUES ($1, $2, $3)`
	// $1, $2, $3 - плейсхолдеры для параметров

	_, err := storage.db.Exec(query, product.Name, product.SBIN, product.ExpiryDate)
	if err != nil {
		// Логируем ошибку, не прерываем выполнение
		log.Printf("Add product error %s: %v", product.Name, err)
	}
}

// GetAll - получение всех товаров из базы данных
// Реализует метод интерфейса ProductStorage
// Возвращает срез всех товаров, хранящихся в таблице products
func (storage *PGStorage) GetAll() []*models.Product {
	// rows - итератор по строкам результата запроса
	rows, err := storage.db.Query(`SELECT name, sbin, expiry_date FROM products`)
	if err != nil {
		log.Printf("query error: %v", err) // логируем ошибку запроса
		return nil
	}

	defer rows.Close() // гарантированное закрытие итератора

	var products []*models.Product // массив для хранения результата

	// Итерация по всем строкам результата
	for rows.Next() {
		var name, sbin string
		var expiryDate time.Time

		// Сканирование полей текущей строки в переменные
		if err := rows.Scan(&name, &sbin, &expiryDate); err != nil {
			log.Printf("scan error: %v", err) // логируем ошибку сканирования
			continue                          // пропускаем проблемную строку
		}

		// Создание модели товара и добавление в результат
		products = append(products, &models.Product{
			Name:       name,
			SBIN:       sbin,
			ExpiryDate: expiryDate,
		})
	}

	return products
}