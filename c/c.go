package c

import (
	"gopkg.in/yaml.v2"
	"io/ioutil"
)

/* Example conf
rifts:
  - pull: 0.0.0.0:13100
    push:
      - URL: 10.1.10.42:13109
        sampleRate: 100
      - URL: 10.1.10.72:13100
        sampleRate: 100
*/

type PullConf struct {
	URL        string `yaml:"URL"`
	SampleRate int    `yaml:"sampleRate"`
}

type Conf struct {
	Rifts []struct {
		Pull string     `yaml:"pull"`
		Name string     `yaml:"name"`
		Push []PullConf `yaml:"push"`
	} `yaml:"rifts"`
}

func GetConf(path string) (c Conf, err error) {
	data, err := ioutil.ReadFile(path)
	if nil != err {
		return c, err
	}
	err = yaml.Unmarshal([]byte(data), &c)
	return c, err
}
