# Storage Configuration Guide

This document provides the environment configuration needed for local and CS3 storage systems.

## Environment Variables

Add these to your `.env` file:

### Common Storage Settings

```bash
# Storage Configuration
FILESYSTEM_TYPE=local         # Options: local, cs3
STORAGE_DIR=storage          # Base storage directory for local storage
```

### Local Storage Settings

For local filesystem storage (development/testing):

```bash
FILESYSTEM_TYPE=local
STORAGE_DIR=storage
```

Files will be stored in:
- `./storage/app/avatars/` - User avatars
- `./storage/app/contents/` - Content files
- `./storage/app/products/` - Product images
- `./storage/app/ads/` - Advertisement images
- `./storage/tmp/` - Temporary uploads

### CS3 (S3-Compatible) Storage Settings

For cloud storage (production):

```bash
# Select CS3 storage
FILESYSTEM_TYPE=cs3

# CS3 Configuration
CS_ENDPOINT=sin1.contabostorage.com              # S3 endpoint URL
CS_BUCKET_CODE=f7fa66e2663f40628d0d1a14d566355e  # Bucket access code
CS_BUCKET=your-bucket-name                       # Bucket name
CS_REGION=us-east-1                              # Region (if applicable)
CS_ACCESS_KEY=your_access_key_here               # Access key ID
CS_SECRET_KEY=your_secret_key_here               # Secret access key
```

## Available Storage Providers

### Contabo Object Storage
```bash
CS_ENDPOINT=sin1.contabostorage.com  # Singapore
# or
CS_ENDPOINT=eu2.contabostorage.com   # Europe
# or
CS_ENDPOINT=usc1.contabostorage.com  # US
```

### AWS S3
```bash
CS_ENDPOINT=s3.amazonaws.com
CS_REGION=us-east-1
CS_BUCKET=your-bucket-name
CS_ACCESS_KEY=your-aws-access-key
CS_SECRET_KEY=your-aws-secret-key
```

### MinIO
```bash
CS_ENDPOINT=your-minio-server.com:9000
CS_BUCKET=your-bucket-name
CS_ACCESS_KEY=your-minio-access-key
CS_SECRET_KEY=your-minio-secret-key
```

### DigitalOcean Spaces
```bash
CS_ENDPOINT=nyc3.digitaloceanspaces.com
CS_REGION=nyc3
CS_BUCKET=your-space-name
CS_ACCESS_KEY=your-spaces-key
CS_SECRET_KEY=your-spaces-secret
```

## Storage API Endpoints

Once configured, the following API endpoints are available:

### Local Storage Endpoints
- `GET /api/v1/storage-local/presigned-url?filename={filename}` - Get presigned URL
- `POST /api/v1/storage-local/uploads` - Upload files (multipart)
- `PUT /api/v1/storage-local/uploads/{file_name}` - Upload single file (binary)
- `PUT /api/v1/storage-local/legitimize-files` - Legitimize uploaded files

### CS3 Storage Endpoints
- `GET /api/v1/storage/presigned-url?filename={filename}` - Get presigned URL
- `PUT /api/v1/storage/legitimize-files` - Legitimize uploaded files

### File Serving Routes (Local Storage Only)
- `GET /storage/avatars/{file_name}` - Serve avatar files
- `GET /storage/contents/{file_name}` - Serve content files
- `GET /storage/products/{file_name}` - Serve product images
- `GET /storage/ads/{file_name}` - Serve advertisement images
- `GET /storage/tmp/{file_name}` - Serve temporary files

## Testing Configuration

### Test Local Storage
```bash
# Set in .env
FILESYSTEM_TYPE=local
STORAGE_DIR=storage

# Test endpoint
curl http://localhost:7789/api/v1/storage-local/presigned-url?filename=test.jpg
```

### Test CS3 Storage
```bash
# Set in .env
FILESYSTEM_TYPE=cs3
# ... (add CS3 credentials)

# Test endpoint
curl http://localhost:7789/api/v1/storage/presigned-url?filename=test.jpg
```

## Directory Structure

For local storage, ensure these directories exist:

```
storage/
├── app/
│   ├── avatars/     # User profile pictures
│   ├── contents/    # Content files
│   ├── products/    # Product images
│   └── ads/         # Advertisement images
├── logs/            # Application logs
└── tmp/             # Temporary uploads
```

Create directories if they don't exist:
```bash
mkdir -p storage/app/{avatars,contents,products,ads}
mkdir -p storage/tmp
mkdir -p storage/logs
```

## Security Recommendations

1. **Never commit `.env` file** - Add to `.gitignore`
2. **Use CS3/S3 for production** - More secure and scalable
3. **Restrict bucket permissions** - Only allow necessary access
4. **Use presigned URLs** - For direct client uploads
5. **Validate file types** - In your application code
6. **Set file size limits** - Prevent abuse
7. **Use HTTPS endpoints** - For CS3/S3 connections

## Troubleshooting

### Local Storage Issues
- **Files not serving**: Check `STORAGE_DIR` path is correct
- **Permission denied**: Ensure storage directory is writable
- **404 errors**: Verify file paths match route patterns

### CS3/S3 Issues
- **Connection failed**: Check `CS_ENDPOINT` is correct
- **Access denied**: Verify `CS_ACCESS_KEY` and `CS_SECRET_KEY`
- **Bucket not found**: Confirm `CS_BUCKET` exists
- **Invalid signature**: Check credentials haven't expired

