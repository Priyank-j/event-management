package redis

import (
	"context"
	"encoding/json"
	"math"
	"os"
	"testing"
	"time"

	"github.com/alicebob/miniredis"
	redis "github.com/go-redis/redis/v8"
)

var mr *miniredis.Miniredis

const nilMessage = "nil"

func TestMain(m *testing.M) {
	var err error
	// create mock redis
	mr, err = miniredis.Run()
	if err != nil {
		panic(err)
	}
	// init the redis package
	Init(false, mr.Addr(), mr.Addr())
	// run the test cases
	code := m.Run()
	// close redis server
	mr.Close()
	// exit with code
	os.Exit(code)
}

func TestSetGetTTL(t *testing.T) {
	ctx := context.Background()
	type testStruct struct {
		key, value string
		ttl        time.Duration
		setError   *error
		getError   *error
	}
	var testCases = []testStruct{
		{"key1", "value1", 0, nil, nil},                   // should have no ttl
		{"key2", "value2", 3 * time.Minute, nil, nil},     // with ttl
		{"key3", "value3", -2, nil, nil},                  // should behave as 0 ttl
		{"key4", "", 0, nil, nil},                         // should work
		{"", "", 0, &ErrorEmptyKey, &ErrorEmptyKey},       // throw error
		{"", "value6", 0, &ErrorEmptyKey, &ErrorEmptyKey}, // throw error
	}
	index := 0
	for _, test := range testCases {
		// set the value
		err := Set(test.key, test.value, test.ttl)
		if test.setError != nil {
			receivedErrorMsg := nilMessage
			if err != nil {
				receivedErrorMsg = err.Error()
			}
			if err != *test.setError {
				t.Errorf("Case %d: Set Error: (expected: %s, got: %s)", index+1, (*test.setError).Error(), receivedErrorMsg)
			}
			continue
		}
		if err != nil {
			t.Errorf("Case %d: Set Error: (expected: %s, got: %s)", index+1, nilMessage, err.Error())
			continue
		}
		// get the value
		value, err := Get(ctx, test.key)
		if test.getError != nil {
			receivedErrorMsg := nilMessage
			if err != nil {
				receivedErrorMsg = err.Error()
			}
			if err != *test.getError {
				t.Errorf("Case %d: Get Error: (expected: %s, got: %s)", index+1, (*test.getError).Error(), receivedErrorMsg)
			}
			continue
		}
		if err != nil {
			t.Errorf("Case %d: Get Error: (expected: %s, got: %s)", index+1, nilMessage, err.Error())
			continue
		}
		if value != test.value {
			t.Errorf("Case %d: Get Mismatch: (expected: %s, got: %s)", index+1, test.value, value)
			continue
		}
		if test.ttl > 0 {
			// verify ttl
			ttlFromRedis, err := GetTTL(ctx, test.key)
			if err != nil {
				t.Errorf("Case %d: TTL Error: (expected: %s, got: %s)", index+1, nilMessage, err.Error())
				continue
			}
			if ttlFromRedis != test.ttl {
				t.Errorf("Case %d: TTL not matching: (expected: %s, got: %s)", index+1, test.ttl, ttlFromRedis)
				continue
			}
			// fast forward ttl
			mr.FastForward(test.ttl)
			// get should give error
			_, err = Get(ctx, test.key)
			if err != redis.Nil {
				errorString := nilMessage
				if err != nil {
					errorString = err.Error()
				}
				t.Errorf("Case %d: Get Post TTL expiry: (expected: %s, got: %s)", index+1, redis.Nil.Error(), errorString)
				continue
			}
		} else {
			// fast forward 24 hours
			mr.FastForward(24 * time.Hour)
			// get should give expected value
			value, err := Get(ctx, test.key)
			if err != nil {
				t.Errorf("Case %d: Get Error No TTL post 24 hours: (expected: %s, got: %s)", index+1, nilMessage, err.Error())
				continue
			}
			if value != test.value {
				t.Errorf("Case %d: Get Mismatch No TTL post 24 hours: (expected: %s, got: %s)", index+1, test.value, value)
				continue
			}
		}
		index++
	}
	// try an extra case of empty key passed in get
	_, err := Get(ctx, "")
	if err != ErrorEmptyKey {
		receivedErrorMsg := nilMessage
		if err != nil {
			receivedErrorMsg = err.Error()
		}
		t.Errorf("Case %d: Get Error empty key: (expected: %s, got: %s)", index+1, ErrorEmptyKey.Error(), receivedErrorMsg)
	}
}

