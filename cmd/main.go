package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"nnyd-back/pb/schemas/protos"
	"nnyd-back/pb/schemas/protos/protosconnect"

	"github.com/bufbuild/connect-go"
	grpcreflect "github.com/bufbuild/connect-grpcreflect-go"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

type UserServer struct{}

func (s *UserServer) User(ctx context.Context, req *connect.Request[protos.UserRequest]) (*connect.Response[protos.UserResponse], error) {
	log.Println("Request headers: ", req.Header())

	if req.Msg.Name == "" || req.Msg.Icon == "" {
		// エラーにステータスコードを追加
		return nil, connect.NewError(connect.CodeInvalidArgument, fmt.Errorf("name and icon is required."))
	}

	responseUser := &protos.User{
		Id:   "dummy",
		Name: req.Msg.GetName(),
		Icon: req.Msg.GetIcon(),
	}
	userResp := &protos.UserResponse{
		User: responseUser,
	}
	resp := connect.NewResponse(userResp)
	return resp, nil
}

// リフレクション設定
func newServeMuxWithReflection() *http.ServeMux {
	mux := http.NewServeMux()
	reflector := grpcreflect.NewStaticReflector(
		"nnyd.user.UserService",
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
	userServer := &UserServer{}

	mux := newServeMuxWithReflection()
	interceptor := newInterCeptors()
	path, handler := protosconnect.NewUserServiceHandler(userServer, interceptor)
	mux.Handle(path, handler)
	http.ListenAndServe(":8080", h2c.NewHandler(mux, &http2.Server{}))
}
