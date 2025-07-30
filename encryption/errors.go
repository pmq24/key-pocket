package encryption

import "fmt"

type ErrNewCipher struct {
	KeyFilePath string
	Err         error
}

func (e ErrNewCipher) Error() string {
	return fmt.Sprintf("Failed to create cipher from %s: %v", e.KeyFilePath, e.Err)
}

type ErrDecrypt struct {
	Err error
}

func (e ErrDecrypt) Error() string {
	return fmt.Sprintf("Failed to encrypt: %v", e.Err)
}
