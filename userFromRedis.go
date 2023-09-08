package user

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"
)

const (
	prefixKey     = "user"
	expiration    = 1 * time.Hour
	offlineStatus = 0
)

type userHelper struct {
	requestExecutor    requestExecutor
	parserBody         iParserBody
	redisClient        iRedis
	userRequestBuilder iUserRequestBuilder
}

type requestExecutor func(ctx context.Context, request *http.Request) (response *http.Response, err error)

type iParserBody interface {
	Parse(inputBody *io.ReadCloser, template interface{}) error
}

type iRedis interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Keys(ctx context.Context, pattern string) ([]string, error)
	IsEqualErrorIsEmpty(err error) bool
	Del(ctx context.Context, keys []string) (int64, error)
}

type iUserRequestBuilder interface {
	ShowUser(ctx context.Context, query string, variables map[string]interface{}) *http.Request
}

func NewUserHelper(
	requestExecutor requestExecutor,
	parserBody iParserBody,
	redisClient iRedis,
	userRequestBuilder iUserRequestBuilder,
) userHelper {
	return userHelper{
		requestExecutor:    requestExecutor,
		parserBody:         parserBody,
		redisClient:        redisClient,
		userRequestBuilder: userRequestBuilder,
	}
}

func (u userHelper) GetUser(ctx context.Context, userId int) (user User, err error) {
	user, err = u.getFromCache(ctx, userId)

	if err != nil {
		return user, err
	}

	return user, err
}

func (u userHelper) GetAllUsers(ctx context.Context) ([]User, error) {
	pattern := fmt.Sprintf("%s:%s", prefixKey, "*")

	keys, err := u.redisClient.Keys(ctx, pattern)

	users := make([]User, 0)

	if err != nil {
		return users, err
	}

	if len(keys) <= 0 {
		return users, nil
	}

	for _, key := range keys {
		user := User{}
		record, err := u.redisClient.Get(ctx, key)

		if err != nil {
			continue
		}

		err = json.Unmarshal([]byte(record), &user)

		users = append(users, user)
	}

	return users, nil
}

func (u userHelper) DeleteFromCache(ctx context.Context, userId int) error {
	key := u.makeKey(userId)

	_, err := u.redisClient.Del(ctx, []string{key})

	return err
}

func (u userHelper) getFromCache(ctx context.Context, userId int) (user User, err error) {
	key := u.makeKey(userId)

	data, err := u.redisClient.Get(ctx, key)

	if err != nil && !u.redisClient.IsEqualErrorIsEmpty(err) {
		return user, err
	}

	if u.redisClient.IsEqualErrorIsEmpty(err) {
		user, err = u.getFromRemote(ctx, userId)

		if err != nil {
			return user, err
		}

		// TODO replace with validate json
		if user.Id == 0 {
			return user, errors.New("invalid user id")
		}

		u.WarmUpUserInCache(ctx, user) // TODO need log error
	} else {
		err = json.Unmarshal([]byte(data), &user)
	}

	return user, err
}

func (u userHelper) getFromRemote(ctx context.Context, userId int) (user User, err error) {
	variables := map[string]interface{}{
		"id": userId,
	}

	request := u.userRequestBuilder.ShowUser(ctx, "", variables)

	response, err := u.requestExecutor(ctx, request)

	if err != nil {
		return user, err
	}

	if response.StatusCode != http.StatusOK {
		return user, errors.New(response.Status)
	}

	data := showUserResponse{}

	err = u.parserBody.Parse(&response.Body, &data)

	if err != nil {
		return user, err
	}

	user.Id = data.Data.ShowUser.User.Id
	user.Avatar = data.Data.ShowUser.User.Avatar
	user.IsOnline = data.Data.ShowUser.User.IsOnline
	user.Teams = data.Data.ShowUser.User.Teams
	user.Username = data.Data.ShowUser.User.Username

	return user, err
}

func (u userHelper) WarmUpUserInCache(ctx context.Context, user User) error {
	data, err := json.Marshal(user)

	if err != nil {
		return err
	}

	key := u.makeKey(user.Id)

	u.redisClient.Set(ctx, key, string(data), expiration)

	return err
}

func (u userHelper) makeKey(userId int) string {
	return fmt.Sprintf("%s:%v", prefixKey, userId)
}
