package make_translation

import (
	"context"
	"io"
	"net/http"
	"time"
)

func NewMakeTranslationHelper(
	requestExecutor func(ctx context.Context, request *http.Request) (response *http.Response, err error),
	parserBody iParserBody,
	staticContentRequestBuilder iStaticContentRequestBuilder,
	redis iRedis,
) IMakeTranslationHelper {
	return makeTranslationHelper{
		requestExecutor:             requestExecutor,
		parserBody:                  parserBody,
		staticContentRequestBuilder: staticContentRequestBuilder,
		redis:                       redis,
	}
}

type iParserBody interface {
	Parse(inputBody *io.ReadCloser, template interface{}) error
}

type iStaticContentRequestBuilder interface {
	ShowStaticContent(ctx context.Context, query string, variables map[string]interface{}) *http.Request
}

type iRedis interface {
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
}

type IMakeTranslationHelper interface {
	Handle(ctx context.Context, key string, locale string) (string, error)
}

type makeTranslationHelper struct {
	requestExecutor             func(ctx context.Context, request *http.Request) (response *http.Response, err error)
	parserBody                  iParserBody
	staticContentRequestBuilder iStaticContentRequestBuilder
	redis                       iRedis
}

func (mt makeTranslationHelper) Handle(ctx context.Context, key string, locale string) (translation string, err error) {
	translation, err = mt.redis.Get(ctx, mt.getTranslationKey(key, locale))

	if err == nil && translation != "" {
		return translation, err
	}

	response, err := mt.showStaticContent(ctx, key, locale)

	if err != nil || response.StatusCode != http.StatusOK {
		return translation, err
	}

	data, err := mt.getResponseDataShowStaticContent(&response.Body)

	if err == nil {
		translation = data.GetFirstValue()
	}

	if translation != "" {
		err = mt.saveTranslation(ctx, translation, key, locale)
	}

	return translation, err
}

func (mt makeTranslationHelper) getTranslationKey(key string, locale string) string {
	return key + "_" + locale
}

func (mt makeTranslationHelper) showStaticContent(ctx context.Context, key string, locale string) (response *http.Response, err error) {
	request := mt.staticContentRequestBuilder.ShowStaticContent(ctx, "", map[string]interface{}{
		"key":    key,
		"locale": locale,
	})

	return mt.requestExecutor(ctx, request)
}

func (mt makeTranslationHelper) getResponseDataShowStaticContent(body *io.ReadCloser) (iShowStaticContentResponse, error) {
	data := newShowStaticContentResponse()

	return data, mt.parserBody.Parse(body, &data)
}

func (mt makeTranslationHelper) saveTranslation(ctx context.Context, translation, key, locale string) error {
	const expiration = 24 * 30 * time.Hour

	return mt.redis.Set(ctx, mt.getTranslationKey(key, locale), translation, expiration)
}
