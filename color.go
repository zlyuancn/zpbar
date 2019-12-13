/*
-------------------------------------------------
   Author :       zlyuan
   date：         2019/5/15
   Description :
-------------------------------------------------
*/

package zpbar

import "fmt"

type ColorType uint8

const (
    ColorDefault = ColorType(iota)      //默认
    ColorRed     = ColorType(iota + 30) //红
    ColorGreen                          //绿
    ColorYellow                         //黄
    ColorBlue                           //蓝
    ColorMagenta                        //紫
    ColorCyan                           //深绿
    ColorWhite                          //白
)

func getColorStartFlag(color ColorType) string {
    return fmt.Sprintf("\x1b[%dm", color)
}

func getColorEndFlag() string {
    return fmt.Sprintf("\x1b[0m")
}
