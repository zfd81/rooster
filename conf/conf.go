package conf

import "strings"

const (
	MySql      string = "mysql"
	Oracle     string = "oracle"
	PostgreSQL string = "PostgreSQL"
)

type Config struct {
	Name     string              `mapstructure:"name"`
	Version  string              `mapstructure:"version"`
	Http     Http                `mapstructure:"http"`
	Pagesqls map[string]*PageSql `mapstructure:"pagesqls"`
}

func (c *Config) PageSql(driverName string, sql string) string {
	return c.Pagesqls[driverName].Sql(sql)
}

type Http struct {
	Port int `mapstructure:"port"`
}

type PageSql struct {
	Template string `mapstructure:"template"`
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
			"select * from ( select rownum num,_init.* from ($sql) _init where rownum <= ${_pageNumber * _pageSize}) where num > ${(_pageNumber - 1) * _pageSize}",
		},
		PostgreSQL: &PageSql{
			"select * from ($sql) _init limit ${_pageSize} offset ${(_pageNumber - 1) * _pageSize}",
		},
	},
}

var globalConf = defaultConf

func GetGlobalConfig() *Config {
	return &globalConf
}
