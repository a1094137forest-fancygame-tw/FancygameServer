package config

import "github.com/spf13/viper"

var (
	RunMode string
	Port    string

	UserServerUrl      string
	CacheServerUrl     string
	LobbyServerUrl     string
	BackendServerUrl   string
	GameServerUrl      string
	GameLogicServerUrl string

	MongoDBName          string
	MongoDBUser          string
	MongoDBPassword      string
	MongoDBMaster        string
	MongoDBAuthMechanism string
	MongoDBAuthSource    string
)

func Initialize() {
	RunMode = viper.GetString("RUN_MODE")
	Port = viper.GetString("PORT")

	UserServerUrl = viper.GetString("USER_SERVER_URL")
	CacheServerUrl = viper.GetString("CACHE_SERVER_URL")
	LobbyServerUrl = viper.GetString("LOBBY_SERVER_URL")
	BackendServerUrl = viper.GetString("BACKEND_SERVER_URL")
	GameServerUrl = viper.GetString("GAME_SERVER_URL")
	GameLogicServerUrl = viper.GetString("GAME_LOGIC_SERVER_URL")

	MongoDBMaster = viper.GetString("MONGO_DB_MASTER")
	MongoDBName = viper.GetString("MONGO_DB_NAME")
	MongoDBUser = viper.GetString("MONGO_DB_USER")
	MongoDBPassword = viper.GetString("MONGO_DB_PASSWORD")
	MongoDBAuthMechanism = viper.GetString("MONGO_DB_AUTHMECHANISM")
	MongoDBAuthSource = viper.GetString("MONGO_DB_AUTHSOURCE")

}