func TestDelete(t *testing.T) {
	ctx := context.Background()
	// case 1: delete existing key
	// save a key
	var keyName = "key1"
	err := Set(keyName, "value1", 5*time.Minute)
	if err != nil {
		t.Errorf("Case %d: Set Error: (expected: %s, got: %s)", 1, nilMessage, err.Error())
	}
	// try getting the key
	_, err = Get(ctx, keyName)
	if err != nil {
		t.Errorf("Case %d: Get Error: (expected: %s, got: %s)", 1, nilMessage, err.Error())
	}
	// delete the key
	err = Delete(ctx, keyName)
	if err != nil {
		t.Errorf("Case %d: Delete Error: (expected: %s, got: %s)", 1, nilMessage, err.Error())
	}
	// try getting the key
	_, err = Get(ctx, keyName)
	if err != redis.Nil {
		errorString := nilMessage
		if err != nil {
			errorString = err.Error()
		}
		t.Errorf("Case %d: Get Post Delete Error: (expected: %s, got: %s)", 1, redis.Nil.Error(), errorString)
	}
	// case 2: delete non-existing key: no error should come
	// delete the key
	err = Delete(ctx, keyName)
	if err != nil {
		t.Errorf("Case %d: Delete Error: (expected: %s, got: %s)", 2, nilMessage, err.Error())
	}
	// case 3: empty key passed in delete: error should come
	err = Delete(ctx, "")
	if err == nil {
		t.Errorf("Case %d: Delete Error: (expected: %s, got: %s)", 3, ErrorEmptyKey.Error(), err.Error())
	}
}

type dummyStruct struct {
	RandomKey1 string
	RandomKey2 int
}

func isEqual(a dummyStruct, b dummyStruct) bool {
	return a.RandomKey1 == b.RandomKey1 && a.RandomKey2 == b.RandomKey2
}

func TestSetGetStructTTL(t *testing.T) {
	ctx := context.Background()
	type testStruct struct {
		key      string
		value    interface{}
		ttl      time.Duration
		setError *error
		getError *error
	}
	var testCases = []testStruct{
		{"key1", dummyStruct{RandomKey1: "key_1", RandomKey2: 2}, 0, nil, nil},                   // should have no ttl
		{"key2", dummyStruct{RandomKey1: "key_2", RandomKey2: -2}, 3 * time.Minute, nil, nil},    // with ttl
		{"key3", dummyStruct{RandomKey1: "key_3", RandomKey2: 3}, -2, nil, nil},                  // should behave as 0 ttl
		{"key4", dummyStruct{RandomKey1: "key_4", RandomKey2: 4}, 0, nil, nil},                   // should work
		{"key5", make(chan int), 0, &ErrorUnsupportedValue, nil},                                 // invalid type specified
		{"key6", math.Inf(1), 0, &ErrorUnsupportedValue, nil},                                    // invalid value specified
		{"", dummyStruct{}, 0, &ErrorEmptyKey, &ErrorEmptyKey},                                   // throw error
		{"", dummyStruct{RandomKey1: "key_8", RandomKey2: 8}, 0, &ErrorEmptyKey, &ErrorEmptyKey}, // throw error
	}
	for index, test := range testCases {
		// set the value
		err := SetStruct(test.key, test.value, test.ttl)
		if test.setError != nil {
			receivedErrorMsg := nilMessage
			if err != nil {
				receivedErrorMsg = err.Error()
			}
			if err != *test.setError {
				t.Errorf("Case %d: Set Error: (expected: %s, got: %s)", index+1, (*test.setError).Error(), receivedErrorMsg)
			}
			continue
		}
		if err != nil {
			t.Errorf("Case %d: Set Error: (expected: %s, got: %s)", index+1, nilMessage, err.Error())
			continue
		}
		// get the value
		valueStr, err := Get(ctx, test.key)
		if test.getError != nil {
			receivedErrorMsg := nilMessage
			if err != nil {
				receivedErrorMsg = err.Error()
			}
			if err != *test.getError {
				t.Errorf("Case %d: Get Error: (expected: %s, got: %s)", index+1, (*test.getError).Error(), receivedErrorMsg)
			}
			continue
		}
		if err != nil {
			t.Errorf("Case %d: Get Error: (expected: %s, got: %s)", index+1, nilMessage, err.Error())
			continue
		}
		// convert valueStr to value
		var value dummyStruct
		err = json.Unmarshal([]byte(valueStr), &value)
		if err != nil {
			t.Errorf("Case %d: Unmarshal Error: (expected: %s, got: %s)", index+1, nilMessage, err.Error())
			continue
		}
		if !isEqual(value, test.value.(dummyStruct)) {
			t.Errorf("Case %d: Get Mismatch: (expected: %+v, got: %+v)", index+1, test.value, value)
			continue
		}
		if test.ttl > 0 {
			// verify ttl
			ttlFromRedis, err := GetTTL(ctx, test.key)
			if err != nil {
				t.Errorf("Case %d: TTL Error: (expected: %s, got: %s)", index+1, nilMessage, err.Error())
				continue
			}
			if ttlFromRedis != test.ttl {
				t.Errorf("Case %d: TTL not matching: (expected: %s, got: %s)", index+1, test.ttl, ttlFromRedis)
				continue
			}
			// fast forward ttl
			mr.FastForward(test.ttl)
			// get should give error
			_, err = Get(ctx, test.key)
			if err != redis.Nil {
				errorString := nilMessage
				if err != nil {
					errorString = err.Error()
				}
				t.Errorf("Case %d: Get Post TTL expiry: (expected: %s, got: %s)", index+1, redis.Nil.Error(), errorString)
				continue
			}
		} else {
			// fast forward 24 hours
			mr.FastForward(24 * time.Hour)
			// get should give expected value
			valueStr, err := Get(ctx, test.key)
			if err != nil {
				t.Errorf("Case %d: Get Error No TTL post 24 hours: (expected: %s, got: %s)", index+1, nilMessage, err.Error())
				continue
			}
			// convert valueStr to value
			var value dummyStruct
			err = json.Unmarshal([]byte(valueStr), &value)
			if err != nil {
				t.Errorf("Case %d: Unmarshal Error: (expected: %s, got: %s)", index+1, nilMessage, err.Error())
				continue
			}
			if !isEqual(value, test.value.(dummyStruct)) {
				t.Errorf("Case %d: Get Mismatch No TTL post 24 hours: (expected: %+v, got: %+v)", index+1, test.value, value)
				continue
			}
		}
	}
}

