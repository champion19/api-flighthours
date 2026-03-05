package idencoder

import (
	"errors"
	"fmt"
	"strings"

	"github.com/champion19/api-flighthours/platform/logger"
	"github.com/google/uuid"
	hashids "github.com/speps/go-hashids/v2"
)

type HashidsEncoder struct {
	hashData *hashids.HashIDData
	logger   logger.Logger
}

type Config struct {
	Secret    string
	MinLength int
}

func NewHashidsEncoder(cfg Config, log logger.Logger) (*HashidsEncoder, error) {
	if cfg.Secret == "" {
		return nil, fmt.Errorf("secret no puede estar vacío")
	}

	if log != nil && cfg.MinLength == 36 {
		log.Warn(logger.LogIDEncoderMinLengthWarn)
	}

	hd := hashids.NewData()
	hd.Salt = cfg.Secret
	hd.MinLength = cfg.MinLength

	hd.Alphabet = "abcdefghijkmnopqrstuvwxyzABCDEFGHJKLMNPQRSTUVWXYZ123456789"

	return &HashidsEncoder{
		hashData: hd,
		logger:   log,
	}, nil
}

func (e *HashidsEncoder) Encode(uuidStr string) (string, error) {

	parsedUUID, err := uuid.Parse(uuidStr)
	if err != nil {
		if e.logger != nil {
			e.logger.Error(logger.LogIDEncoderInvalidUUID, "error", err, "uuid", uuidStr)
		}
		return "", err
	}

	uuidBytes := parsedUUID[:]

	numbers := make([]int, 0, 8)
	for i := 0; i < len(uuidBytes); i += 2 {
		num := int(uuidBytes[i])<<8 | int(uuidBytes[i+1])
		numbers = append(numbers, num)
	}

	h, err := hashids.NewWithData(e.hashData)
	if err != nil {
		if e.logger != nil {
			e.logger.Error(logger.LogIDEncoderHashidsCreate, "error", err)
		}
		return "", err
	}

	encoded, err := h.Encode(numbers)
	if err != nil {
		if e.logger != nil {
			e.logger.Error(logger.LogIDEncoderEncodingError, "error", err, "uuid", uuidStr)
		}
		return "", err
	}

	return encoded, nil
}

func (e *HashidsEncoder) Decode(encoded string) (string, error) {
	if encoded == "" {
		err := errors.New("ID ofuscado no puede estar vacío")
		e.logError(logger.LogIDEncoderEmptyID, err, "")
		return "", err
	}

	h, err := hashids.NewWithData(e.hashData)
	if err != nil {
		e.logError(logger.LogIDEncoderHashidsCreate, err, "")
		return "", err
	}

	numbers, err := h.DecodeWithError(encoded)
	if err != nil {
		e.logError(logger.LogIDEncoderDecodingError, err, encoded)
		return "", err
	}

	if len(numbers) != 8 {
		err := errors.New("ID ofuscado tiene formato incorrecto")
		e.logError(logger.LogIDEncoderInvalidFormat, err, encoded)
		return "", err
	}

	uuidBytes := make([]byte, 16)
	for i, num := range numbers {
		uuidBytes[i*2] = byte(num >> 8)
		uuidBytes[i*2+1] = byte(num & 0xFF)
	}

	parsedUUID, err := uuid.FromBytes(uuidBytes)
	if err != nil {
		e.logError(logger.LogIDEncoderUUIDError, err, encoded)
		return "", err
	}

	return parsedUUID.String(), nil
}

// logError logs an error if the logger is available.
func (e *HashidsEncoder) logError(msg string, err error, encoded string) {
	if e.logger == nil {
		return
	}
	if encoded != "" {
		e.logger.Error(msg, "error", err, "encoded", encoded)
	} else {
		e.logger.Error(msg, "error", err)
	}
}

func (e *HashidsEncoder) MustEncode(uuidStr string) string {
	encoded, err := e.Encode(uuidStr)
	if err != nil {
		if e.logger != nil {
			e.logger.Error(logger.LogIDEncoderEncodingError, "error", err, "uuid", uuidStr)
		}
	}
	return encoded
}

func (e *HashidsEncoder) IsValidEncoded(encoded string) bool {
	_, err := e.Decode(encoded)
	return err == nil
}
func IsUUID(str string) bool {
	str = strings.ToLower(str)
	_, err := uuid.Parse(str)
	return err == nil
}
