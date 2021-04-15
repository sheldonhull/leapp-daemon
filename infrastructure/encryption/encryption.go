package encryption

import (
  "crypto/aes"
  "crypto/cipher"
  "crypto/rand"
  "fmt"
  "github.com/denisbrodbeck/machineid"
  "io"
  "leapp_daemon/infrastructure/http/http_error"
)

type Encryption struct {}

func(encryption *Encryption) Encrypt(plainText string) (string, error) {
	key, err := getAesKey()
	if err != nil {
		return "", err
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonce := make([]byte, gcm.NonceSize())
	_, err = io.ReadFull(rand.Reader, nonce)
	if err != nil {
		return "", err
	}

	encryptedText := gcm.Seal(nonce, nonce, []byte(plainText), nil)
	return string(encryptedText), nil
}

func(encryption *Encryption) Decrypt(encryptedText string) (string, error) {
	key, err := getAesKey()
	if err != nil {
		return "", err
	}

	c, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	gcm, err := cipher.NewGCM(c)
	if err != nil {
		return "", err
	}

	nonceSize := gcm.NonceSize()
	encryptedTextByteSlice := []byte(encryptedText)

	if len(encryptedTextByteSlice) <= nonceSize {
		return "", http_error.NewInternalServerError(fmt.Errorf("encrypted file's length is shorter than the expected nonce size"))
	}

	nonce, ciphertext := encryptedTextByteSlice[:nonceSize], encryptedTextByteSlice[nonceSize:]

	plainText, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plainText), nil
}

// TODO: move encryption key to domain layer
// TODO: read random-generated key from keychain
func getAesKey() ([]byte, error) {
  machineId, err := getMachineId()
  if err != nil { return nil, err }
  machineIdRuneSlice := []rune(machineId)
  machineIdRuneSlice = append(machineIdRuneSlice[0:8], machineIdRuneSlice[8+1:]...)
  machineIdRuneSlice = append(machineIdRuneSlice[0:12], machineIdRuneSlice[12+1:]...)
  machineIdRuneSlice = append(machineIdRuneSlice[0:16], machineIdRuneSlice[16+1:]...)
  machineIdRuneSlice = append(machineIdRuneSlice[0:20], machineIdRuneSlice[20+1:]...)
  return []byte(string(machineIdRuneSlice)), nil
}

func getMachineId() (string, error) {
  id, err := machineid.ID()
  if err != nil { return "", err }
  return id, nil
}
