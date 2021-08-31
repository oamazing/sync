package update

type Client interface {
	Write(string)       // 写入
	Remove(string)      // 删除
	List() []string     // 文件列表
	Download(string)    // 下载
	Downloads([]string) // 批量下载
	Close()
}
