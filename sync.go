package main

import (
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/oamazing/sync/update"
	"github.com/oamazing/sync/update/qiniu_client"
	"gopkg.in/yaml.v3"
)

func init() {
	bs, err := ioutil.ReadFile("./config.yml")
	if err != nil {
		log.Panic(err)
	}
	if err = yaml.Unmarshal(bs, &conf); err != nil {
		log.Panic(err)
	}
}

type Config struct {
	BasePath string `yaml:"base_path"`
	Storage  string `yaml:"storage"`
	Qiniu    Qiniu  `yaml:"qiniu"`
}

type Qiniu struct {
	Ak string `yaml:"ak"`
	Sk string `yaml:"sk"`
}

var (
	watcher *fsnotify.Watcher
	conf    Config
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
