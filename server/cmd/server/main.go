package main

import (
	"context"
	"flag"
	"github.com/go-kratos/kratos/v2"
	"github.com/go-kratos/kratos/v2/log"
	"github.com/go-kratos/kratos/v2/transport/grpc"
	"github.com/go-kratos/kratos/v2/transport/http"
	"os"
	"vgpu/internal/conf"
	"vgpu/internal/database"

	_ "go.uber.org/automaxprocs"
)

// go build -ldflags "-X main.Version=x.y.z"
var (
	// Name is the name of the compiled software.
	Name string
	// Version is the version of the compiled software.
	Version string
	// flagconf is the config flag.
	flagconf string

	id, _ = os.Hostname()
)

func init() {
	flag.StringVar(&flagconf, "conf", "../../config/config.yaml", "config path, eg: -conf config.yaml")
}

func main() {
	flag.Parse()
	var ctx = context.Background()

	if err := initDatabase(); err != nil {
		log.Errorf("数据库初始化失败: %v", err)
		os.Exit(1)
	}

	app, cleanup, err := initApp(flagconf, ctx)
	if err != nil {
		panic(err)
	}
	defer cleanup()

	// start and wait for stop signal
	if err := app.Run(); err != nil {
		panic(err)
	}
}

func newApp(ctx context.Context, logger log.Logger, gs *grpc.Server, hs *http.Server) *kratos.App {
	return kratos.New(
		kratos.Context(ctx),
		kratos.ID(id),
		kratos.Name(Name),
		kratos.Version(Version),
		kratos.Metadata(map[string]string{}),
		kratos.Logger(logger),
		kratos.Server(
			gs,
			hs,
		),
	)
}

func getNodeSelectors(c *conf.Bootstrap) map[string]string {
	return c.NodeSelectors
}

func initDatabase() error {
	config, err := database.LoadConfig(flagconf)
	log.Infof("config: %+v", config)
	if err != nil {
		log.Errorf("Failed to load config: %v", err)
		return err
	}
	database.InitDB(&config.Database)
	log.Infof("初始化%s成功", config.Database.Driver)
	return nil
}
