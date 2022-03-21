package app

import (
	"context"
	"fmt"
	"github.com/Shemistan/Blog/docs"
	"github.com/Shemistan/Blog/internal/app/repo"
	"github.com/jmoiron/sqlx"
	"github.com/mvrilo/go-redoc"
	"log"
	"net"
	"net/http"
	"os"
	"runtime"
	"sync"

	"github.com/BurntSushi/toml"
	blog_v1 "github.com/Shemistan/Blog/internal/app/api/blog.v1"
	"github.com/Shemistan/Blog/internal/app/config"
	blog_system "github.com/Shemistan/Blog/internal/app/service/blog.system"
	"github.com/Shemistan/Blog/internal/app/service/logger"
	pb "github.com/Shemistan/Blog/pkg/blog.v1"
	gateway_runtime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	_ "github.com/jackc/pgx/stdlib"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	grpcServer *grpc.Server
	mux        *gateway_runtime.ServeMux

	db     *sqlx.DB
	dbInfo string

	loggerService *logger.Service

	BlogSystemService blog_system.IBlogSystemService

	stdOut          *logrus.Logger
	stdErr          *logrus.Logger
	outFile         *os.File
	errFile         *os.File
	loggerFormatter *logrus.TextFormatter
	stdOutLevel     *logrus.Level
	stdErrLevel     *logrus.Level

	reDoc redoc.Redoc

	appConfig *config.Config
}

func NewApp(ctx context.Context, configPath string) (*App, error) {
	runtime.GOMAXPROCS(runtime.NumCPU())

	a := &App{}

	a.initConfig(configPath)
	a.initDB()
	a.initLogger()
	a.initReDoc()
	a.initGRPCServer()
	if err := a.initHTTPServer(ctx); err != nil {
		return nil, err
	}

	return a, nil
}

func (a *App) Run() error {
	wg := sync.WaitGroup{}
	wg.Add(3)

	go func() {
		defer func() {
			a.Close()
			wg.Done()
		}()
		log.Fatal(a.runGRPC())
	}()

	go func() {
		defer wg.Done()
		log.Fatal(a.runHTTP())
	}()

	go func() {
		defer wg.Done()
		log.Fatal(a.runDocumentation())
	}()

	wg.Wait()
	return nil
}

func (a *App) initConfig(configPath string) {
	_, err := toml.DecodeFile(configPath, a.GetConfig())
	if err != nil {
		log.Println("Can not find config file, using default values:", err)
		a.SetConfigDefaultParams(a.GetConfig())
	}
}

func (a *App) initLogger() {
	_ = a.GetStdOut()
	_ = a.GetStdErr()
	if a.appConfig.LoggerConf.WriteLoggerInfoInFile {
		a.SetFileToOutLogger(a.GetOutFile())
		a.SetFileToErrLogger(a.GetErrFile())
	}
	a.SetLoggerFormatter(a.GetLoggerFormatter())
	a.SetStdOutLoggerLevel(a.GetOutLoggerLevel())
	a.SetStdErrLoggerLevel(a.GetErrLoggerLevel())
}

func (a *App) initGRPCServer() {
	a.grpcServer = grpc.NewServer()
	pb.RegisterBlogV1Server(
		a.grpcServer,
		&blog_v1.Blog{
			BlogService: a.GetBlogSystemService(),
		},
	)
}

func (a *App) initDB() {
	var err error
	a.SetDBInfo(a.GetDBInfo())
	a.db, err = sqlx.Open("pgx", a.dbInfo)
	if err != nil {
		log.Fatal("failed to opening connection to db:", err.Error())
	}
}

func (a *App) GetLoggerService() *logger.Service {
	if a.loggerService == nil {
		a.loggerService = logger.NewLoggerService(a.GetStdOut(), a.GetStdErr())
	}

	return a.loggerService
}

func (a *App) initHTTPServer(ctx context.Context) error {
	// Параметр отключает тег omitempty из структуры для json ответа, что бы отображались zero value параметры в response
	a.mux = gateway_runtime.NewServeMux(
		gateway_runtime.WithMarshalerOption(
			gateway_runtime.MIMEWildcard, &gateway_runtime.JSONPb{OrigName: true, EmitDefaults: true},
		),
		gateway_runtime.WithIncomingHeaderMatcher(a.GetCustomMatcherHeaders()),
	)

	opts := []grpc.DialOption{grpc.WithTransportCredentials(insecure.NewCredentials())}
	err := pb.RegisterBlogV1HandlerFromEndpoint(ctx, a.mux, a.appConfig.Server.GRPCPort, opts)

	if err != nil {
		return err
	}

	return nil
}

func (a *App) initReDoc() {
	a.reDoc = docs.Initialize()
}

func (a *App) runGRPC() error {
	listener, err := net.Listen("tcp", a.appConfig.Server.GRPCPort)
	if err != nil {
		return err
	}

	a.loggerService.Info("GRPC server running on port:", a.appConfig.Server.GRPCPort)

	return a.grpcServer.Serve(listener)
}

func (a *App) runHTTP() error {

	a.loggerService.Info("HTTP server running on port:", a.appConfig.Server.HTTPPort)

	return http.ListenAndServe(a.appConfig.Server.HTTPPort, a.mux)
}

func (a *App) runDocumentation() error {
	a.loggerService.Info("Swagger documentation running on port:", a.appConfig.Server.DocsPort)

	return http.ListenAndServe(a.appConfig.Server.DocsPort, a.reDoc.Handler())
}

