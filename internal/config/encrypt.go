package config

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"golang.org/x/crypto/scrypt"
	"io"
)

// 加密密钥和盐值（在实际应用中应该从环境变量或配置文件中获取）
var (
	encryptionKey = []byte("godbmodeler-encryption-key-12345")
	salt          = []byte("godbmodeler-salt-12345")
)

// EncryptPassword 加密密码
func EncryptPassword(password string) (string, error) {
	if password == "" {
		return "", nil
	}

	// 使用scrypt生成密钥
	key, err := scrypt.Key(encryptionKey, salt, 32768, 8, 1, 32)
	if err != nil {
		return "", err
	}

	// 创建加密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 创建GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 创建随机数
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", err
	}

	// 加密
	ciphertext := gcm.Seal(nonce, nonce, []byte(password), nil)

	// 编码为Base64
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// DecryptPassword 解密密码
func DecryptPassword(encrypted string) (string, error) {
	if encrypted == "" {
		return "", nil
	}

	// 解码Base64
	ciphertext, err := base64.StdEncoding.DecodeString(encrypted)
	if err != nil {
		return "", err
	}

	// 使用scrypt生成密钥
	key, err := scrypt.Key(encryptionKey, salt, 32768, 8, 1, 32)
	if err != nil {
		return "", err
	}

	// 创建加密块
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", err
	}

	// 创建GCM模式
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", err
	}

	// 检查长度
	if len(ciphertext) < gcm.NonceSize() {
		return "", errors.New("密文太短")
	}

	// 提取nonce
	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]

	// 解密
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", err
	}

	return string(plaintext), nil
}

// EncryptConnectionPassword 加密连接配置中的密码
func EncryptConnectionPassword(config ConnectionConfig) (ConnectionConfig, error) {
	if config.Password == "" {
		return config, nil
	}

	// 加密密码
	encrypted, err := EncryptPassword(config.Password)
	if err != nil {
		return config, err
	}

	// 更新配置
	config.Password = encrypted

	return config, nil
}

// DecryptConnectionPassword 解密连接配置中的密码
func DecryptConnectionPassword(config ConnectionConfig) (ConnectionConfig, error) {
	if config.Password == "" {
		return config, nil
	}

	// 解密密码
	decrypted, err := DecryptPassword(config.Password)
	if err != nil {
		return config, err
	}

	// 更新配置
	config.Password = decrypted

	return config, nil
}
