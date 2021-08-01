/*
为兼容linux系统
godotenv在windows下加载env file会将环境变量名转为大写
这里在linux下将变量名也转为大写, 以便配合github.com/timest/env库使用
*/

package gcommon

import (
	"os"
	"strings"

	"github.com/timest/env"
	"github.com/joho/godotenv"
)

var (
	// EnvFill env Fill func
	EnvFill = env.Fill
	// EnvIgnorePrefix  env IgnorePrefix func
	EnvIgnorePrefix = env.IgnorePrefix
)

func readFile(filename string) (envMap map[string]string, err error) {
	file, err := os.Open(filename)
	if err != nil {
		return
	}
	defer file.Close()

	return godotenv.Parse(file)
}

// LoadEnvFile load env file
func LoadEnvFile(filename string, overload bool) error {
	envMap, err := readFile(filename)
	if err != nil {
		return err
	}

	currentEnv := map[string]bool{}
	rawEnv := os.Environ()
	for _, rawEnvLine := range rawEnv {
		key := strings.Split(rawEnvLine, "=")[0]
		currentEnv[key] = true
	}

	for key, value := range envMap {
		key = strings.ToUpper(key)
		if !currentEnv[key] || overload {
			os.Setenv(key, value)
		}
	}

	return nil
}

// LoadEnvFiles load env files
func LoadEnvFiles(filenames ...string) error {
	// 如果没有指定fileame 则加载当前下的.env文件
	if len(filenames) == 0 {
		filenames = append(filenames, ".env")
	}
	for _, filename := range filenames {
		err := LoadEnvFile(filename, false)
		if err != nil {
			return err
		}
	}
	return nil
}

// OverLoadEnvFiles over load env files
// 以覆盖的形式加载
func OverLoadEnvFiles(filenames ...string) error {
	// 如果没有指定fileame 则加载当前下的.env文件
	if len(filenames) == 0 {
		filenames = append(filenames, ".env")
	}
	for _, filename := range filenames {
		err := LoadEnvFile(filename, true)
		if err != nil {
			return err
		}
	}
	return nil
}