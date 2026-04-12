package parser

import (
	"bufio"   // для построчного чтения файла
	"fmt"     
	"os"      // для работы с файловой системой
	"strings" 
	"go_storage/internal/models"
)

// ParseProductsFromFile - функция для парсинга товаров из текстового файла
// Формат файла: каждая строка содержит Name;SBIN;Date
// Параметры:
// - path: путь к файлу
// Возвращает:
// - срез указателей на Product при успешном чтении
// - ошибку при проблемах с открытием/чтением файла
func ParseProductsFromFile(path string) ([]*models.Product, error) {
	// Открытие файла
	file, err := os.Open(path)
	if err != nil {
		return nil, fmt.Errorf("Error: %w", err)
	}

	defer file.Close() // гарантированное закрытие файла при выходе из функции

	var Products []*models.Product // массив для хранения успешно созданных товаров
	lineNum := 1                   // счётчик строк 
	scanner := bufio.NewScanner(file) // создание сканера для построчного чтения

	// Построчное чтение файла
	for scanner.Scan() {
		line := scanner.Text() // получение текущей строки

		// Пропуск пустых строк
		if strings.TrimSpace(line) == "" {
			lineNum++
			continue
		}

		// Разделение строки на части по символу ";"
		parts := strings.Split(line, ";")

		// Проверка: должно быть ровно 3 части name, sbin, date
		if len(parts) != 3 {
			fmt.Printf("[Строка: %d] Недостаточно данных для сохранения продукта", lineNum)
			lineNum++
			continue
		}

		// Удаление лишних пробелов из каждой части
		name := strings.TrimSpace(parts[0])
		sbin := strings.TrimSpace(parts[1])
		date := strings.TrimSpace(parts[2])

		// Создание товара с валидацией
		product, err := models.NewProduct(name, sbin, date)

		if err != nil {
			// Если валидация не пройдена - вывод ошибки, без прерыва чтения
			fmt.Printf("Товар %s на строке %d отклонён. Error: %v\n", name, lineNum, err)
		} else {
			// Успешно созданный товар добавляем в результат
			Products = append(Products, product)
		}
		lineNum++
	}

	// Проверка на ошибки сканирования
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("Ошибка чтения файла %w", err)
	}

	return Products, nil
}