package test

import (
	"fmt"
	"go_test_rabbitmq/config"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestConfigDB(t *testing.T) {
	conf := config.NewAppConfig()
	testCases := []struct {
		name          string
		key           string
		expectedValue any
	}{
		{name: "successfully", key: "rabbitmq.protocol", expectedValue: "amqp"},
		{name: "successfully", key: "rabbitmq.host", expectedValue: "localhost"},
		{name: "successfully", key: "rabbitmq.port", expectedValue: "15672"},
		{name: "successfully", key: "rabbitmq.user", expectedValue: "guest"},
		{name: "successfully", key: "rabbitmq.password", expectedValue: "guest"},
		{name: "successfully", key: "app.env", expectedValue: "testing"},
		{name: "failure not exist key", key: "rabbitmq.protocols", expectedValue: nil},
		{name: "failure not exist key", key: "rabbitmq.hosts", expectedValue: nil},
		{name: "failure not exist key", key: "rabbitmq.ports", expectedValue: nil},
		{name: "failure not exist key", key: "rabbitmq.users", expectedValue: nil},
		{name: "failure not exist key", key: "rabbitmq.passwords", expectedValue: nil},
	}

	for i, testCase := range testCases {
		i++
		t.Run(fmt.Sprintf("Case %d : %s", i, testCase.name), func(t *testing.T) {
			require.IsType(t,testCase.expectedValue,conf.Get(testCase.key))
		})
	}
}
