package gokalshi

import (
	"crypto"
	"crypto/ed25519"
	cryptorand "crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
	"io"
	"os"
	"strconv"
	"strings"
	"time"
)

// Credentials holds the key used for API authentication. Kalshi takes two key types (docs.kalshi.com/getting_started/
// api_keys, checked 2026-11-11): Ed25519, the recommended one since 2026-09-24 and the web app's default since
// 2026-10-01, and RSA. Exactly one of PrivateKey and Ed25519Key is set; Ed25519Key wins if both are.
type Credentials struct {
	KeyID string
	// PrivateKey is an RSA key: requests are signed with RSA-PSS, SHA-256, salt length equal to the hash.
	PrivateKey *rsa.PrivateKey
	// Ed25519Key is an Ed25519 key: requests are signed with Ed25519 (RFC 8032) over the same text, no hash first.
	Ed25519Key ed25519.PrivateKey
	// Rand is the random source for PSS signing. Defaults to crypto/rand.Reader. Ed25519 needs none.
	Rand io.Reader
}

// LoadCredentials reads a private key (Ed25519 or RSA) from a PEM file and returns Credentials.
func LoadCredentials(keyID, privateKeyPath string) (*Credentials, error) {
	pemData, err := os.ReadFile(privateKeyPath)
	if err != nil {
		return nil, &AuthError{Op: "load_credentials", Err: fmt.Errorf("failed to read private key file %s: %w", privateKeyPath, err)}
	}
	creds, err := parseCredentials(keyID, pemData)
	if err != nil {
		return nil, &AuthError{Op: "load_credentials", Err: fmt.Errorf("%s: %w", privateKeyPath, err)}
	}
	return creds, nil
}

// LoadCredentialsFromPEM parses a private key (Ed25519 or RSA) from a PEM-encoded string and returns Credentials. This
// is useful when the key is stored in an environment variable rather than a file (e.g. CI/CD secrets).
func LoadCredentialsFromPEM(keyID, pemString string) (*Credentials, error) {
	creds, err := parseCredentials(keyID, []byte(pemString))
	if err != nil {
		return nil, &AuthError{Op: "load_credentials_pem", Err: fmt.Errorf("PEM string: %w", err)}
	}
	return creds, nil
}

// parseCredentials picks the signing algorithm from the parsed key, never from the PEM header: a PKCS#8 block
// ("BEGIN PRIVATE KEY") holds an Ed25519 key or an RSA key alike, and a PKCS#1 block ("BEGIN RSA PRIVATE KEY") is RSA.
func parseCredentials(keyID string, pemData []byte) (*Credentials, error) {
	block, _ := pem.Decode(pemData)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}
	key, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		pkcs1Key, pkcs1Err := x509.ParsePKCS1PrivateKey(block.Bytes)
		if pkcs1Err != nil {
			return nil, fmt.Errorf("failed to parse private key: %w (also tried PKCS1: %v)", err, pkcs1Err)
		}
		key = pkcs1Key
	}
	switch k := key.(type) {
	case *rsa.PrivateKey:
		return &Credentials{KeyID: keyID, PrivateKey: k, Rand: cryptorand.Reader}, nil
	case ed25519.PrivateKey:
		return &Credentials{KeyID: keyID, Ed25519Key: k}, nil
	default:
		return nil, fmt.Errorf("private key is %T; Kalshi takes RSA or Ed25519", key)
	}
}

// RequestHeaders generates the authentication headers required by the Kalshi API.
func (c *Credentials) RequestHeaders(method, path string) (map[string]string, error) {
	timestampMs := strconv.FormatInt(time.Now().UnixMilli(), 10)

	// Strip query parameters from path for signing
	pathOnly := strings.SplitN(path, "?", 2)[0]
	msgString := timestampMs + method + pathOnly

	signature, err := c.SignText(msgString)
	if err != nil {
		return nil, &AuthError{Op: "request_headers", Err: fmt.Errorf("failed to sign request: %w", err)}
	}

	return map[string]string{
		"Content-Type":            "application/json",
		"KALSHI-ACCESS-KEY":       c.KeyID,
		"KALSHI-ACCESS-SIGNATURE": signature,
		"KALSHI-ACCESS-TIMESTAMP": timestampMs,
	}, nil
}

// SignText signs text with the credentials' key and returns the standard base64 signature Kalshi expects in
// KALSHI-ACCESS-SIGNATURE: Ed25519 over the text itself, or RSA-PSS SHA-256 (SignPSSText).
func (c *Credentials) SignText(text string) (string, error) {
	if c.Ed25519Key != nil {
		if len(c.Ed25519Key) != ed25519.PrivateKeySize {
			return "", &AuthError{Op: "sign", Err: fmt.Errorf("Ed25519 key is %d bytes, want %d", len(c.Ed25519Key), ed25519.PrivateKeySize)}
		}
		return base64.StdEncoding.EncodeToString(ed25519.Sign(c.Ed25519Key, []byte(text))), nil
	}
	return c.SignPSSText(text)
}

// SignPSSText signs the given text using RSA-PSS with SHA-256.
func (c *Credentials) SignPSSText(text string) (string, error) {
	if c.PrivateKey == nil {
		return "", &AuthError{Op: "sign", Err: fmt.Errorf("no RSA private key")}
	}
	message := []byte(text)
	hash := sha256.Sum256(message)

	signature, err := rsa.SignPSS(
		c.randReader(),
		c.PrivateKey,
		crypto.SHA256,
		hash[:],
		&rsa.PSSOptions{
			SaltLength: rsa.PSSSaltLengthEqualsHash,
		},
	)
	if err != nil {
		return "", &AuthError{Op: "sign", Err: fmt.Errorf("RSA-PSS sign failed: %w", err)}
	}

	return base64.StdEncoding.EncodeToString(signature), nil
}

// NewCredentials creates Credentials from a pre-loaded RSA private key.
// Use this when you have already parsed the key yourself.
func NewCredentials(keyID string, key *rsa.PrivateKey) *Credentials {
	return &Credentials{
		KeyID:      keyID,
		PrivateKey: key,
		Rand:       cryptorand.Reader,
	}
}

// NewEd25519Credentials creates Credentials from a pre-loaded Ed25519 private key.
func NewEd25519Credentials(keyID string, key ed25519.PrivateKey) *Credentials {
	return &Credentials{KeyID: keyID, Ed25519Key: key}
}

func (c *Credentials) randReader() io.Reader {
	if c.Rand != nil {
		return c.Rand
	}
	return cryptorand.Reader
}
