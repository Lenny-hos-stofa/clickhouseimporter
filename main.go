package main

import (
	"bufio"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"time"

	"github.com/ClickHouse/clickhouse-go"
	"github.com/jessevdk/go-flags"
)

func main() {

	if _, err := parser.Parse(); err != nil {

		if flagsErr, ok := err.(*flags.Error); ok && flagsErr.Type == flags.ErrHelp {
			os.Exit(0)
		} else {
			os.Exit(1)
		}
	}

	debug("Loading configfile")
	readFile(&cfg)

	debug("Loading Env")
	readEnv(&cfg)

	debug("Connecting to Clickhouse")
	connect, err := sql.Open("clickhouse", "tcp://"+cfg.Database.Hostname+":"+cfg.Database.Port+"?username="+cfg.Database.Username+"&password="+cfg.Database.Password+"&database="+cfg.Database.Dbname+"&read_timeout=10&write_timeout=20")
	checkErr(err)
	debug("Connection established")

	if err := connect.Ping(); err != nil {
		if exception, ok := err.(*clickhouse.Exception); ok {
			fmt.Printf("[%d] %s \n%s\n", exception.Code, exception.Message, exception.StackTrace)
		} else {
			fmt.Println(err)
		}
		return
	}
	debug("Database seems to be alive")

	debug("Loading Audit log file")

	file, err := os.Open(options.Auditinput)
	if err != nil {
		log.Fatal(err)
	}
	debug("Audit log filed opened")

	defer file.Close()

	debug("Processing lines")

	reader := bufio.NewReader(file)
	for {

		debug("Begin insert transaction")

		tx, err := connect.Begin()
		checkErr(err)

		debug("Preparing")

		stmt, err := tx.Prepare("insert into " + cfg.Database.Tablename + " (time, name,record,command_class,connection_id,status,sqltext,user,host,os_user,ip,db,dbserver) values (?,?,?,?,?,?,?,?,?,?,?,?,?)")
		checkErr(err)

		debug("Parsing lines. Not outputting for each line, simply too much info")

		for i := 0; i <= cfg.Misc.Batchsize; i++ {
			auditline, err := reader.ReadBytes('\n')
			if err != nil {
				if err == io.EOF {
					fmt.Println("And now, the time is near!")
					if i > 0 {
						checkErr(tx.Commit())
					}
					os.Exit(0)
				}
				log.Println(err)
				continue
			}

			var unjson AuditEntry
			if err = json.Unmarshal(auditline, &unjson); err != nil {
				fmt.Println(err)
				continue
			}

			var format string
			if match, _ := regexp.MatchString(`UTC$`, unjson.AuditRecord.Timestamp); match {
				format = "2006-01-02T15:04:05 UTC"
			} else {
				format = "2006-01-02T15:04:05Z"
			}
			timestamp, err := time.Parse(format, unjson.AuditRecord.Timestamp)
			if _, err := stmt.Exec(
				timestamp,
				unjson.AuditRecord.Name,
				unjson.AuditRecord.Record,
				unjson.AuditRecord.CommandClass,
				unjson.AuditRecord.ConnectionID,
				unjson.AuditRecord.Status,
				unjson.AuditRecord.Sqltext,
				unjson.AuditRecord.User,
				unjson.AuditRecord.Host,
				unjson.AuditRecord.OsUser,
				unjson.AuditRecord.IP,
				unjson.AuditRecord.Db,
				options.Dbserver); err != nil {
				log.Fatal(err)
			}
		}
		debug("Committing")

		checkErr(tx.Commit())
	}
}
