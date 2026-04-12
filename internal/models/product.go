package models

import (
	"errors"
	"regexp"
	//	"strings"
	"time"
)

// sbinRegex - регулярное выражение для валидации SBIN
// Допускает только цифры длиной от 12 до 20 символов
var sbinRegex = regexp.MustCompile(`^[0-9]{12,20}$`)

// Product - структура, представляющая товар
// Поля:
// - Name: название товара
// - SBIN: серийный номер/штрих-код
// - ExpiryDate: срок годности товара
type Product struct {
	Name       string
	SBIN       string
	ExpiryDate time.Time
}

// NewProduct - конструктор товара с валидацией
// Параметры:
// - name: название товара
// - sbin: серийный номер/штрих-код
// - date: дата в формате "DD.MM.YYYY"
// Возвращает:
// - указатель на Product при успешном создании
// - или ошибку при невалидных данных
func NewProduct(name, sbin, date string) (*Product, error) {
	// Проверка формата SBIN
	if !sbinRegex.MatchString(sbin) {
		return nil, errors.New("Некорректная длина SBIN")
	}

	layout := "02.01.2006"
	expire, error := time.Parse(layout, date)
	if error != nil {
		return nil, errors.New("Не получилось спарсить дату!")
	}

	// Проверка, что срок годности не истёк
	if expire.Before(time.Now()) {
		return nil, errors.New("Товар с истёкшим сроком годности не может быть добавлен")
	}

	return &Product{
		Name:       name,
		SBIN:       sbin,
		ExpiryDate: expire,
	}, nil
}

// DateToString - метод для форматированного вывода даты
func (p *Product) DateToString() string {
	return p.ExpiryDate.Format("02.01.2006")
}
