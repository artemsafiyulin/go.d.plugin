// SPDX-License-Identifier: GPL-3.0-or-later

package nvidia_smi

import (
	"time"

	"github.com/netdata/go.d.plugin/agent/module"
	"github.com/netdata/go.d.plugin/pkg/web"
)

func init() {
	module.Register("nvidia_smi", module.Creator{
		Defaults: module.Defaults{
			Disabled:    true,
			UpdateEvery: 10,
		},
		Create: func() module.Module { return New() },
	})
}

func New() *NvidiaSMI {
	return &NvidiaSMI{
		Config: Config{
			Timeout:      web.Duration{Duration: time.Second * 5},
			UseCSVFormat: true,
		},
		binName: "nvidia-smi",
		charts:  &module.Charts{},
		gpus:    make(map[string]bool),
	}

}

type Config struct {
	Timeout      web.Duration
	BinaryPath   string `yaml:"binary_path"`
	UseCSVFormat bool   `yaml:"use_csv_format"`
}

type (
	NvidiaSMI struct {
		module.Base
		Config `yaml:",inline"`

		charts *module.Charts

		binName string
		exec    nvidiaSMI

		gpuQueryProperties []string

		gpus map[string]bool
	}
	nvidiaSMI interface {
		queryGPUInfoXML() ([]byte, error)
		queryGPUInfoCSV(properties []string) ([]byte, error)
		queryHelpQueryGPU() ([]byte, error)
	}
)

func (nv *NvidiaSMI) Init() bool {
	if nv.exec == nil {
		smi, err := nv.initNvidiaSMIExec()
		if err != nil {
			nv.Error(err)
			return false
		}
		nv.exec = smi
	}

	return true
}

func (nv *NvidiaSMI) Check() bool {
	return len(nv.Collect()) > 0
}

func (nv *NvidiaSMI) Charts() *module.Charts {
	return nv.charts
}

func (nv *NvidiaSMI) Collect() map[string]int64 {
	mx, err := nv.collect()
	if err != nil {
		nv.Error(err)
	}

	if len(mx) == 0 {
		return nil
	}
	return mx
}

func (nv *NvidiaSMI) Cleanup() {}
