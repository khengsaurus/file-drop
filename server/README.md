## AWS S3

### Bucket policy

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "S3 public GET",
      "Effect": "Allow",
      "Principal": "*",
      "Action": "s3:GetObject",
      "Resource": "<resource-arn>/*"
    }
  ]
}
```

### CORS policy

```json
[
  {
    "AllowedHeaders": ["*"],
    "AllowedMethods": ["GET", "HEAD", "POST", "PUT"],
    "AllowedOrigins": ["<allowed-origin>"],
    "ExposeHeaders": []
  }
]
```

### Lifecycle rule

- Rule config: rule scope - Apply to all objects in the bucket
- Lifecycle rule actions:
  - Expire current versions of objects
  - Permanently delete noncurrent versions of objects
  - Delete expired object delete markers or incomplete multipart uploads
- Expire current versions: 1 day after object creation
- Permanently delete noncurrent versions of objects: 1 day after objects become noncurrent
- Delete expired object delete markers or incomplete multipart uploads: Delete incomplete multipart uploads, 1 day

### IAM policy

```json
{
  "Version": "2012-10-17",
  "Statement": [
    {
      "Sid": "VisualEditor0",
      "Effect": "Allow",
      "Action": [
        "s3:PutObject",
        "s3:PutObjectAcl",
        "s3:GetObject",
        "s3:DeleteObject"
      ],
      "Resource": "<resource-arn>/*"
    },
    {
      "Sid": "VisualEditor1",
      "Effect": "Allow",
      "Action": ["s3:ListBucket"],
      "Resource": "<resource-arn>"
    }
  ]
}
```

## Running services with Docker

Requirements: [LS CLI](https://github.com/localstack/awscli-local), Docker Compose, Python 3

```bash
# Run services in docker:
> docker compose -f docker-compose-services.yml up

# Init LocalStack S3:
> python3 scripts/ls_s3_setup.py

# List LocalStack S3 buckets
> awslocal s3 ls

# List contents of LocalStack S3 bucket
> awslocal s3 ls <bucket-name>

# Show tags of LocalStack S3 object
> awslocal s3api get-object-tagging --bucket <bucket-name> --key <object-key>
```

<hr/>

## Issues

- Adding header `x-amz-tagging` works in LS but results in 403 for AWS

- [Localstack docs](https://docs.localstack.cloud/localstack/persistence-mechanism/): "_...please note that persistence in LocalStack, as currently intended, is a Pro only feature..._" 😕

## Reference

- [AWS-CLI docs](https://docs.localstack.cloud/integrations/aws-cli/#aws-cli)
- [Setting up LS S3](https://alojea.com/how-to-create-an-aws-local-bucket/)
