package service

import (
	"context"
	"fmt"
	"log"
	"nnyd-back/controller"
	"nnyd-back/db"
	protosv1 "nnyd-back/pb/schemas/protos/v1"

	"connectrpc.com/connect"
)

type TotpServer struct{}

func (ps *TotpServer) GenerateTotpKey(ctx context.Context, req *connect.Request[protosv1.GenerateTotpKeyRequest]) (*connect.Response[protosv1.GenerateTotpKeyResponse], error) {
	conn := db.GetDB()
	uc := &controller.TotpController{}
	resp, err := uc.GenerateTotpKeyController(conn, req.Msg)
	if err != nil {
		log.Println(err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed generate totp secret"))
	}
	return connect.NewResponse(resp), nil
}

func (ps *TotpServer) VerifyTotp(ctx context.Context, req *connect.Request[protosv1.VerifyTotpRequest]) (*connect.Response[protosv1.VerifyTotpResponse], error) {
	conn := db.GetDB()
	uc := &controller.TotpController{}
	resp, err := uc.VerifyTotpController(conn, req.Msg)

	if err != nil {
		log.Println(err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed token verify"))
	}

	return connect.NewResponse(resp), nil
}

func (ps *TotpServer) ValidateTotp(ctx context.Context, req *connect.Request[protosv1.ValidateTotpRequest]) (*connect.Response[protosv1.ValidateTotpResponse], error) {
	conn := db.GetDB()
	uc := &controller.TotpController{}
	log.Println(req.Msg)
	resp, err := uc.ValidateTotpController(conn, req.Msg)

	if err != nil {
		log.Println(err)
		return nil, connect.NewError(connect.CodeInternal, fmt.Errorf("failed token verify"))
	}

	return connect.NewResponse(resp), nil
}
