package qiniu_client

import (
	"context"
	"log"
	"path/filepath"
	"sync"
	"time"

	"github.com/oamazing/sync/config"
	"github.com/qiniu/go-sdk/v7/auth"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
)

type QiniuClient struct {
	Config  storage.Config
	Bucket  string `c:"存储空间"`
	Token   string `c:"上传凭证"`
	mux     *sync.RWMutex
	auth    *auth.Credentials
	ticker  *time.Ticker
	ch      chan bool
	manager *storage.BucketManager
}

func NewQiniuClient() *QiniuClient {
	conf := config.GetConfig().Qiniu
	if conf.Bucket == `` {
		log.Panic(`config file error`)
	}
	cfg := storage.Config{}
	// 空间对应的机房
	cfg.Zone = &storage.ZoneHuanan
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
	manager := storage.NewBucketManager(mac, &cfg)
	client := &QiniuClient{
		Config:  cfg,
		auth:    mac,
		mux:     &sync.RWMutex{},
		Token:   token,
		ticker:  time.NewTicker(1 * time.Hour),
		Bucket:  conf.Bucket,
		manager: manager,
	}
	client.ch = make(chan bool, 1)
	go func(client *QiniuClient) {
		for {
			select {
			case <-client.ticker.C:
				client.mux.Lock()
				log.Println("refresh token")
				client.Token = policy.UploadToken(mac)
				log.Println(client.Token)
				client.mux.Unlock()
			case <-client.ch:
				log.Println("close")
				return
			}
		}
	}(client)
	return client
}

func (qiniu *QiniuClient) Write(relpath, file string) {

	uploader := storage.NewFormUploader(&qiniu.Config)
	_, fileName := filepath.Split(file)
	if relpath != `` {
		fileName = relpath + `/` + fileName
	}
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": fileName,
		},
	}
	// loader.Put(context.Background(), nil, qiniu.Token, "")
	if err := uploader.PutFile(context.Background(), &storage.PutRet{}, qiniu.Token, fileName, file, &putExtra); err != nil {
		log.Printf("upload error %s", err)
		return
	}
	log.Printf("write file %s success", file)
}
func (qiniu *QiniuClient) Remove(relpath, file string) {
	_, fileName := filepath.Split(file)
	if relpath != `` {
		fileName = relpath + `/` + fileName
	}
	if err := qiniu.manager.Delete(qiniu.Bucket, fileName); err != nil {
		log.Printf("remove file %s error %s", fileName, err)
		return
	}
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
	qiniu.ticker.Stop()
	qiniu.ch <- true
}
