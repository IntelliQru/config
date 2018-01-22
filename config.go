package config

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"io/ioutil"
	"os"
	"path"

	"github.com/DisposaBoy/JsonConfigReader"
	"github.com/hashicorp/vault/api"
)

type VaultConnection struct {
	Address    string
	Token      string
	Path       string
	ConfigName string
}

type Config struct {
	data    map[string]interface{}
	workDir string
}

func NewConfig() (*Config, error) {

	workDir, err := os.Getwd()
	if err != nil {
		return nil, err
	}

	config := &Config{
		data:    make(map[string]interface{}),
		workDir: workDir,
	}

	return config, nil
}

func (c *Config) AddFromVault(conn *VaultConnection) error {

	vaultClientConfig := api.DefaultConfig()
	client, err := api.NewClient(vaultClientConfig)
	if err != nil {
		return err
	}

	client.SetAddress(conn.Address)
	client.SetToken(conn.Token)
	secret, err := client.Logical().Read(conn.Path)
	if err != nil {
		return err
	}

	if secret == nil {
		return errors.New("Empty vault result.")
	}

	for key, value := range secret.Data {
		keyName := conn.ConfigName + key
		c.data[keyName] = value
	}

	return nil
}

func (c *Config) Str(key string) (res string) {

	res, _ = c.data[key].(string)
	return
}

func (c *Config) Bool(key string) (res bool) {

	res, _ = c.data[key].(bool)
	return
}

func (c *Config) Int(key string) (res int) {

	f, ok := c.data[key].(float64)
	if ok {
		res = int(f)
	}
	return
}

func (c *Config) Int64(key string) (res int64) {

	f, ok := c.data[key].(float64)
	if ok {
		res = int64(f)
	}
	return
}

func (c *Config) Float64(key string) (res float64) {

	res, _ = c.data[key].(float64)
	return
}

func (c *Config) Array(key string) (res []interface{}) {

	res, _ = c.data[key].([]interface{})
	return
}

func (c *Config) ArrayStr(key string) (res []string) {

	res = make([]string, 0)
	arr, ok := c.data[key].([]interface{})
	if !ok {
		return
	}

	for _, val := range arr {
		str, ok := val.(string)
		if ok {
			res = append(res, str)
		}
	}

	return
}

func (c *Config) Map(key string) (res map[string]interface{}) {

	res, _ = c.data[key].(map[string]interface{})
	return
}

func (c *Config) MapStr(key string) (res map[string]string) {

	res = make(map[string]string, 0)
	m, ok := c.data[key].(map[string]interface{})

	if !ok {
		return
	}

	for key, val := range m {
		if strVal, ok := val.(string); ok {
			res[key] = strVal
		}
	}

	return
}

func (c *Config) ReadConfig() error {

	pathToConfig := path.Join(c.workDir, "config.json")
	configData, err := ioutil.ReadFile(pathToConfig)
	if err != nil {
		return err
	}

	reader := JsonConfigReader.New(bytes.NewReader(configData))
	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, reader); err != nil {
		return err
	}

	return json.Unmarshal(buf.Bytes(), &c.data)
}
