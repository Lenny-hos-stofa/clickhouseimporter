package main

var cfg Config
var options Options

//Options read from CLI
type Options struct {
	Configfile string `short:"c" long:"configfile" description:"/path/to/configfile.yml" required:"true"`
	Auditinput string `short:"a" long:"auditinput" description:"/path/to/auditinput.log" required:"true"`
	Dbserver   string `short:"d" long:"dbserver" description:"origin of the logfile (hostname of databaseserver)" required:"true"`
	Verbose    bool   `short:"v" long:"verbose" description:"Show verbose debug information"`
}

//Config A struct for holding config parameters
type Config struct {
	Database struct {
		Hostname  string `yaml:"hostname" envconfig:"DB_HOSTNAME"`
		Dbname    string `yaml:"dbname" envconfig:"DB_DBNAME"`
		Tablename string `yaml:"tablename" envconfig:"TABLENAME"`
		Username  string `yaml:"user" envconfig:"DB_USERNAME"`
		Password  string `yaml:"pass" envconfig:"DB_PASSWORD"`
		Port      string `yaml:"port" envconfig:"DB_PORT"`
	} `yaml:"database"`
	Misc struct {
		Batchsize int `yaml:"batchsize" envconfig:"BATCHSIZE"`
	} `yaml:"misc"`
}

//AuditEntry the struct for a jsonline
type AuditEntry struct {
	AuditRecord struct {
		Name         string `json:"name"`
		Record       string `json:"record"`
		Timestamp    string `json:"timestamp"`
		CommandClass string `json:"command_class"`
		ConnectionID string `json:"connection_id"`
		Status       int    `json:"status"`
		Sqltext      string `json:"sqltext"`
		User         string `json:"user"`
		Host         string `json:"host"`
		OsUser       string `json:"os_user"`
		IP           string `json:"ip"`
		Db           string `json:"db"`
	} `json:"audit_record"`
}
