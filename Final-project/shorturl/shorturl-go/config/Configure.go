package config

import (
	"errors"
	"strconv"

	"github.com/ewangplay/config"
)

type Configure struct {
	ConfigureMap map[string]string
}

func NewConfigure() (*Configure, error) {
	config := &Configure{}

	config.ConfigureMap = make(map[string]string)
	err := config.ParseConfigure()
	if err != nil {
		return nil, err
	}

	return config, nil
}

func (this *Configure) loopConfigure(sectionName string, cfg *config.Config) error {

	if cfg.HasSection(sectionName) {
		section, err := cfg.SectionOptions(sectionName)
		if err == nil {
			for _, v := range section {
				options, err := cfg.String(sectionName, v)
				if err == nil {
					this.ConfigureMap[v] = options
				}
			}

			return nil
		}
		return errors.New("Parse Error")
	}

	return errors.New("No Section")
}

func (this *Configure) ParseConfigure() error {
	this.ConfigureMap["port"] = "9000"
	this.ConfigureMap["counter"] = "redis"
	this.ConfigureMap["redishost"] = "short-redis"
	this.ConfigureMap["redisport"] = "6379"
	this.ConfigureMap["status"] = "true"

	return nil
}

func (this *Configure) GetPort() (int, error) {

	portstr, ok := this.ConfigureMap["port"]
	if ok == false {
		return 9090, errors.New("No Port set, use default")
	}

	port, err := strconv.Atoi(portstr)
	if err != nil {
		return 9090, err
	}

	return port, nil
}

func (this *Configure) GetRedisHost() (string, error) {
	redishost, ok := this.ConfigureMap["redishost"]

	if ok == false {
		return "127.0.0.1", errors.New("No redishost,use defualt")
	}

	return redishost, nil
}

func (this *Configure) GetRedisPort() (string, error) {
	redisport, ok := this.ConfigureMap["redisport"]

	if ok == false {
		return "6379", errors.New("No redisport,use defualt")
	}

	return redisport, nil
}

func (this *Configure) GetRedisStatus() bool {

	status, ok := this.ConfigureMap["status"]
	if ok == false {
		return true
	}

	if status == "true" {
		return true
	}
	return false

}

func (this *Configure) GetCounterType() string {

	count_type, ok := this.ConfigureMap["counter"]
	if ok == false {
		return "inner"
	}

	return count_type

}
