package enc

import "fmt"

type ErrNewEncryptor struct {
	KeyFilePath string
	Err         error
}

func (e ErrNewEncryptor) Error() string {
	return fmt.Sprintf("Failed to create cipher from %s: %v", e.KeyFilePath, e.Err)
}

type ErrDecrypt struct {
	Err error
}

func (e ErrDecrypt) Error() string {
	return fmt.Sprintf("Failed to encrypt: %v", e.Err)
}
