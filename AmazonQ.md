# SpurtCMS Image Upload Endpoint Documentation

## Overview
SpurtCMS provides an endpoint for uploading base64-encoded images that can be used with the built-in image editor. This document explains how to use this endpoint and how it handles storage based on your configuration.

## Endpoint Details

- **URL**: `/uploadb64image`
- **Method**: POST
- **Content-Type**: application/x-www-form-urlencoded

## Request Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| imagedata | string | Base64-encoded image data |

## Response Format

```json
{
  "imagepath": "https://your-domain.com/image-resize?name=path/to/image.jpg",
  "imgname": "image.jpg"
}
```

## Error Response

```json
{
  "Error": "Failed to upload image. Please verify your storage configuration and try again."
}
```

## Storage Configuration

SpurtCMS supports multiple storage backends:

1. **Local Storage**: Files are stored in the local filesystem
2. **AWS S3**: Files are stored in Amazon S3 buckets
3. **Azure Storage**: (Coming soon)

The system automatically detects which storage type to use based on your configuration in the database.

## Implementation Details

The endpoint uses the `UploadImage` function from the `storagecontroller` package, which:

1. Checks the selected storage type from the database
2. Routes the upload to the appropriate storage handler
3. Returns the URL path to the uploaded image

## Integration with the Editor

The image upload URL is dynamically configured in the entry editor pages. The system:

1. Reads the `BASE_URL` from environment variables (e.g., `http://localhost:8082/`)
2. Appends `uploadb64image` to create the full upload endpoint URL
3. Passes this URL to the editor component as a JSON configuration object:

```json
{
  "path": "http://localhost:8082/uploadb64image", 
  "payload": "imagedata", 
  "success": {
    "imagepath": "imagepath", 
    "imagename": "imagename"
  }
}
```

This ensures that the image upload endpoint is always relative to the current domain where SpurtCMS is running.

## Example Usage

```javascript
// Example using fetch API
const uploadImage = async (base64Data) => {
  const formData = new FormData();
  formData.append('imagedata', base64Data);
  
  const response = await fetch('https://your-domain.com/uploadb64image', {
    method: 'POST',
    body: formData
  });
  
  const result = await response.json();
  return result.imagepath;
};
```

## Security Considerations

- The endpoint should only be accessible to authenticated users
- Input validation is performed to ensure only valid image data is processed
- File size limits should be configured in your web server

## Troubleshooting

If you encounter issues with image uploads:

1. Check your storage configuration in the database
2. Verify that the appropriate credentials are set in your environment variables
3. Check server logs for detailed error messages
4. Ensure the `BASE_URL` in your .env file is correctly set to your domain
