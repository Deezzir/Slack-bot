package blob

import (
	"bytes"
	"context"
	"slack-bot/pkg/config"
	"slack-bot/pkg/utils"

	"github.com/Azure/azure-sdk-for-go/sdk/storage/azblob"
	"github.com/patrickmn/go-cache"
)

type File struct {
	Filename string
	Data     string
	Desc     string
}

func GetBlobFiles(ctx context.Context, container string) []*File {
	var files []*File
	pager := config.AzContainerClient.ListBlobsFlat(&azblob.ContainerListBlobsFlatOptions{
		Include: []azblob.ListBlobsIncludeItem{azblob.ListBlobsIncludeItemMetadata},
	})

	for pager.NextPage(ctx) {
		page := pager.PageResponse()
		for _, blob := range page.ListBlobsFlatSegmentResponse.Segment.BlobItems {
			tmp := &File{}

			tmp.Filename = *blob.Name
			if desc := blob.Metadata["desc"]; desc != nil {
				tmp.Desc = *desc
			}

			files = append(files, tmp)
		}
	}

	if err := pager.Err(); err != nil {
		utils.ErrorLogger.Printf("Failed to get blob files - %s\n", err.Error())
	}

	return files
}

func GetBlobFile(ctx context.Context, containerName, filename string) (*File, bool) {
	var file *File = &File{}

	if data, ok := config.FileCache.Get(filename); ok {
		utils.InfoLogger.Printf("File '%s' is cached, returning\n", filename)
		return data.(*File), ok
	}

	blobClient, err := config.AzContainerClient.NewBlockBlobClient(filename)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to create blob client - %s\n", err.Error())
		return file, false
	}
	res, err := blobClient.Download(ctx, nil)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to download blob file -\n %s\n", err.Error())
		return file, true
	}
	data := &bytes.Buffer{}
	reader := res.Body(&azblob.RetryReaderOptions{})
	defer reader.Close()
	_, err = data.ReadFrom(reader)
	if err != nil {
		utils.ErrorLogger.Printf("Failed to read from blob file -\n %s\n", err.Error())
		return file, false
	}

	utils.InfoLogger.Printf("Caching '%s' file\n", filename)
	file.Filename = filename
	file.Data = data.String()

	config.FileCache.Set(filename, file, cache.DefaultExpiration)
	return file, true
}
