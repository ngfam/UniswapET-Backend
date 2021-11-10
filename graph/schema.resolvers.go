package graph

// This file will be automatically regenerated based on the schema, any resolver implementations
// will be copied through when generating and any unknown code will be moved to the end.

import (
	"context"
	"errors"
	"fmt"
	"internal/links"
	"internal/pairs"
	"internal/tokens"
	"internal/users"
	"log"
	"pkg/jwt"
	"strconv"

	"github.com/ngfam/uniswap69420/auth"
	"github.com/ngfam/uniswap69420/graph/generated"
	"github.com/ngfam/uniswap69420/graph/model"
)

func (r *mutationResolver) CreateLink(ctx context.Context, input model.NewLink) (*model.Link, error) {
	var link links.Link

	link.Title = input.Title
	link.Address = input.Address
	linkID := link.Save()

	return &model.Link{ID: strconv.FormatInt(linkID, 10), Title: link.Title, Address: link.Address}, nil
}

func (r *mutationResolver) CreateUser(ctx context.Context, input model.NewUser) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password

	err := user.Create()
	if err != nil {
		return "", errors.New("username has already existed")
	}

	userID, _ := users.GetUserIdByUsername(user.Username)
	users.IncreaseBalance(userID, "tether", 1000000)

	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) Login(ctx context.Context, input model.Login) (string, error) {
	var user users.User
	user.Username = input.Username
	user.Password = input.Password
	correct := user.Authenticate()
	if !correct {
		// 1
		return "", &users.WrongUsernameOrPasswordError{}
	}
	token, err := jwt.GenerateToken(user.Username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) RefreshToken(ctx context.Context, input model.RefreshTokenInput) (string, error) {
	username, err := jwt.ParseToken(input.Token)
	if err != nil {
		return "", fmt.Errorf("access denied")
	}
	token, err := jwt.GenerateToken(username)
	if err != nil {
		return "", err
	}
	return token, nil
}

func (r *mutationResolver) Swap(ctx context.Context, inToken string, outToken string, inAmount float64) (bool, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return false, errors.New("access denied")
	}
	_, intermediaryToken := pairs.GetBestExchangeRate(inToken, outToken, inAmount)
	return pairs.Swap(user.ID, inToken, outToken, inAmount, intermediaryToken), nil
}

func (r *queryResolver) Links(ctx context.Context) ([]*model.Link, error) {
	pendle := "pendle"
	tether := "tether"
	var resultLinks []*model.Link
	_, intermediaryToken := pairs.GetBestExchangeRate(pendle, tether, 100000)
	done := pairs.Swap(1, "pendle", "tether", 1, intermediaryToken)
	log.Println(done)
	return resultLinks, nil

	var dbLinks []links.Link
	dbLinks = links.GetAll()

	for _, link := range dbLinks {
		resultLinks = append(resultLinks, &model.Link{ID: link.ID, Title: link.Title, Address: link.Address})
	}
	return resultLinks, nil
}

func (r *queryResolver) Tokens(ctx context.Context) ([]*model.Token, error) {
	var resultTokens []*model.Token
	var dbTokens []tokens.Token
	dbTokens = tokens.GetAll()

	for _, token := range dbTokens {
		resultTokens = append(resultTokens, &model.Token{ID: token.ID, Name: token.Name, TotalSupply: token.TotalSupply, IconURL: token.IconURL, Price: token.Price})
	}
	return resultTokens, nil
}

func (r *queryResolver) TokenSearch(ctx context.Context, prefix string) ([]*model.Token, error) {
	var resultTokens []*model.Token
	var dbTokens []tokens.Token
	dbTokens = tokens.GetByPrefix(prefix)

	for _, token := range dbTokens {
		resultTokens = append(resultTokens, &model.Token{ID: token.ID, Name: token.Name, TotalSupply: token.TotalSupply, IconURL: token.IconURL, Price: token.Price})
	}
	return resultTokens, nil
}

func (r *queryResolver) GetBestExchangeRate(ctx context.Context, inToken string, outToken string, inAmount float64) (float64, error) {
	bestRate, _ := pairs.GetBestExchangeRate(inToken, outToken, inAmount)
	return bestRate, nil
}

func (r *queryResolver) GetUserBalance(ctx context.Context, token string) (float64, error) {
	user := auth.ForContext(ctx)
	if user == nil {
		return 0, errors.New("access denied")
	}
	userBalance, _ := users.GetUserBalance(user.ID, token)
	return userBalance.Balance, nil
}

// Mutation returns generated.MutationResolver implementation.
func (r *Resolver) Mutation() generated.MutationResolver { return &mutationResolver{r} }

// Query returns generated.QueryResolver implementation.
func (r *Resolver) Query() generated.QueryResolver { return &queryResolver{r} }

type mutationResolver struct{ *Resolver }
type queryResolver struct{ *Resolver }
