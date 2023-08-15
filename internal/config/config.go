package config

import (
	"fmt"
	"os"
	"strings"

	"github.com/ghkadim/highload_architect/internal/app/mysql"
	"github.com/ghkadim/highload_architect/internal/logger"
	"github.com/ghkadim/highload_architect/internal/models"
)

func Get[T any](variable string, defaultValue T) T {
	valueStr := os.Getenv(variable)
	if valueStr == "" {
		return defaultValue
	}

	var value T
	err := parseValue(valueStr, &value)
	if err != nil {
		logger.Infof("Failed to parse env variable %s, return defaultValue %v: %v",
			variable, defaultValue, err)
		return defaultValue
	}
	return value
}

func parseValue(valueStr string, value any) error {
	switch val := value.(type) {
	case *mysql.DedicatedShardID:
		for _, kv := range strings.Split(valueStr, ",") {
			kvArr := strings.Split(kv, ":")
			if len(kvArr) != 2 {
				return fmt.Errorf("failed to parse key value %s", kv)
			}
			var userID models.UserID
			err := parseValue(kvArr[0], &userID)
			if err != nil {
				return err
			}
			(*val)[userID] = kvArr[1]
		}
	default:
		_, err := fmt.Sscan(valueStr, val)
		if err != nil {
			return fmt.Errorf("failed to parse value '%s' to type %T: %w", valueStr, value, err)
		}
	}
	return nil
}
