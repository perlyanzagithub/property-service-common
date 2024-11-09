package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/perlyanzagithub/property-service-common/dtos"
	"gorm.io/gorm"
	"io"
	"reflect"
)

func ConvertToDTO[T any](data interface{}) ([]T, error) {
	var dtos []T

	jsonData, err := json.Marshal(data)
	if err != nil {
		return nil, fmt.Errorf("error marshaling data to JSON: %v", err)
	}

	if jsonData[0] == '[' {
		if err := json.Unmarshal(jsonData, &dtos); err != nil {
			return nil, fmt.Errorf("error unmarshaling JSON to DTO slice: %v", err)
		}
	} else {
		var dto T
		if err := json.Unmarshal(jsonData, &dto); err != nil {
			return nil, fmt.Errorf("error unmarshaling JSON to DTO: %v", err)
		}
		dtos = append(dtos, dto)
	}

	return dtos, nil
}
func ConvertToDTOs[T any](data []interface{}) ([]T, error) {
	var dtos []T

	for _, item := range data {
		jsonData, err := json.Marshal(item)
		if err != nil {
			return nil, fmt.Errorf("error marshaling item to JSON: %v", err)
		}

		var dto T
		if err := json.Unmarshal(jsonData, &dto); err != nil {
			return nil, fmt.Errorf("error unmarshaling JSON to DTO: %v", err)
		}
		dtos = append(dtos, dto)
	}

	return dtos, nil
}
func TotalPage(totalData int64, size int) int64 {
	return (totalData + int64(size) - 1) / int64(size)
}
func loadSecretKey() ([]byte, error) {
	secretKey := "secretkeyproperty"
	return []byte(secretKey), nil
}

func createCipherBlock() (cipher.Block, error) {
	key, err := loadSecretKey()
	if err != nil {
		return nil, err
	}
	return aes.NewCipher(key)
}

func EncryptAES(plainText string) (string, error) {
	block, err := createCipherBlock()
	if err != nil {
		return "", err
	}

	cipherText := make([]byte, aes.BlockSize+len(plainText))
	iv := cipherText[:aes.BlockSize]

	if _, err := io.ReadFull(rand.Reader, iv); err != nil {
		return "", err
	}

	stream := cipher.NewCFBEncrypter(block, iv)
	stream.XORKeyStream(cipherText[aes.BlockSize:], []byte(plainText))

	return base64.StdEncoding.EncodeToString(cipherText), nil
}

func DecryptAES(cipherText string) (string, error) {
	block, err := createCipherBlock()
	if err != nil {
		return "", err
	}

	cipherTextBytes, err := base64.StdEncoding.DecodeString(cipherText)
	if err != nil {
		return "", err
	}

	if len(cipherTextBytes) < aes.BlockSize {
		return "", errors.New("cipher text too short")
	}

	iv := cipherTextBytes[:aes.BlockSize]
	cipherTextBytes = cipherTextBytes[aes.BlockSize:]

	stream := cipher.NewCFBDecrypter(block, iv)
	stream.XORKeyStream(cipherTextBytes, cipherTextBytes)

	return string(cipherTextBytes), nil
}
func CopyDTOToModel(dto interface{}, model interface{}) error {
	dtoVal := reflect.ValueOf(dto).Elem()
	modelVal := reflect.ValueOf(model).Elem()

	for i := 0; i < dtoVal.NumField(); i++ {
		dtoField := dtoVal.Type().Field(i)
		modelField := modelVal.FieldByName(dtoField.Name)

		if modelField.IsValid() && modelField.CanSet() {
			modelField.Set(dtoVal.Field(i))
		}
	}
	return nil
}
func ApplyFilters[T any](db *gorm.DB, request dtos.RequestDTO, out *[]T) (int64, int64, error) {
	var totalData int64

	if request.Filter != nil {
		for key, value := range request.Filter {
			db = db.Where(fmt.Sprintf("%s LIKE ?", key), "%"+value.(string)+"%")
		}
	}

	if err := db.Count(&totalData).Error; err != nil {
		return 0, 0, err
	}
	if err := db.Order(request.OrderBy).
		Offset((request.Page - 1) * request.Size).
		Limit(request.Size).
		Find(out).Error; err != nil {
		return 0, 0, err
	}

	return totalData, TotalPage(totalData, request.Size), nil
}
