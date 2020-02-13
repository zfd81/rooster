package conf

import "strings"

const (
	MySql      string = "mysql"
	Oracle     string = "oracle"
	PostgreSQL string = "PostgreSQL"
)

type Config struct {
	Name     string              `toml:"name"`
	Version  string              `toml:"version"`
	Http     Http                `toml:"http"`
	Pagesqls map[string]*PageSql `toml:"pagesqls"`
}

func (c *Config) PageSql(driverName string, sql string) string {
	return c.Pagesqls[driverName].Sql(sql)
}

type Http struct {
	Port int `toml:"port"`
}

type PageSql struct {
	Template string `toml:"template"`
}

func (p *PageSql) Sql(sql string) string {
	return strings.Replace(p.Template, "$sql", strings.TrimSpace(sql), 1)
}

var defaultConf = Config{
	Name:    "Rooster",
	Version: "1.0.0",
	Http: Http{
		Port: 8143,
	},
	Pagesqls: map[string]*PageSql{
		MySql: &PageSql{
			"select * from ($sql) _init limit ${(_pageNumber - 1) * _pageSize} , ${_pageSize}",
		},
		Oracle: &PageSql{
			"select * from ( select rownum num,init.* from ($sql) init where rownum <= ${_pageNumber * _pageSize}) where num > ${(_pageNumber - 1) * _pageSize}",
		},
		PostgreSQL: &PageSql{
			"select * from ($sql) init limit ${_pageSize} offset ${(_pageNumber - 1) * _pageSize}",
		},
	},
}

var globalConf = defaultConf

func GetGlobalConfig() *Config {
	return &globalConf
}
