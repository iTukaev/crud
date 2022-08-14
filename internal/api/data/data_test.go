package data

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	userMockPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/mock"
	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
	errorsPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo/customerrors"
	"gitlab.ozon.dev/iTukaev/homework/pkg/adaptor"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
	apiMockPkg "gitlab.ozon.dev/iTukaev/homework/pkg/mock"
)

func TestDataApi_UserCreate(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	ctx := context.Background()

	cases := []struct {
		name   string
		err    error
		expErr error
		expRes *pb.UserCreateResponse
	}{
		{
			name:   "success",
			err:    nil,
			expErr: nil,
			expRes: &pb.UserCreateResponse{},
		},
		{
			name:   "failed, user already exists",
			err:    errorsPkg.ErrUserAlreadyExists,
			expErr: status.Error(codes.AlreadyExists, errorsPkg.ErrUserAlreadyExists.Error()),
			expRes: nil,
		},
		{
			name:   "failed, deadline exceeded",
			err:    errorsPkg.ErrTimeout,
			expErr: status.Error(codes.DeadlineExceeded, errorsPkg.ErrTimeout.Error()),
			expRes: nil,
		},
		{
			name:   "failed, unexpected error",
			err:    errorsPkg.ErrUnexpected,
			expErr: status.Error(codes.Internal, errorsPkg.ErrUnexpected.Error()),
			expRes: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockUser := userMockPkg.NewMockInterface(ctl)
			userCtl := New(mockUser)

			gomock.InOrder(
				mockUser.EXPECT().Create(gomock.Any(), models.User{}).
					Return(c.err).MaxTimes(1),
			)
			res, err := userCtl.UserCreate(ctx, &pb.UserCreateRequest{})

			require.ErrorIs(t, err, c.expErr)
			assert.Equal(t, c.expRes, res)
		})
	}
}

func TestDataApi_UserUpdate(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	ctx := context.Background()

	cases := []struct {
		name   string
		err    error
		expErr error
		expRes *pb.UserUpdateResponse
	}{
		{
			name:   "success",
			err:    nil,
			expErr: nil,
			expRes: &pb.UserUpdateResponse{},
		},
		{
			name:   "failed, user not found",
			err:    errorsPkg.ErrUserNotFound,
			expErr: status.Error(codes.InvalidArgument, errorsPkg.ErrUserNotFound.Error()),
			expRes: nil,
		},
		{
			name:   "failed, deadline exceeded",
			err:    errorsPkg.ErrTimeout,
			expErr: status.Error(codes.DeadlineExceeded, errorsPkg.ErrTimeout.Error()),
			expRes: nil,
		},
		{
			name:   "failed, unexpected error",
			err:    errorsPkg.ErrUnexpected,
			expErr: status.Error(codes.Internal, errorsPkg.ErrUnexpected.Error()),
			expRes: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockUser := userMockPkg.NewMockInterface(ctl)
			userCtl := New(mockUser)

			gomock.InOrder(
				mockUser.EXPECT().Update(gomock.Any(), models.User{}).
					Return(c.err).MaxTimes(1),
			)
			res, err := userCtl.UserUpdate(ctx, &pb.UserUpdateRequest{})

			require.ErrorIs(t, err, c.expErr)
			assert.Equal(t, c.expRes, res)
		})
	}
}

func TestDataApi_UserDelete(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	ctx := context.Background()

	cases := []struct {
		name   string
		err    error
		expErr error
		expRes *pb.UserDeleteResponse
	}{
		{
			name:   "success",
			err:    nil,
			expErr: nil,
			expRes: &pb.UserDeleteResponse{},
		},
		{
			name:   "failed, user not found",
			err:    errorsPkg.ErrUserNotFound,
			expErr: status.Error(codes.InvalidArgument, errorsPkg.ErrUserNotFound.Error()),
			expRes: nil,
		},
		{
			name:   "failed, deadline exceeded",
			err:    errorsPkg.ErrTimeout,
			expErr: status.Error(codes.DeadlineExceeded, errorsPkg.ErrTimeout.Error()),
			expRes: nil,
		},
		{
			name:   "failed, unexpected error",
			err:    errorsPkg.ErrUnexpected,
			expErr: status.Error(codes.Internal, errorsPkg.ErrUnexpected.Error()),
			expRes: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockUser := userMockPkg.NewMockInterface(ctl)
			userCtl := New(mockUser)

			gomock.InOrder(
				mockUser.EXPECT().Delete(gomock.Any(), models.User{}.Name).
					Return(c.err).MaxTimes(1),
			)
			res, err := userCtl.UserDelete(ctx, &pb.UserDeleteRequest{})

			require.ErrorIs(t, err, c.expErr)
			assert.Equal(t, c.expRes, res)
		})
	}
}

