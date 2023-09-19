package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"nnyd-back/config"
	"nnyd-back/db"
	protosv1 "nnyd-back/pb/schemas/protos/v1"
	"nnyd-back/pb/schemas/protos/v1/protosv1connect"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"github.com/rs/cors"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type UserServer struct{}

func (s *UserServer) CreateUser(ctx context.Context, req *connect.Request[protosv1.CreateUserRequest]) (*connect.Response[protosv1.CreateUserResponse], error) {
	log.Println("Request headers: ", req.Header())

	if req.Msg.Name == "" || req.Msg.Icon == "" {
		// エラーにステータスコードを追加
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("name and icon is required."))
	}

	conn := db.GetDB()

	u := db.Users{
		Uid:  req.Msg.Name,
		Name: req.Msg.Name,
		Icon: req.Msg.Icon,
	}

	if err := conn.Create(&u).Error; err != nil {
		resp := connect.NewError(connect.CodeInternal, err)
		log.Fatal(err)
		return nil, resp
	}

	responseUser := &protosv1.User{
		Id:   "dummy",
		Name: req.Msg.GetName(),
		Icon: req.Msg.GetIcon(),
	}
	userResp := &protosv1.CreateUserResponse{
		User: responseUser,
	}
	resp := connect.NewResponse(userResp)
	return resp, nil
}

// リフレクション設定
func newServeMuxWithReflection() *http.ServeMux {
	mux := http.NewServeMux()
	reflector := grpcreflect.NewStaticReflector(
		"schemas.protos.v1.UserService",
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
	userServer := &UserServer{}

	if config.ENV == "develop" {
		db.Init()
		defer db.Close()
		db.AutoMigration()
	}

	mux := newServeMuxWithReflection()
	interceptor := newInterCeptors()
	path, handler := protosv1connect.NewUserServiceHandler(userServer, interceptor)
	mux.Handle(path, handler)

	portStr := ":" + config.PORT
	http.ListenAndServe(portStr, cors.AllowAll().Handler(h2c.NewHandler(mux, &http2.Server{})))
}
