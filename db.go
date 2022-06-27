package main

import (
	"database/sql"
	"encoding/json"
    _ "github.com/go-sql-driver/mysql"
)

// 定义一个全局对象db
var db *sql.DB

// 定义一个初始化数据库的函数
func initDB() (err error) {
	mysqlHost := MysqlServerSetting.HostIp
	mysqlUser := MysqlServerSetting.User
	mysqlPassword := MysqlServerSetting.Password
	mysqlDb := MysqlServerSetting.Db
	mysqlPort := MysqlServerSetting.Port

    dsn := mysqlUser + ":" + mysqlPassword + "@tcp(" + mysqlHost + ":" + mysqlPort + ")/" + mysqlDb + "?charset=utf8&parseTime=True"
    // 不会校验账号密码是否正确
    db, err = sql.Open("mysql", dsn)
    if err != nil {
        return err
    }
    // 尝试与数据库建立连接（校验dsn是否正确）
    err = db.Ping()
    if err != nil {
        return err
    }
    return nil
}

func getJSON(sqlString string) (string, error) {
	stmt, err := db.Prepare(sqlString)
	if err != nil {
		return "", err
	}
	defer stmt.Close()
	rows, err := stmt.Query()
	if err != nil {
		return "", err
	}
	defer rows.Close()
	columns, err := rows.Columns()
	if err != nil {
		return "", err
	}
	count := len(columns)
	tableData := make([]map[string]interface{}, 0)
	values := make([]interface{}, count)
	valuePtrs := make([]interface{}, count)
	for rows.Next() {
		for i := 0; i < count; i++ {
			valuePtrs[i] = &values[i] 
		}
		rows.Scan(valuePtrs...)
		entry := make(map[string]interface{})
		for i, col := range columns {
			var v interface{}
			val := values[i] 
			b, ok := val.([]byte)
			if ok {
				v = string(b)
			} else {
				v = val
			}
			entry[col] = v
		}
		tableData = append(tableData, entry)
	}
	jsonData, err := json.Marshal(tableData)
	if err != nil {
		return "", err
	}
	return string(jsonData), nil
}
	