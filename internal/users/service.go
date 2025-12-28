package users

import (
	"context"
	"fmt"

	repo "github.com/szuryanailham/expense-tracker/internal/adapters/sqlc"
	"github.com/szuryanailham/expense-tracker/internal/auth"
	"github.com/szuryanailham/expense-tracker/internal/env"
)

type Service interface {
	Register(ctx context.Context, arg repo.CreateUserParams) (RegisterResult, error)
	Login(ctx context.Context, arg LoginParams) (LoginResult, error)
}




type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier)Service{
	return &svc{repo: repo}
}

func(s *svc)Register(ctx context.Context,arg repo.CreateUserParams)(RegisterResult,error){
	hashedPassword, err := auth.HashPassword(arg.Password)
	if err != nil {
		return RegisterResult{},err
	}
	
	user,err := s.repo.CreateUser(ctx, repo.CreateUserParams{
	FirstName: arg.FirstName,
	LastName: arg.LastName,
	Email: arg.Email,
	Password: hashedPassword,
	})
	 if err != nil {
        return RegisterResult{}, err
    }

	token, err := auth.CreateJWT([]byte(env.GetString("JWT_SECRET","")), user.ID.String())

	 if err != nil {
        return RegisterResult{}, err
    }

	return RegisterResult{
        UserID: user.ID.String(),
        Token:  token,
    }, nil
}

func(s *svc)Login(ctx context.Context, arg LoginParams)(LoginResult, error){
	user, err  := s.repo.FindUserByEmail(ctx ,arg.Email)
	if err != nil {
		return LoginResult{}, nil
	}

	if !auth.ComparePassword(user.Password, []byte(arg.Password)){
		  return LoginResult{}, fmt.Errorf("email atau password salah")
	}

	token, err := auth.CreateJWT([]byte(env.GetString("JWT_SECRET","")), user.ID.String())

	if err != nil {
		return LoginResult{}, err
	}

	return LoginResult{
		UserID: user.ID.String(),
		Token: token,
	}, nil

}
