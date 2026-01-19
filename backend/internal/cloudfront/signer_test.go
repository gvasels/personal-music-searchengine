package cloudfront

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// generateTestPrivateKey creates an RSA private key for testing.
func generateTestPrivateKey(t *testing.T) []byte {
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	require.NoError(t, err)

	privateKeyBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	pemBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	return pem.EncodeToMemory(pemBlock)
}

func TestNewSigner_Success(t *testing.T) {
	privateKeyPEM := generateTestPrivateKey(t)

	signer, err := NewSigner("d123456789.cloudfront.net", "KEYPAIRID123", privateKeyPEM)

	require.NoError(t, err)
	assert.NotNil(t, signer)
	assert.Equal(t, "d123456789.cloudfront.net", signer.domain)
	assert.Equal(t, "KEYPAIRID123", signer.keyPairID)
}

func TestNewSigner_NormalizeDomain(t *testing.T) {
	privateKeyPEM := generateTestPrivateKey(t)

	tests := []struct {
		input    string
		expected string
	}{
		{"https://d123.cloudfront.net/", "d123.cloudfront.net"},
		{"http://d123.cloudfront.net", "d123.cloudfront.net"},
		{"d123.cloudfront.net", "d123.cloudfront.net"},
	}

	for _, test := range tests {
		signer, err := NewSigner(test.input, "KEYPAIRID", privateKeyPEM)
		require.NoError(t, err)
		assert.Equal(t, test.expected, signer.domain)
	}
}

func TestNewSigner_InvalidPEM(t *testing.T) {
	_, err := NewSigner("d123.cloudfront.net", "KEYPAIRID", []byte("invalid"))

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "decode PEM")
}

func TestGenerateSignedURL_ValidSignature(t *testing.T) {
	ctx := context.Background()
	privateKeyPEM := generateTestPrivateKey(t)
	signer, err := NewSigner("d123.cloudfront.net", "KEYPAIRID123", privateKeyPEM)
	require.NoError(t, err)

	url, err := signer.GenerateSignedURL(ctx, "audio/user-123/track-456/file.mp3", time.Hour)

	require.NoError(t, err)
	assert.Contains(t, url, "https://d123.cloudfront.net/audio/user-123/track-456/file.mp3")
	assert.Contains(t, url, "Expires=")
	assert.Contains(t, url, "Signature=")
	assert.Contains(t, url, "Key-Pair-Id=KEYPAIRID123")
}

func TestGenerateSignedURL_Expiry(t *testing.T) {
	ctx := context.Background()
	privateKeyPEM := generateTestPrivateKey(t)
	signer, err := NewSigner("d123.cloudfront.net", "KEYPAIRID123", privateKeyPEM)
	require.NoError(t, err)

	before := time.Now().Add(time.Hour).Unix()
	url, err := signer.GenerateSignedURL(ctx, "test.mp3", time.Hour)
	after := time.Now().Add(time.Hour).Unix()

	require.NoError(t, err)

	// Extract expires value
	parts := strings.Split(url, "Expires=")
	require.Len(t, parts, 2)
	expiresStr := strings.Split(parts[1], "&")[0]

	var expires int64
	_, err = fmt.Sscanf(expiresStr, "%d", &expires)
	require.NoError(t, err)

	// Expires should be within the time range
	assert.GreaterOrEqual(t, expires, before)
	assert.LessOrEqual(t, expires, after)
}

func TestSignStreamURL_HLSPath(t *testing.T) {
	ctx := context.Background()
	privateKeyPEM := generateTestPrivateKey(t)
	signer, err := NewSigner("d123.cloudfront.net", "KEYPAIRID", privateKeyPEM)
	require.NoError(t, err)

	url, err := signer.SignStreamURL(ctx, "user-123", "track-456", time.Hour)

	require.NoError(t, err)
	assert.Contains(t, url, "hls/user-123/track-456/master.m3u8")
}

func TestSignDownloadURL_MediaPath(t *testing.T) {
	ctx := context.Background()
	privateKeyPEM := generateTestPrivateKey(t)
	signer, err := NewSigner("d123.cloudfront.net", "KEYPAIRID", privateKeyPEM)
	require.NoError(t, err)

	url, err := signer.SignDownloadURL(ctx, "audio/user-123/track-456/original.mp3", time.Hour)

	require.NoError(t, err)
	assert.Contains(t, url, "audio/user-123/track-456/original.mp3")
}

func TestEncodeBase64ForURL(t *testing.T) {
	// Test that URL-safe encoding works
	input := []byte{0xfb, 0xfe, 0x3f} // Will contain +, /, = in standard base64
	result := encodeBase64ForURL(input)

	// Should not contain URL-unsafe characters
	assert.NotContains(t, result, "+")
	assert.NotContains(t, result, "/")
	assert.NotContains(t, result, "=")

	// Should contain CloudFront-specific replacements
	// + -> -, / -> ~, = -> _
}

func TestGetDomain(t *testing.T) {
	privateKeyPEM := generateTestPrivateKey(t)
	signer, err := NewSigner("d123.cloudfront.net", "KEYPAIRID", privateKeyPEM)
	require.NoError(t, err)

	assert.Equal(t, "d123.cloudfront.net", signer.GetDomain())
}
