package main

import (
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/oamazing/sync/config"
	"github.com/oamazing/sync/update"
	"github.com/oamazing/sync/update/qiniu_client"
)

var (
	watcher *fsnotify.Watcher
	client  update.Client
)

const (
	qiniu  = "qiniu"
	alioss = "alioss"
)

func main() {
	var err error
	watcher, err = fsnotify.NewWatcher()
	if err != nil {
		log.Panic(err)
	}
	defer watcher.Close()
	conf := config.GetConfig()
	switch conf.Storage {
	case qiniu:
		client = qiniu_client.NewQiniuClient()
	default:
		log.Panic("not found client")
	}
	go func() {
		for {
			select {
			case event, ok := <-watcher.Events:
				if !ok {
					return
				}
				handlerEvent(event, client)
			case err, ok := <-watcher.Errors:
				if !ok {
					return
				}
				log.Fatalln(err)
			}
		}
	}()
	addListener(conf.BasePath)
	defer func() {
		log.Println("quit")
	}()
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)
	<-ch
	client.Close()
}

func addListener(path string) {
	err := watcher.Add(path)
	if err != nil {
		log.Panic(err)
	}
	log.Printf("listener dir %s", path)
}

func handlerEvent(event fsnotify.Event, client update.Client) {

	fname := event.Name
	// finfo, err := os.Stat(fname)
	// if err != nil {
	// 	log.Printf("sync: get file info err %s", err)
	// }
	// if finfo.IsDir() {
	// 	// 如果是文件夹，那么就监听

	// }
	if event.Op&fsnotify.Create == fsnotify.Create {
		// 创建操作
		client.Write(fname)
	} else if event.Op&fsnotify.Remove == fsnotify.Remove {
		// 删除文件
		client.Remove(fname)
	}
}
