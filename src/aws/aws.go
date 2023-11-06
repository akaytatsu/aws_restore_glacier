package aws

import (
	"app/utils"
	"bufio"
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/vbauerster/mpb/v7"
	"github.com/vbauerster/mpb/v7/decor"
)

const concurrency = 200

type AWSOperations struct {
	BucketName string
	Region     string
	PartialKey string
	client     *s3.Client
}

type BucketItem struct {
	Key          string
	Owner        string
	StorageClass string
	LastModified string
	Size         string
	IsRestoring  string
}

type RestoreItemData struct {
	Key     string
	Err     string
	Success bool
}

func (u *AWSOperations) Init(aws_access_key string, aws_secret_key string) {

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(aws_access_key, aws_secret_key, "")),
		config.WithRegion(u.Region),
	)
	if err != nil {
		panic("failed to load config, " + err.Error())
	}

	u.client = s3.NewFromConfig(cfg)

}

func (u *AWSOperations) processObjects(objs []types.Object, dataResponse *[]BucketItem, sem chan struct{}, wg *sync.WaitGroup) {
	for _, object := range objs {
		sem <- struct{}{}
		wg.Add(1)

		go func(obj types.Object) {
			defer wg.Done()
			defer func() { <-sem }()

			if obj.StorageClass == "DEEP_ARCHIVE" || obj.StorageClass == "GLACIER" || obj.StorageClass == "GLACIER_IR" {

				// dataResponse = append(dataResponse, obj)

				var key, onwer, storeClass, lastModified, size string
				var isRestoring string

				if obj.RestoreStatus != nil {
					isRestoring = fmt.Sprintf("%t", obj.RestoreStatus.IsRestoreInProgress)
				}

				if obj.Key != nil {
					key = *obj.Key
				}

				if obj.Owner != nil {
					onwer = *obj.Owner.DisplayName
				}

				storeClass = string(obj.StorageClass)

				if obj.LastModified != nil {
					lastModified = obj.LastModified.String()
				}

				size = fmt.Sprintf("%d", obj.Size)

				*dataResponse = append(*dataResponse, BucketItem{
					Key:          key,
					Owner:        onwer,
					StorageClass: storeClass,
					LastModified: lastModified,
					Size:         size,
					IsRestoring:  isRestoring,
				})

				// size = fmt.Sprintf("%d", obj.Size)

				// // writer.WriteString(
				// // 	fmt.Sprintf("%s;%s;%s;%s;%s\n", key, onwer, storeClass, lastModified, size),
				// // )
			}
		}(object)
	}
}

func (u *AWSOperations) processAllObjects(objs []types.Object, dataResponse *[]BucketItem, sem chan struct{}, wg *sync.WaitGroup) {
	for _, object := range objs {
		sem <- struct{}{}
		wg.Add(1)

		go func(obj types.Object) {
			defer wg.Done()
			defer func() { <-sem }()

			var key, onwer, storeClass, lastModified, size string

			if obj.Key != nil {
				key = *obj.Key
			}

			if obj.Owner != nil {
				onwer = *obj.Owner.DisplayName
			}

			storeClass = string(obj.StorageClass)

			if obj.LastModified != nil {
				lastModified = obj.LastModified.String()
			}

			isRestoring := ""

			if obj.RestoreStatus != nil {
				isRestoring = fmt.Sprintf("%t", obj.RestoreStatus.IsRestoreInProgress)
			}

			size = fmt.Sprintf("%d", obj.Size)

			*dataResponse = append(*dataResponse, BucketItem{
				Key:          key,
				Owner:        onwer,
				StorageClass: storeClass,
				LastModified: lastModified,
				Size:         size,
				IsRestoring:  isRestoring,
			})

		}(object)
	}
}

func (u *AWSOperations) ListObjects(pathToFile string) (objects []BucketItem, err error) {

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(u.BucketName),
		Prefix: aws.String(u.PartialKey),
	}

	file, err := utils.GenerateFile("list", u.BucketName)
	if err != nil {
		fmt.Printf("Erro ao criar arquivo: %v\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	writer.WriteString(
		fmt.Sprintf("%s;%s;%s;%s;%s;%s\n", "key", "owner", "storageClass", "lastModified", "size", "isRestoring"),
	)

	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	println("Listando objetos...")

	paginator := s3.NewListObjectsV2Paginator(u.client, input)

	for paginator.HasMorePages() {

		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			fmt.Printf("Erro ao listar objetos: %v\n", err)
			return objects, err
		}

		u.processObjects(page.Contents, &objects, sem, &wg)
	}

	wg.Wait()

	for _, object := range objects {
		writer.WriteString(
			fmt.Sprintf("%s;%s;%s;%s;%s;%s\n", object.Key, object.Owner, object.StorageClass, object.LastModified, object.Size, object.IsRestoring),
		)
	}

	err = writer.Flush()
	if err != nil {
		fmt.Printf("Erro ao finalizar escrita no arquivo: %v\n", err)
	}

	println("Total de objetos: ", len(objects))

	return objects, nil
}

