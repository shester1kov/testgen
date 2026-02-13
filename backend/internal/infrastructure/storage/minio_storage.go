package storage

import (
	"context"
	"fmt"
	"io"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.uber.org/zap"
)

var noSuchKeyError string = "NoSuchKey"

// MinIOStorage - реализация интерфейса Storage для MinIO
// хранит все необходимое для работы с MinIO
type MinIOStorage struct {
	client     *minio.Client // клиент для взаимодействия с API
	bucketName string        // имя бакета для сохранения файлов
	logger     *zap.Logger   // логгер для отслеживания операций
}

// NewMinIOStorage - конструктор (фабричная функция) для создания MinIOStorage
// Параметры:
//   - endpoint: адрес MinIO сервера (например: "localhost:9000")
//   - accessKey: ключ доступа (аналог логина)
//   - secretKey: секретный ключ (аналог пароля)
//   - bucketName: имя bucket для хранения файлов
//   - useSSL: использовать ли HTTPS (false для локальной разработки, true для production)
//   - logger: логгер для записи событий
//
// Возвращает:
//   - *MinIOStorage: настроенный экземпляр хранилища
//   - error: ошибку, если не удалось подключиться к MinIO
func NewMinIOStorage(
	endpoint string,
	accessKey string,
	secretKey string,
	bucketName string,
	useSSL bool,
	logger *zap.Logger,
) (*MinIOStorage, error) {
	// Создаём клиент для подключения к MinIO
	// minio.New() принимает endpoint и опции конфигурации
	minioClient, err := minio.New(endpoint, &minio.Options{
		// Credentials - учётные данные для аутентификации
		// NewStaticV4 - использует статические ключи (не временные токены)
		Creds: credentials.NewStaticV4(accessKey, secretKey, ""),
		// Secure - использовать ли HTTPS/TLS
		Secure: useSSL,
	})

	if err != nil {
		return nil, fmt.Errorf("failed to create MinIO client: %w", err)
	}

	// Шаг 2: Проверяем существование bucket и создаём его, если нужно
	ctx := context.Background()

	// BucketExists() проверяет существует ли бакет с таким именем
	exists, err := minioClient.BucketExists(ctx, bucketName)
	if err != nil {
		return nil, fmt.Errorf("failed to check bucket existence: %w", err)
	}

	// Если bucket не существует - создаём его
	if !exists {
		err = minioClient.MakeBucket(ctx, bucketName, minio.MakeBucketOptions{
			// Region можно указать, если нужно (для AWS S3)
			// Для MinIO обычно не требуется
		})
		if err != nil {
			return nil, fmt.Errorf("failed to create bucket '%s': %w", bucketName, err)
		}

		logger.Info(
			"MinIO bucket created successfully",
			zap.String("bucket", bucketName),
		)
	} else {
		logger.Info(
			"MinIO bucket created successfully",
			zap.String("bucket", bucketName),
		)
	}

	// Шаг 3: Возвращаем настроенный экземпляр MinIOStorage
	return &MinIOStorage{
		client:     minioClient,
		bucketName: bucketName,
		logger:     logger,
	}, nil
}

// Upload реализует метод Upload из интерфейса Storage
// Загружает файл в MinIO
func (s *MinIOStorage) Upload(
	ctx context.Context,
	key string,
	reader io.Reader,
	size int64,
	contentType string,
) error {
	// Логируем начало операции
	s.logger.Info(
		"Uploading file to MinIO",
		zap.String("bucket", s.bucketName),
		zap.String("key", key),
		zap.Int64("size", size),
		zap.String("content_type", contentType),
	)

	// Настройки для загрузки объекта
	opts := minio.PutObjectOptions{
		// ContentType - MIME тип файла (важно для правильного отображения в браузере)
		ContentType: contentType,
		// UserMetadata - произвольные метаданные, которые можно прикрепить к объекту
		// Будут храниться вместе с файлом и доступны при скачивании
		UserMetadata: map[string]string{
			"uploaded_at": time.Now().Format(time.RFC3339),
		},

		// PartSize - размер части для multipart upload (для больших файлов)
		// -1 означает автоматический выбор размера
		// MinIO сам разделит файл на части, если он большой
		PartSize: 0, // 0 - автоматический выбор
	}

	// PutObject загружает файл в MinIO
	// Параметры:
	//   - ctx: контекст (для таймаутов и отмены)
	//   - bucketName: имя bucket
	//   - key: ключ объекта (путь в bucket)
	//   - reader: поток данных файла
	//   - size: размер файла (-1 если неизвестен, но лучше указывать)
	//   - opts: опции загрузки
	uploadInfo, err := s.client.PutObject(
		ctx,
		s.bucketName,
		key,
		reader,
		size,
		opts,
	)
	if err != nil {
		s.logger.Error(
			"Failed to upload file to MinIO",
			zap.String("key", key),
			zap.Error(err),
		)
		return fmt.Errorf("failed to upload to MinIO: %w", err)
	}

	// Логируем успешную загрузку с информацией о результате
	s.logger.Info(
		"File uploaded to MinIO successfully",
		zap.String("key", key),
		zap.Int64("size", uploadInfo.Size),
		zap.String("etag", uploadInfo.ETag), // ETag - хеш файла (для проверки целостности)
	)

	return nil

}

