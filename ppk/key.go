package ppk

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"gopkg.in/src-d/go-git.v4/plumbing/transport/ssh"
)

// GetKey gets a private key string and returns a *ssh.PublicKeys and an error
func GetKey(key string) (*ssh.PublicKeys, error) {

	publicKey, err := ssh.NewPublicKeys("git", []byte(key), "")
	if err != nil {
		return nil, errors.New("creating ssh auth method: " + err.Error())
	}

	return publicKey, nil

}

// CreateKey creates an rsa keypair and returns them as a byte array
func CreateKey() ([]byte, []byte, error) {

	bitSize := 4096

	// Generate RSA key.
	key, err := rsa.GenerateKey(rand.Reader, bitSize)
	if err != nil {
		return nil, nil, err
	}

	// Extract public component.
	pub := key.Public()

	// Encode private key to PKCS#1 ASN.1 PEM.
	privateKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(key),
		},
	)

	// Encode public key to PKCS#1 ASN.1 PEM.
	pubKeyPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: x509.MarshalPKCS1PublicKey(pub.(*rsa.PublicKey)),
		},
	)

	return privateKeyPEM, pubKeyPEM, nil

}

// CreateSecret creates a secret with a key passed as an argument
func CreateSecret(privateKey []byte, publicKey []byte) *corev1.Secret {

	secret := corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "airflow-deploy-keys",
			Namespace: "dev-data",
		},
		Data: map[string][]byte{
			"privateKey": privateKey,
			"publicKey":  publicKey,
		},
		Type: "Opaque",
	}
	return &secret
}