func TestDataApi_UserGet(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	ctx := context.Background()

	cases := []struct {
		name   string
		err    error
		expErr error
		expRes *pb.UserGetResponse
	}{
		{
			name:   "success",
			err:    nil,
			expErr: nil,
			expRes: &pb.UserGetResponse{User: adaptor.ToUserPbModel(models.User{})},
		},
		{
			name:   "failed, user not found",
			err:    errorsPkg.ErrUserNotFound,
			expErr: status.Error(codes.InvalidArgument, errorsPkg.ErrUserNotFound.Error()),
			expRes: nil,
		},
		{
			name:   "failed, deadline exceeded",
			err:    errorsPkg.ErrTimeout,
			expErr: status.Error(codes.DeadlineExceeded, errorsPkg.ErrTimeout.Error()),
			expRes: nil,
		},
		{
			name:   "failed, unexpected error",
			err:    errorsPkg.ErrUnexpected,
			expErr: status.Error(codes.Internal, errorsPkg.ErrUnexpected.Error()),
			expRes: nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockUser := userMockPkg.NewMockInterface(ctl)
			userCtl := New(mockUser)

			gomock.InOrder(
				mockUser.EXPECT().Get(gomock.Any(), models.User{}.Name).
					Return(models.User{}, c.err).MaxTimes(1),
			)
			res, err := userCtl.UserGet(ctx, &pb.UserGetRequest{})

			require.ErrorIs(t, err, c.expErr)
			assert.Equal(t, c.expRes, res)
		})
	}
}

func TestDataApi_UserList(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	ctx := context.Background()

	cases := []struct {
		name   string
		err    error
		expErr error
		expRes *pb.UserListResponse
	}{
		{
			name:   "success",
			err:    nil,
			expErr: nil,
			expRes: &pb.UserListResponse{Users: adaptor.ToUserListPbModel([]models.User{})},
		},
		{
			name:   "failed, deadline exceeded",
			err:    errorsPkg.ErrTimeout,
			expErr: status.Error(codes.DeadlineExceeded, errorsPkg.ErrTimeout.Error()),
			expRes: &pb.UserListResponse{},
		},
		{
			name:   "failed, unexpected error",
			err:    errorsPkg.ErrUnexpected,
			expErr: status.Error(codes.Internal, errorsPkg.ErrUnexpected.Error()),
			expRes: &pb.UserListResponse{},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockUser := userMockPkg.NewMockInterface(ctl)
			userCtl := New(mockUser)

			gomock.InOrder(
				mockUser.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return([]models.User{}, c.err).MaxTimes(1),
			)
			res, err := userCtl.UserList(ctx, &pb.UserListRequest{})

			require.ErrorIs(t, err, c.expErr)
			assert.Equal(t, c.expRes, res)
		})
	}
}

func TestDataApi_UserAllList(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	cases := []struct {
		name    string
		listErr error
		sendErr error
		expErr  error
		first   []models.User
		second  []models.User
		toSend  *pb.UserAllListResponse
	}{
		{
			name:    "success",
			listErr: nil,
			sendErr: nil,
			expErr:  nil,
			first:   []models.User{{}, {}},
			second:  []models.User{},
			toSend:  &pb.UserAllListResponse{Users: adaptor.ToUserListPbModel([]models.User{{}, {}})},
		},
		{
			name:    "failed, List unexpected error",
			listErr: errorsPkg.ErrUnexpected,
			sendErr: nil,
			expErr:  status.Error(codes.Internal, errorsPkg.ErrUnexpected.Error()),
			first:   []models.User{{}, {}},
			second:  []models.User{},
			toSend:  &pb.UserAllListResponse{Users: adaptor.ToUserListPbModel([]models.User{{}, {}})},
		},
		{
			name:    "failed, Send unexpected error",
			listErr: nil,
			sendErr: errorsPkg.ErrUnexpected,
			expErr:  status.Error(codes.Internal, errorsPkg.ErrUnexpected.Error()),
			first:   []models.User{{}, {}},
			second:  []models.User{},
			toSend:  &pb.UserAllListResponse{Users: adaptor.ToUserListPbModel([]models.User{{}, {}})},
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockUser := userMockPkg.NewMockInterface(ctl)
			mockStream := apiMockPkg.NewMockUser_UserAllListServer(ctl)
			userCtl := New(mockUser)

			gomock.InOrder(
				mockStream.EXPECT().Context().Times(1),
				mockUser.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(c.first, c.listErr).Times(1),
				mockStream.EXPECT().Send(c.toSend).
					Return(c.sendErr).MaxTimes(1),
				mockStream.EXPECT().Context().MaxTimes(1),
				mockUser.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(c.second, c.listErr).MaxTimes(1),
			)
			err := userCtl.UserAllList(&pb.UserAllListRequest{}, mockStream)

			require.ErrorIs(t, err, c.expErr)
		})
	}
}
