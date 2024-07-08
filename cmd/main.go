package main

import (
	"context"
	"fliqt/pkg/common"
	"fliqt/pkg/common/db"
	handler "fliqt/pkg/http"
	"fliqt/pkg/repo"
	"fliqt/pkg/svc"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/charmbracelet/log"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func init() {
	var cfgName string
	cfgName = "app-local"

	env := os.Getenv("ENV")
	log.Info("ENV:", env)

	if env == "stg" {
		cfgName = "app-stg"
		gin.SetMode(gin.ReleaseMode)
	}

	viper.SetConfigName(cfgName)
	viper.SetConfigType("yaml")
	viper.AddConfigPath("./config")
	err := viper.ReadInConfig()
	if err != nil {
		log.Error("[main] Read config failed ===>")
		panic(err)
	}

	log.SetLevel(log.DebugLevel)
}

func main() {
	router := gin.Default()

	store, err := redis.NewStore(10, "tcp", viper.GetString("redis.address"), viper.GetString("redis.password"), []byte(common.FLIQT_CONST))
	if err != nil {
		log.Error("[main] Init Redis failed ===>")
		panic(err)
	}
	store.Options(sessions.Options{MaxAge: 3600})

	readDB, writeDB, err := db.InitDatabases()
	if err != nil {
		log.Error("[main] Init DB failed ===>")
		panic(err)
	}
	log.Info("[main] Init DB complete...")

	userRepo := repo.NewUserRepository(readDB, writeDB)
	userSvc := svc.NewUserService(userRepo)

	authSvc := svc.NewAuthService(userRepo)

	handler.NewHandler(router, store, userSvc, authSvc)

	httpServer := &http.Server{
		Addr:         fmt.Sprintf(":%s", viper.GetString("http.port")),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		Handler:      router,
	}

	go func() {
		log.Infof("[main] Listening and serving HTTP on %s", httpServer.Addr)
		err := httpServer.ListenAndServe()
		if err != nil {
			log.Error("[main] http server listen failed ===>")
			panic(err)
		}
	}()

	stopChan := make(chan os.Signal, 1)
	signal.Notify(stopChan, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	<-stopChan
	log.Info("[main] shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		log.Error("[main] http server shutdown error ===>")
		panic(err)
	} else {
		log.Info("[main] gracefully stopped")
	}
}
