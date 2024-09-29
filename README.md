# AWS S3 Signed URL Generator

## Project Purpose

The purpose of this project is to generate signed URLs for AWS S3 storage with two key enhancements over the standard method:

1. The ability to specify not only the expiration time but also the **start date** from which the validity of the link begins. In the default method, the start time is determined by the current time, but this project allows you to define the exact moment (past or future) when the signed URL becomes valid. This provides the flexibility to generate links that can be used at any desired time, including both past and future dates.

2. The project introduces the concept of **sliding window URL generation**, which is perhaps its most powerful feature. With this, you can generate the same signed URL for repeated access within a predefined sliding window period. The key advantage here is that for the same bucket and object (key), the URL remains identical throughout the window duration.

### Key Benefits:
- **Consistent URL generation in sliding windows**: This is particularly useful for systems employing local caching, leading to significant savings in the number of calls made to S3. For example, if a local cache stores both the URL and its corresponding content, future accesses to that content within the sliding window will generate the same URL. This consistency allows the cache to serve the stored content without making repeated calls to S3. If the generated URLs were different each time, the cache would be unable to reuse the content, negating its benefits.

- **Custom Start Date for Signed URLs**: Generate signed URLs valid at any specific time, be it in the past or future, providing more control over when the links are valid for access.

## Features

### 1. Custom Start Date for Signed URLs
Unlike the typical behavior where the URL validity starts immediately upon generation, this project allows you to define a custom start time, providing more flexibility in scheduling access to your S3 resources.

### 2. Sliding Window Signed URL Generation
This feature allows you to generate the same signed URL for repeated access within a sliding window. By specifying the window duration and the residual time, the system can transition to the next time slot early. The slot will always start in advance by the residual time specified.

- **Benefit for caching**: By generating identical URLs during the sliding window, local caches can store and reuse both the URL and the content without making redundant calls to S3.

## How to Use

1. Specify the desired `bucket`, `key`, `start time`, and `expiration time` for generating a custom signed URL.
2. Optionally, configure the `sliding window` and `residual time` to generate consistent URLs for repeated access across time slots.
3. Use the generated signed URLs to access your S3 objects securely and flexibly.

# Examples

## Basic S3 URL Signing

```go
// Example 1: Basic S3 URL Signing
s3UrlSigner := NewS3UrlSigner("<awsIdAccessKey>", "<awsSecretAccessKey>", "<awsRegion>", "<bucket>")
url := s3UrlSigner.GetS3SignedUrl("/images/test1.jpg", time.Unix(1727633084, 0), 4*time.Hour)
fmt.Println(url)
```

## S3 URL Signing with Sliding Window

```go
// Example 2: S3 URL Signing with Sliding Window
s3UrlSignerSlidingWindow := NewS3UrlSignerSlidingWindow(s3UrlSigner, 48*time.Hour, 4*time.Hour)
url := s3UrlSignerSlidingWindow.GetS3SignedUrl("/images/test1.jpg")

fmt.Println(url)
```

## License

This project is licensed under the Apache License 2.0.
