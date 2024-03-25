package ezlibp2p

import (
	"crypto/rand"
	"os"
	"path/filepath"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
)

const DEFAULT_ID_FILE_NAME = "libp2p_id"

//

func PersistentIdentity() (libp2p.Option, error) {
	currentDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	filePath := filepath.Join(currentDir, DEFAULT_ID_FILE_NAME)
	return PersistentIdentityFilePath(filePath)
}

func PersistentIdentityFilePath(filePath string) (libp2p.Option, error) {
	_, err := os.Stat(filePath)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if os.IsNotExist(err) {
		privKey, _, err := crypto.GenerateEd25519Key(rand.Reader)
		if err != nil {
			return nil, err
		}

		err = writePrivateKey(filePath, privKey)
		if err != nil {
			return nil, err
		}

		return libp2p.Identity(privKey), nil
	}

	privKey, err := readPrivateKey(filePath)
	if err != nil {
		return nil, err
	}

	return libp2p.Identity(privKey), nil
}

func readPrivateKey(fileName string) (crypto.PrivKey, error) {
	keyStringBytes, err := os.ReadFile(fileName)
	if err != nil {
		return nil, err
	}

	keyBytes, err := crypto.ConfigDecodeKey(string(keyStringBytes))
	if err != nil {
		return nil, err
	}

	return crypto.UnmarshalPrivateKey(keyBytes)
}

func writePrivateKey(fileName string, key crypto.PrivKey) error {
	keyBytes, err := crypto.MarshalPrivateKey(key)
	if err != nil {
		return err
	}

	keyString := crypto.ConfigEncodeKey(keyBytes)
	return os.WriteFile(fileName, []byte(keyString), os.ModePerm)
}
