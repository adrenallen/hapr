package helpers

import (
	"database/sql"
	"fmt"
	"log"
	"reflect"
	"strings"

	_ "github.com/lib/pq"
)

var globalDB *sql.DB

func NewDatabaseConnection() *sql.DB {
	if globalDB != nil {
		return globalDB
	}
	dbConfig, err := getDatabaseConfig()
	if err != nil {
		log.Fatal(err)
		return nil
	}
	connStr := fmt.Sprintf("user=%v dbname=happy sslmode=disable password=%v host=%v port=%v", dbConfig.User, dbConfig.Pass, dbConfig.Server, dbConfig.Port)
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		log.Fatal(err)
	}

	globalDB = db
	return globalDB
}

var databaseConfig *DatabaseConfig

func getDatabaseConfig() (*DatabaseConfig, error) {
	if databaseConfig != nil {
		return databaseConfig, nil
	}
	databaseConfig := &DatabaseConfig{}

	server, err := GetConfigValue("db.server")
	if err != nil {
		return nil, err
	}
	port, err := GetConfigValue("db.port")
	if err != nil {
		return nil, err
	}
	user, err := GetConfigValue("db.user")
	if err != nil {
		return nil, err
	}
	pass, err := GetConfigValue("db.pass")
	if err != nil {
		return nil, err
	}

	databaseConfig.User = user
	databaseConfig.Server = server
	databaseConfig.Port = port
	databaseConfig.Pass = pass

	return databaseConfig, nil
}

type DatabaseConfig struct {
	Server string
	Port   string
	User   string
	Pass   string
}

const columnTagName = "column"
const encryptTagName = "encrypted"

//Gets the select columns for the provided model
func GetSQLSelectForModel(model interface{}) string {
	return GetSQLSelectForModelWithTableAlias(model, "")
}

func GetSQLSelectForModelWithTableAlias(model interface{}, alias string) string {

	columns := []string{}
	t := reflect.TypeOf(model)
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		column := field.Tag.Get(columnTagName)

		if _, ok := field.Tag.Lookup(encryptTagName); ok {
			column = GetDecryptSQLString(column, alias)
		} else if len(alias) > 0 {
			column = getColumnWithAlias(column, alias)
		}

		columns = append(columns, column)
	}

	return strings.Join(columns, ", ")
}

//returns the string required to encrypt the provided string in a sql query
func GetEncryptSQLString(item string) string {
	dbEncryptKey, err := GetConfigValue("db.encryption_key")
	if err != nil {
		log.Panic(err)
	}
	return fmt.Sprintf(`encrypt(%s::bytea, '%s'::bytea, 'aes'::text)`, item, dbEncryptKey)
}

//Returns the string required to decrypt the provided column ina sql query
func GetDecryptSQLString(column string, alias string) string {
	dbEncryptKey, err := GetConfigValue("db.encryption_key")
	if err != nil {
		log.Panic(err)
	}
	return fmt.Sprintf(`convert_from(decrypt(%s::bytea, '%s', 'aes'), 'SQL_ASCII') as %s`, getColumnWithAlias(column, alias), dbEncryptKey, column)
}

func getColumnWithAlias(column string, alias string) string {
	if len(alias) > 0 {
		return fmt.Sprintf("%s.%s", alias, column)
	}
	return column
}
