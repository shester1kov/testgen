package storage

import (
	"context"
	"io"
	"time"
)

// Storage - интерфейс для работы с хранилищем файлов
// этот интерфейс определяет контракт (набор методов), которые должны реализовывать все реализации хранилищ
// (MinIO, AWS S3...)

type Storage interface {
	// Upload загружает файл
	// Параметры:
	//	- ctx: контекст для отмены операций и таймаутов
	//	- key: уникальный ключ объекта
	//	- reader: поток данных файла
	//	- size: размер файла в байтах
	//	- contentType: MIME-тип файла, например "application/pdf"
	// Возвращает:
	//	- error: ошибку, если загрузка не удалась
	Upload(ctx context.Context, key string, reader io.Reader, size int64, contentType string) error

	// Download скачивает файл из хранилища
	// Параметры:
	//   - ctx: контекст
	//   - key: ключ объекта в хранилище
	// Возвращает:
	//   - io.ReadCloser: поток данных файла (ReadCloser = Reader + Close метод)
	//                    ВАЖНО: вызывающая сторона должна вызвать Close() после чтения!
	//   - error: ошибку, если файл не найден или произошла ошибка
	Download(ctx context.Context, key string) (io.ReadCloser, error)

	// Delete удаляет файл из хранилища
	// Параметры:
	//   - ctx: контекст
	//   - key: ключ объекта для удаления
	// Возвращает:
	//   - error: ошибку, если удаление не удалось
	Delete(ctx context.Context, key string) error

	// Exists проверяет существование файла в хранилище
	// Параметры:
	//   - ctx: контекст
	//   - key: ключ объекта
	// Возвращает:
	//   - bool: true если файл существует, false если нет
	//   - error: ошибку, если произошла проблема при проверке (не путать с "файл не найден")
	Exists(ctx context.Context, key string) (bool, error)

	// GetPresignedURL генерирует временную подписанную ссылку для прямого доступа к файлу
	// Это нужно для того, чтобы пользователь мог скачать файл напрямую из MinIO,
	// без прохождения через backend (экономия ресурсов и трафика)
	// Параметры:
	//   - ctx: контекст
	//   - key: ключ объекта
	//   - expiry: время жизни ссылки (например: 1 час, 24 часа)
	// Возвращает:
	//   - string: подписанный URL, действительный в течение expiry времени
	//   - error: ошибку, если генерация не удалась
	GetPresignedURL(ctx context.Context, key string, expiry time.Duration) (string, error)
}
