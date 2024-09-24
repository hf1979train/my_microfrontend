package myenv

type Env struct {
	Server Server
	Mysql  Mysql
	Gmail  Gmail
}

type Server struct {
	Host string
	Port string
}

type Mysql struct {
	User             string
	Password         string
	Host             string
	Port             string
	TableNameBudgeet string
}

type Gmail struct {
	PathCredentials string
	PathToken       string
}
