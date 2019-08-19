# zpbar
单行色彩进度条

## 获得zpbar
`go get -u github.com/zlyuancn/zpbar`

## 示例
```go
package main

import (
    "github.com/zlyuancn/zpbar"
    "time"
)

func main() {
    p := zpbar.NewPbar(
        zpbar.WithTotal(100),
    )

    p.Start()
    for i := 0; i < 100; i++ {
        p.Done()
        time.Sleep(0.05e9)
    }
    p.Close()
}
```
