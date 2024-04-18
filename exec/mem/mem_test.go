package mem

import (
	"fmt"
	"testing"
)

func TestCgroupv2Memory(t *testing.T) {
	total, i, err := cgroupV2AvailableAndTotal("/sys/fs/cgroup/", "ram", 5449, true)
	if err != nil {
		t.Error(err)
	}
	fmt.Println(total, i)
}
