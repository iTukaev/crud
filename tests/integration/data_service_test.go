//go:build integration
// +build integration

package integration

import (
	"github.com/jackc/pgx/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"gitlab.ozon.dev/iTukaev/homework/internal/pkg/core/user/models"
	"gitlab.ozon.dev/iTukaev/homework/pkg/adaptor"
	pb "gitlab.ozon.dev/iTukaev/homework/pkg/api"
	pbModels "gitlab.ozon.dev/iTukaev/homework/pkg/api/models"
	"gitlab.ozon.dev/iTukaev/homework/tests/integration/fixtures"
)

func (s *repositorySuite) TestData_UserCreate() {
	cases := []struct {
		name   string
		user   *pbModels.User
		expRes *pb.UserCreateResponse
		expErr codes.Code
	}{
		{
			name:   "success with User1",
			user:   adaptor.ToUserPbModel(*fixtures.User1),
			expRes: &pb.UserCreateResponse{User: nil},
			expErr: codes.OK,
		},
		{
			name:   "failed, user already exists",
			user:   adaptor.ToUserPbModel(*fixtures.ExistedUser1),
			expRes: nil,
			expErr: codes.AlreadyExists,
		},
		{
			name:   "failed, internal error",
			user:   nil,
			expRes: nil,
			expErr: codes.Internal,
		},
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			resp, errExpr := s.user.UserCreate(s.ctx, &pb.UserCreateRequest{User: c.user})
			st, _ := status.FromError(errExpr)

			s.Require().Equal(c.expErr, st.Code())
			s.Require().Equal(c.expRes, resp)

			if errExpr == nil {
				user, err := s.getUser(c.user.GetName())
				s.Require().NoError(err)
				s.Assert().Equal(c.user.GetName(), user.Name)
			}
		})
	}
}

func (s *repositorySuite) TestData_UserUpdate() {
	cases := []struct {
		name   string
		user   *pbModels.User
		expRes *pb.UserUpdateResponse
		expErr codes.Code
	}{
		{
			name:   "success with ExistedUser",
			user:   adaptor.ToUserPbModel(*fixtures.ExistedUser1.PasswordSet("987654321")),
			expRes: &pb.UserUpdateResponse{User: nil},
			expErr: codes.OK,
		},
		{
			name:   "failed, user not found",
			user:   adaptor.ToUserPbModel(*fixtures.User1),
			expRes: nil,
			expErr: codes.InvalidArgument,
		},
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			resp, errExpr := s.user.UserUpdate(s.ctx, &pb.UserUpdateRequest{
				Name: c.user.GetName(),
				Profile: &pbModels.Profile{
					Password: &c.user.Password,
					Email:    &c.user.Email,
					FullName: &c.user.FullName,
				},
			})
			st, _ := status.FromError(errExpr)

			s.Require().Equal(c.expErr, st.Code())
			s.Require().Equal(c.expRes, resp)

			if errExpr == nil {
				user, err := s.getUser(c.user.GetName())
				s.Require().NoError(err)
				s.Assert().Equal(c.user.GetPassword(), user.Password)
			}
		})
	}
}

func (s *repositorySuite) TestData_UserDelete() {
	cases := []struct {
		name   string
		user   string
		expRes *pb.UserDeleteResponse
		expErr codes.Code
	}{
		{
			name:   "success with ExistedUser",
			user:   fixtures.ExistedUser1.Name,
			expRes: &pb.UserDeleteResponse{},
			expErr: codes.OK,
		},
		{
			name:   "failed, user not found",
			user:   fixtures.User1.Name,
			expRes: nil,
			expErr: codes.InvalidArgument,
		},
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			resp, errExpr := s.user.UserDelete(s.ctx, &pb.UserDeleteRequest{Name: c.user})
			st, _ := status.FromError(errExpr)

			s.Require().Equal(c.expErr, st.Code())
			s.Require().Equal(c.expRes, resp)

			if errExpr == nil {
				_, err := s.getUser(c.user)
				s.Require().ErrorIs(err, pgx.ErrNoRows)
			}
		})
	}
}

func (s *repositorySuite) TestData_UserGet() {
	cases := []struct {
		name   string
		user   string
		expRes *pb.UserGetResponse
		expErr codes.Code
	}{
		{
			name:   "success with ExistedUser",
			user:   fixtures.ExistedUser1.Name,
			expRes: &pb.UserGetResponse{User: adaptor.ToUserPbModel(*fixtures.ExistedUser1)},
			expErr: codes.OK,
		},
		{
			name:   "failed, user not found",
			user:   fixtures.User1.Name,
			expRes: nil,
			expErr: codes.InvalidArgument,
		},
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			resp, errExpr := s.user.UserGet(s.ctx, &pb.UserGetRequest{Name: c.user})
			st, _ := status.FromError(errExpr)

			s.Require().Equal(c.expErr, st.Code())
			s.Assert().Equal(c.expRes, resp)
		})
	}
}

func (s *repositorySuite) TestData_UserList() {
	order := false
	limit := uint64(2)
	offset := uint64(0)

	cases := []struct {
		name   string
		expRes *pb.UserListResponse
		expErr codes.Code
	}{
		{
			name: "success with ExistedUser",
			expRes: &pb.UserListResponse{
				Users: adaptor.ToUserListPbModel(
					[]models.User{
						*fixtures.ExistedUser1,
						*fixtures.ExistedUser2,
					},
				),
			},
			expErr: codes.OK,
		},
	}

	for _, c := range cases {
		s.Run(c.name, func() {
			resp, errExpr := s.user.UserList(s.ctx, &pb.UserListRequest{
				Order:  order,
				Limit:  limit,
				Offset: offset,
			})
			st, _ := status.FromError(errExpr)

			s.Require().Equal(c.expErr, st.Code())
			s.Assert().Equal(c.expRes, resp)
		})
	}
}

func (s *repositorySuite) getUser(name string) (models.User, error) {
	var user models.User
	row := s.db.QueryRow(s.ctx, selectUser, name)
	err := row.Scan(&user.Name, &user.Password, &user.Email, &user.FullName, &user.CreatedAt)
	return user, err
}
