package main

import (
	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"github.com/subosito/gotenv"
	"os"
	"todo"
	"todo/pkg/handler"
	"todo/pkg/repository"
	"todo/pkg/service"
)

func main() {
	log.SetFormatter(new(log.JSONFormatter))
	if err := initConfig(); err != nil {
		log.Fatalf("erroc occured while initizialising configs %s ", err.Error())
	}
	if err := gotenv.Load(); err != nil {
		log.Fatalf("error loading env variables %s", err.Error())
	}
	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode")})
	if err != nil {
		log.Fatalf("failed to initialize db: %s", err.Error())
	}
	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := new(todo.Server)
	if err := srv.Run(viper.GetString("port"), handlers.InitRoutes()); err != nil {
		log.Fatalf("error occured while running http server: %s", err.Error())
	}

}

func initConfig() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