func (u *AWSOperations) ListAllObjects(pathToFile string) (objects []BucketItem, err error) {

	input := &s3.ListObjectsV2Input{
		Bucket: aws.String(u.BucketName),
		Prefix: aws.String(u.PartialKey),
	}

	file, err := utils.GenerateFile("list_all", u.BucketName)
	if err != nil {
		fmt.Printf("Erro ao criar arquivo: %v\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	writer.WriteString(
		fmt.Sprintf("%s;%s;%s;%s;%s;%s\n", "key", "owner", "storageClass", "lastModified", "size", "isRestoring"),
	)

	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	println("Listando objetos...")

	paginator := s3.NewListObjectsV2Paginator(u.client, input)

	for paginator.HasMorePages() {

		page, err := paginator.NextPage(context.TODO())
		if err != nil {
			fmt.Printf("Erro ao listar objetos: %v\n", err)
			return objects, err
		}

		u.processAllObjects(page.Contents, &objects, sem, &wg)
	}

	wg.Wait()

	for _, object := range objects {
		writer.WriteString(
			fmt.Sprintf("%s;%s;%s;%s;%s;%s\n", object.Key, object.Owner, object.StorageClass, object.LastModified, object.Size, object.IsRestoring),
		)
	}

	err = writer.Flush()
	if err != nil {
		fmt.Printf("Erro ao finalizar escrita no arquivo: %v\n", err)
	}

	println("Total de objetos: ", len(objects))

	return objects, nil
}

func (u *AWSOperations) RestoreObjects(objects *[]BucketItem) {

	var response []RestoreItemData = make([]RestoreItemData, 0)

	file, err := utils.GenerateFile("restore", u.BucketName)
	if err != nil {
		fmt.Printf("Erro ao criar arquivo: %v\n", err)
		return
	}
	defer file.Close()

	writer := bufio.NewWriter(file)

	writer.WriteString(
		fmt.Sprintf("%s;%s;%s;\n", "key", "Success", "Error"),
	)

	println("Restaurando objetos...")

	sem := make(chan struct{}, concurrency)
	var wg sync.WaitGroup

	p := mpb.New(mpb.WithWaitGroup(&wg))

	totalItems := len(*objects)

	bar := p.AddBar(int64(totalItems),
		mpb.PrependDecorators(
			decor.CountersNoUnit("%d/%d", decor.WCSyncWidth),
		),
		mpb.AppendDecorators(
			decor.Percentage(decor.WCSyncWidth),
		),
	)

	for _, object := range *objects {
		key := object.Key

		sem <- struct{}{}
		wg.Add(1)

		go u.processRestoreObject(key, &wg, sem, bar, &response)
	}

	wg.Wait()

	p.Wait()

	for _, object := range response {
		writer.WriteString(
			fmt.Sprintf("%s;%t;%s\n", object.Key, object.Success, object.Err),
		)
	}

	err = writer.Flush()
	if err != nil {
		fmt.Printf("Erro ao finalizar escrita no arquivo: %v\n", err)
	}
}

func (u *AWSOperations) processRestoreObject(key string, wg *sync.WaitGroup, sem chan struct{}, bar *mpb.Bar, response *[]RestoreItemData) {
	defer wg.Done()

	input := &s3.RestoreObjectInput{
		Bucket: aws.String(u.BucketName),
		Key:    aws.String(key),
		RestoreRequest: &types.RestoreRequest{
			Days: 30,
		},
	}

	_, err := u.client.RestoreObject(context.TODO(), input)

	bar.Increment()

	if err != nil {
		*response = append(*response, RestoreItemData{
			Key:     key,
			Err:     err.Error(),
			Success: false,
		})
	}

	*response = append(*response, RestoreItemData{
		Key:     key,
		Err:     "",
		Success: true,
	})

	defer func() { <-sem }()
}
