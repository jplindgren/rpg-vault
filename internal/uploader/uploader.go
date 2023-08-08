package uploader

import (
	"encoding/base64"
	"fmt"
	"strings"

	"github.com/jplindgren/rpg-vault/internal/clients"
)

func UploadCoverImage(s3 *clients.S3ClientWrapper, base64Image string, destination string) (string, error) {
	if base64Image == "" {
		return "", nil
	}

	// _, err := url.ParseRequestURI(cover)
	// isURL := err == nil
	// if !isURL { //if it was not a url, it was not uploaded yet
	// 	imgBytes, err := base64.StdEncoding.DecodeString(cover)
	// 	if err != nil {
	// 		return fmt.Errorf("error decoding base64 image: %w", err)
	// 	}
	// 	destinationPath := fmt.Sprintf(publishImageKeyFormat, userId)
	// 	err = uploadFileToS3(s3, imgBytes, destinationPath)
	// 	return err
	// 	//https://my-bucket.s3-ap-southeast-2.amazonaws.com/foo/bar.txt
	//   s3://rpg-go/f572a37c-33a7-4b5b-85d9-86cb596d2edb/world/cover.png
	//   https://rpg-go.s3.sa-east-1.amazonaws.com/f572a37c-33a7-4b5b-85d9-86cb596d2edb/world/cover.png
	// }
	//return nil

	b64data := base64Image[strings.IndexByte(base64Image, ',')+1:]
	imgBytes, err := base64.StdEncoding.DecodeString(b64data)
	if err != nil {
		return "", fmt.Errorf("error decoding base64 image: %w", err)
	}

	url, err := s3.Upload(imgBytes, destination)
	if err != nil {
		return "", err
	}

	return url, nil
}
