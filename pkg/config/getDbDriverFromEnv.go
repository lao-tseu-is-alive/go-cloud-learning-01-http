package config

import (
	"fmt"
	"os"
)

//GetDbDriverFromEnv returns the DB driver based on the value of environment variables :
//  DB_DRIVER : string containing the type of storage to use one of (memory|postgres)
func GetDbDriverFromEnv(defaultDbDriver string) string {
	dbDriver := defaultDbDriver
	val, exist := os.LookupEnv("DB_DRIVER")
	if exist {
		dbDriver = val
	}
	return fmt.Sprintf("%s", dbDriver)
}
