package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	core_auth "github.com/alisupurov/todoApp-golang/internal/core/auth"
	core_config "github.com/alisupurov/todoApp-golang/internal/core/config"
	core_logger "github.com/alisupurov/todoApp-golang/internal/core/logger"
	core_pgx_pool "github.com/alisupurov/todoApp-golang/internal/core/repository/postgres/pool/pgx"
	core_goredis_pool "github.com/alisupurov/todoApp-golang/internal/core/repository/redis/pool/goredis"
	core_http_middleware "github.com/alisupurov/todoApp-golang/internal/core/transport/http/middleware"
	core_http_server "github.com/alisupurov/todoApp-golang/internal/core/transport/http/server"
	auth_postgres_repository "github.com/alisupurov/todoApp-golang/internal/features/auth/repository/postgres"
	auth_service "github.com/alisupurov/todoApp-golang/internal/features/auth/service"
	auth_transport_http "github.com/alisupurov/todoApp-golang/internal/features/auth/transport/http"
	statistics_cached_repository "github.com/alisupurov/todoApp-golang/internal/features/statistics/repository/cached"
	statistics_postgres_repository "github.com/alisupurov/todoApp-golang/internal/features/statistics/repository/postgres"
	statistics_service "github.com/alisupurov/todoApp-golang/internal/features/statistics/service"
	statistics_transport_http "github.com/alisupurov/todoApp-golang/internal/features/statistics/transport/http"
	tasks_postgres_repository "github.com/alisupurov/todoApp-golang/internal/features/tasks/repository/postgres"
	tasks_service "github.com/alisupurov/todoApp-golang/internal/features/tasks/service"
	tasks_transport_http "github.com/alisupurov/todoApp-golang/internal/features/tasks/transport/http"
	users_postgres_repository "github.com/alisupurov/todoApp-golang/internal/features/users/repository/postgres"
	users_service "github.com/alisupurov/todoApp-golang/internal/features/users/service"
	users_transport_http "github.com/alisupurov/todoApp-golang/internal/features/users/transport/http"
	web_fs_repository "github.com/alisupurov/todoApp-golang/internal/features/web/repository/file_system"
	web_service "github.com/alisupurov/todoApp-golang/internal/features/web/service"
	web_transport_http "github.com/alisupurov/todoApp-golang/internal/features/web/transport/http"
	"go.uber.org/zap"

	_ "github.com/alisupurov/todoApp-golang/docs"
)

// @title    Golang To do API
// @Description  Todo Application REST-API schema
// @version  1.0
// @BasePath /api/v1
// @host 127.0.0.1:5050
func main() {
	cfg := core_config.NewConfigMust()
	time.Local = cfg.TimeZone

	fmt.Println("Hello, todoApp!")
	ctx, cancel := signal.NotifyContext(
		context.Background(),
		syscall.SIGINT, syscall.SIGTERM,
	)
	defer cancel()

	logger, err := core_logger.NewLogger(core_logger.NewConfigMust())
	if err != nil {
		fmt.Println("failed to init application logger:", err)
		os.Exit(1)
	}
	defer logger.Close()

	logger.Debug("application time zone:", zap.Any("zone", time.Local))

	logger.Debug("initializing postgres connection pool")
	pool, err := core_pgx_pool.NewPool(
		ctx,
		core_pgx_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init postgres connection pool", zap.Error(err))
	}
	defer pool.Close()

	logger.Debug("initializing redis connection pool")
	redisPool, err := core_goredis_pool.NewPool(
		ctx,
		core_goredis_pool.NewConfigMust(),
	)
	if err != nil {
		logger.Fatal("failed to init redis connection pool", zap.Error(err))
	}
	defer redisPool.Close()

	jwtConfig := core_auth.NewConfigMust()
	jwtMiddleware := core_http_middleware.JWT(jwtConfig)

	withJWT := func(routes []core_http_server.Route) []core_http_server.Route {
		for i := range routes {
			routes[i].Middleware = append(routes[i].Middleware, jwtMiddleware)
		}
		return routes
	}

	logger.Debug("initializing feature", zap.String("feature", "auth"))
	authRepository := auth_postgres_repository.NewAuthRepository(pool)
	authService := auth_service.NewAuthService(authRepository, jwtConfig)
	authTransportHTTP := auth_transport_http.NewAuthHTTPHandler(authService)

	logger.Debug("initializing feature", zap.String("feature", "users"))
	usersRepository := users_postgres_repository.NewUsersRepository(pool)
	usersService := users_service.NewUsersService(usersRepository)
	usersTransportHTTP := users_transport_http.NewUsersHTTPHandler(usersService)

	logger.Debug("initializing feature", zap.String("feature", "tasks"))
	tasksRepository := tasks_postgres_repository.NewTasksRepository(pool)
	tasksService := tasks_service.NewTasksService(tasksRepository)
	tasksTransportHTTP := tasks_transport_http.NewTasksHTTPHandler(tasksService)

	logger.Debug("initializing feature", zap.String("feature", "statistics"))
	statisticsRepository := statistics_postgres_repository.NewStatisticsRepository(pool)
	cachedStatisticsRepository := statistics_cached_repository.NewCachedRepository(redisPool, statisticsRepository)
	statisticsService := statistics_service.NewStatisticsService(cachedStatisticsRepository)
	statisticsTransportHTTP := statistics_transport_http.NewStatisticsHTTPHandler(statisticsService)

	logger.Debug("initializing feature", zap.String("feature", "web"))
	webRepository := web_fs_repository.NewWebRepository()
	webService := web_service.NewWebService(webRepository)
	webTransportHTTP := web_transport_http.NewWebHTTPHandler(webService)

	logger.Debug("initializing HTTP server")
	httpServer := core_http_server.NewHTTPServer(
		core_http_server.NewConfigMust(),
		logger,
		core_http_middleware.CORS(),
		core_http_middleware.RequestID(),
		core_http_middleware.Logger(logger),
		core_http_middleware.Trace(),
		core_http_middleware.Panic(),
	)

	// Все роуты под /api/v1.
	// Auth роуты — публичные (без JWT).
	// Остальные — защищённые (с JWT в Middleware каждого Route).
	apiRouter := core_http_server.NewApiVersionRouter(core_http_server.ApiVersion1)
	apiRouter.RegisterRoutes(authTransportHTTP.Routes()...)
	apiRouter.RegisterRoutes(withJWT(usersTransportHTTP.Routes())...)
	apiRouter.RegisterRoutes(withJWT(tasksTransportHTTP.Routes())...)
	apiRouter.RegisterRoutes(withJWT(statisticsTransportHTTP.Routes())...)

	httpServer.RegisterApiRouters(apiRouter)
	httpServer.RegisterRoutes(webTransportHTTP.Routes()...)

	httpServer.RegisterSwagger()

	if err := httpServer.Run(ctx); err != nil {
		logger.Error("HTTP server run error", zap.Error(err))
	}
}
