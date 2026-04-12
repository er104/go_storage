package main

import (
	"flag"
	"fmt"

	// "go_storage/internal/models"
	"go_storage/internal/parser"
	"go_storage/internal/postgres"

	// "go_storage/internal/storage"
	"log"
)

// main отвечает за:
// Парсинг аргументов командной строки
// Подключение к PostgreSQL
// Чтение товаров из файла
// Сохранение товаров в БД
// Вывод списка сохранённых товаров
func main() {
	// rep := storage.NewRepository[*models.Product]()

	// Определение флагов командной строки
	fileNamePTR := flag.String("file", "data.txt", "Передайте адрес data файла") // путь к файлу с данными
	pgHost := flag.String("pg-host", "localhost", "Postgre host")                // хост PostgreSQL
	pgPort := flag.String("pg-port", "", "Postgre port")                         // порт PostgreSQL
	pgUser := flag.String("pg-user", "postgres", "Postgre user")                 // пользователь PostgreSQL
	pgPass := flag.String("pg-pass", "", "Postgre password")                     // пароль PostgreSQL
	pgName := flag.String("pg-db", "productdb", "Postgre database name")         // имя базы данных

	flag.Parse() // парсинг переданных флагов

	// Формирование конфигурации для подключения к БД
	cfg := postgres.Config{
		Host:     *pgHost,
		Port:     *pgPort,
		User:     *pgUser,
		Password: *pgPass,
		DBName:   *pgName,
	}

	// Создание репозитория для работы с PostgreSQL
	repo, err := postgres.NewPGStorage(cfg)
	if err != nil {
		log.Fatal("Dont connected to DB %v", err) // фатальная ошибка при подключении
	}
	fmt.Println("Connected to Postgre. Table created successfully")

	// Получение пути к файлу из флага
	currentFile := *fileNamePTR

	// var repos storage.productStorage = storage.

	// Вывод информации о файле для чтения
	fmt.Printf("Система настроена на чтение из файла: %s", currentFile)
	fmt.Printf("Чтение из файла (%s)", currentFile)

	// Парсинг товаров из файла
	Products, err := parser.ParseProductsFromFile(currentFile)

	if err != nil {
		log.Fatalf("Fatal Error: %v\n", err) // фатальная ошибка при чтении файла
	}

	// Добавление каждого товара в базу данных
	for _, product := range Products {
		repo.Add(product)
	}

	// Получение всех товаров из базы данных
	items := repo.GetAll()

	// Вывод количества отгруженных товаров
	fmt.Printf("На склад отгружено %d товаров \n", len(items))

	// Детальный вывод информации о каждом товаре
	for _, item := range items {
		fmt.Printf("-Name: %10s | SBIN: %20s | Годен до %s \n", item.Name, item.SBIN, item.DateToString())
	}

	fmt.Scanln()
}
