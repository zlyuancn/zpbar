/*
-------------------------------------------------
   Author :       Zhang Fan
   date：         2019/5/14
   Description :
-------------------------------------------------
*/

package zpbar

import (
    "os"
)

func printOs(a string) {
    _, _ = os.Stdout.WriteString(a)
}
