package services

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"crypto/x509/pkix"
	"fmt"
	"log"
	"math/big"
	"os"
	"path/filepath"
	"time"

	"github.com/pavlo-v-chernykh/keystore-go/v4"
)

// GenerateKeystore creates a new Android keystore for a user in pure Go
func GenerateKeystore(userID string, password string) (string, error) {
	log.Printf("[KeystoreService] Starting generation for user: %s\n", userID)

	keystoreDir := filepath.Join("..", "lib", "keystores")
	log.Printf("[KeystoreService] Keystore directory: %s\n", keystoreDir)

	if _, err := os.Stat(keystoreDir); os.IsNotExist(err) {
		log.Printf("[KeystoreService] Directory doesn't exist, creating: %s\n", keystoreDir)
		if err := os.MkdirAll(keystoreDir, 0755); err != nil {
			log.Printf("[KeystoreService] ERROR creating directory: %v\n", err)
			return "", fmt.Errorf("failed to create keystore directory: %v", err)
		}
	}

	keystorePath := filepath.Join(keystoreDir, fmt.Sprintf("%s.keystore", userID))
	log.Printf("[KeystoreService] Target path: %s\n", keystorePath)

	if _, err := os.Stat(keystorePath); err == nil {
		log.Printf("[KeystoreService] Keystore already exists at %s, skipping generation\n", keystorePath)
		return keystorePath, nil
	}

	// Generate RSA private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return "", fmt.Errorf("failed to generate private key: %v", err)
	}

	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to marshal private key: %v", err)
	}

	serialNumberLimit := new(big.Int).Lsh(big.NewInt(1), 128)
	serialNumber, err := rand.Int(rand.Reader, serialNumberLimit)
	if err != nil {
		return "", fmt.Errorf("failed to generate serial number: %v", err)
	}

	tmpl := x509.Certificate{
		SerialNumber: serialNumber,
		Subject: pkix.Name{
			CommonName:         "RenaBuilder",
			Organization:       []string{"Rena"},
			OrganizationalUnit: []string{"Development"},
			Locality:           []string{"Global"},
			Province:           []string{"Global"},
			Country:            []string{"US"},
		},
		NotBefore:             time.Now().Add(-5 * time.Minute),
		NotAfter:              time.Now().Add(3650 * 24 * time.Hour),
		KeyUsage:              x509.KeyUsageDigitalSignature | x509.KeyUsageKeyEncipherment,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageClientAuth, x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	certDER, err := x509.CreateCertificate(rand.Reader, &tmpl, &tmpl, &privateKey.PublicKey, privateKey)
	if err != nil {
		return "", fmt.Errorf("failed to create certificate: %v", err)
	}

	keystoreStore := keystore.New()
	entry := keystore.PrivateKeyEntry{
		CreationTime: time.Now(),
		PrivateKey:   privateKeyBytes,
		CertificateChain: []keystore.Certificate{
			{
				Type:    "X.509",
				Content: certDER,
			},
		},
	}

	if err := keystoreStore.SetPrivateKeyEntry("rena_key", entry, []byte(password)); err != nil {
		return "", fmt.Errorf("failed to set private key entry: %v", err)
	}

	keystoreFile, err := os.Create(keystorePath)
	if err != nil {
		return "", fmt.Errorf("failed to create keystore file: %v", err)
	}
	defer keystoreFile.Close()

	if err := keystoreStore.Store(keystoreFile, []byte(password)); err != nil {
		return "", fmt.Errorf("failed to store keystore: %v", err)
	}

	log.Printf("[KeystoreService] Successfully generated keystore at %s\n", keystorePath)
	return keystorePath, nil
}
