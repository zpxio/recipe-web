/*
 * Copyright 2022 zpxio (Jeff Sharpe)
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

package config

import (
	"github.com/apex/log"
	"gopkg.in/yaml.v2"
	"os"
	"path/filepath"
)

type Config struct {
	BaseDir string
	Server  struct {
		LocalIP string `yaml:"localip"`
		Port    int    `yaml:"port"`
	} `yaml:"server"`
	Content struct {
	} `yaml:"content"`
}

func Load(path string) (Config, error) {
	realpath := filepath.Clean(path)
	c := Config{
		BaseDir: filepath.Dir(realpath),
	}

	configData, err := os.Open(realpath)
	if err != nil {
		return c, err
	}
	defer func() {
		err := configData.Close()
		if err != nil {
			log.Errorf("failed to close config file: %s", err)
		}
	}()

	decoder := yaml.NewDecoder(configData)
	err = decoder.Decode(&c)
	if err != nil {
		return c, err
	}

	return c, nil
}
