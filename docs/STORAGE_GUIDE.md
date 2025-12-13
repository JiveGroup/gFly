# gFly Framework: Storage System Guide

## Table of Contents
1. [Overview](#overview)
2. [Storage Providers](#storage-providers)
3. [Setup and Configuration](#setup-and-configuration)
4. [Basic Storage Operations](#basic-storage-operations)
5. [File Upload Patterns](#file-upload-patterns)
6. [Presigned URLs](#presigned-urls)
7. [File Legitimization](#file-legitimization)
8. [Storage API Endpoints](#storage-api-endpoints)
9. [Transformers and Public URLs](#transformers-and-public-urls)
10. [Best Practices](#best-practices)

---

## Overview

The gFly framework provides a **unified storage abstraction** that supports multiple storage backends through a consistent interface. This allows you to seamlessly switch between different storage providers (local filesystem, S3-compatible storage, etc.) without changing your application code.

### Key Features

- **Multiple Storage Backends**: Support for local storage and S3-compatible cloud storage (CS3)
- **Unified Interface**: Single API for all storage operations regardless of backend
- **Presigned URLs**: Generate secure, time-limited upload/download URLs
- **File Legitimization**: Convert temporary uploaded files to permanent storage locations
- **URL Generation**: Automatic public URL generation based on storage type
- **Type Safety**: Full Go type support with proper error handling

### Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                Application Code                             │
│         storage.Instance() / storage.Disk(type)             │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│            Storage Abstraction Layer                        │
│         github.com/gflydev/storage                          │
└──────────────────────┬──────────────────────────────────────┘
                       │
         ┌─────────────┴─────────────┐
         │                           │
         ▼                           ▼
┌──────────────────┐      ┌──────────────────┐
│  Local Storage   │      │   CS3 Storage    │
│   (filesystem)   │      │ (S3-compatible)  │
└──────────────────┘      └──────────────────┘
```

---

## Storage Providers

### 1. Local Storage (`github.com/gflydev/storage/local`)

**Use Cases:**
- Development and testing
- Small-scale applications
- Applications with simple file storage needs
- Fast local file access

**Characteristics:**
- Files stored on local filesystem
- Direct file system access
- No external dependencies
- URLs served through application server

**Storage Path Pattern:**
```
/storage/app/{directory}/{filename}
```

### 2. CS3 Storage (`github.com/gflydev/storage/cs3`)

**Use Cases:**
- Production environments
- Scalable cloud storage
- CDN integration
- Distributed applications
- Large file storage

**Characteristics:**
- S3-compatible cloud storage (Contabo, AWS S3, MinIO, etc.)
- Presigned URL support for direct uploads
- Public URL generation
- Scalable and highly available

**Storage Path Pattern:**
```
https://{endpoint}/{bucket_code}:{bucket}/{path}/{filename}
```

---

## Setup and Configuration

### 1. Register Storage Providers

In your `cmd/web/main.go`, register the storage providers you want to use:

```go
import (
    "github.com/gflydev/storage"
    storageLocal "github.com/gflydev/storage/local"
    storageCS3 "github.com/gflydev/storage/cs3"
)

func main() {
    // Register storage providers
    storage.Register(storageLocal.Type, storageLocal.New())
    storage.Register(storageCS3.Type, storageCS3.New())

    // ... rest of your application setup
}
```

### 2. Environment Configuration

Configure storage settings in your `.env` file:

**Common Settings:**
```bash
# Storage Configuration
FILESYSTEM_TYPE=cs3           # Options: local, cs3
STORAGE_DIR=storage           # Base storage directory for local storage
```

**Local Storage Settings:**
```bash
STORAGE_DIR=storage
```

**CS3 (S3-Compatible) Storage Settings:**
```bash
CS_ENDPOINT=sin1.contabostorage.com    # S3 endpoint
CS_BUCKET_CODE=f7fa66e2663f40628d0d1a14d566355e  # Bucket access code
CS_BUCKET=thietngon                     # Bucket name
CS_REGION=us-east-1                     # Region (if applicable)
CS_ACCESS_KEY=your_access_key           # Access key
CS_SECRET_KEY=your_secret_key           # Secret key
```

### 3. Storage Types

The framework uses storage type constants:

```go
// Local storage type
storageLocal.Type  // "local"

// CS3 storage type
storageCS3.Type    // "cs3"
```

---

## Basic Storage Operations

### Get Storage Instance

**Default Storage (from environment):**
```go
import "github.com/gflydev/storage"

// Get the default storage instance based on FILESYSTEM_TYPE
fs := storage.Instance()
```

**Specific Storage Type:**
```go
// Get local storage explicitly
localFS := storage.Disk(storageLocal.Type)

// Get CS3 storage explicitly
cs3FS := storage.Disk(storageCS3.Type)
```

### Generate Public URLs

Convert a storage path to a publicly accessible URL:

```go
fs := storage.Instance()

// Path: "avatars/user-123.jpg"
// Returns: Full public URL based on storage type
publicURL := fs.Url("avatars/user-123.jpg")

// Local: http://yourapp.com/storage/avatars/user-123.jpg
// CS3: https://sin1.contabostorage.com/bucket:name/avatars/user-123.jpg
```

### Check if URL is Absolute

```go
import "github.com/gflydev/core"

if strings.HasPrefix(path, core.SchemaHTTP) {
    // Already an absolute URL (http:// or https://)
    return path
}

// Generate URL using storage
return fs.Url(path)
```

---

## File Upload Patterns

### Pattern 1: Multipart Form Upload (Local Storage)

Used for traditional form-based file uploads.

**Controller Example:**
```go
type UploadApi struct {
    core.Api
}

func NewUploadApi() *UploadApi {
    return &UploadApi{}
}

// Handle processes file upload from multipart form
// @Summary Upload files via form
// @Tags Storage
// @Accept multipart/form-data
// @Produce json
// @Param files formData file true "Files to upload"
// @Success 200 {object} http.Success[[]response.UploadedFile]
// @Router /api/v1/storage-local/uploads [post]
func (h *UploadApi) Handle(c *core.Ctx) error {
    // Files are handled by gFly storage module
    // Implementation in github.com/gflydev/modules/storage/api

    // Returns list of uploaded files with paths
    return c.Success(uploadedFiles)
}
```

### Pattern 2: Direct Binary Upload

Upload a file directly as request body:

```go
// PUT /api/v1/storage-local/uploads/{file_name}
// Content-Type: application/octet-stream
// Body: <binary file data>

func (h *UploadFileApi) Handle(c *core.Ctx) error {
    fileName := c.PathVal("file_name")

    // Read file from request body
    fileData := c.Body()

    // Save to storage
    // Implementation handles saving

    return c.Success(uploadedFile)
}
```

### Pattern 3: Presigned URL Upload (CS3)

Client-side direct upload to cloud storage (recommended for production).

**Step 1: Get Presigned URL**
```go
// Client requests presigned URL
// GET /api/v1/storage/presigned-url?filename=avatar.png

type PresignURLRequest struct {
    Filename string `json:"filename" validate:"required"`
}

func (h *PresignedURLApi) Handle(c *core.Ctx) error {
    var req dto.PreSignURL

    // Parse and validate
    if err := c.Parse(&req); err != nil {
        return c.JSON(400, err)
    }

    // Generate presigned URLs
    // Returns both upload URL and final file URL
    return c.Success(response.PresignedURL{
        UploadURL: "https://s3.../upload?signature=...",
        FileURL:   "https://s3.../avatars/file.jpg",
    })
}
```

**Step 2: Client Uploads to Presigned URL**
```javascript
// Frontend code (JavaScript example)
const response = await fetch('/api/v1/storage/presigned-url?filename=avatar.png');
const { upload_url, file_url } = await response.json();

// Upload file directly to S3
await fetch(upload_url, {
    method: 'PUT',
    body: fileBlob,
    headers: { 'Content-Type': 'image/png' }
});

// Use file_url in your application
console.log('File available at:', file_url);
```

**Step 3: Legitimize the File**
```javascript
// Tell backend the file is ready and should be moved from temp to permanent location
await fetch('/api/v1/storage/legitimize-files', {
    method: 'PUT',
    body: JSON.stringify({
        files: [{
            file: file_url,
            name: 'avatar.png',
            dir: 'avatars'
        }]
    })
});
```

---

## Presigned URLs

Presigned URLs provide secure, time-limited access to cloud storage for uploads and downloads.

### Benefits

- **Direct Upload**: Files upload directly to cloud storage, not through your server
- **Reduced Server Load**: No file data passes through your application
- **Faster Uploads**: Direct connection to cloud storage
- **Secure**: Time-limited and scoped permissions

### Generate Presigned URL

```go
import "github.com/gflydev/storage"

fs := storage.Disk(storageCS3.Type)

// Generate presigned URL for upload
uploadURL, fileURL, err := fs.PresignedURL(
    "avatars/user-123.jpg",  // File path in storage
    "PUT",                    // HTTP method (PUT for upload, GET for download)
    15 * time.Minute,         // URL validity duration
)

if err != nil {
    return err
}

// uploadURL: Temporary URL for uploading
// fileURL: Permanent URL to access the file after upload
```

### DTO for Presigned URL Request

```go
// internal/dto/upload_dto.go
type PreSignURL struct {
    Filename string `json:"filename" validate:"required,lte=255" doc:"The name of the file to be uploaded."`
}
```

### Response Structure

```go
// internal/http/response/presign_url_response.go
type PresignedURL struct {
    UploadURL string `json:"upload_url" example:"https://s3.amazonaws.com/bucket/upload?signature=abc123" doc:"The URL to upload the file to."`
    FileURL   string `json:"file_url" example:"https://s3.amazonaws.com/bucket/file.jpg" doc:"The URL to access the uploaded file."`
}
```

---

## File Legitimization

File legitimization is the process of converting temporary uploaded files to permanent storage locations with proper validation.

### Purpose

When files are uploaded (especially via presigned URLs), they initially exist in a temporary location. Legitimization:
- Moves files from temporary to permanent locations
- Validates file existence
- Generates proper public URLs
- Confirms successful upload

### Two Approaches to Legitimization

There are **two ways** to legitimize uploaded files in gFly:

1. **Via API Endpoint** (Client-side approach)
   - Frontend calls `/api/v1/storage/legitimize-files` endpoint
   - Used when client handles file upload directly
   - Suitable for SPA/frontend applications

2. **Via Utility Function** (Server-side approach)
   - Backend calls `utils.LegitimizeUploadedFile()` in service layer
   - Used when backend processes file uploads
   - No additional HTTP request needed
   - Suitable for server-side file handling

### Legitimization Flow

```
┌─────────────────────────────────────────────────────────────┐
│  1. Client uploads to presigned URL                         │
│     File stored in: /tmp/avatar.png                         │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│  2. Client calls legitimize-files endpoint                  │
│     Provides: file URL, name, target directory              │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│  3. Backend validates and moves file                        │
│     File moved to: /avatars/avatar.png                      │
│     Generates public URL                                    │
└──────────────────────┬──────────────────────────────────────┘
                       │
                       ▼
┌─────────────────────────────────────────────────────────────┐
│  4. Returns legitimized file information                    │
│     Public URL ready for use                                │
└─────────────────────────────────────────────────────────────┘
```

### Legitimization Request DTO

```go
// internal/dto/upload_dto.go
type LegitimizeItem struct {
    File          string `json:"file" validate:"required,lte=255" doc:"The file URL pointing to the uploaded file."`
    Name          string `json:"name" validate:"required,lte=255" doc:"The display name of the file."`
    Dir           string `json:"dir" validate:"required,lte=255" doc:"The directory where the file is stored."`
    LegitimizeURL string `json:"legitimize_url" doc:"The public legitimize URL, created by the backend."`
}

type LegitimizeFile struct {
    Files []LegitimizeItem `json:"files" validate:"required" doc:"The list of files to be legitimized."`
}
```

### Approach 1: Via API Endpoint (Client-Side)

Use this approach when the frontend handles file uploads directly.

**Client Example:**
```javascript
// After uploading file via presigned URL
const legitimizeResponse = await fetch('/api/v1/storage/legitimize-files', {
    method: 'PUT',
    headers: { 'Content-Type': 'application/json' },
    body: JSON.stringify({
        files: [
            {
                file: 'https://s3.../tmp/avatar.png',
                name: 'avatar.png',
                dir: 'avatars'
            }
        ]
    })
});

const { files } = await legitimizeResponse.json();
// files[0].legitimize_url contains the permanent public URL
```

### Approach 2: Via Utility Function (Server-Side)

Use this approach when the backend processes file uploads in the service layer. This is more efficient as it doesn't require an additional HTTP request.

**Utility Function Implementation:**

The `utils.LegitimizeUploadedFile()` function is located in `pkg/utils/storage_utils.go`:

```go
// internal/services/user_services.go
const (
    UploadAvatarDir = "avatars"
)

// UpdateUserProfile updates user profile information including avatar
func UpdateUserProfile(user *models.User, updateProfileDto dto.UpdateProfile) (*models.User, error) {
    // ... other profile updates ...

    // Handle avatar upload if provided
    if updateProfileDto.Avatar != "" && updateProfileDto.Avatar != user.Avatar.String {
        // Legitimize the uploaded file (move from temp to permanent location)
        legitimizedFile := utils.LegitimizeUploadedFile(updateProfileDto.Avatar, UploadAvatarDir)

        // Only update if legitimization was successful
        if legitimizedFile != "" {
            user.Avatar = dbNull.String(legitimizedFile)
        }
    }

    // Update timestamp
    user.UpdatedAt = dbNull.TimeNow()

    // Save updated user
    err := mb.UpdateModel(user)
    if err != nil {
        log.Errorf("Error while updating user profile: %v", err)
        return nil, errors.New("Failed to update user profile")
    }

    return user, nil
}
```

**How It Works:**

1. **Client uploads file** to presigned URL and gets temporary URL (e.g., `/tmp/avatar_123.jpg`)
2. **Client sends temporary URL** to backend in the update profile request
3. **Service layer calls** `utils.LegitimizeUploadedFile()` with the temp URL and target directory
4. **Function automatically detects** storage type (CS3 or local) from URL pattern
5. **File is moved** from temporary location to permanent directory (e.g., `tmp/avatar_123.jpg` → `avatars/avatar_123.jpg`)
6. **Returns proxy URL** that can be stored in the database
7. **Database is updated** with the legitimized file URL

**When to Use Each Approach:**

| Scenario | Approach | Reason |
|----------|----------|--------|
| Frontend uploads file directly | API Endpoint | Client needs to legitimize after upload |
| Form-based file upload | Utility Function | Backend handles the entire flow |
| User profile update with avatar | Utility Function | Service layer processes the update |
| Bulk file imports | Utility Function | Server-side batch processing |
| SPA with direct S3 upload | API Endpoint | Frontend manages upload lifecycle |
| Background job processing | Utility Function | No HTTP context available |

---

## Storage API Endpoints

The framework provides two sets of storage endpoints based on the storage type.

### Local Storage Endpoints

**Base Path:** `/api/v1/storage-local`

```go
// Get presigned URL (local storage)
GET /api/v1/storage-local/presigned-url?filename={filename}

// Upload files via multipart form
POST /api/v1/storage-local/uploads
Content-Type: multipart/form-data

// Upload single file via binary body
PUT /api/v1/storage-local/uploads/{file_name}
Content-Type: application/octet-stream

// Legitimize uploaded files
PUT /api/v1/storage-local/legitimize-files
Content-Type: application/json
```

### CS3 Storage Endpoints

**Base Path:** `/api/v1/storage`

```go
// Get presigned URL (CS3 cloud storage)
GET /api/v1/storage/presigned-url?filename={filename}

// Legitimize uploaded files
PUT /api/v1/storage/legitimize-files
Content-Type: application/json
```

### Endpoint Registration

In `internal/http/routes/api_routes.go`:

```go
import (
    storageApi "github.com/gflydev/modules/storage/api"
    storageCS3Api "github.com/gflydev/modules/storagecs3/api"
)

func ApiRoutes(f core.IFly) {
    // ...

    /* ============================ Storage Group (Local) ============================ */
    apiRouter.Group("/storage-local", func(uploadGroup *core.Group) {
        uploadGroup.GET("/presigned-url", storageApi.NewPresignedURLApi())
        uploadGroup.POST("/uploads", storageApi.NewUploadApi())
        uploadGroup.PUT("/uploads/{file_name}", storageApi.NewUploadFileApi())
        uploadGroup.PUT("/legitimize-files", storageApi.NewLegitimizeFileApi())
    })

    /* ============================ Storage Group (CS3) ============================ */
    apiRouter.Group("/storage", func(uploadGroup *core.Group) {
        uploadGroup.GET("/presigned-url", storageCS3Api.NewPresignedURLApi())
        uploadGroup.PUT("/legitimize-files", storageCS3Api.NewLegitimizeFileApi())
    })
}
```

### Serving Local Files

For local storage, you need routes to serve the uploaded files:

```go
// internal/http/routes/web_routes.go
func WebRoutes(f core.IFly) {
    // Serve uploaded files for Local storage
    f.GET("/storage/avatars/{file_name}", page.NewStoragePage("app/avatars"))
    f.GET("/storage/contents/{file_name}", page.NewStoragePage("app/contents"))
    f.GET("/storage/products/{file_name}", page.NewStoragePage("app/products"))
    f.GET("/storage/ads/{file_name}", page.NewStoragePage("app/ads"))
    f.GET("/storage/tmp/{file_name}", page.NewStoragePage("tmp"))
}
```

**Storage Page Controller:**
```go
// internal/http/controllers/page/storage_page.go
type StoragePage struct {
    BasePage
    path string
}

func NewStoragePage(path string) *StoragePage {
    return &StoragePage{path: path}
}

func (h *StoragePage) Handle(c *core.Ctx) error {
    // ./storage/tmp/file_name.ext
    filePath := fmt.Sprintf("%s/%s/%s",
        os.Getenv("STORAGE_DIR"),
        h.path,
        c.PathVal("file_name"))

    return c.File(filePath)
}
```

### Proxy for CS3 Storage (Optional)

If you want to proxy CS3 URLs through your application:

```go
// internal/http/controllers/page/objects_proxy_page.go
type ObjectsProxyPage struct {
    core.Page
}

func (m *ObjectsProxyPage) Handle(c *core.Ctx) error {
    // Forward: /objects/contents/article-default-img.jpg
    // To: https://sin1.contabostorage.com/bucket:name/contents/article-default-img.jpg

    remainingPath := strings.TrimPrefix(c.Path(), "/objects")

    targetURL := fmt.Sprintf("https://%s/%s:%s%s",
        utils.Getenv("CS_ENDPOINT", ""),
        utils.Getenv("CS_BUCKET_CODE", ""),
        utils.Getenv("CS_BUCKET", ""),
        remainingPath,
    )

    return c.ProxyStream(targetURL)
}

// Register route
f.GET("/objects/{path:*}", page.NewObjectsProxyPage())
```

---

## Transformers and Public URLs

Transformers convert domain models to API responses, handling URL generation for uploaded files.

### Generic File Transformer

```go
// internal/http/transformers/generic_transformer.go
import (
    "github.com/gflydev/core"
    "github.com/gflydev/storage"
    "strings"
)

// PublicUploadedFile converts media URL to public URL
func PublicUploadedFile(mediaUrl string) string {
    fs := storage.Instance()

    // Already an absolute URL
    if strings.HasPrefix(mediaUrl, core.SchemaHTTP) {
        return mediaUrl
    }

    // Auto generate with current storage type
    return fs.Url(mediaUrl)
}
```

### User Avatar Transformer

```go
// internal/http/transformers/user_transformer.go
func PublicAvatar(avatar string) *string {
    if avatar == "" {
        return nil
    }

    fs := storage.Instance()

    // Absolute URL
    if strings.HasPrefix(avatar, core.SchemaHTTP) {
        return &avatar
    }

    // Generate public URL
    avatar = fs.Url(avatar)
    return &avatar
}

func ToUserResponse(user models.User) response.User {
    return response.User{
        ID:        user.ID,
        Email:     user.Email,
        FirstName: dbNull.StringNil(user.FirstName),
        LastName:  dbNull.StringNil(user.LastName),
        Avatar:    PublicAvatar(user.Avatar.String),  // Transform avatar URL
        CreatedAt: user.CreatedAt,
        // ... other fields
    }
}
```

### Uploaded File Response Transformer

```go
// internal/http/transformers/storage_transformer.go
func ToUploadedFileResponse(model core.UploadedFile) response.UploadedFile {
    return response.UploadedFile{
        Field: model.Field,
        Name:  model.Name,
        Path:  model.Path,
        Size:  model.Size,
    }
}
```

### Usage in Controllers

```go
type GetUserApi struct {
    core.Api
}

func (h *GetUserApi) Handle(c *core.Ctx) error {
    userID := c.PathValInt("id")

    // Get user from service
    user, err := services.GetUserByID(userID)
    if err != nil {
        return c.Error(http.Error{
            Code:    "USER_NOT_FOUND",
            Message: "User not found",
        })
    }

    // Transform to response (avatar URL is converted automatically)
    userResponse := transformers.ToUserResponse(user)

    return c.JSON(userResponse)
}
```

---

## Best Practices

### 1. Storage Type Selection

**Use Local Storage:**
- Development and testing environments
- Simple applications with low storage needs
- Quick prototyping
- Single-server deployments

**Use CS3/Cloud Storage:**
- Production environments
- Scalable applications
- Applications requiring CDN
- Multi-server or distributed deployments
- Large file storage requirements

### 2. Environment-Based Configuration

Use environment variables to switch storage backends:

```go
// Development (.env.development)
FILESYSTEM_TYPE=local
STORAGE_DIR=storage

// Production (.env.production)
FILESYSTEM_TYPE=cs3
CS_ENDPOINT=sin1.contabostorage.com
CS_BUCKET=myapp-production
```

### 3. Upload Strategy

**For Small Files (< 5MB):**
- Use multipart form upload
- Simple implementation
- Good for user profiles, avatars

**For Large Files (> 5MB):**
- Use presigned URLs
- Direct client-to-cloud upload
- Reduces server load
- Faster upload times

### 4. Directory Organization

Organize files by type in separate directories:

```
storage/app/
├── avatars/        # User avatars
├── products/       # Product images
├── contents/       # Article/content images
├── ads/            # Advertisement banners
└── documents/      # User documents
```

### 5. URL Generation Best Practices

Always use transformers for URL generation:

```go
// ✅ GOOD: Use transformer
avatar := transformers.PublicAvatar(user.Avatar.String)

// ❌ BAD: Hardcode URL generation
avatar := "https://myapp.com/storage/avatars/" + user.Avatar.String
```

### 6. Error Handling

Always handle storage errors properly:

```go
fs := storage.Instance()

url, err := fs.Url(filePath)
if err != nil {
    log.Error("Failed to generate URL", "error", err, "path", filePath)
    return c.Error(http.Error{
        Code:    "STORAGE_ERROR",
        Message: "Failed to generate file URL",
    })
}
```

### 7. File Validation

Validate files before upload:

```go
type UploadRequest struct {
    Filename string `json:"filename" validate:"required,lte=255"`
}

// Check file extension
allowedExtensions := []string{".jpg", ".jpeg", ".png", ".gif"}
ext := filepath.Ext(req.Filename)
if !slices.Contains(allowedExtensions, strings.ToLower(ext)) {
    return c.Error(http.Error{
        Code:    "INVALID_FILE_TYPE",
        Message: "File type not allowed",
    })
}

// Check file size (if uploading via form)
const maxFileSize = 10 * 1024 * 1024 // 10MB
if fileSize > maxFileSize {
    return c.Error(http.Error{
        Code:    "FILE_TOO_LARGE",
        Message: "File size exceeds 10MB limit",
    })
}
```

### 8. Secure File Names

Sanitize file names to prevent security issues:

```go
import "github.com/gflydev/core/utils"

// Generate safe filename
safeFilename := utils.Slug(originalFilename)

// Or add unique identifier
timestamp := time.Now().Unix()
uniqueFilename := fmt.Sprintf("%d-%s", timestamp, safeFilename)
```

### 9. Storage Instance Caching

Cache storage instances when making multiple operations:

```go
// ✅ GOOD: Cache instance
fs := storage.Instance()
url1 := fs.Url(path1)
url2 := fs.Url(path2)
url3 := fs.Url(path3)

// ❌ LESS EFFICIENT: Multiple instance calls
url1 := storage.Instance().Url(path1)
url2 := storage.Instance().Url(path2)
url3 := storage.Instance().Url(path3)
```

### 10. Testing Storage Operations

For unit tests, mock the storage interface:

```go
func TestFileUpload(t *testing.T) {
    // Setup mock storage
    mockFS := &mockStorage{
        urlFunc: func(path string) string {
            return "http://test.com/" + path
        },
    }

    // Use mock in tests
    // ...
}
```

### 11. Documentation

Always document storage-related fields in DTOs and responses:

```go
type ProductRequest struct {
    Image string `json:"image" validate:"required,url" doc:"Product image URL from storage upload"`
}

type ProductResponse struct {
    Image string `json:"image" example:"https://cdn.com/products/item.jpg" doc:"Public URL to product image"`
}
```

### 12. Cleanup Temporary Files

For temporary uploads, implement cleanup:

```go
// After processing, clean up temp files
if strings.Contains(filePath, "/tmp/") {
    defer func() {
        if err := fs.Delete(filePath); err != nil {
            log.Error("Failed to cleanup temp file", "path", filePath, "error", err)
        }
    }()
}
```

---

## Common Patterns

### Complete Upload Flow Example

**1. Define Constants**
```go
// internal/services/user_services.go
const (
    UploadAvatarDir = "avatars"
)
```

**2. Service Method for Avatar Update**
```go
func UpdateUserAvatar(userID int, avatarURL string) error {
    // Legitimize the uploaded file
    legitimizedURL := utils.LegitimizeUploadedFile(avatarURL, UploadAvatarDir)

    // Update user record
    user := models.User{
        ID:     userID,
        Avatar: dbNull.String(legitimizedURL),
    }

    err := mb.Model(&user).
        Where("id", mb.Eq, userID).
        Update(mb.Fields{"avatar": legitimizedURL})

    return err
}
```

**3. Controller**
```go
type UpdateProfileApi struct {
    core.Api
}

func (h *UpdateProfileApi) Handle(c *core.Ctx) error {
    var req request.UpdateProfile

    // Parse request
    if err := http.ProcessData(c, &req); err != nil {
        return err
    }

    userID := c.Locals("user_id").(int)

    // Update avatar if provided
    if req.Avatar != nil {
        if err := services.UpdateUserAvatar(userID, *req.Avatar); err != nil {
            return c.Error(http.Error{
                Code:    "UPDATE_FAILED",
                Message: "Failed to update avatar",
            })
        }
    }

    // Get updated user
    user, _ := services.GetUserByID(userID)

    // Return response with transformed avatar URL
    return c.JSON(transformers.ToUserResponse(user))
}
```

### Storage Type Detection

```go
func GetStorageTypeFromURL(fileURL string) string {
    bucketPath := fmt.Sprintf("/%s:%s/",
        utils.Getenv("CS_BUCKET_CODE", ""),
        utils.Getenv("CS_BUCKET", ""),
    )

    if strings.Contains(fileURL, bucketPath) {
        return storageCS3.Type
    }

    return storageLocal.Type
}
```

---

## Troubleshooting

### Issue: URLs not generating correctly

**Problem:** `fs.Url()` returns empty or incorrect URLs

**Solution:**
- Check `FILESYSTEM_TYPE` environment variable
- Verify storage provider is registered in `main.go`
- For CS3, ensure all `CS_*` environment variables are set

### Issue: Files not accessible after upload

**Problem:** Uploaded files return 404

**Local Storage:**
- Verify `STORAGE_DIR` is correct
- Check file permissions on storage directory
- Ensure storage routes are registered in `web_routes.go`

**CS3 Storage:**
- Verify bucket permissions (should be public or have proper ACLs)
- Check `CS_ENDPOINT`, `CS_BUCKET_CODE`, and `CS_BUCKET` are correct
- Test presigned URL expiration time

### Issue: Presigned URLs expire too quickly

**Problem:** Upload fails due to expired presigned URL

**Solution:**
- Increase presigned URL validity duration
- Implement client-side retry logic
- Consider chunked uploads for large files

### Issue: CORS errors with CS3 uploads

**Problem:** Browser blocks direct uploads to S3

**Solution:**
Configure CORS on your S3 bucket:
```xml
<CORSConfiguration>
    <CORSRule>
        <AllowedOrigin>https://yourapp.com</AllowedOrigin>
        <AllowedMethod>PUT</AllowedMethod>
        <AllowedMethod>POST</AllowedMethod>
        <AllowedMethod>GET</AllowedMethod>
        <AllowedHeader>*</AllowedHeader>
    </CORSRule>
</CORSConfiguration>
```

## Additional Resources

- [STORAGE_CONFIG.md](STORAGE_CONFIG.md) - Complete storage configuration
- [gFly Storage Documentation](https://www.gfly.dev/docs/storage)
- [S3 API Documentation](https://docs.aws.amazon.com/s3/)

---

## Summary

The gFly storage system provides:

1. **Unified Interface**: Single API for multiple storage backends
2. **Flexibility**: Easy switching between local and cloud storage
3. **Security**: Presigned URLs for secure, direct uploads
4. **Performance**: Reduced server load with client-side uploads
5. **Type Safety**: Full Go type support with error handling
6. **Best Practices**: Built-in patterns for URL generation and file management

By following this guide, you can implement robust file storage in your gFly applications that scales from development to production seamlessly.