func TestSetGetLongTTLStruct(t *testing.T) {
	ctx := context.Background()
	type testStruct struct {
		key      string
		value    dummyStruct
		setError *error
		getError *error
	}
	var testCases = []testStruct{
		{"key1", dummyStruct{RandomKey1: "key_1", RandomKey2: 2}, nil, nil},                   // should work
		{"", dummyStruct{}, &ErrorEmptyKey, &ErrorEmptyKey},                                   // throw error
		{"", dummyStruct{RandomKey1: "key_6", RandomKey2: 6}, &ErrorEmptyKey, &ErrorEmptyKey}, // throw error
	}
	for index, test := range testCases {
		// set the value
		err := SetStructWithLongTTL(test.key, test.value)
		if test.setError != nil {
			receivedErrorMsg := nilMessage
			if err != nil {
				receivedErrorMsg = err.Error()
			}
			if err != *test.setError {
				t.Errorf("Case %d: Set Error: (expected: %s, got: %s)", index+1, (*test.setError).Error(), receivedErrorMsg)
			}
			continue
		}
		if err != nil {
			t.Errorf("Case %d: Set Error: (expected: %s, got: %s)", index+1, nilMessage, err.Error())
			continue
		}
		// get the value
		valueStr, err := Get(ctx, test.key)
		if test.getError != nil {
			receivedErrorMsg := nilMessage
			if err != nil {
				receivedErrorMsg = err.Error()
			}
			if err != *test.getError {
				t.Errorf("Case %d: Get Error: (expected: %s, got: %s)", index+1, (*test.getError).Error(), receivedErrorMsg)
			}
			continue
		}
		if err != nil {
			t.Errorf("Case %d: Get Error: (expected: %s, got: %s)", index+1, nilMessage, err.Error())
			continue
		}
		// convert valueStr to value
		var value dummyStruct
		err = json.Unmarshal([]byte(valueStr), &value)
		if err != nil {
			t.Errorf("Case %d: Unmarshal Error: (expected: %s, got: %s)", index+1, nilMessage, err.Error())
			continue
		}
		if !isEqual(value, test.value) {
			t.Errorf("Case %d: Get Mismatch: (expected: %+v, got: %+v)", index+1, test.value, value)
			continue
		}
		// verify ttl
		ttlFromRedis, err := GetTTL(ctx, test.key)
		if err != nil {
			t.Errorf("Case %d: TTL Error: (expected: %s, got: %s)", index+1, nilMessage, err.Error())
			continue
		}
		if ttlFromRedis != LongRedisTTL {
			t.Errorf("Case %d: TTL not matching: (expected: %s, got: %s)", index+1, LongRedisTTL, ttlFromRedis)
			continue
		}
		// fast forward ttl
		mr.FastForward(LongRedisTTL)
		// get should give error
		_, err = Get(ctx, test.key)
		if err != redis.Nil {
			errorString := nilMessage
			if err != nil {
				errorString = err.Error()
			}
			t.Errorf("Case %d: Get Post TTL expiry: (expected: %s, got: %s)", index+1, redis.Nil.Error(), errorString)
			continue
		}
	}
}
