// Scry Info.  All rights reserved.
// license that can be found in the license file.

package dot

import (
	"gopkg.in/yaml.v3"
	"os"
	"regexp"

	"github.com/scryinfo/scryg/sutils/skit"
)

const (

	//SconfigTypeID scofig dot type id
	SconfigTypeID = "484ef01d-3c04-4517-a643-2d776a9ae758"
)

var reVar = regexp.MustCompile(`^\${(\w+)}$`)

type StringFromEnv string

func (e *StringFromEnv) UnmarshalYAML(value *yaml.Node) error {
	var s string
	if err := value.Decode(&s); err != nil {
		return err
	}
	if match := reVar.FindStringSubmatch(s); len(match) > 0 {
		*e = StringFromEnv(os.Getenv(match[1]))
	} else {
		*e = StringFromEnv(s)
	}
	return nil
}

//SConfig config belongs to one component Dot, but it is so basic, every Dot need it, so define it in dot.go file
//S represents scryinfo config this name is used frequently, so add s to distinguish it
type SConfig interface {
	//RootPath root path
	RootPath()
	//Config file path
	ConfigPath() string
	//Without path, only file name
	ConfigFile() string
	//Whether key existing
	Key(key string) bool
	//If no config or config is empty, return nil
	Map() (m map[string]interface{}, err error)

	//Priorly bring data from RAM and operate, if RAM is nil then check whether original config file existing
	Unmarshal(s interface{}) error
	//Analyze key is corresponding type
	UnmarshalKey(key string, obj interface{}) error

	Marshal(data []byte) error

	//If no corresponding key or data type cannot be converted, must pay attention to default value,
	//so add "Def" before function to notify default value
	DefInterface(key string, def interface{}) interface{}
	DefArray(key string, def []interface{}) []interface{}
	DefMap(key string, def map[string]interface{}) map[string]interface{}
	DefString(key string, def string) string
	DefInt32(key string, def int32) int32
	DefUint32(key string, def uint32) uint32
	DefInt64(key string, def int64) int64
	DefUint64(key string, def uint64) uint64
	DefBool(key string, def bool) bool
	DefFloat32(key string, def float32) float32
	DefFloat64(key string, def float64) float64
}

//UnMarshalConfig unmarshal config
func UnMarshalConfig(conf []byte, obj interface{}) (err error) {
	err = nil
	if conf != nil {
		err = yaml.Unmarshal(conf, obj)
		//err = json.Unmarshal(conf, obj)
	} else {
		err = SError.Parameter
	}
	return err
}

//MarshalConfig marshal config
func MarshalConfig(lconf *LiveConfig) (conf []byte, err error) {
	conf = nil
	err = nil

	if lconf != nil {
		if !skit.IsNil(lconf.Config) {
			return yaml.Marshal(lconf.Config)
			//return json.Marshal(lconf.Config)
		}
	}
	return conf, err
}