func (a *App) Close() {
	err := a.outFile.Close()
	if err != nil {
		a.stdErr.Error("failed to closing out file", err.Error())
	}

	err = a.errFile.Close()
	if err != nil {
		a.stdErr.Error("failed to closing err file", err.Error())
	}

	err = a.db.Close()
	if err != nil {
		a.stdErr.Error("failed to closing DB", err.Error())
	}
}

func (a *App) GetCustomMatcherHeaders() gateway_runtime.HeaderMatcherFunc {
	return func(key string) (string, bool) {
		switch key {
		case "authorization":
			return key, true
		default:
			return key, false
		}
	}
}

func (a *App) GetBlogSystemService() blog_system.IBlogSystemService {
	if a.BlogSystemService == nil {
		blogRepo := repo.NewRepo(*a.db)
		a.BlogSystemService = blog_system.NewBlogSystemService(
			a.GetLoggerService(),
			a.GetConfig(),
			blogRepo,
		)
	}

	return a.BlogSystemService
}

func (a *App) GetStdOut() *logrus.Logger {
	if a.stdOut == nil {
		a.stdOut = logrus.New()
	}

	return a.stdOut
}

func (a *App) GetStdErr() *logrus.Logger {
	if a.stdErr == nil {
		a.stdErr = logrus.New()
	}

	return a.stdErr
}

func (a *App) GetOutFile() *os.File {
	if a.outFile == nil {
		var err error
		a.outFile, err = os.OpenFile(a.appConfig.LoggerConf.StdoutFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			log.Println("failed to open out log file:", err.Error())
			return nil
		}
	}

	return a.outFile
}

func (a *App) GetErrFile() *os.File {
	if a.errFile == nil {
		var err error
		a.errFile, err = os.OpenFile(a.appConfig.LoggerConf.StderrFileName, os.O_WRONLY|os.O_APPEND|os.O_CREATE, 0644)
		if err != nil {
			log.Println("filed to open error log file:", err.Error())
			return nil
		}
	}

	return a.errFile
}

func (a *App) GetLoggerFormatter() *logrus.TextFormatter {
	if a.loggerFormatter == nil {
		a.loggerFormatter = new(logrus.TextFormatter)
		a.loggerFormatter.TimestampFormat = a.appConfig.LoggerConf.TimestampFormat
		a.loggerFormatter.FullTimestamp = a.appConfig.LoggerConf.FullTimestamp
	}

	return a.loggerFormatter
}

func (a *App) GetOutLoggerLevel() *logrus.Level {
	if a.stdOutLevel == nil {
		level, err := logrus.ParseLevel(a.appConfig.LoggerConf.StderrLoggerLevelValue)
		if err != nil {
			log.Println("failed to setting out logger level: ", err.Error())
			return nil
		}
		a.stdOutLevel = &level
	}

	return a.stdOutLevel
}

func (a *App) GetErrLoggerLevel() *logrus.Level {
	if a.stdErrLevel == nil {
		level, err := logrus.ParseLevel(a.appConfig.LoggerConf.StdoutLoggerLevelValue)
		if err != nil {
			log.Println("failed to setting err logger level: ", err.Error())
			return nil
		}
		a.stdErrLevel = &level
	}

	return a.stdErrLevel
}

func (a *App) GetConfig() *config.Config {
	if a.appConfig == nil {
		a.appConfig = config.NewConfig()
	}

	return a.appConfig
}

func (a *App) GetDBInfo() string {
	if a.dbInfo == "" {
		a.dbInfo = fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
			a.appConfig.DB.Host,
			a.appConfig.DB.Port,
			a.appConfig.DB.User,
			a.appConfig.DB.Password,
			a.appConfig.DB.DBName,
			a.appConfig.DB.SSLMode,
		)
	}

	return a.dbInfo
}

func (a *App) SetFileToOutLogger(file *os.File) {
	a.stdOut.SetOutput(file)
}

func (a *App) SetFileToErrLogger(file *os.File) {
	a.stdErr.SetOutput(file)
}

func (a *App) SetLoggerFormatter(loggerFormatter *logrus.TextFormatter) {
	a.stdOut.SetFormatter(loggerFormatter)
	a.stdErr.SetFormatter(loggerFormatter)
}

func (a *App) SetStdOutLoggerLevel(level *logrus.Level) {
	a.stdOut.SetLevel(*level)
}

func (a *App) SetStdErrLoggerLevel(level *logrus.Level) {
	a.stdErr.SetLevel(*level)
}

func (a *App) SetDBInfo(dbInfo string) {
	a.dbInfo = dbInfo
}

func (a *App) SetConfigDefaultParams(appConfig *config.Config) {
	log.Println("setting default params")
	appConfig.Server.GRPCPort = config.GrpcPort
	appConfig.Server.HTTPPort = config.HttpPort

	appConfig.LoggerConf.StderrFileName = config.StdoutFileName
	appConfig.LoggerConf.StderrFileName = config.StderrFileName
	appConfig.LoggerConf.WriteLoggerInfoInFile = config.WriteLoggerInfoInFile
	appConfig.LoggerConf.TimestampFormat = config.TimestampFormat
	appConfig.LoggerConf.StdoutLoggerLevelValue = config.StdoutLoggerLevelValue
	appConfig.LoggerConf.StderrLoggerLevelValue = config.StderrLoggerLevelValue
	appConfig.LoggerConf.FullTimestamp = config.FullTimestamp

	appConfig.DB.Host = config.DBHost
	appConfig.DB.Port = config.DBPort
	appConfig.DB.User = config.DBUser
	appConfig.DB.Password = config.DBPassword
	appConfig.DB.DBName = config.DBName
	appConfig.DB.SSLMode = config.DBSSLMode
}
