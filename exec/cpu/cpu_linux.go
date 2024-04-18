/*
 * Copyright 1999-2020 Alibaba Group Holding Ltd.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package cpu

import (
	"context"
	"fmt"
	"github.com/chaosblade-io/chaosblade-exec-os/exec"
	"github.com/chaosblade-io/chaosblade-spec-go/channel"
	"github.com/chaosblade-io/chaosblade-spec-go/log"
	"github.com/containerd/cgroups"
	cgroupsv2 "github.com/containerd/cgroups/v2"
	"github.com/shirou/gopsutil/cpu"
	"strconv"
	"time"
)

func getUsed(ctx context.Context, percpu bool, cpuIndex int) float64 {

	pid := ctx.Value(channel.NSTargetFlagName)
	cpuCount := ctx.Value("cpuCount").(int)

	if pid != nil {
		p, err := strconv.Atoi(pid.(string))
		if err != nil {
			log.Fatalf(ctx, "get cpu usage fail, %s", err.Error())
		}

		cgroupRoot := ctx.Value("cgroup-root")
		if cgroupRoot == "" {
			cgroupRoot = "/sys/fs/cgroup/"
		}

		log.Debugf(ctx, "get cpu useage by cgroup, root path: %s", cgroupRoot)
		var used float64
		if cgroups.Mode() == cgroups.Unified {
			// Adapt to cgroup v2
			used, err = cgroupV2Used(cgroupRoot.(string), p, cpuCount)
		} else {
			used, err = cgroupV1Used(cgroupRoot.(string), p, cpuCount)
		}
		if err != nil {
			log.Fatalf(ctx, "get cpu usage fail, %s", err.Error())
		}
		return used
	}

	totalCpuPercent, err := cpu.Percent(time.Second, percpu)
	if err != nil {
		log.Fatalf(ctx, "get cpu usage fail, %s", err.Error())
	}
	if percpu {
		if cpuIndex > len(totalCpuPercent) {
			log.Fatalf(ctx, "illegal cpu index %d", cpuIndex)
		}
		return totalCpuPercent[cpuIndex]
	}
	return totalCpuPercent[0]
}

func cgroupV1Used(cgroupRoot string, pid, cpuCount int) (float64, error) {
	cgroup, err := cgroups.Load(exec.Hierarchy(cgroupRoot), exec.PidPath(pid))

	if err != nil {
		return 0, fmt.Errorf("load cgroup error, %v", err)
	}
	stats, err := cgroup.Stat(cgroups.IgnoreNotExist)
	if err != nil {
		return 0, fmt.Errorf("load cgroup stat error, %v", err)
	}

	pre := float64(stats.CPU.Usage.Total) / float64(time.Second)
	time.Sleep(time.Second)
	nextStats, err := cgroup.Stat(cgroups.IgnoreNotExist)
	if err != nil {
		return 0, fmt.Errorf("get cpu usage fail, %s", err.Error())
	}
	next := float64(nextStats.CPU.Usage.Total) / float64(time.Second)
	return ((next - pre) * 100) / float64(cpuCount), nil

}
func cgroupV2Used(cgroupRoot string, pid, cpuCount int) (float64, error) {
	group, err := cgroupsv2.PidGroupPath(pid)
	if err != nil {
		return 0, err
	}
	cgroup, err := cgroupsv2.LoadManager(cgroupRoot, group)

	if err != nil {
		return 0, fmt.Errorf("load cgroup error, %v", err)
	}
	stats, err := cgroup.Stat()
	if err != nil {
		return 0, fmt.Errorf("load cgroup stat error, %v", err)
	}
	pre := float64(stats.CPU.UsageUsec) / float64(time.Second)
	time.Sleep(time.Second)
	nextStats, err := cgroup.Stat()
	if err != nil {
		return 0, fmt.Errorf("get cpu usage fail, %s", err.Error())
	}
	next := float64(nextStats.CPU.UserUsec) / float64(time.Second)
	return ((next - pre) * 100) / float64(cpuCount), nil

}
