package service

import (
	"context"
	protosv1 "nnyd-back/pb/schemas/protos/v1"

	"connectrpc.com/connect"
)

type PostServer struct{}

func (ps *PostServer) CreatePost(context.Context, *connect.Request[protosv1.CreatePostRequest]) (*connect.Response[protosv1.CreatePostResponse], error) {
	// mock
	resp := &protosv1.CreatePostResponse{
		Post: nil,
	}
	return connect.NewResponse(resp), nil
}

func (ps *PostServer) GetPost(context.Context, *connect.Request[protosv1.GetPostRequest]) (*connect.Response[protosv1.GetPostResponse], error) {
	// mock
	resp := &protosv1.GetPostResponse{
		Post: nil,
	}
	return connect.NewResponse(resp), nil
}

func (ps *PostServer) GetPosts(context.Context, *connect.Request[protosv1.GetPostsRequest]) (*connect.Response[protosv1.GetPostsResponse], error) {
	// mock
	resp := &protosv1.GetPostsResponse{
		Posts: []*protosv1.Post{},
	}
	return connect.NewResponse(resp), nil
}

func (ps *PostServer) DeletePost(context.Context, *connect.Request[protosv1.DeletePostRequest]) (*connect.Response[protosv1.DeletePostResponse], error) {
	// mock
	resp := &protosv1.DeletePostResponse{
		Status: true,
	}
	return connect.NewResponse(resp), nil
}

func (ps *PostServer) GetReplies(context.Context, *connect.Request[protosv1.GetRepliesRequest]) (*connect.Response[protosv1.GetRepliesResponse], error) {
	// mock
	resp := &protosv1.GetRepliesResponse{
		Replies: []*protosv1.Post{},
	}
	return connect.NewResponse(resp), nil
}

func (ps *PostServer) CreateFavorite(context.Context, *connect.Request[protosv1.CreateFavoriteRequest]) (*connect.Response[protosv1.CreateFavoriteResponse], error) {
	// mock
	resp := &protosv1.CreateFavoriteResponse{
		Status: true,
	}
	return connect.NewResponse(resp), nil
}

func (ps *PostServer) DeleteFavorite(context.Context, *connect.Request[protosv1.DeleteFavoriteRequest]) (*connect.Response[protosv1.DeleteFavoriteResponse], error) {
	// mock
	resp := &protosv1.DeleteFavoriteResponse{
		Status: true,
	}
	return connect.NewResponse(resp), nil
}
