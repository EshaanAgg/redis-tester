package internal

import (
	"github.com/codecrafters-io/tester-utils/logger"
	"github.com/go-redis/redis"
)

type RunnableTestCase interface {
	Run(client *redis.Client, logger *logger.Logger) error
}

func runTestCases(
	testCases []RunnableTestCase,
	client *redis.Client,
	logger *logger.Logger,
) error {
	for _, testCase := range testCases {
		if err := testCase.Run(client, logger); err != nil {
			return err
		}
	}
	return nil
}
