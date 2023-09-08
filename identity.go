package ezlibp2p

import (
	"crypto/rand"
	"os"

	"github.com/libp2p/go-libp2p"
	"github.com/libp2p/go-libp2p/core/crypto"
)

const DEFAULT_ID_FILE_NAME = "libp2p_id"

func PersistantIdentity() (libp2p.Option, error) {
	return PersistantIdentityFileName(DEFAULT_ID_FILE_NAME)
}

func PersistantIdentityFileName(fileName string) (libp2p.Option, error) {
	_, err := os.Stat(fileName)
	if err != nil && !os.IsNotExist(err) {
		return nil, err
	}

	if os.IsNotExist(err) {
		privKey, _, err := crypto.GenerateEd25519Key(rand.Reader)
		if err != nil {
			return nil, err
		}

		err = writePrivKey(fileName, privKey)
		if err != nil {
			return nil, err
		}

		return libp2p.Identity(privKey), nil
	}

	privKey, err := readPrivKey(fileName)
	if err != nil {
		return nil, err
	}

	return libp2p.Identity(privKey), nil

}

func readPrivKey(fileName string) (crypto.PrivKey, error) {
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

func writePrivKey(fileName string, key crypto.PrivKey) error {
	keyBytes, err := crypto.MarshalPrivateKey(key)
	if err != nil {
		return err
	}

	keyString := crypto.ConfigEncodeKey(keyBytes)
	return os.WriteFile(fileName, []byte(keyString), os.ModePerm)
}
