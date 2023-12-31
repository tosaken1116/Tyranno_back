// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: schemas/protos/v1/post.proto

package protosv1connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
	http "net/http"
	v1 "nnyd-back/pb/schemas/protos/v1"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion0_1_0

const (
	// PostServiceName is the fully-qualified name of the PostService service.
	PostServiceName = "schemas.protos.v1.PostService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// PostServiceCreatePostProcedure is the fully-qualified name of the PostService's CreatePost RPC.
	PostServiceCreatePostProcedure = "/schemas.protos.v1.PostService/CreatePost"
	// PostServiceGetPostProcedure is the fully-qualified name of the PostService's GetPost RPC.
	PostServiceGetPostProcedure = "/schemas.protos.v1.PostService/GetPost"
	// PostServiceGetPostsProcedure is the fully-qualified name of the PostService's GetPosts RPC.
	PostServiceGetPostsProcedure = "/schemas.protos.v1.PostService/GetPosts"
	// PostServiceDeletePostProcedure is the fully-qualified name of the PostService's DeletePost RPC.
	PostServiceDeletePostProcedure = "/schemas.protos.v1.PostService/DeletePost"
	// PostServiceGetRepliesProcedure is the fully-qualified name of the PostService's GetReplies RPC.
	PostServiceGetRepliesProcedure = "/schemas.protos.v1.PostService/GetReplies"
	// PostServiceCreateFavoriteProcedure is the fully-qualified name of the PostService's
	// CreateFavorite RPC.
	PostServiceCreateFavoriteProcedure = "/schemas.protos.v1.PostService/CreateFavorite"
	// PostServiceDeleteFavoriteProcedure is the fully-qualified name of the PostService's
	// DeleteFavorite RPC.
	PostServiceDeleteFavoriteProcedure = "/schemas.protos.v1.PostService/DeleteFavorite"
	// PostServiceGetMyFavoritePostsProcedure is the fully-qualified name of the PostService's
	// GetMyFavoritePosts RPC.
	PostServiceGetMyFavoritePostsProcedure = "/schemas.protos.v1.PostService/GetMyFavoritePosts"
	// PostServiceGetOthersFavoritePostsProcedure is the fully-qualified name of the PostService's
	// GetOthersFavoritePosts RPC.
	PostServiceGetOthersFavoritePostsProcedure = "/schemas.protos.v1.PostService/GetOthersFavoritePosts"
	// PostServiceGetUsersFavoritedPostProcedure is the fully-qualified name of the PostService's
	// GetUsersFavoritedPost RPC.
	PostServiceGetUsersFavoritedPostProcedure = "/schemas.protos.v1.PostService/GetUsersFavoritedPost"
)

// PostServiceClient is a client for the schemas.protos.v1.PostService service.
type PostServiceClient interface {
	CreatePost(context.Context, *connect.Request[v1.CreatePostRequest]) (*connect.Response[v1.CreatePostResponse], error)
	GetPost(context.Context, *connect.Request[v1.GetPostRequest]) (*connect.Response[v1.GetPostResponse], error)
	GetPosts(context.Context, *connect.Request[v1.GetPostsRequest]) (*connect.Response[v1.GetPostsResponse], error)
	DeletePost(context.Context, *connect.Request[v1.DeletePostRequest]) (*connect.Response[v1.DeletePostResponse], error)
	GetReplies(context.Context, *connect.Request[v1.GetRepliesRequest]) (*connect.Response[v1.GetRepliesResponse], error)
	CreateFavorite(context.Context, *connect.Request[v1.CreateFavoriteRequest]) (*connect.Response[v1.CreateFavoriteResponse], error)
	DeleteFavorite(context.Context, *connect.Request[v1.DeleteFavoriteRequest]) (*connect.Response[v1.DeleteFavoriteResponse], error)
	GetMyFavoritePosts(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.GetPostsResponse], error)
	GetOthersFavoritePosts(context.Context, *connect.Request[v1.GetUserRequest]) (*connect.Response[v1.GetPostsResponse], error)
	GetUsersFavoritedPost(context.Context, *connect.Request[v1.GetPostRequest]) (*connect.Response[v1.GetUsersResponse], error)
}

