package sample

import (
	"testing"
	"time"
	"fmt"
	"github.com/stretchr/testify/assert"
)

type TestType int

const (
	normal   TestType = iota
	abnormal
)

func TestCreateRedisClient(t *testing.T) {

	var err error
	redisClient, err := CreateRedisClient()
	if err != nil {
		t.Errorf("%v", err)
	}
	type redisRecord struct {
		key        string
		value      string
		expiration time.Duration
	}

	defaultExpiration := time.Duration(5) * time.Minute

	testCases := []struct {
		testType TestType
		title    string
		expect   redisRecord
	}{
		{
			testType: normal,
			title:    "通常値が取得できるか",
			expect: redisRecord{
				key:        "key",
				value:      "value",
				expiration: defaultExpiration,
			},
		},
		{
			testType: normal,
			title:    "keyが空文字",
			expect: redisRecord{
				key:        "",
				value:      "blank value",
				expiration: defaultExpiration,
			},
		},
		{
			testType: normal,
			title:    "valueが空文字",
			expect: redisRecord{
				key:        "blank key",
				value:      "",
				expiration: defaultExpiration,
			},
		},
	}

	for _, testCase := range testCases {

		err = redisClient.Set(
			testCase.expect.key,
			testCase.expect.value,
			testCase.expect.expiration,
		).Err()
		if err != nil {
			t.Error(err)
		}
		actualValue, err := redisClient.Get(testCase.expect.key).Result()
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, testCase.expect.value, actualValue, fmt.Sprintf("[%s]: value of redis record", testCase.title))
	}

}
