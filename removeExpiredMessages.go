import (
	"context"
	"encoding/json"
	"fmt"
	"strings"
	"time"
)

const (
	poolLimit  = 5000
	count      = 10
	limitDays  = 30
	expiration = 0
)

var excludedKeysPrefixes = []string{
	"countUnread_",
	"count_limit_message_",
	"limit_message_"
}

type removeExpiredMessagesInCacheCommand struct {
	redisClient iRedisClient
}

type iRedisClient interface {
	Scan(ctx context.Context, cursor uint64, match string, count int64) ([]string, uint64, error)
	Set(ctx context.Context, key string, value string, expiration time.Duration) error
	Get(ctx context.Context, key string) (string, error)
}

func MakeRemoveExpiredMessagesInCacheCommand(redisClient iRedisClient) removeExpiredMessagesInCacheCommand {
	return removeExpiredMessagesInCacheCommand{
		redisClient: redisClient,
	}
}

func (rem removeExpiredMessagesInCacheCommand) Run(ctx context.Context) {
	fmt.Println("Started removeExpiredMessagesFromCache command")
	now := time.Now()
	numberOfScans := 0

	var cursor uint64
	var err error

	for {
		var keys []string

		keys, cursor, err = rem.redisClient.Scan(ctx, cursor, "*", count)

		if err != nil {
			fmt.Println(err)
			return
		}

		numberOfScans++

		for _, key := range keys {

			if !rem.doesKeyShouldBeExcluded(key) {
				fmt.Println("Skipped key:", key)
				continue
			}

			record, err := rem.redisClient.Get(ctx, key)

			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			messages := in_app.Messages{}

			err = json.Unmarshal([]byte(record), &messages)

			if err != nil {
				fmt.Println(err.Error())
				continue
			}

			position := 0

			for _, message := range messages {
				messageCreatedAt := time.Unix(message.CreatedAt, 0)
				diffDays := int(now.Sub(messageCreatedAt).Hours() / 24)

				if diffDays < limitDays {
					break
				}

				position++
			}

			if position > 0 {
				messages = messages[position:]

				value, err := json.Marshal(messages)

				if err != nil {
					fmt.Println(err.Error())
					continue
				}

				err = rem.redisClient.Set(ctx, key, string(value), expiration)

				if err != nil {
					fmt.Println(err.Error())
					continue
				}

				fmt.Println("Updated key:", key)
			}
		}

		if numberOfScans >= poolLimit {
			fmt.Println("Wait for connections to die")
			numberOfScans = 0
			time.Sleep(time.Second)
		}

		if cursor == 0 {
			break
		}
	}

	fmt.Println("Finished removeExpiredMessagesFromCache command")
}

func (rem removeExpiredMessagesInCacheCommand) doesKeyShouldBeExcluded(key string) bool {
	for _, exlcudedKeyPrefix := range excludedKeysPrefixes {
		if strings.Contains(key, exlcudedKeyPrefix) {
			return false
		}
	}

	return true
}
