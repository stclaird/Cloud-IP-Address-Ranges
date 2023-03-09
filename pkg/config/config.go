package config

import (
	"log"

	"github/stclaird/cloudIPtoDB/pkg/ipfile"

	"github.com/spf13/viper"
)

type Config struct {
	Dbfile      string                //The name of the database file.
	Dbdir       string                //The database output directory
	Downloaddir string                //The directory where cloud IP files are downloaed to.
	Ipfiles     []ipfile.DownloadFile //The details of the vendor IP Files to be converted to SQL such as their URLs
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

	dbdir := addTrailingSlash(viper.GetString("dbdir"))
	downloaddir := addTrailingSlash(viper.GetString("downloaddir"))

	newConfig := Config{
		Dbfile:      viper.GetString("dbfile"),
		Dbdir:       dbdir,
		Downloaddir: downloaddir,
		Ipfiles:     ipfiles,
	}

	return newConfig

}
