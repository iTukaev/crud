package delete

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
	apiMockPkg "gitlab.ozon.dev/iTukaev/homework/pkg/mock"
)

var (
	ErrUnexpected = errors.New("unexpected error")

	userName = "Ivan"
)

func TestAddCommand_Process(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()
	ctx := context.Background()

	cases := []struct {
		name    string
		args    string
		expText string
		expErr  error
	}{
		{
			name:    "success",
			args:    userName,
			expText: fmt.Sprintf("user [%s] deleted", userName),
			expErr:  nil,
		},
		{
			name:    "failed, UserDelete returns error",
			args:    userName,
			expText: "internal error",
			expErr:  ErrUnexpected,
		},
		{
			name:    "failed, UserDelete returns specific error",
			args:    userName,
			expText: "error message",
			expErr:  status.Error(codes.Internal, "error message"),
		},
		{
			name:    "failed, invalid arguments",
			args:    "",
			expText: "invalid arguments",
			expErr:  nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockClient := apiMockPkg.NewMockUserClient(ctl)
			deleteCommand := New(mockClient)

			gomock.InOrder(
				mockClient.EXPECT().UserDelete(ctx, &pb.UserDeleteRequest{
					Name: c.args,
				}).Return(&pb.UserDeleteResponse{}, c.expErr).MaxTimes(1),
			)
			text := deleteCommand.Process(ctx, c.args)

			assert.Equal(t, c.expText, text)
		})
	}
}

func TestAddCommand_Name(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	t.Run("success", func(t *testing.T) {
		mockClient := apiMockPkg.NewMockUserClient(ctl)
		deleteCommand := New(mockClient)

		name := deleteCommand.Name()

		assert.Equal(t, deleteName, name)
	})
}

func TestAddCommand_Description(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	t.Run("success", func(t *testing.T) {
		mockClient := apiMockPkg.NewMockUserClient(ctl)
		deleteCommand := New(mockClient)

		description := deleteCommand.Description()

		assert.Equal(t, deleteDescription, description)
	})
}
