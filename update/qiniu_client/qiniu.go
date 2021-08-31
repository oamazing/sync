package qiniu_client

import (
	"fmt"
	"log"
	"sync"
	"time"

	"github.com/oamazing/sync/config"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type QiniuClient struct {
	Config storage.Config
	Bucket string `c:"存储空间"`
	Token  string `c:"上传凭证"`
	mux    *sync.RWMutex
	auth   *auth.Credentials
	ticker *time.Ticker
}

func NewQiniuClient() *QiniuClient {
	conf := config.GetConfig().Qiniu
	if conf.Bucket == `` {
		log.Panic(`config file error`)
	}
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuadong
	// 是否使用https域名
	cfg.UseHTTPS = false
	// 上传是否使用CDN上传加速
	cfg.UseCdnDomains = false
	mac := qbox.NewMac(conf.Ak, conf.Sk)
	policy := storage.PutPolicy{
		Scope:   conf.Bucket,
		Expires: 7200,
	}
	token := policy.UploadToken(mac)
	client := &QiniuClient{
		Config: cfg,
		auth:   mac,
		mux:    &sync.RWMutex{},
		Token:  token,
		ticker: time.NewTicker(10 * time.Second),
		Bucket: conf.Bucket,
	}
	go func(client *QiniuClient) {
		for {
			select {
			case <-client.ticker.C:
				client.mux.Lock()
				log.Println("refresh token")
				client.Token = policy.UploadToken(mac)
				client.mux.Unlock()
			}
		}
	}(client)
	return client
}

func (qiniu *QiniuClient) Write(file string) {
	// loader := storage.NewFormUploader(&qiniu.Config)
	// loader.PutFile()
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

func (qiniu *QiniuClient) Close() {
	fmt.Println("close client")
	qiniu.ticker.Stop()
}
