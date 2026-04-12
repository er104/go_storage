package storage

// Repository - хранилище для любого типа данных
// Позволяет переиспользовать логику хранения для разных структур
// T - тип элементов, которые будут храниться
type Repository[T any] struct {
	items []T // массив для хранения элементов типа T
}

// NewRepository - конструктор репозитория
// Создаёт и возвращает указатель на новый Repository с инициализированным массивом
func NewRepository[T any]() *Repository[T] {
	return &Repository[T]{
		items: make([]T, 0), // создаём пустой массив нужного типа
	}
}

// Add - добавление элемента в репозиторий
// Параметр item - элемент типа T, который нужно добавить
func (r *Repository[T]) Add(item T) {
	r.items = append(r.items, item) // добавляем в конец массива
}

// GetAll - получение всех элементов из репозитория
// Возвращает срез всех сохранённых элементов типа T
func (r *Repository[T]) GetAll() []T {
	return r.items
}