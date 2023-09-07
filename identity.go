package ezlibp2p

import (
	"errors"
	"os"

	"github.com/libp2p/go-libp2p/core/crypto"
)

// PersistantIdentity stores ID in the file system.
type PersistantIdentity struct {
	fileName string
}

const DEFAULT_ID_FILE_NAME = "libp2p_id"

func NewPersistantIdentity() *PersistantIdentity {
	return NewPersistantIdentityWithFileName(DEFAULT_ID_FILE_NAME)
}

func NewPersistantIdentityWithFileName(fileName string) *PersistantIdentity {
	return &PersistantIdentity{
		fileName: fileName,
	}
}

var ErrNoIDFileFound = errors.New("No ID file found.")

// ReadPrivKey returns the previously written PrivKey, if any.
func (p *PersistantIdentity) ReadPrivKey() (crypto.PrivKey, error) {
	_, err := os.Stat(p.fileName)
	if os.IsNotExist(err) {
		return nil, ErrNoIDFileFound
	}

	keyStringBytes, err := os.ReadFile(p.fileName)
	if err != nil {
		return nil, err
	}

	keyBytes, err := crypto.ConfigDecodeKey(string(keyStringBytes))
	if err != nil {
		return nil, err
	}

	return crypto.UnmarshalPrivateKey(keyBytes)
}

// WritePrivKey writes new or overrides the previous PrivKey.
func (p *PersistantIdentity) WritePrivKey(key crypto.PrivKey) error {
	keyBytes, err := crypto.MarshalPrivateKey(key)
	if err != nil {
		return err
	}

	keyString := crypto.ConfigEncodeKey(keyBytes)
	return os.WriteFile(p.fileName, []byte(keyString), os.ModePerm)
}
