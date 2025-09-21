package entity

type MysqlDBConnOption struct {
	URL                 string
	MaxIdleConn         string
	MaxOpenConn         string
	MaxLifetimeInMinute string
}