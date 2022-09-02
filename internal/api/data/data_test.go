package data

import (
	"context"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	errorsPkg "gitlab.ozon.dev/iTukaev/homework/internal/customerrors"
	userMockPkg "gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/mock"
	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
	"gitlab.ozon.dev/iTukaev/homework/pkg/adaptor"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
	loggerPkg "gitlab.ozon.dev/iTukaev/homework/pkg/logger"
	apiMockPkg "gitlab.ozon.dev/iTukaev/homework/pkg/mock"
)

func TestDataApi_UserAllList(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	ctx := context.WithValue(context.Background(), "meta", "meta")

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
			userCtl := New(mockUser, loggerPkg.NewFatal())

			gomock.InOrder(
				mockStream.EXPECT().Context().Return(ctx).Times(2),
				mockUser.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(c.first, c.listErr).Times(1),
				mockStream.EXPECT().Send(c.toSend).
					Return(c.sendErr).MaxTimes(1),
				mockStream.EXPECT().Context().Return(ctx).MaxTimes(1),
				mockUser.EXPECT().List(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
					Return(c.second, c.listErr).MaxTimes(1),
			)
			err := userCtl.UserAllList(&pb.UserAllListRequest{}, mockStream)

			require.ErrorIs(t, err, c.expErr)
		})
	}
}
