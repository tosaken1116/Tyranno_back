package main

import (
	"context"
	"fmt"
	"net/http"
	"nnyd-back/config"
	"nnyd-back/db"
	"nnyd-back/pb/schemas/protos/v1/protosv1connect"
	"nnyd-back/service"
	"nnyd-back/utils"
	"strings"
	"time"

	"connectrpc.com/connect"
	"connectrpc.com/grpcreflect"
	"github.com/dgrijalva/jwt-go"
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
		"schemas.protos.v1.AuthService",
	)
	mux.Handle(grpcreflect.NewHandlerV1(reflector))
	mux.Handle(grpcreflect.NewHandlerV1Alpha(reflector))
	return mux
}

// インターセプタ設定
func newInterCeptors() connect.Option {
	interceptor := func(next connect.UnaryFunc) connect.UnaryFunc {
		return connect.UnaryFunc(func(ctx context.Context, req connect.AnyRequest) (connect.AnyResponse, error) {
			auth := req.Header().Get("Authorization")
			if auth == "" {
				return next(ctx, req)
			}
			authArray := strings.Split(auth, " ")
			if len(authArray) < 2 {
				return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("please set valid Authorization Header"))
			}
			if authArray[0] != "Bearer" {
				return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("please set valid Authorization Header"))
			}
			token := authArray[1]
			authProvider := req.Header().Get("AuthProvider")
			switch authProvider {
			case "firebase":
				client := utils.GetFirebaseAuthClient()
				userInfo, err := client.VerifyIDToken(ctx, token)
				if err != nil {
					return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("failed verifying firebase token"))
				}
				ctx = context.WithValue(ctx, config.FIREBASE_ID, userInfo.UID)
			case "origin":
				t, err := jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
					if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
						return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
					}

					return []byte(config.JST_SECRET_KEY), nil
				})

				claims, ok := t.Claims.(jwt.MapClaims)
				if !ok || !t.Valid {
					return nil, connect.NewError(connect.CodeUnauthenticated, err)
				}
				user_id := string(claims["user_id"].(string))
				exp := int64(claims["exp"].(float64))

				if time.Now().Unix() > exp {
					return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("token is expired"))
				}
				ctx = context.WithValue(ctx, config.USER_ID, user_id)
			default:
				return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("please set valid AuthProvider Header"))
			}
			return next(ctx, req)
		})
	}
	return connect.WithInterceptors(connect.UnaryInterceptorFunc(interceptor))
}

func main() {
	config.LoadConfig()
	utils.InitFirebaseAuthClient()

	db.Init()
	defer db.Close()
	db.AutoMigration()

	mux := newServeMuxWithReflection()
	interceptor := newInterCeptors()
	mux.Handle(protosv1connect.NewUserServiceHandler(&service.UserServer{}, interceptor))
	mux.Handle(protosv1connect.NewPostServiceHandler(&service.PostServer{}, interceptor))
	mux.Handle(protosv1connect.NewAuthServiceHandler(&service.AuthServer{}, interceptor))

	portStr := ":" + config.PORT
	debugCors := cors.New(cors.Options{
		AllowedOrigins: []string{"*"},
		AllowedMethods: []string{
			http.MethodOptions,
			http.MethodHead,
			http.MethodGet,
			http.MethodPost,
			http.MethodPut,
			http.MethodPatch,
			http.MethodDelete,
		},
		AllowedHeaders:     []string{"*"},
		AllowCredentials:   true,
		OptionsPassthrough: true,
		Debug:              true,
	})
	http.ListenAndServe(portStr, debugCors.Handler(h2c.NewHandler(mux, &http2.Server{})))
}
