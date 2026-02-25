package config

import (
	"log"

	"github.com/stclaird/Cloud-IP-Address-Ranges/pkg/ipfile"

	"github.com/spf13/viper"
)

type DatabaseConfig struct {
	Type     string `mapstructure:"type"`     // "sqlite" or "duckdb"
	Filename string `mapstructure:"filename"` // Database filename
}

type Config struct {
	Dbfile      string                //The name of the database file (deprecated - use Databases)
	Dbdir       string                //The database output directory
	Downloaddir string                //The directory where cloud IP files are downloaded to.
	Ipfiles     []ipfile.DownloadFile //The details of the vendor IP Files to be converted to SQL such as their URLs
	Databases   []DatabaseConfig      //List of databases to generate (SQLite, DuckDB, or both)
}

func addTrailingSlash(strIn string) string {
	//Return a string with trailing slash, return as-is if one already exists.
	lb := strIn[len(strIn)-1:]

	if lb == "/" {
		return strIn
	}
	return strIn + "/"
}

func NewConfig() Config {
	viper.AddConfigPath("./config")
	viper.SetConfigType("yaml")
	if err := viper.ReadInConfig(); err != nil {
		log.Fatalf("Error reading config file: %s", err)
	}

	var ipfiles []ipfile.DownloadFile
	viper.UnmarshalKey("ipfiles", &ipfiles)

	var databases []DatabaseConfig
	viper.UnmarshalKey("databases", &databases)

	dbdir := addTrailingSlash(viper.GetString("dbdir"))
	downloaddir := addTrailingSlash(viper.GetString("downloaddir"))

	// Backward compatibility: if no databases specified, use legacy dbfile
	dbfile := viper.GetString("dbfile")
	if len(databases) == 0 && dbfile != "" {
		databases = []DatabaseConfig{
			{Type: "sqlite", Filename: dbfile},
		}
	}

	newConfig := Config{
		Dbfile:      dbfile,
		Dbdir:       dbdir,
		Downloaddir: downloaddir,
		Ipfiles:     ipfiles,
		Databases:   databases,
	}

	return newConfig

}
