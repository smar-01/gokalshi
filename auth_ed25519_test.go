package gokalshi

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// pkcs8PEM is key as a "BEGIN PRIVATE KEY" block: the header Kalshi's Ed25519 keys and PKCS#8 RSA keys share.
func pkcs8PEM(t *testing.T, key any) string {
	t.Helper()
	der, err := x509.MarshalPKCS8PrivateKey(key)
	require.NoError(t, err)
	return string(pem.EncodeToMemory(&pem.Block{Type: "PRIVATE KEY", Bytes: der}))
}

// The pre-sign text is timestamp + method + path without the query (docs.kalshi.com/getting_started/api_keys).
func presign(t *testing.T, h map[string]string, method, path string) []byte {
	t.Helper()
	ts := h["KALSHI-ACCESS-TIMESTAMP"]
	_, err := strconv.ParseInt(ts, 10, 64)
	require.NoError(t, err)
	return []byte(ts + method + strings.SplitN(path, "?", 2)[0])
}

func TestEd25519_SignAndVerifyWithPublicKey(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)

	creds, err := LoadCredentialsFromPEM("ed-key", pkcs8PEM(t, priv))
	require.NoError(t, err)
	require.Nil(t, creds.PrivateKey, "an Ed25519 key must not be read as RSA")
	require.NotNil(t, creds.Ed25519Key)

	path := "/trade-api/v2/portfolio/orders?status=resting&subaccount=1"
	h, err := creds.RequestHeaders("GET", path)
	require.NoError(t, err)
	assert.Equal(t, "ed-key", h["KALSHI-ACCESS-KEY"])

	sig, err := base64.StdEncoding.DecodeString(h["KALSHI-ACCESS-SIGNATURE"])
	require.NoError(t, err)
	assert.Len(t, sig, ed25519.SignatureSize)
	assert.True(t, ed25519.Verify(pub, presign(t, h, "GET", path), sig), "signature must verify against the public key")
	assert.False(t, ed25519.Verify(pub, []byte(h["KALSHI-ACCESS-TIMESTAMP"]+"GET"+path), sig), "the query must not be signed")
}

func TestEd25519_FromFileAndDeterministic(t *testing.T) {
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	p := filepath.Join(t.TempDir(), "kalshi.key")
	require.NoError(t, os.WriteFile(p, []byte(pkcs8PEM(t, priv)), 0o600))

	creds, err := LoadCredentials("ed-key", p)
	require.NoError(t, err)
	a, err := creds.SignText("1700000000000GET/trade-api/v2/portfolio/balance")
	require.NoError(t, err)
	b, err := NewEd25519Credentials("ed-key", priv).SignText("1700000000000GET/trade-api/v2/portfolio/balance")
	require.NoError(t, err)
	assert.Equal(t, a, b, "Ed25519 signatures are deterministic (RFC 8032)")
}

func TestRSA_PKCS8AndPKCS1_SignAndVerifyWithPublicKey(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	pkcs1 := string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: x509.MarshalPKCS1PrivateKey(key)}))
	for name, pemText := range map[string]string{"pkcs8": pkcs8PEM(t, key), "pkcs1": pkcs1} {
		t.Run(name, func(t *testing.T) {
			creds, err := LoadCredentialsFromPEM("rsa-key", pemText)
			require.NoError(t, err)
			require.NotNil(t, creds.PrivateKey)
			require.Nil(t, creds.Ed25519Key)

			path := "/trade-api/v2/portfolio/balance?subaccount=3"
			h, err := creds.RequestHeaders("GET", path)
			require.NoError(t, err)
			sig, err := base64.StdEncoding.DecodeString(h["KALSHI-ACCESS-SIGNATURE"])
			require.NoError(t, err)
			digest := sha256.Sum256(presign(t, h, "GET", path))
			require.NoError(t, rsa.VerifyPSS(&key.PublicKey, crypto.SHA256, digest[:], sig, &rsa.PSSOptions{SaltLength: rsa.PSSSaltLengthEqualsHash}))
		})
	}
}

// fixedReader replays the same bytes, so two RSA-PSS signatures with the same salt can be compared byte for byte.
type fixedReader struct{ r io.Reader }

func (f fixedReader) Read(p []byte) (int, error) { return f.r.Read(p) }

// RSA keys keep working byte for byte: SignText (the path RequestHeaders now takes) produces exactly what
// SignPSSText, the only signer before Ed25519, produced for the same key, text and random salt.
func TestRSA_SignTextIsSignPSSTextByteForByte(t *testing.T) {
	key, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)
	salt := bytes.Repeat([]byte{7}, 4096)
	text := "1700000000000POST/trade-api/v2/portfolio/events/orders"

	before := &Credentials{KeyID: "k", PrivateKey: key, Rand: fixedReader{bytes.NewReader(salt)}}
	after := &Credentials{KeyID: "k", PrivateKey: key, Rand: fixedReader{bytes.NewReader(salt)}}
	a, err := before.SignPSSText(text)
	require.NoError(t, err)
	b, err := after.SignText(text)
	require.NoError(t, err)
	assert.Equal(t, a, b)
}

func TestKeyType_FromParsedKeyNotHeader(t *testing.T) {
	// An Ed25519 key wrapped in a block labelled "RSA PRIVATE KEY" is still Ed25519: the label is never read.
	_, priv, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	der, err := x509.MarshalPKCS8PrivateKey(priv)
	require.NoError(t, err)
	mislabelled := string(pem.EncodeToMemory(&pem.Block{Type: "RSA PRIVATE KEY", Bytes: der}))
	creds, err := LoadCredentialsFromPEM("k", mislabelled)
	require.NoError(t, err)
	assert.NotNil(t, creds.Ed25519Key)
	assert.Nil(t, creds.PrivateKey)

	// A key type Kalshi does not take is refused with its type named.
	ec, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	require.NoError(t, err)
	_, err = LoadCredentialsFromPEM("k", pkcs8PEM(t, ec))
	require.Error(t, err)
	assert.Contains(t, err.Error(), "ecdsa")
}

func TestSignText_NoKey(t *testing.T) {
	_, err := (&Credentials{KeyID: "k"}).SignText("x")
	assert.Error(t, err)
}

// The WebSocket handshake signs the same way (ws_client.go uses RequestHeaders).
func TestEd25519_WebSocketPathSigned(t *testing.T) {
	pub, priv, err := ed25519.GenerateKey(rand.Reader)
	require.NoError(t, err)
	creds := NewEd25519Credentials("k", priv)
	h, err := creds.RequestHeaders("GET", "/trade-api/ws/v2")
	require.NoError(t, err)
	sig, _ := base64.StdEncoding.DecodeString(h["KALSHI-ACCESS-SIGNATURE"])
	assert.True(t, ed25519.Verify(pub, presign(t, h, "GET", "/trade-api/ws/v2"), sig))
	ts, _ := strconv.ParseInt(h["KALSHI-ACCESS-TIMESTAMP"], 10, 64)
	assert.InDelta(t, time.Now().UnixMilli(), ts, 5000)
}