// Download реализует метод Download из интерфейса Storage
// Скачивает файл из MinIO и возвращает поток для чтения
func (s *MinIOStorage) Download(ctx context.Context, key string) (io.ReadCloser, error) {
	s.logger.Info(
		"Download file from MinIO",
		zap.String("bucket", s.bucketName),
		zap.String("key", key),
	)

	// GetObject возвращает объект из MinIO
	// Возвращает *minio.Object, который реализует io.ReadCloser
	object, err := s.client.GetObject(
		ctx,
		s.bucketName,
		key,
		minio.GetObjectOptions{}, // опции скачивания (пока пустые)
	)
	if err != nil {
		s.logger.Error(
			"Failed to get object from MinIO",
			zap.String("key", key),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to download from MinIO: %w", err)
	}

	// Проверяем, что объект действительно существует
	// GetObject не возвращает ошибку сразу, ошибка возникает при попытке чтения
	// Поэтому делаем Stat() для проверки
	_, err = object.Stat()
	if err != nil {
		object.Close()

		// Проверяем, является ли ошибка "объект не найден"
		errResponse := minio.ToErrorResponse(err)
		if errResponse.Code == noSuchKeyError {
			s.logger.Warn(
				"Object not found in MinIO",
				zap.String("key", key),
			)
			return nil, fmt.Errorf("object not found: %s", key)
		}

		s.logger.Error(
			"Failed to stat object in MinIO",
			zap.String("key", key),
			zap.Error(err),
		)
		return nil, fmt.Errorf("failed to verify object: %w", err)
	}

	s.logger.Info(
		"File downloaded from MinIO successfully",
		zap.String("key", key),
	)

	// Возвращаем object (он реализует io.ReadCloser)
	// Вызывающая сторона должна вызвать Close() после использования
	return object, nil
}

// Delete реализует метод Delete из интерфейса Storage
// Удаляет файл из MinIO
func (s *MinIOStorage) Delete(ctx context.Context, key string) error {
	s.logger.Info(
		"Deleting file from MinIO",
		zap.String("key", key),
	)

	// RemoveObject удаляет объект из bucket
	err := s.client.RemoveObject(
		ctx,
		s.bucketName,
		key,
		minio.RemoveObjectOptions{}, // опции удаления
	)
	if err != nil {
		s.logger.Error(
			"Failed to delete file from MinIO",
			zap.String("key", key),
			zap.Error(err),
		)
		return fmt.Errorf("failed to delete MinIO: %w", err)
	}

	s.logger.Info(
		"File deleted from MinIO successfully",
		zap.String("key", key),
	)

	return nil

}

// Exists реализует метод Exists из интерфейса Storage
// Проверяет существование файла в MinIO
func (s *MinIOStorage) Exists(ctx context.Context, key string) (bool, error) {
	// StatObject возвращает метаданные объекта
	// Если объект не существует - вернёт ошибку
	_, err := s.client.StatObject(
		ctx,
		s.bucketName,
		key,
		minio.GetObjectOptions{},
	)

	if err != nil {
		// Преобразуем ошибку MinIO в структуру для анализа
		errResponse := minio.ToErrorResponse(err)

		// Если ошибка "NoSuchKey" - объект не существует (это нормально, не ошибка)
		if errResponse.Code == noSuchKeyError {
			return false, nil
		}

		// Любая другая ошибка - это проблема
		s.logger.Error(
			"Failed to check object existence in MinIO",
			zap.String("key", key),
			zap.Error(err),
		)
		return false, fmt.Errorf("failed to check existence: %w", err)

	}

	// объект сущетвует
	return true, nil

}

// GetPresignedURL реализует метод GetPresignedURL из интерфейса Storage
// Генерирует временную подписанную ссылку для прямого доступа к файлу
func (s *MinIOStorage) GetPresignedURL(
	ctx context.Context,
	key string,
	expiry time.Duration,
) (string, error) {
	s.logger.Info(
		"Generating presigned URL",
		zap.String("bucket", s.bucketName),
		zap.String("key", key),
		zap.Duration("expiry", expiry),
	)

	// PresignedGetObject генерирует подписанный URL для скачивания
	// Параметры:
	//   - ctx: контекст
	//   - bucketName: имя bucket
	//   - key: ключ объекта
	//   - expiry: время жизни ссылки (максимум 7 дней для AWS S3)
	//   - reqParams: дополнительные параметры запроса (обычно nil)
	url, err := s.client.PresignedGetObject(
		ctx,
		s.bucketName,
		key,
		expiry,
		nil, // // Request parameters (например, можно указать response-content-type)
	)
	if err != nil {
		s.logger.Error(
			"Failed to generate presigned URL",
			zap.String("key", key),
			zap.Error(err),
		)
		return "", fmt.Errorf("failed to generate presigned URL: %w", err)
	}

	s.logger.Info(
		"Presigned URL generated successfully",
		zap.String("key", key),
		zap.String("url", url.String()),
	)

	return url.String(), nil
}
