package storage

import (
	"go_storage/internal/models"
)

// ProductStorage - интерфейс для работы с хранилищем товаров
// Определяет контракт, которому должны следовать все реализации хранилищ
type ProductStorage interface {
	Add(product *models.Product)      // добавление товара
	GetAll() []*models.Product        // получение всех товаров
}

// Storage -  реализация хранилища товаров 
type Storage struct {
	products []*models.Product // массив для хранения товаров
}

// NewStorage - конструктор хранилища
// Возвращает экземпляр, реализующий интерфейс ProductStorage
func NewStorage() ProductStorage {
	return &Storage{
		products: make([]*models.Product, 0), // инициализация пустого массива
	}
}

// Add - добавление товара в хранилище
func (storage *Storage) Add(product *models.Product) {
	storage.products = append(storage.products, product)
}

// GetAll - получение всех товаров из хранилища
// Возвращает срез всех сохранённых товаров
func (storage *Storage) GetAll() []*models.Product {
	return storage.products
}