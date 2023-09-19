package main

import (
	"context"
	"net/http"
	"nnyd-back/config"
	"nnyd-back/db"
	"nnyd-back/pb/schemas/protos/v1/protosv1connect"
	"nnyd-back/service"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

// リフレクション設定
func newServeMuxWithReflection() *http.ServeMux {
	mux := http.NewServeMux()
	reflector := grpcreflect.NewStaticReflector(
		"schemas.protos.v1.UserService",
		"schemas.protos.v1.PostService",
		"schemas.protos.v1.TotpService",
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
	return mux
}

// インターセプタ設定
func newInterCeptors() connect.Option {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			// ここでヘッダをセットするなど色々処理を書ける
			// req.Header().Set("hoge", "fuga")
			return next(ctx, req)
		})
	}
	return connect.WithInterceptors(connect.UnaryInterceptorFunc(interceptor))
}

func main() {
	config.LoadConfig()

	db.Init()
	defer db.Close()
	db.AutoMigration()

	mux := newServeMuxWithReflection()
	interceptor := newInterCeptors()
	mux.Handle(protosv1connect.NewUserServiceHandler(&service.UserServer{}, interceptor))
	mux.Handle(protosv1connect.NewPostServiceHandler(&service.PostServer{}, interceptor))
	mux.Handle(protosv1connect.NewTotpServiceHandler(&service.TotpServer{}, interceptor))

	portStr := ":" + config.PORT
	http.ListenAndServe(portStr, cors.AllowAll().Handler(h2c.NewHandler(mux, &http2.Server{})))
}
