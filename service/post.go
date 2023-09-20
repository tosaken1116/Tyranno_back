package service

import (
	"context"
	"fmt"
	"log"
	"nnyd-back/config"
	"nnyd-back/controller"
	"nnyd-back/db"
	protosv1 "nnyd-back/pb/schemas/protos/v1"

	"connectrpc.com/connect"
	"google.golang.org/protobuf/types/known/emptypb"
)

var (
	pc  = &controller.PostController{}
	fac = &controller.FavoriteController{}
)

type PostServer struct{}

func (ps *PostServer) CreatePost(ctx context.Context, req *connect.Request[protosv1.CreatePostRequest]) (*connect.Response[protosv1.CreatePostResponse], error) {
	user_id := ctx.Value(config.USER_ID).(string)

	conn := db.GetDB()
	if _, err := uc.GetUserById(conn, user_id); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	resp, err := pc.CreatePost(conn, req.Msg, user_id)
	if err != nil {
		log.Println(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(resp), nil
}

func (ps *PostServer) GetPost(ctx context.Context, req *connect.Request[protosv1.GetPostRequest]) (*connect.Response[protosv1.GetPostResponse], error) {
	user_id := ctx.Value(config.USER_ID).(string)
	var p_user_id *string = nil

	conn := db.GetDB()
	if user_id != "" {
		if _, err := uc.GetUserById(conn, user_id); err != nil {
			return nil, connect.NewError(connect.CodeUnauthenticated, err)
		}
		p_user_id = &user_id
	}
	resp, err := pc.GetPostByID(conn, req.Msg.Id, p_user_id)
	if err != nil {
		log.Println(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(resp), nil
}

func (ps *PostServer) GetPosts(ctx context.Context, req *connect.Request[protosv1.GetPostsRequest]) (*connect.Response[protosv1.GetPostsResponse], error) {
	user_id := ctx.Value(config.USER_ID).(string)
	var p_user_id *string = nil

	conn := db.GetDB()
	if user_id != "" {
		if _, err := uc.GetUserById(conn, user_id); err != nil {
			return nil, connect.NewError(connect.CodeUnauthenticated, err)
		}
		p_user_id = &user_id
	}

	resp, err := pc.GetPosts(conn, p_user_id)
	if err != nil {
		log.Println(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(resp), nil
}

func (ps *PostServer) DeletePost(ctx context.Context, req *connect.Request[protosv1.DeletePostRequest]) (*connect.Response[protosv1.DeletePostResponse], error) {
	user_id := ctx.Value(config.USER_ID).(string)

	conn := db.GetDB()
	userResp, err := uc.GetUserById(conn, user_id)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	postResp, err := pc.GetPostByID(conn, req.Msg.Id, nil)
	if err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}
	if postResp.Post.User.DisplayId != userResp.User.DisplayId {
		return nil, connect.NewError(connect.CodePermissionDenied, fmt.Errorf("only delete your post"))
	}

	resp, err := pc.DeletePost(conn, req.Msg.Id)
	if err != nil {
		log.Println(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(resp), nil
}

func (ps *PostServer) GetReplies(ctx context.Context, req *connect.Request[protosv1.GetRepliesRequest]) (*connect.Response[protosv1.GetRepliesResponse], error) {
	user_id := ctx.Value(config.USER_ID).(string)
	var p_user_id *string = nil

	conn := db.GetDB()
	if user_id != "" {
		if _, err := uc.GetUserById(conn, user_id); err != nil {
			return nil, connect.NewError(connect.CodeUnauthenticated, err)
		}
		p_user_id = &user_id
	}
	if _, err := pc.GetPostByID(conn, req.Msg.ReplyAt, nil); err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	resp, err := pc.GetReplies(conn, req.Msg.ReplyAt, p_user_id)
	if err != nil {
		log.Println(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(resp), nil
}

func (ps *PostServer) CreateFavorite(ctx context.Context, req *connect.Request[protosv1.CreateFavoriteRequest]) (*connect.Response[protosv1.CreateFavoriteResponse], error) {
	user_id := ctx.Value(config.USER_ID).(string)

	conn := db.GetDB()
	if _, err := uc.GetUserById(conn, user_id); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}
	if _, err := pc.GetPostByID(conn, req.Msg.FavoriteAt, nil); err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	resp, err := fac.CreateFavorite(conn, user_id, req.Msg.FavoriteAt)
	if err != nil {
		log.Println(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(resp), nil
}

func (ps *PostServer) DeleteFavorite(ctx context.Context, req *connect.Request[protosv1.DeleteFavoriteRequest]) (*connect.Response[protosv1.DeleteFavoriteResponse], error) {
	user_id := ctx.Value(config.USER_ID).(string)

	conn := db.GetDB()
	if _, err := uc.GetUserById(conn, user_id); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}
	if _, err := pc.GetPostByID(conn, req.Msg.FavoriteAt, nil); err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	resp, err := fac.DeleteFavorite(conn, user_id, req.Msg.FavoriteAt)
	if err != nil {
		log.Println(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(resp), nil
}

func (ps *PostServer) GetMyFavoritePosts(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[protosv1.GetPostsResponse], error) {
	user_id := ctx.Value(config.USER_ID).(string)
	conn := db.GetDB()
	if _, err := uc.GetUserById(conn, user_id); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	resp, err := fac.GetFavoritePosts(conn, user_id, &user_id)
	if err != nil {
		log.Println(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(resp), nil
}

func (ps *PostServer) GetOthersFavoritePosts(ctx context.Context, req *connect.Request[protosv1.GetUserRequest]) (*connect.Response[protosv1.GetPostsResponse], error) {
	my_user_id := ctx.Value(config.USER_ID).(string)
	var p_my_user_id *string = nil

	conn := db.GetDB()
	if my_user_id != "" {
		if _, err := uc.GetUserById(conn, my_user_id); err != nil {
			return nil, connect.NewError(connect.CodeUnauthenticated, err)
		}
		p_my_user_id = &my_user_id
	}
	user_id, err := uc.GetIdFromDisplayId(conn, req.Msg.DisplayId)
	if err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}

	resp, err := fac.GetFavoritePosts(conn, user_id, p_my_user_id)
	if err != nil {
		log.Println(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(resp), nil
}

func (ps *PostServer) GetUsersFavoritedPost(ctx context.Context, req *connect.Request[protosv1.GetPostRequest]) (*connect.Response[protosv1.GetUsersResponse], error) {
	user_id := ctx.Value(config.USER_ID).(string)
	conn := db.GetDB()
	if _, err := uc.GetUserById(conn, user_id); err != nil {
		return nil, connect.NewError(connect.CodeUnauthenticated, err)
	}
	if _, err := pc.GetPostByID(conn, req.Msg.Id, nil); err != nil {
		return nil, connect.NewError(connect.CodeNotFound, err)
	}

	resp, err := fac.GetFavoritingUser(conn, req.Msg.Id)
	if err != nil {
		log.Println(err)
		return nil, connect.NewError(connect.CodeInternal, err)
	}
	return connect.NewResponse(resp), nil
}