// NewPostServiceClient constructs a client for the schemas.protos.v1.PostService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewPostServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) PostServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &postServiceClient{
		createPost: connect.NewClient[v1.CreatePostRequest, v1.CreatePostResponse](
			httpClient,
			baseURL+PostServiceCreatePostProcedure,
			opts...,
		),
		getPost: connect.NewClient[v1.GetPostRequest, v1.GetPostResponse](
			httpClient,
			baseURL+PostServiceGetPostProcedure,
			opts...,
		),
		getPosts: connect.NewClient[v1.GetPostsRequest, v1.GetPostsResponse](
			httpClient,
			baseURL+PostServiceGetPostsProcedure,
			opts...,
		),
		deletePost: connect.NewClient[v1.DeletePostRequest, v1.DeletePostResponse](
			httpClient,
			baseURL+PostServiceDeletePostProcedure,
			opts...,
		),
		getReplies: connect.NewClient[v1.GetRepliesRequest, v1.GetRepliesResponse](
			httpClient,
			baseURL+PostServiceGetRepliesProcedure,
			opts...,
		),
		createFavorite: connect.NewClient[v1.CreateFavoriteRequest, v1.CreateFavoriteResponse](
			httpClient,
			baseURL+PostServiceCreateFavoriteProcedure,
			opts...,
		),
		deleteFavorite: connect.NewClient[v1.DeleteFavoriteRequest, v1.DeleteFavoriteResponse](
			httpClient,
			baseURL+PostServiceDeleteFavoriteProcedure,
			opts...,
		),
		getMyFavoritePosts: connect.NewClient[emptypb.Empty, v1.GetPostsResponse](
			httpClient,
			baseURL+PostServiceGetMyFavoritePostsProcedure,
			opts...,
		),
		getOthersFavoritePosts: connect.NewClient[v1.GetUserRequest, v1.GetPostsResponse](
			httpClient,
			baseURL+PostServiceGetOthersFavoritePostsProcedure,
			opts...,
		),
		getUsersFavoritedPost: connect.NewClient[v1.GetPostRequest, v1.GetUsersResponse](
			httpClient,
			baseURL+PostServiceGetUsersFavoritedPostProcedure,
			opts...,
		),
	}
}

// postServiceClient implements PostServiceClient.
type postServiceClient struct {
	createPost             *connect.Client[v1.CreatePostRequest, v1.CreatePostResponse]
	getPost                *connect.Client[v1.GetPostRequest, v1.GetPostResponse]
	getPosts               *connect.Client[v1.GetPostsRequest, v1.GetPostsResponse]
	deletePost             *connect.Client[v1.DeletePostRequest, v1.DeletePostResponse]
	getReplies             *connect.Client[v1.GetRepliesRequest, v1.GetRepliesResponse]
	createFavorite         *connect.Client[v1.CreateFavoriteRequest, v1.CreateFavoriteResponse]
	deleteFavorite         *connect.Client[v1.DeleteFavoriteRequest, v1.DeleteFavoriteResponse]
	getMyFavoritePosts     *connect.Client[emptypb.Empty, v1.GetPostsResponse]
	getOthersFavoritePosts *connect.Client[v1.GetUserRequest, v1.GetPostsResponse]
	getUsersFavoritedPost  *connect.Client[v1.GetPostRequest, v1.GetUsersResponse]
}

// CreatePost calls schemas.protos.v1.PostService.CreatePost.
func (c *postServiceClient) CreatePost(ctx context.Context, req *connect.Request[v1.CreatePostRequest]) (*connect.Response[v1.CreatePostResponse], error) {
	return c.createPost.CallUnary(ctx, req)
}

// GetPost calls schemas.protos.v1.PostService.GetPost.
func (c *postServiceClient) GetPost(ctx context.Context, req *connect.Request[v1.GetPostRequest]) (*connect.Response[v1.GetPostResponse], error) {
	return c.getPost.CallUnary(ctx, req)
}

// GetPosts calls schemas.protos.v1.PostService.GetPosts.
func (c *postServiceClient) GetPosts(ctx context.Context, req *connect.Request[v1.GetPostsRequest]) (*connect.Response[v1.GetPostsResponse], error) {
	return c.getPosts.CallUnary(ctx, req)
}

// DeletePost calls schemas.protos.v1.PostService.DeletePost.
func (c *postServiceClient) DeletePost(ctx context.Context, req *connect.Request[v1.DeletePostRequest]) (*connect.Response[v1.DeletePostResponse], error) {
	return c.deletePost.CallUnary(ctx, req)
}

// GetReplies calls schemas.protos.v1.PostService.GetReplies.
func (c *postServiceClient) GetReplies(ctx context.Context, req *connect.Request[v1.GetRepliesRequest]) (*connect.Response[v1.GetRepliesResponse], error) {
	return c.getReplies.CallUnary(ctx, req)
}

