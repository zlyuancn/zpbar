/*
-------------------------------------------------
   Author :       zlyuan
   date：         2019/8/19
   Description :
-------------------------------------------------
*/

package zpbar

import "fmt"

func timeFmt(cost int64) string {
    var h, m, s int64
    h = cost / 3600
    m = (cost % 3600) / 60
    s = cost % 60

    if cost < 3600 {
        return fmt.Sprintf("%02d:%02d", m, s)
    }
    return fmt.Sprintf("%02d:%02d:%02d", h, m, s)
}
