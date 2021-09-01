package main

import (
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/signal"
	"path/filepath"
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
	var (
		currentFiles = []string{}
		osFiles      = []string{}
		syncFiles    = []string{}
	)
	// 读取目录中的所有文件和子目录
	files, err := ioutil.ReadDir(conf.BasePath)
	if err != nil {
		panic(err)
	}
	// 获取文件，并输出它们的名字
	for _, file := range files {
		if !file.IsDir() {
			currentFiles = append(currentFiles, file.Name())
		}
	}
	osFiles = client.List()
	for _, osFile := range osFiles {
		if !ContainerString(osFile, currentFiles) {
			syncFiles = append(syncFiles, osFile)
		}
	}
	// log.Printf("os files %+v", osFiles)
	download("http://files.yyang.xin", osFiles)
	client.Downloads(syncFiles)
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
	finfo, err := os.Stat(fname)
	if err != nil {
		log.Printf("sync: get file info err %s", err)

	} else {
		if finfo.IsDir() {
			// 如果是文件夹，那么就监听
			return
		}
	}
	relpath, err := filepath.Rel(config.GetConfig().BasePath, filepath.Dir(fname))
	if err != nil {
		log.Println("get real path error")
		return
	}
	if relpath == `.` {
		relpath = ``
	}
	if event.Op&fsnotify.Create == fsnotify.Create {
		// 创建操作
		client.Write(relpath, fname)

	} else if event.Op&fsnotify.Remove == fsnotify.Remove {
		// 删除文件
		client.Remove(relpath, fname)
	}
}

func ContainerString(key string, s []string) bool {
	for _, v := range s {
		if v == key {
			return true
		}
	}
	return false
}

func download(url string, fileNames []string) {
	for _, name := range fileNames {
		// files = append(files,)
		url = url + `/` + name
		f, err := os.Create(filepath.Join(config.GetConfig().BasePath, name))
		if err != nil {
			log.Printf("create file err: %s", err)
		}
		resp, err := http.Get(url)
		if err != nil {
			log.Printf("get file err: %s", err)
		}
		_, err = io.Copy(f, resp.Body)
		if err != nil {
			log.Printf("copy data err: %s", err)
		}
		log.Printf("sync file %s", name)
		f.Close()
		resp.Body.Close()
	}
	// log.Println(files)
}
