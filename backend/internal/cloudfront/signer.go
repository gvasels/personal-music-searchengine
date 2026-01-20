// Package cloudfront provides CloudFront URL signing capabilities.
package cloudfront

import (
	"context"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"errors"
	"fmt"
	"strings"
	"time"
)

// Expiration bounds for signed URLs
const (
	MinExpiration = 5 * time.Minute   // Minimum URL expiration
	MaxExpiration = 7 * 24 * time.Hour // Maximum URL expiration (7 days)
)

// ErrExpirationTooShort is returned when the expiration duration is less than MinExpiration.
var ErrExpirationTooShort = errors.New("expiration duration too short (minimum 5 minutes)")

// ErrExpirationTooLong is returned when the expiration duration exceeds MaxExpiration.
var ErrExpirationTooLong = errors.New("expiration duration too long (maximum 7 days)")

// Signer provides CloudFront URL signing operations.
type Signer struct {
	domain     string
	keyPairID  string
	privateKey *rsa.PrivateKey
}

// NewSigner creates a new CloudFront signer.
func NewSigner(domain, keyPairID string, privateKeyPEM []byte) (*Signer, error) {
	block, _ := pem.Decode(privateKeyPEM)
	if block == nil {
		return nil, fmt.Errorf("failed to decode PEM block")
	}

	var privateKey *rsa.PrivateKey
	var err error

	switch block.Type {
	case "RSA PRIVATE KEY":
		privateKey, err = x509.ParsePKCS1PrivateKey(block.Bytes)
	case "PRIVATE KEY":
		key, e := x509.ParsePKCS8PrivateKey(block.Bytes)
		if e != nil {
			return nil, fmt.Errorf("failed to parse PKCS8 private key: %w", e)
		}
		var ok bool
		privateKey, ok = key.(*rsa.PrivateKey)
		if !ok {
			return nil, fmt.Errorf("private key is not RSA")
		}
	default:
		return nil, fmt.Errorf("unsupported PEM block type: %s", block.Type)
	}

	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	// Normalize domain - remove trailing slash and protocol if present
	domain = strings.TrimPrefix(domain, "https://")
	domain = strings.TrimPrefix(domain, "http://")
	domain = strings.TrimSuffix(domain, "/")

	return &Signer{
		domain:     domain,
		keyPairID:  keyPairID,
		privateKey: privateKey,
	}, nil
}

// GenerateSignedURL generates a CloudFront signed URL with a canned policy.
// The expiry duration must be between MinExpiration (5 minutes) and MaxExpiration (7 days).
func (s *Signer) GenerateSignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	// Validate expiration bounds
	if expiry < MinExpiration {
		return "", ErrExpirationTooShort
	}
	if expiry > MaxExpiration {
		return "", ErrExpirationTooLong
	}

	expiresAt := time.Now().Add(expiry).Unix()

	// Build the URL
	key = strings.TrimPrefix(key, "/")
	url := fmt.Sprintf("https://%s/%s", s.domain, key)

	// Create canned policy
	policy := fmt.Sprintf(`{"Statement":[{"Resource":"%s","Condition":{"DateLessThan":{"AWS:EpochTime":%d}}}]}`, url, expiresAt)

	// Sign the policy
	signature, err := s.signPolicy(policy)
	if err != nil {
		return "", fmt.Errorf("failed to sign policy: %w", err)
	}

	// Build signed URL with canned policy
	signedURL := fmt.Sprintf("%s?Expires=%d&Signature=%s&Key-Pair-Id=%s",
		url,
		expiresAt,
		encodeBase64ForURL(signature),
		s.keyPairID,
	)

	return signedURL, nil
}

// SignStreamURL generates a signed URL for HLS streaming.
func (s *Signer) SignStreamURL(ctx context.Context, userID, trackID string, expiry time.Duration) (string, error) {
	key := fmt.Sprintf("hls/%s/%s/master.m3u8", userID, trackID)
	return s.GenerateSignedURL(ctx, key, expiry)
}

// SignDownloadURL generates a signed URL for direct file download.
func (s *Signer) SignDownloadURL(ctx context.Context, s3Key string, expiry time.Duration) (string, error) {
	return s.GenerateSignedURL(ctx, s3Key, expiry)
}

// SignCoverArtURL generates a signed URL for cover art.
func (s *Signer) SignCoverArtURL(ctx context.Context, s3Key string, expiry time.Duration) (string, error) {
	return s.GenerateSignedURL(ctx, s3Key, expiry)
}

// signPolicy signs a policy using RSA-SHA1.
func (s *Signer) signPolicy(policy string) ([]byte, error) {
	hash := sha1.Sum([]byte(policy))
	signature, err := rsa.SignPKCS1v15(rand.Reader, s.privateKey, crypto.SHA1, hash[:])
	if err != nil {
		return nil, fmt.Errorf("failed to sign: %w", err)
	}
	return signature, nil
}

// encodeBase64ForURL encodes bytes to URL-safe base64.
func encodeBase64ForURL(data []byte) string {
	encoded := base64.StdEncoding.EncodeToString(data)
	// CloudFront requires URL-safe base64 with specific replacements
	encoded = strings.ReplaceAll(encoded, "+", "-")
	encoded = strings.ReplaceAll(encoded, "=", "_")
	encoded = strings.ReplaceAll(encoded, "/", "~")
	return encoded
}

// GetDomain returns the CloudFront distribution domain.
func (s *Signer) GetDomain() string {
	return s.domain
}
