/*
Copyright 2023 cuisongliu@qq.com.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package system

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/containers/buildah"
	"github.com/containers/storage/pkg/homedir"

	"github.com/labring/sealos/pkg/constants"
)

type envSystemConfig struct{}

func Get(key string) (string, error) {
	cfg, err := globalConfig.getValueOrDefault(key)
	if err != nil {
		return "", err
	}
	return cfg.DefaultValue, nil
}

func GetConfig(key string) (*ConfigOption, error) {
	return globalConfig.getValueOrDefault(key)
}

var globalConfig *envSystemConfig

func init() {
	globalConfig = &envSystemConfig{}
}

type ConfigOption struct {
	Key          string
	Description  string
	DefaultValue string
	OSEnv        string
}

var configOptions = []ConfigOption{
	{
		Key:          PromptConfigKey,
		Description:  "toggle interactive prompting in the terminal.",
		DefaultValue: "enabled",
	},
	{
		Key:          RuntimeRootConfigKey,
		Description:  "root directory for runtime actions.",
		DefaultValue: filepath.Join(homedir.Get(), fmt.Sprintf(".%s", constants.AppName)),
	},
	{
		Key:          DataRootConfigKey,
		Description:  "cluster root directory for remote.",
		DefaultValue: filepath.Join("/var/lib", constants.AppName),
	},
	{
		Key:          BuildahFormatConfigKey,
		Description:  "`format` of the image manifest and metadata.",
		DefaultValue: buildah.OCI,
		OSEnv:        BuildahFormatConfigKey,
	},
	{
		Key:          ScpChecksumConfigKey,
		Description:  "whether to check the md5sum value is consistent during the copy process",
		DefaultValue: "yes",
	},
	{
		Key:          RegistrySyncExperimentalConfigKey,
		Description:  "enable registry sync experimental feature",
		DefaultValue: "false",
	},
}

const (
	PromptConfigKey                   = "PROMPT"
	RuntimeRootConfigKey              = "RUNTIME_ROOT"
	DataRootConfigKey                 = "DATA_ROOT"
	BuildahFormatConfigKey            = "BUILDAH_FORMAT"
	ScpChecksumConfigKey              = "SCP_CHECKSUM"
	RegistrySyncExperimentalConfigKey = "REGISTRY_SYNC_EXPERIMENTAL"
)

func (*envSystemConfig) getValueOrDefault(key string) (*ConfigOption, error) {
	for _, option := range configOptions {
		if option.Key == key {
			if option.OSEnv == "" {
				option.OSEnv = strings.ReplaceAll(strings.ToUpper(constants.AppName+"_"+option.Key), "-", "_")
			}
			if value, ok := os.LookupEnv(option.OSEnv); ok {
				option.DefaultValue = value
			}
			return &option, nil
		}
	}
	return nil, fmt.Errorf("not found config key %s", key)
}

func ConfigOptions() []ConfigOption {
	return configOptions
}
