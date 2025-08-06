package database

import (
	"encoding/json"
	"fmt"
	"github.com/go-kratos/kratos/v2/log"
	"gopkg.in/yaml.v3"
	"os"
	"strings"
	"sync"
)

var (
	configData map[string]interface{}
	loadOnce   sync.Once
	loadErr    error
	configPath string
)

// InitConfigPath 设置配置路径（可选）
func InitConfigPath(path string) {
	configPath = path
}

// loadYAML 加载 YAML 到 map
func loadYAML(filePath string) (map[string]interface{}, error) {
	yamlFile, err := os.ReadFile(filePath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var raw map[string]interface{}
	err = yaml.Unmarshal(yamlFile, &raw)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshal yaml: %v", err)
	}

	b, _ := json.MarshalIndent(raw, "", "  ")
	log.Info("loadYAML: ", string(b))
	converted, ok := toStringKeyMap(raw).(map[string]interface{})
	if !ok {
		return nil, fmt.Errorf("failed to convert config to map[string]interface{}")
	}

	return converted, nil
}

// GetConfig 获取全局配置 map，只加载一次
func GetConfig() (map[string]interface{}, error) {
	loadOnce.Do(func() {
		if configPath == "" {
			configPath = "config.yaml"
		}
		configData, loadErr = loadYAML(configPath)
	})
	return configData, loadErr
}

// Get 获取嵌套配置，例如 Get("database.driver")
func Get(key string) (interface{}, error) {
	cfg, err := GetConfig()
	if err != nil {
		return nil, err
	}

	keys := strings.Split(key, ".")
	var val interface{} = cfg

	for _, k := range keys {
		m, ok := val.(map[string]interface{})
		if !ok {
			return nil, fmt.Errorf("invalid path: %s", key)
		}
		val, ok = m[k]
		if !ok {
			return nil, fmt.Errorf("key not found: %s", key)
		}
	}
	return val, nil
}

func toStringKeyMap(i interface{}) interface{} {
	switch x := i.(type) {
	case map[interface{}]interface{}:
		m2 := map[string]interface{}{}
		for k, v := range x {
			keyStr := fmt.Sprintf("%v", k)
			m2[keyStr] = toStringKeyMap(v)
		}
		return m2
	case []interface{}:
		for i, v := range x {
			x[i] = toStringKeyMap(v)
		}
	}
	return i
}
