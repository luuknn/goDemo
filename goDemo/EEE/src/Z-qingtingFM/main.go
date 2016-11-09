package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"path"
	"sync"
)

//翻唱FM的URL：http://www.qingting.fm/#/vchannels/136962/programs/5659080
//翻唱FM的ajax地址:http://www.qingting.fm/s/vchannels/136962/programs/5745196/ajax

var QingTingAjaxAddrFmt = "http://www.qingting.fm/s/vchannels/%d/programs/%d/ajax"

//Go语言中json的解析，如果没有固定的输入格式，尽量使用通用格式 map [] interface{}
//因为GO语言导出结构体默认使用首字母大写，而且需要和实际的格式进行对应，还不如使用map+interface+[]
//type QingTingSongSubItem struct {
//     Name       string           "歌曲名称"
//     ParentName string           "专辑名称"
//     QTSongs    []string        "歌曲URL"
//     ParentID   int              "父目录ID"
//     Duration   int              "时长"
//     Type       string           "类型"
//     ID         int              "类型"
//     Thumb      string           "缩略图"
//}
//
//type QingTingSongItem struct {
//     PlayInfo QingTingSongSubItem
//}
//
//type QingTingPlayInfo QingTingSongItem
//
//func (this *QingTingPlayInfo) Download(dirname string) (bool, error) {
//     return false, nil
//}
//
//func (this *QingTingPlayInfo) String() string{
//     return fmt.Sprintf("歌曲名称:%#v 专辑名称:%#v 歌曲URL地址:%#v 父目录ID:%#v 时长:%#v 类型:%#v 缩略图:%#v",
//            this.PlayInfo.Name,
//            this.PlayInfo.ParentName,
//            this.PlayInfo.QTSongs,
//            this.PlayInfo.ParentID,
//            this.PlayInfo.Duration,
//            this.PlayInfo.Type,
//            this.PlayInfo.ID,
//            this.PlayInfo.Thumb)
//}

type QingTingPlayInfo map[string]interface{}

func Download(this QingTingPlayInfo, dirname string, wg *sync.WaitGroup) (bool, error) {
	defer wg.Done()

	//pthis := (map[string]interface{})(this)
	name := this["name"].(string)
	url := "http://od.qingting.fm" + this["urls"].([]interface{})[0].(string)
	filename := path.Join(dirname, name+".m4a")
	fmt.Printf("开始下载:%s 到：%s\n", url, filename)

	res, err := http.Get(url)
	if err != nil {
		fmt.Printf("下载错误：%s\n", err)
		return false, err
	}

	defer res.Body.Close()
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Printf("下载失败:%s\n", err)
		return false, err
	}

	err2 := ioutil.WriteFile(filename, data, 0666)
	if err2 != nil {
		fmt.Printf("下载失败:%s\n", err2)
		return false, err2
	}

	fmt.Printf("!!!下载成功:%s 到：%s\n", url, filename)
	return true, nil
}

//-----------------------------------------------------------------
type QingTingFM struct {
	Channel int "频道"
	ID      int "频道ID"
}

func (this *QingTingFM) Crawl() ([]QingTingPlayInfo, error) {
	jsonbyte, err := this.readJson()
	if err != nil {
		return nil, err
	}

	if jsonbyte == nil {
		return nil, errors.New("返回了无效数据")
	}

	//go语言和json的映射关系 :对应结构体数据 []对应数组
	pis := make([]QingTingPlayInfo, 0)
	jerr := json.Unmarshal(jsonbyte, &pis)
	if jerr != nil {
		return nil, jerr
	}

	downloadDir := "music_go"
	os.MkdirAll(downloadDir, 0666)
	//for _,pi := range pis{
	//     fmt.Println(pi["name"])
	//}
	var wg sync.WaitGroup
	//连接缓存最大数目
	maxIndex := len(pis) - 1
	cacheMax := 5
	for i, pi := range pis {
		wg.Add(1)
		cacheMax++
		//wg必须是传引用，这里必须要限制异步个数，否则与服务器的连接过多，服务器会拒绝，从而导致下载失败
		go Download(pi, downloadDir, &wg)
		//使用缓存
		if cacheMax > 5 || maxIndex == i {
			wg.Wait()
			cacheMax = 0
		}
	}

	return pis, nil
}

func (this *QingTingFM) readJson() ([]byte, error) {
	//ajaxAddr := "http://www.golangtc.com/t/533c0ad2320b520cc400004f"
	//http.Get自动gzip解压缩数据流
	res, err := http.Get(this.getAajxAddr())
	if err != nil {
		return nil, err
	}

	readerCloser := res.Body
	//读取之后需要关闭
	defer readerCloser.Close()

	if res.StatusCode != 200 {
		return nil, errors.New("返回错误状态码")
	}

	bytes := make([]byte, 0, 1024)
	//buffer必须是数组，或者可以使用ioutil ReadAll来读取
	var buffera [128]byte
	buffer := buffera[:]

	for {
		readed, rerr := readerCloser.Read(buffer)
		if readed == 0 && rerr == nil {
			return nil, errors.New("readed=0 rerr=nil")
		}

		//切片...，自动将数组解压成很多个参数输入
		//实现原理和...，变参支持的原理一样
		//这句应该在 ==io.EOF之前，如果写到后面，则会丢失最后部分的数据
		bytes = append(bytes, buffer[:readed]...)

		//数据流读取完毕
		if rerr == io.EOF {
			break
		}

		if rerr != nil {
			return nil, rerr
		}

	}

	//将bytes转换成string
	return bytes, nil
}

func (this *QingTingFM) getAajxAddr() string {
	return fmt.Sprintf(QingTingAjaxAddrFmt, this.Channel, this.ID)
}

//-----------------------------------------------------------------
func main() {
	qtfm := QingTingFM{Channel: 105674, ID: 110000}
	qtfm.Crawl()
}
