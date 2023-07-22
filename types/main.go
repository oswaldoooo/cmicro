// common types between different process
package types

type NetData_one struct {
	Path    string `json:"path" yaml:"path" xml:"path"`
	Header  string `json:"header" yaml:"header" xml:"header"`
	Content string `json:"content" yaml:"content" xml:"content"`
}

type NetData_two struct {
	Code    int               `json:"code" yaml:"code" xml:"code"`
	Args    map[string]string `json:"args" yaml:"args" xml:"args"`
	Header  string            `json:"header" yaml:"header" xml:"header"`
	Content string            `json:"content" yaml:"content" xml:"content"`
}

type DataBase_Info struct {
	Name     string `json:"name" yaml:"name" xml:"name"`
	Address  string `json:"address" yaml:"address" xml:"address"`
	Port     int    `json:"port" yaml:"port" xml:"port"`
	UserName string `json:"username" yaml:"username" xml:"username"`
	Passwd   string `json:"password" yaml:"password" xml:"password"`
	Db       string `json:"db" yaml:"db" xml:"db"`
	PoolSize uint   `json:"poolsize" yaml:"poolsize" xml:"poolsize"`
}

type NetService struct {
	Address string `json:"address" yaml:"address" xml:"address"`
	Type    int8   `json:"type" yaml:"type" xml:"type"`
}
