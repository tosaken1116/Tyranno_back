package service

import (
	"context"
	"fmt"
	"log"
	"nnyd-back/config"
	"nnyd-back/controller"
	"nnyd-back/db"
	protosv1 "nnyd-back/pb/schemas/protos/v1"
	"time"

	"connectrpc.com/connect"
	"github.com/dgrijalva/jwt-go"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	ac = &controller.AuthController{}
)

type AuthServer struct{}

func (as *AuthServer) SignIn(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[protosv1.SignInResponse], error) {
	firebase_id := ctx.Value(config.FIREBASE_ID)
	if firebase_id == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("verifying failed"))
	}
	firebase_id_str := firebase_id.(string)
	conn := db.GetDB()
	user_id, err := ac.CheckVerifyTotp(conn, firebase_id_str)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	claims := jwt.MapClaims{
		"user_id": user_id,
		"exp":     time.Now().Add(3000 * time.Minute).Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	accessToken, err := token.SignedString([]byte(config.JST_SECRET_KEY))
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	resp := &protosv1.SignInResponse{
		Token: accessToken,
	}
	return connect.NewResponse(resp), nil
}

func (as *AuthServer) SignOut(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[protosv1.SignOutResponse], error) {
	user_id := ctx.Value(config.USER_ID)
	if user_id == nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, fmt.Errorf("missing credential"))
	}
	user_id_str := user_id.(string)
	conn := db.GetDB()
	if _, err := uc.GetUserById(conn, user_id_str); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}
	resp, err := ac.SignOut(conn, user_id_str)
	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(resp), nil
}

func (as *AuthServer) GenerateTotpKey(ctx context.Context, req *connect.Request[protosv1.GenerateTotpKeyRequest]) (*connect.Response[protosv1.GenerateTotpKeyResponse], error) {
	conn := db.GetDB()
	resp, err := ac.GenerateTotpKeyController(conn, req.Msg)
	if err != nil {
		log.Println(err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed generate totp secret"))
	}
	return connect.NewResponse(resp), nil
}

func (as *AuthServer) VerifyTotp(ctx context.Context, req *connect.Request[protosv1.VerifyTotpRequest]) (*connect.Response[protosv1.VerifyTotpResponse], error) {
	conn := db.GetDB()
	resp, err := ac.VerifyTotpController(conn, req.Msg)

	if err != nil {
		log.Println(err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed token verify"))
	}

	return connect.NewResponse(resp), nil
}

func (as *AuthServer) ValidateTotp(ctx context.Context, req *connect.Request[protosv1.ValidateTotpRequest]) (*connect.Response[protosv1.ValidateTotpResponse], error) {
	conn := db.GetDB()
	log.Println(req.Msg)
	resp, err := ac.ValidateTotpController(conn, req.Msg)

	if err != nil {
		log.Println(err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed token verify"))
	}

	return connect.NewResponse(resp), nil
}
