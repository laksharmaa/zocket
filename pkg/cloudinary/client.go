package cloudinary

import (
    "context"
    "github.com/cloudinary/cloudinary-go/v2"
    "github.com/cloudinary/cloudinary-go/v2/api/uploader"
)

type Client struct {
    cld *cloudinary.Cloudinary
}

func NewCloudinaryClient(cloudName, apiKey, apiSecret string) (*Client, error) {
    cld, err := cloudinary.NewFromParams(cloudName, apiKey, apiSecret)
    if err != nil {
        return nil, err
    }

    return &Client{
        cld: cld,
    }, nil
}

func (c *Client) UploadAndCompressImage(ctx context.Context, imageURL string) (string, error) {
    // Upload params with compression settings
    params := uploader.UploadParams{
        ResourceType: "auto",
        // Use Transformation parameter for quality and format optimization
        Transformation: "q_auto:low,f_auto",
    }

    // Upload the image
    result, err := c.cld.Upload.Upload(ctx, imageURL, params)
    if err != nil {
        return "", err
    }

    return result.SecureURL, nil
}