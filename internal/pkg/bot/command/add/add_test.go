package add

import (
	"context"
	"fmt"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
	errorsPkg "gitlab.ozon.dev/iTukaev/homework/internal/repo/customerrors"
	"gitlab.ozon.dev/iTukaev/homework/pkg/adaptor"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
	apiMockPkg "gitlab.ozon.dev/iTukaev/homework/pkg/mock"
)

var (
	user = models.User{
		Name:     "Ivan",
		Password: "123",
		Email:    "ivan@email.com",
		FullName: "IvanDummy",
	}
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
			args:    fmt.Sprintf("%s %s %s %s", user.Name, user.Password, user.Email, user.FullName),
			expText: fmt.Sprintf("user [%s] added", user.Name),
			expErr:  nil,
		},
		{
			name:    "failed, UserCreate returns error",
			args:    fmt.Sprintf("%s %s %s %s", user.Name, user.Password, user.Email, user.FullName),
			expText: "internal error",
			expErr:  errorsPkg.ErrUnexpected,
		},
		{
			name:    "failed, UserCreate returns specific error",
			args:    fmt.Sprintf("%s %s %s %s", user.Name, user.Password, user.Email, user.FullName),
			expText: "error message",
			expErr:  status.Error(codes.Internal, "error message"),
		},
		{
			name:    "failed, invalid arguments",
			args:    fmt.Sprintf("%s %s %s", user.Name, user.Password, user.Email),
			expText: "invalid arguments",
			expErr:  nil,
		},
	}

	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			mockClient := apiMockPkg.NewMockUserClient(ctl)
			addCommand := New(mockClient)

			gomock.InOrder(
				mockClient.EXPECT().UserCreate(ctx, &pb.UserCreateRequest{
					User: adaptor.ToUserPbModel(user),
				}).Return(&pb.UserCreateResponse{}, c.expErr).MaxTimes(1),
			)
			text := addCommand.Process(ctx, c.args)

			assert.Equal(t, c.expText, text)
		})
	}
}

func TestAddCommand_Name(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	t.Run("success", func(t *testing.T) {
		mockClient := apiMockPkg.NewMockUserClient(ctl)
		addCommand := New(mockClient)

		name := addCommand.Name()

		assert.Equal(t, addName, name)
	})
}

func TestAddCommand_Description(t *testing.T) {
	ctl := gomock.NewController(t)
	defer ctl.Finish()

	t.Run("success", func(t *testing.T) {
		mockClient := apiMockPkg.NewMockUserClient(ctl)
		addCommand := New(mockClient)

		description := addCommand.Description()

		assert.Equal(t, addDescription, description)
	})
}
