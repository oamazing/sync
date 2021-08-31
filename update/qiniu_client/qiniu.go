package qiniu_client

import (
	"log"

	"github.com/qiniu/go-sdk/v7/storage"
)

type QiniuClient struct {
}

func NewQiniuClient() *QiniuClient {
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuadong
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	return &QiniuClient{}
}

func (qiniu *QiniuClient) Write(file string) {
	log.Printf("write file %s", file)
}
func (qiniu *QiniuClient) Remove(file string) {
	log.Printf("remove file %s", file)
}
func (qiniu *QiniuClient) List() []string {
	return []string{}
}
func (qiniu *QiniuClient) Download(string) {

}
func (qiniu *QiniuClient) Downloads([]string) {

}
