package database

import (
	"database/sql"
	"fmt"
	"log"
	"strings"

	_ "github.com/go-sql-driver/mysql"
)

type DBModel struct {
	DBEngine *sql.DB
	DBInfo   *DBInfo
}

type DBInfo struct {
	DBType   string
	Host     string
	UserName string
	Password string
	Charset  string
}

func NewDBModel(info *DBInfo) *DBModel {
	return &DBModel{DBInfo: info}
}

func (m *DBModel) ConnectDB(db string) error {
	var err error
	s := "%s:%s@tcp(%s)/" + db + "?charset=%s&parseTime=True&loc=Local"
	dsn := fmt.Sprintf(
		s,
		m.DBInfo.UserName,
		m.DBInfo.Password,
		m.DBInfo.Host,
		m.DBInfo.Charset,
	)
	m.DBEngine, err = sql.Open(m.DBInfo.DBType, dsn)
	if err != nil {
		return err
	}

	return nil
}


func (m *DBModel) GeRows() string {
	//rows, err := m.DBEngine.Query("select * from book")
	//查询数据库中所有的表
	//rows, err := m.DBEngine.Query("select * from INFORMATION_SCHEMA.TABLES")

	//rows, err := m.DBEngine.Query("SELECT * FROM INFORMATION_SCHEMA.COLUMNS")

	//查询fullbook数据库中所有的表
	//query:="SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE  TABLE_SCHEMA='fullbook' AND TABLE_NAME='book'"

	query := "SELECT * FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_SCHEMA='fullbook'"
	//query:="select * from book"
	rows, err := m.DBEngine.Query(query)
	if err != nil {
		panic(err.Error())
	}

	defer rows.Close()

	columns, err := rows.Columns()
	if err != nil {
		panic(err.Error())
	}

	values := make([]sql.RawBytes, len(columns))

	scanArgs := make([]interface{}, len(values))
	for i := range values {
		scanArgs[i] = &values[i]
	}
	list := "["

	for rows.Next() {
		err = rows.Scan(scanArgs...)
		if err != nil {
			fmt.Println("log:", err)
			panic(err.Error())
		}

		row := "{"
		var value string
		for i, col := range values {
			if col == nil {
				value = "NULL"
			} else {
				value = string(col)
			}

			columName := strings.ToLower(columns[i])

			cell := fmt.Sprintf(`"%v":"%v"`, columName, value)
			row = row + cell + ","
		}
		row = row[0 : len(row)-1]
		row += "}"
		list = list + row + ","

	}
	list = list[0 : len(list)-1]
	list += "]"

	return list
}

var username string
var password string
var host string
var charset string
var dbType string
var dbName string
var tableName string

func GetRun() {
	username = "root"
	password = "root"
	host = "127.0.0.1:3306"
	charset = "utf8mb4"
	dbType = "mysql"
	dbName = "fullbook"
	tableName = "book"
	dbInfo := &DBInfo{
		DBType:   dbType,
		Host:     host,
		UserName: username,
		Password: password,
		Charset:  charset,
	}
	dbModel := NewDBModel(dbInfo)
	err := dbModel.ConnectDB("fullbook")
	if err != nil {
		log.Fatalf("db.Connect err: %v", err)
	}
	rows := dbModel.GeRows()
	if err != nil {
		log.Fatalf("dbModel.GetColumns err: %v", err)
	}
	//jsonData, err := json.Marshal(columns)
	fmt.Println("cccccc", rows)
}
