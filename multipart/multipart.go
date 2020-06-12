package multipart

import (
	"apicore"
	"bytes"
	"crypto/md5"
	"fmt"
	"io"
	"os"
)

var MaxCatchSize = 1024 * 1024 * 1

func init() {
	apicore.AddMiddleware(func() apicore.MiddleWare {
		return &Multipart{}
	})
}

type Multipart struct{}

func (m *Multipart) Before(ctx apicore.Context) {
	var catches []byte
	form, err := ctx.Raw().MultipartReader()
	if err != nil {
		return
	}
	for {
		item, err := form.NextPart()
		if err != nil {
			break
		}
		// 读到文件
		if item.FileName() != "" {
			// 初始化缓存
			if catches == nil {
				catches = make([]byte, MaxCatchSize)
			}
			// 尝试读入缓存中
			n, err := _copy(catches, item)
			if err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println(n)
			if n < MaxCatchSize {
				fmt.Println("catches")
				data := make([]byte, n)
				copy(data, catches[:n])
				reader := bytes.NewReader(data)
				ctx.SetValue(item.FormName(), &File{name: item.FileName(), ReadCloser: &BufferFile{reader}})
			}
			// 超过缓存，写入文件
			if n == MaxCatchSize {
				fmt.Println("tofile")
				tempName := md5V3(item.FileName())
				f, err := os.Create(tempName)
				if err != nil {
					fmt.Println(err)
					continue
				}
				io.Copy(f, bytes.NewReader(catches))
				io.Copy(f, item)
				f.Seek(0, io.SeekStart)
				ctx.SetValue(item.FormName(), &File{name: item.FileName(), ReadCloser: &IOFile{name: tempName, File: f}})
			}
		}
	}
}

func _copy(dst []byte, reader io.Reader) (int, error) {
	var (
		err   error
		n     int
		total int
	)
	for {
		n, err = reader.Read(dst[total:])
		total += n
		if err == io.EOF {
			return total, nil
		}
		if err != nil {
			return 0, err
		}
		if total == MaxCatchSize {
			return total, nil
		}
	}
}

func md5V3(str string) string {
	w := md5.New()
	io.WriteString(w, str)
	md5str := fmt.Sprintf("%x", w.Sum(nil))
	return md5str
}

func (m *Multipart) After(ctx apicore.Context) {
	return
}

func (m Multipart) Index() int {
	return 10
}
