# CloudFront Package - CLAUDE.md

## Overview

CloudFront URL signing package for generating signed URLs for media streaming and downloads. Implements RSA-SHA1 canned policy signing compatible with AWS CloudFront.

## File Descriptions

| File | Purpose |
|------|---------|
| `signer.go` | CloudFront URL signer implementation |
| `signer_test.go` | Unit tests for URL signing |

## Key Types

### Signer
Provides CloudFront URL signing operations.
```go
type Signer struct {
    domain     string
    keyPairID  string
    privateKey *rsa.PrivateKey
}
```

## Functions

| Function | Signature | Description |
|----------|-----------|-------------|
| `NewSigner` | `func NewSigner(domain, keyPairID string, privateKeyPEM []byte) (*Signer, error)` | Creates new CloudFront signer |
| `GenerateSignedURL` | `func (s *Signer) GenerateSignedURL(ctx, key string, expiry time.Duration) (string, error)` | Generates signed URL with canned policy |
| `SignStreamURL` | `func (s *Signer) SignStreamURL(ctx, userID, trackID string, expiry time.Duration) (string, error)` | Signs HLS master playlist URL |
| `SignDownloadURL` | `func (s *Signer) SignDownloadURL(ctx, s3Key string, expiry time.Duration) (string, error)` | Signs direct download URL |
| `SignCoverArtURL` | `func (s *Signer) SignCoverArtURL(ctx, s3Key string, expiry time.Duration) (string, error)` | Signs cover art URL |
| `GetDomain` | `func (s *Signer) GetDomain() string` | Returns the CloudFront domain |

## Usage Example

```go
import (
    "github.com/gvasels/personal-music-searchengine/internal/cloudfront"
)

func main() {
    // Load private key from Secrets Manager or file
    privateKeyPEM := loadPrivateKey()

    signer, err := cloudfront.NewSigner(
        "d123456789.cloudfront.net",
        "KEYPAIRID123",
        privateKeyPEM,
    )
    if err != nil {
        log.Fatal(err)
    }

    // Sign a stream URL (24 hour expiry)
    url, err := signer.SignStreamURL(ctx, "user-123", "track-456", 24*time.Hour)

    // Sign a download URL
    downloadURL, err := signer.SignDownloadURL(ctx, "audio/user-123/track-456/file.mp3", 24*time.Hour)
}
```

## URL Structure

Signed URLs follow CloudFront canned policy format:
```
https://{domain}/{key}?Expires={unix_timestamp}&Signature={url_safe_base64}&Key-Pair-Id={key_pair_id}
```

### HLS Stream URLs
```
https://d123.cloudfront.net/hls/{userId}/{trackId}/master.m3u8?Expires=...
```

### Download URLs
```
https://d123.cloudfront.net/audio/{userId}/{trackId}/original.mp3?Expires=...
```

## Security

- RSA-SHA1 signature algorithm (CloudFront requirement)
- URL-safe Base64 encoding with CloudFront-specific character replacements:
  - `+` → `-`
  - `/` → `~`
  - `=` → `_`
- Private key should be stored in AWS Secrets Manager
- Key pair ID is public and can be stored in environment variables

## Dependencies

| Package | Purpose |
|---------|---------|
| `crypto/rsa` | RSA key handling |
| `crypto/x509` | Key parsing |
| `encoding/pem` | PEM block decoding |
| `crypto/sha1` | Policy hashing |

No external dependencies - uses standard library only.

## Testing

Tests use dynamically generated RSA keys:
- Key generation with `rsa.GenerateKey`
- Signature verification via URL structure checks
- Domain normalization tests
- PEM parsing error handling
