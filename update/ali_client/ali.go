package ali_client

import (
	"log"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/oamazing/sync/config"
)

type AliClient struct {
	client *oss.Client
	bucket *oss.Bucket
}

func NewAliClient() *AliClient {
	conf := config.GetAli()
	client, err := oss.New(conf.Region, conf.Key, conf.Secret)
	if err != nil {
		log.Panic("init aliyun oss client failed")
	}
	bucket, err := client.Bucket(conf.Bucket)
	if err != nil {
		log.Panic("get bucket failed")
	}
	return &AliClient{
		client: client,
		bucket: bucket,
	}

}

func (ali *AliClient) Write(realpath string, name string) {
	if err := ali.bucket.PutObjectFromFile(name, realpath); err != nil {
		log.Println("upload file error")
	}
	log.Printf("write file %s success", name)
}
func (ali *AliClient) Remove(string, string) {

}
func (ali *AliClient) List() []string {
	return nil
}
func (ali *AliClient) Download(string) {

}
func (ali *AliClient) Downloads([]string) {

}
func (ali *AliClient) Close() {

}