// CreateFavorite calls schemas.protos.v1.PostService.CreateFavorite.
func (c *postServiceClient) CreateFavorite(ctx context.Context, req *connect.Request[v1.CreateFavoriteRequest]) (*connect.Response[v1.CreateFavoriteResponse], error) {
	return c.createFavorite.CallUnary(ctx, req)
}

// DeleteFavorite calls schemas.protos.v1.PostService.DeleteFavorite.
func (c *postServiceClient) DeleteFavorite(ctx context.Context, req *connect.Request[v1.DeleteFavoriteRequest]) (*connect.Response[v1.DeleteFavoriteResponse], error) {
	return c.deleteFavorite.CallUnary(ctx, req)
}

// GetMyFavoritePosts calls schemas.protos.v1.PostService.GetMyFavoritePosts.
func (c *postServiceClient) GetMyFavoritePosts(ctx context.Context, req *connect.Request[emptypb.Empty]) (*connect.Response[v1.GetPostsResponse], error) {
	return c.getMyFavoritePosts.CallUnary(ctx, req)
}

// GetOthersFavoritePosts calls schemas.protos.v1.PostService.GetOthersFavoritePosts.
func (c *postServiceClient) GetOthersFavoritePosts(ctx context.Context, req *connect.Request[v1.GetUserRequest]) (*connect.Response[v1.GetPostsResponse], error) {
	return c.getOthersFavoritePosts.CallUnary(ctx, req)
}

// GetUsersFavoritedPost calls schemas.protos.v1.PostService.GetUsersFavoritedPost.
func (c *postServiceClient) GetUsersFavoritedPost(ctx context.Context, req *connect.Request[v1.GetPostRequest]) (*connect.Response[v1.GetUsersResponse], error) {
	return c.getUsersFavoritedPost.CallUnary(ctx, req)
}

// PostServiceHandler is an implementation of the schemas.protos.v1.PostService service.
type PostServiceHandler interface {
	CreatePost(context.Context, *connect.Request[v1.CreatePostRequest]) (*connect.Response[v1.CreatePostResponse], error)
	GetPost(context.Context, *connect.Request[v1.GetPostRequest]) (*connect.Response[v1.GetPostResponse], error)
	GetPosts(context.Context, *connect.Request[v1.GetPostsRequest]) (*connect.Response[v1.GetPostsResponse], error)
	DeletePost(context.Context, *connect.Request[v1.DeletePostRequest]) (*connect.Response[v1.DeletePostResponse], error)
	GetReplies(context.Context, *connect.Request[v1.GetRepliesRequest]) (*connect.Response[v1.GetRepliesResponse], error)
	CreateFavorite(context.Context, *connect.Request[v1.CreateFavoriteRequest]) (*connect.Response[v1.CreateFavoriteResponse], error)
	DeleteFavorite(context.Context, *connect.Request[v1.DeleteFavoriteRequest]) (*connect.Response[v1.DeleteFavoriteResponse], error)
	GetMyFavoritePosts(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.GetPostsResponse], error)
	GetOthersFavoritePosts(context.Context, *connect.Request[v1.GetUserRequest]) (*connect.Response[v1.GetPostsResponse], error)
	GetUsersFavoritedPost(context.Context, *connect.Request[v1.GetPostRequest]) (*connect.Response[v1.GetUsersResponse], error)
}

// NewPostServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewPostServiceHandler(svc PostServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	postServiceCreatePostHandler := connect.NewUnaryHandler(
		PostServiceCreatePostProcedure,
		svc.CreatePost,
		opts...,
	)
	postServiceGetPostHandler := connect.NewUnaryHandler(
		PostServiceGetPostProcedure,
		svc.GetPost,
		opts...,
	)
	postServiceGetPostsHandler := connect.NewUnaryHandler(
		PostServiceGetPostsProcedure,
		svc.GetPosts,
		opts...,
	)
	postServiceDeletePostHandler := connect.NewUnaryHandler(
		PostServiceDeletePostProcedure,
		svc.DeletePost,
		opts...,
	)
	postServiceGetRepliesHandler := connect.NewUnaryHandler(
		PostServiceGetRepliesProcedure,
		svc.GetReplies,
		opts...,
	)
	postServiceCreateFavoriteHandler := connect.NewUnaryHandler(
		PostServiceCreateFavoriteProcedure,
		svc.CreateFavorite,
		opts...,
	)
	postServiceDeleteFavoriteHandler := connect.NewUnaryHandler(
		PostServiceDeleteFavoriteProcedure,
		svc.DeleteFavorite,
		opts...,
	)
	postServiceGetMyFavoritePostsHandler := connect.NewUnaryHandler(
		PostServiceGetMyFavoritePostsProcedure,
		svc.GetMyFavoritePosts,
		opts...,
	)
	postServiceGetOthersFavoritePostsHandler := connect.NewUnaryHandler(
		PostServiceGetOthersFavoritePostsProcedure,
		svc.GetOthersFavoritePosts,
		opts...,
	)
	postServiceGetUsersFavoritedPostHandler := connect.NewUnaryHandler(
		PostServiceGetUsersFavoritedPostProcedure,
		svc.GetUsersFavoritedPost,
		opts...,
	)
	return "/schemas.protos.v1.PostService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case PostServiceCreatePostProcedure:
			postServiceCreatePostHandler.ServeHTTP(w, r)
		case PostServiceGetPostProcedure:
			postServiceGetPostHandler.ServeHTTP(w, r)
		case PostServiceGetPostsProcedure:
			postServiceGetPostsHandler.ServeHTTP(w, r)
		case PostServiceDeletePostProcedure:
			postServiceDeletePostHandler.ServeHTTP(w, r)
		case PostServiceGetRepliesProcedure:
			postServiceGetRepliesHandler.ServeHTTP(w, r)
		case PostServiceCreateFavoriteProcedure:
			postServiceCreateFavoriteHandler.ServeHTTP(w, r)
		case PostServiceDeleteFavoriteProcedure:
			postServiceDeleteFavoriteHandler.ServeHTTP(w, r)
		case PostServiceGetMyFavoritePostsProcedure:
			postServiceGetMyFavoritePostsHandler.ServeHTTP(w, r)
		case PostServiceGetOthersFavoritePostsProcedure:
			postServiceGetOthersFavoritePostsHandler.ServeHTTP(w, r)
		case PostServiceGetUsersFavoritedPostProcedure:
			postServiceGetUsersFavoritedPostHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedPostServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedPostServiceHandler struct{}

func (UnimplementedPostServiceHandler) CreatePost(context.Context, *connect.Request[v1.CreatePostRequest]) (*connect.Response[v1.CreatePostResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("schemas.protos.v1.PostService.CreatePost is not implemented"))
}

func (UnimplementedPostServiceHandler) GetPost(context.Context, *connect.Request[v1.GetPostRequest]) (*connect.Response[v1.GetPostResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("schemas.protos.v1.PostService.GetPost is not implemented"))
}

func (UnimplementedPostServiceHandler) GetPosts(context.Context, *connect.Request[v1.GetPostsRequest]) (*connect.Response[v1.GetPostsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("schemas.protos.v1.PostService.GetPosts is not implemented"))
}

func (UnimplementedPostServiceHandler) DeletePost(context.Context, *connect.Request[v1.DeletePostRequest]) (*connect.Response[v1.DeletePostResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("schemas.protos.v1.PostService.DeletePost is not implemented"))
}

func (UnimplementedPostServiceHandler) GetReplies(context.Context, *connect.Request[v1.GetRepliesRequest]) (*connect.Response[v1.GetRepliesResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("schemas.protos.v1.PostService.GetReplies is not implemented"))
}

func (UnimplementedPostServiceHandler) CreateFavorite(context.Context, *connect.Request[v1.CreateFavoriteRequest]) (*connect.Response[v1.CreateFavoriteResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("schemas.protos.v1.PostService.CreateFavorite is not implemented"))
}

func (UnimplementedPostServiceHandler) DeleteFavorite(context.Context, *connect.Request[v1.DeleteFavoriteRequest]) (*connect.Response[v1.DeleteFavoriteResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("schemas.protos.v1.PostService.DeleteFavorite is not implemented"))
}

func (UnimplementedPostServiceHandler) GetMyFavoritePosts(context.Context, *connect.Request[emptypb.Empty]) (*connect.Response[v1.GetPostsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("schemas.protos.v1.PostService.GetMyFavoritePosts is not implemented"))
}

func (UnimplementedPostServiceHandler) GetOthersFavoritePosts(context.Context, *connect.Request[v1.GetUserRequest]) (*connect.Response[v1.GetPostsResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("schemas.protos.v1.PostService.GetOthersFavoritePosts is not implemented"))
}

func (UnimplementedPostServiceHandler) GetUsersFavoritedPost(context.Context, *connect.Request[v1.GetPostRequest]) (*connect.Response[v1.GetUsersResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("schemas.protos.v1.PostService.GetUsersFavoritedPost is not implemented"))
}
