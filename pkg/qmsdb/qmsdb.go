package qmsdb


import (
	"database/sql"
	"fmt"
	"github.com/jmoiron/sqlx"
	_ "github.com/wendal/go-oci8"
	//_ "github.com/godror/godror"
	"os"
	"log"
)
var qmsdb *sqlx.DB

func Init() func() {
	os.Setenv("NLS_LANG","SIMPLIFIED CHINESE_CHINA.UTF8")

	var err error
	// [username/password@sid] or [username/password@host:port/service_name] for oracle10 and higher
	qmsdb, err = sqlx.Open("oci8", "qms/system@MG")
	//qmsdb, err = sqlx.Open("oci8", "qms/system@172.27.36.1:1521/qmsdb")
	//qmsdb, err := sql.Open("godror", `user="qms" password="system" connectString="172.27.36.1:1521/qmsdb" libDir=""`)
	if err != nil {
		panic(err.Error())
		log.Fatalln("Initializing qms database failed: ", fmt.Sprintf("%s", err.Error()))
	}

	cleanFunc := func() {
		err := qmsdb.Close()
		if err != nil {
			log.Fatalln("Qms database close error: ", fmt.Sprintf("%s", err.Error()))
		}
	}
	return cleanFunc

}

func Test() {

	rows, err := qmsdb.Query("SELECT count(*) FROM HSM2_TEMPO_DATA")
	if err != nil {
		log.Fatal(err)
	}
	for rows.Next() {
		var name string
		rows.Scan(&name)
		log.Printf("Name = %s, len=%d", name, len(name))
	}
	rows.Close()
}

func Select(dest interface{}, query string, args ...interface{}) error {
	return qmsdb.Select(dest, query, args...)
}

func Get(dest interface{}, query string, args ...interface{}) error {
	return qmsdb.Get(dest, query, args...)
}

func Query(query string, args ...interface{}) (*sql.Rows, error) {
	return qmsdb.Query(query, args...)
}