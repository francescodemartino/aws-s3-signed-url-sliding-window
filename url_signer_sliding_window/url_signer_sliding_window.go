package url_signer_sliding_window

import (
	"github.com/francescodemartino/aws-s3-signed-url-sliding-window/url_signer"
	"time"
)

type S3UrlSignerSlidingWindow struct {
	s3UrlSigner  *url_signer.S3UrlSigner
	slidingTime  time.Duration
	residualTime time.Duration
}

func NewS3UrlSignerSlidingWindow(s3UrlSigner *url_signer.S3UrlSigner, slidingTime time.Duration, residualTime time.Duration) *S3UrlSignerSlidingWindow {
	return &S3UrlSignerSlidingWindow{
		s3UrlSigner:  s3UrlSigner,
		slidingTime:  slidingTime,
		residualTime: residualTime,
	}
}

func (s S3UrlSignerSlidingWindow) GetS3SignedUrl(key string) string {
	now := time.Now().Unix()
	slidingTime := int64(s.slidingTime.Seconds())
	residualTime := int64(s.residualTime.Seconds())
	slotTime := ((now+residualTime)/slidingTime)*slidingTime - residualTime
	slotTimeEnd := ((now+residualTime)/slidingTime)*slidingTime + slidingTime
	expiration := slotTimeEnd - slotTime
	start := time.Unix(slotTime, 0).UTC()

	return s.s3UrlSigner.GetS3SignedUrl(key, start, time.Duration(expiration)*time.Second)
}
