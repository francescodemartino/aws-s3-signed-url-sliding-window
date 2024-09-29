package url_signer

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"net/url"
	"strconv"
	"time"
)

type S3UrlSigner struct {
	awsIdAccessKey     string
	awsSecretAccessKey string
	awsRegion          string
	bucket             string
}

func NewS3UrlSigner(awsIdAccessKey string, awsSecretAccessKey string, awsRegion string, bucket string) *S3UrlSigner {
	return &S3UrlSigner{
		awsIdAccessKey:     awsIdAccessKey,
		awsSecretAccessKey: awsSecretAccessKey,
		awsRegion:          awsRegion,
		bucket:             bucket,
	}
}

func (s S3UrlSigner) GetS3SignedUrl(key string, start time.Time, expiration time.Duration) string {
	canonicalRequest := s.getCanonicalRequest(key, start, expiration)
	canonicalRequestHash := sha256Hex(canonicalRequest)
	stringToSign := s.getStringToSign(canonicalRequestHash, start)
	signature := s.getSignature(stringToSign, start)
	return s.getFullUrl(key, start, expiration, signature)
}

func (s S3UrlSigner) getFullUrl(key string, start time.Time, expiration time.Duration, signature string) string {
	return "https://" + s.bucket + ".s3." + s.awsRegion + ".amazonaws.com" +
		key +
		"?X-Amz-Algorithm=AWS4-HMAC-SHA256" +
		"&X-Amz-Credential=" + url.QueryEscape(s.getCredential(start)) +
		"&X-Amz-Date=" + start.Format("20060102T150405Z") +
		"&X-Amz-Expires=" + strconv.Itoa(int(expiration.Seconds())) +
		"&X-Amz-SignedHeaders=host" +
		"&X-Amz-Signature=" + signature
}

func (s S3UrlSigner) getHost() string {
	return s.bucket + ".s3." + s.awsRegion + ".amazonaws.com"
}

func (s S3UrlSigner) getCredential(start time.Time) string {
	return s.awsIdAccessKey + "/" + start.Format("20060102") + "/" + s.awsRegion + "/s3/aws4_request"
}

func (s S3UrlSigner) getCanonicalRequest(key string, start time.Time, expiration time.Duration) string {
	host := "host:" + s.getHost()
	return "GET\n" +
		key + "\n" +
		s.getCanonicalQuery(start, expiration) + "\n" +
		host + "\n\n" +
		"host\n" +
		"UNSIGNED-PAYLOAD"
}

func (s S3UrlSigner) getCanonicalQuery(start time.Time, expiration time.Duration) string {
	return "X-Amz-Algorithm=AWS4-HMAC-SHA256&X-Amz-Credential=" +
		url.QueryEscape(s.getCredential(start)) +
		"&X-Amz-Date=" + start.Format("20060102T150405Z") + "&X-Amz-Expires=" + strconv.Itoa(int(expiration.Seconds())) + "&X-Amz-SignedHeaders=host"
}

func (s S3UrlSigner) getStringToSign(canonicalRequestHash string, start time.Time) string {
	return "AWS4-HMAC-SHA256\n" +
		start.Format("20060102T150405Z") + "\n" +
		start.Format("20060102") + "/" + s.awsRegion + "/s3/aws4_request\n" +
		canonicalRequestHash
}

func (s S3UrlSigner) getSignature(stringToSign string, start time.Time) string {
	kDate := hmacSHA256([]byte("AWS4"+s.awsSecretAccessKey), []byte(start.Format("20060102")))
	kRegion := hmacSHA256(kDate, []byte(s.awsRegion))
	kService := hmacSHA256(kRegion, []byte("s3"))
	kSigning := hmacSHA256(kService, []byte("aws4_request"))
	signature := hmacSHA256(kSigning, []byte(stringToSign))
	return hex.EncodeToString(signature)
}

func hmacSHA256(key, data []byte) []byte {
	h := hmac.New(sha256.New, key)
	h.Write(data)
	return h.Sum(nil)
}

func sha256Hex(data string) string {
	hash := sha256.New()
	hash.Write([]byte(data))
	return hex.EncodeToString(hash.Sum(nil))
}
