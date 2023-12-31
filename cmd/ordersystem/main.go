package main

import (
	"database/sql"
	"fmt"
	"log"
	"net"
	"net/http"

	graphql_handler "github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/nimbo1999/20-CleanArch/configs"
	"github.com/nimbo1999/20-CleanArch/internal/event/handler"
	"github.com/nimbo1999/20-CleanArch/internal/infra/database"
	"github.com/nimbo1999/20-CleanArch/internal/infra/graph"
	"github.com/nimbo1999/20-CleanArch/internal/infra/grpc/pb"
	"github.com/nimbo1999/20-CleanArch/internal/infra/grpc/service"
	"github.com/nimbo1999/20-CleanArch/internal/infra/web/webserver"
	"github.com/nimbo1999/20-CleanArch/pkg/events"
	"github.com/streadway/amqp"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// mysql
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	configs, err := configs.LoadConfig(".")
	if err != nil {
		panic(err)
	}

	dbSourceUrl := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?multiStatements=true", configs.DBUser, configs.DBPassword, configs.DBHost, configs.DBPort, configs.DBName)
	db, err := sql.Open(configs.DBDriver, dbSourceUrl)
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	migration := database.NewMigration(db, "mysql")
	err = migration.Up()

	if err != nil {
		log.Fatal(err)
	}

	rabbitMQChannel := getRabbitMQChannel()

	eventDispatcher := events.NewEventDispatcher()
	eventDispatcher.Register("OrderCreated", &handler.OrderCreatedHandler{
		RabbitMQChannel: rabbitMQChannel,
	})

	createOrderUseCase := NewCreateOrderUseCase(db, eventDispatcher)
	listOrderUseCase := NewListOrderUseCase(db)

	webServer := webserver.NewWebServer(configs.WebServerPort)
	webOrderHandler := NewWebOrderHandler(db, eventDispatcher)
	err = webServer.AddHandler("/order", &webserver.WebServerHandler{
		Handler: webOrderHandler.Create,
		Method:  http.MethodPost,
	})
	if err != nil {
		panic(err)
	}
	err = webServer.AddHandler("/order", &webserver.WebServerHandler{
		Handler: webOrderHandler.List,
		Method:  http.MethodGet,
	})
	if err != nil {
		panic(err)
	}
	fmt.Println("Starting web server on port", configs.WebServerPort)
	go webServer.Start()

	grpcServer := grpc.NewServer()
	createOrderService := service.NewOrderService(*createOrderUseCase, *listOrderUseCase)
	pb.RegisterOrderServiceServer(grpcServer, createOrderService)
	reflection.Register(grpcServer)

	fmt.Println("Starting gRPC server on port", configs.GRPCServerPort)
	lis, err := net.Listen("tcp", fmt.Sprintf(":%s", configs.GRPCServerPort))
	if err != nil {
		panic(err)
	}
	go grpcServer.Serve(lis)

	srv := graphql_handler.NewDefaultServer(graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
		CreateOrderUseCase: *createOrderUseCase,
		ListOrderUseCase:   *listOrderUseCase,
	}}))
	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	fmt.Println("Starting GraphQL server on port", configs.GraphQLServerPort)
	http.ListenAndServe(":"+configs.GraphQLServerPort, nil)
}

func getRabbitMQChannel() *amqp.Channel {
	conn, err := amqp.Dial("amqp://guest:guest@localhost:5672/")
	if err != nil {
		panic(err)
	}
	ch, err := conn.Channel()
	if err != nil {
		panic(err)
	}
	return ch
}
