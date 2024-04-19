/**
 * @Author:      Adam wu
 * @Description:
 * @File:        disk_test.go
 * @Version:     1.0.0
 * @Date:        2024/4/19
 */

package disk

import (
	"fmt"
	"github.com/chaosblade-io/chaosblade-exec-os/exec/util"
	"testing"
)

func TestDisk(t *testing.T) {
	fmt.Println(util.IsDir("/Users/10000134/code/chaosblade-exec-os/README_CN.md"))
}
