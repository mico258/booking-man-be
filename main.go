package main

import (
	"fmt"
	"time"

	"github.com/booking-man-be/config"
	"github.com/booking-man-be/handler"
	"github.com/booking-man-be/lib/logger"
	"github.com/booking-man-be/lib/server"
	userPb "github.com/booking-man-be/proto/user"
	"github.com/booking-man-be/user"
	"github.com/gomodule/redigo/redis"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	gormLogger "gorm.io/gorm/logger"
	"gorm.io/gorm/schema"
)

func main() {
	cfg := config.Get()
	// init db
	db := initDB(cfg)
	redis := connectRedis(cfg)

	// init repo
	userRepository := user.NewRepository(db, redis)

	// init service
	userService := user.NewService(userRepository)

	// TODO change port to config
	svc := server.NewService(
		server.GRPCPort("9090"),
		server.RESTPort("80"),
	)

	svc.Init()

	// init handler
	userHandler := handler.NewUserHandler(userService)

	// register handler to grpc and rest
	userPb.RegisterUserServer(svc.Server(), userHandler)
	svc.RegisterRESTHandler(userPb.RegisterUserHandler)

	if err := <-svc.RunServers(); err != nil {
		logger.Fatal(err)
	}

}

func initDB(cfg config.Config) *gorm.DB {
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=%s&parseTime=%+v&loc=%s", cfg.MysqlUsername, cfg.MysqlPassword, cfg.MysqlHost, cfg.MysqlPort, cfg.MysqlDBName, cfg.MysqlCharset, cfg.MysqlParseTime, cfg.MysqlLoc)
	// open mysql connecion
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		Logger: gormLogger.Default.LogMode(gormLogger.LogLevel(cfg.MysqlLogMode)),
		NamingStrategy: schema.NamingStrategy{
			SingularTable: true,
		},
	})
	if err != nil {
		logger.Panic(err)
	}
	// set configuration pooling connection
	mysqlDb, _ := db.DB()
	mysqlDb.SetMaxOpenConns(cfg.MysqlMaxOpenConnection)
	mysqlDb.SetConnMaxLifetime(time.Duration(cfg.MysqlMaxLifetimeConnection) * time.Minute)
	mysqlDb.SetMaxIdleConns(cfg.MysqlMaxIdleConnection)

	return db

}

func connectRedis(cfg config.Config) *redis.Pool {
	// construct URL from host and port
	URL := fmt.Sprintf("%s:%d", cfg.RedisHost, cfg.RedisPort)
	// create redis pooling
	Redis := &redis.Pool{
		MaxIdle:         cfg.RedisMaxIdle,
		MaxActive:       cfg.RedisMaxActive,
		Wait:            cfg.RedisWait,
		MaxConnLifetime: time.Duration(cfg.RedisMaxConnLifetime) * time.Minute,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", URL)
			if err != nil {
				return nil, err
			}
			if len(cfg.RedisPassword) > 0 {
				if _, err := c.Do("AUTH", cfg.RedisPassword); err != nil {
					return nil, err
				}
			}
			return c, nil
		},
	}
	// test initialized connection pooling using PING command
	conn := Redis.Get()
	_, err := conn.Do("PING")
	if err != nil {
		logger.Panicf("[ERR] Redis connection failed, %s", err.Error())
	}
	defer conn.Close()
	return Redis
}
