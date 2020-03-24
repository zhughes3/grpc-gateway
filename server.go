package main

import (
	"context"
	"database/sql"
	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	"github.com/soheilhy/cmux"
	"github.com/zhughes3/grpc/models"
	"google.golang.org/grpc"
	"log"
	"net"
	"net/http"
	"time"
)


type server struct {
	name string
	lis net.Listener
	grpcServer *grpc.Server
	httpServer *http.Server
	db *db
}

func NewServer(cfg *serverConfig, l net.Listener, d *sql.DB) *server {
	return &server{
		name: cfg.httpHost + ":" + cfg.httpPort,
		lis: l,
		db: &db{
			db: d,
		},
	}
}

func (s *server) Start() error {
	var err error

	ctx := context.Background()
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	// tcpMuxer
	tcpMux := cmux.New(s.lis)

	// connection dispatcher rules
	grpcL := tcpMux.MatchWithWriters(cmux.HTTP2MatchHeaderFieldPrefixSendSettings("content-type", "application/grpc"))
	httpL := tcpMux.Match(cmux.HTTP1Fast())

	//init gprc
	s.grpcServer, err = s.withGRPC(ctx)
	if err != nil {
		log.Fatal("Unable to init gRPC server")
		return err
	}

	s.httpServer, err = prepareHTTP(ctx, s.name)
	if err != nil {
		log.Fatalln("unable to init http server")
		return err
	}

	// start servers
	go func() {
		if err := s.grpcServer.Serve(grpcL); err != nil {
			log.Fatalln("unable to start grpc server")
		}
	}()

	go func() {
		if err := s.httpServer.Serve(httpL); err != nil {
			log.Fatalln("unable to start http server")
		}
	}()

	return tcpMux.Serve()
}

func (s *server) withGRPC(ctx context.Context) (*grpc.Server, error) {
	grpcServer := grpc.NewServer()
	models.RegisterPostServiceServer(grpcServer, s)
	return grpcServer, nil
}

func prepareHTTP(ctx context.Context, name string) (*http.Server, error) {
	router := http.NewServeMux()

	//gateway
	gw, err := prepareGateway(name, ctx)
	if err != nil {
		log.Fatalln("unable to init grpc gateway")
		return nil, err
	}
	router.Handle("/", gw)

	return &http.Server{
		Addr: name,
		Handler: router,
		ReadTimeout: 5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout: 120 * time.Second,
	}, nil
}

func prepareGateway(target string, ctx context.Context) (http.Handler, error) {
	opts := []grpc.DialOption{
		grpc.WithInsecure(),
	}
	conn, err := grpc.Dial(target, opts...)
	if err != nil {
		log.Fatalf("fail to dial: %v", err)
		return nil, err
	}

	gwMux := runtime.NewServeMux(
		runtime.WithMarshalerOption(runtime.MIMEWildcard, &runtime.JSONPb{OrigName: true, EmitDefaults: true}),
		runtime.WithProtoErrorHandler(runtime.DefaultHTTPProtoErrorHandler),
	)

	err = models.RegisterPostServiceHandler(ctx, gwMux, conn)
	if err != nil {
		return nil, err
	}

	return gwMux, nil
}
