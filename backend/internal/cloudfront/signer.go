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

// SignedURLOptions contains options for generating signed URLs.
type SignedURLOptions struct {
	// ForceDownload adds response-content-disposition header to force browser download
	ForceDownload bool
	// FileName is the filename to use when ForceDownload is true
	FileName string
}

// GenerateSignedURL generates a CloudFront signed URL with a canned policy.
// The expiry duration must be between MinExpiration (5 minutes) and MaxExpiration (7 days).
func (s *Signer) GenerateSignedURL(ctx context.Context, key string, expiry time.Duration) (string, error) {
	return s.GenerateSignedURLWithOptions(ctx, key, expiry, nil)
}

// GenerateSignedURLWithOptions generates a CloudFront signed URL with options.
func (s *Signer) GenerateSignedURLWithOptions(ctx context.Context, key string, expiry time.Duration, opts *SignedURLOptions) (string, error) {
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
	baseURL := fmt.Sprintf("https://%s/%s", s.domain, key)

	// Add query parameters if needed
	resourceURL := baseURL
	queryParams := ""
	if opts != nil && opts.ForceDownload {
		filename := opts.FileName
		if filename == "" {
			// Extract filename from key
			parts := strings.Split(key, "/")
			filename = parts[len(parts)-1]
		}
		// URL-encode the filename for Content-Disposition header
		encodedFilename := strings.ReplaceAll(filename, " ", "%20")
		encodedFilename = strings.ReplaceAll(encodedFilename, "\"", "%22")
		queryParams = fmt.Sprintf("response-content-disposition=attachment%%3B%%20filename%%3D%%22%s%%22", encodedFilename)
		resourceURL = baseURL + "?" + queryParams
	}

	// Create canned policy - use URL with query params for signature
	policy := fmt.Sprintf(`{"Statement":[{"Resource":"%s","Condition":{"DateLessThan":{"AWS:EpochTime":%d}}}]}`, resourceURL, expiresAt)

	// Sign the policy
	signature, err := s.signPolicy(policy)
	if err != nil {
		return "", fmt.Errorf("failed to sign policy: %w", err)
	}

	// Build signed URL
	var signedURL string
	if queryParams != "" {
		signedURL = fmt.Sprintf("%s?%s&Expires=%d&Signature=%s&Key-Pair-Id=%s",
			baseURL,
			queryParams,
			expiresAt,
			encodeBase64ForURL(signature),
			s.keyPairID,
		)
	} else {
		signedURL = fmt.Sprintf("%s?Expires=%d&Signature=%s&Key-Pair-Id=%s",
			baseURL,
			expiresAt,
			encodeBase64ForURL(signature),
			s.keyPairID,
		)
	}

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

// GenerateSignedDownloadURL generates a signed URL with Content-Disposition: attachment header
// to force the browser to download the file instead of playing it.
func (s *Signer) GenerateSignedDownloadURL(ctx context.Context, key string, expiry time.Duration, filename string) (string, error) {
	opts := &SignedURLOptions{
		ForceDownload: true,
		FileName:      filename,
	}
	return s.GenerateSignedURLWithOptions(ctx, key, expiry, opts)
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
