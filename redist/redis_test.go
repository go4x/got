package redist

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
)

// TestMockRedis tests the MockRedis function
func TestMockRedis(t *testing.T) {
	client, mock := MockRedis()

	if client == nil {
		t.Error("MockRedis should return a non-nil client")
	}

	if mock == nil {
		t.Error("MockRedis should return a non-nil mock")
	}

	// Test basic mock functionality
	mock.ExpectPing().SetVal("PONG")

	ctx := context.Background()
	result, err := client.Ping(ctx).Result()
	if err != nil {
		t.Errorf("Ping should not return error, got: %v", err)
	}

	if result != "PONG" {
		t.Errorf("Expected PONG, got: %s", result)
	}

	if err := mock.ExpectationsWereMet(); err != nil {
		t.Errorf("Mock expectations were not met: %v", err)
	}
}

// TestNewMiniRedis tests the NewMiniRedis function
func TestNewMiniRedis(t *testing.T) {
	client := NewMiniRedis()

	if client == nil {
		t.Error("NewMiniRedis should return a non-nil client")
	}

	// Test basic operations
	ctx := context.Background()

	// Test SET operation
	err := client.Set(ctx, "test-key", "test-value", 0).Err()
	if err != nil {
		t.Errorf("Set operation should not return error, got: %v", err)
	}

	// Test GET operation
	val, err := client.Get(ctx, "test-key").Result()
	if err != nil {
		t.Errorf("Get operation should not return error, got: %v", err)
	}

	if val != "test-value" {
		t.Errorf("Expected 'test-value', got: %s", val)
	}

	// Test DEL operation
	err = client.Del(ctx, "test-key").Err()
	if err != nil {
		t.Errorf("Del operation should not return error, got: %v", err)
	}

	// Test GET after DEL (should return error)
	_, err = client.Get(ctx, "test-key").Result()
	if err != redis.Nil {
		t.Errorf("Get after delete should return redis.Nil, got: %v", err)
	}
}

// TestRedisOperations tests various Redis operations
func TestRedisOperations(t *testing.T) {
	client := NewMiniRedis()
	ctx := context.Background()

	// Test string operations
	t.Run("String Operations", func(t *testing.T) {
		key := "string-test"
		value := "hello world"

		// SET
		err := client.Set(ctx, key, value, time.Hour).Err()
		if err != nil {
			t.Errorf("Set failed: %v", err)
		}

		// GET
		result, err := client.Get(ctx, key).Result()
		if err != nil {
			t.Errorf("Get failed: %v", err)
		}
		if result != value {
			t.Errorf("Expected %s, got %s", value, result)
		}

		// EXISTS
		exists, err := client.Exists(ctx, key).Result()
		if err != nil {
			t.Errorf("Exists failed: %v", err)
		}
		if exists != 1 {
			t.Errorf("Expected key to exist, got %d", exists)
		}
	})

	// Test hash operations
	t.Run("Hash Operations", func(t *testing.T) {
		key := "hash-test"

		// HSET
		err := client.HSet(ctx, key, "field1", "value1").Err()
		if err != nil {
			t.Errorf("HSet failed: %v", err)
		}

		// HGET
		val, err := client.HGet(ctx, key, "field1").Result()
		if err != nil {
			t.Errorf("HGet failed: %v", err)
		}
		if val != "value1" {
			t.Errorf("Expected 'value1', got %s", val)
		}

		// HGETALL
		hash, err := client.HGetAll(ctx, key).Result()
		if err != nil {
			t.Errorf("HGetAll failed: %v", err)
		}
		if hash["field1"] != "value1" {
			t.Errorf("Expected hash field1 to be 'value1', got %s", hash["field1"])
		}
	})

	// Test list operations
	t.Run("List Operations", func(t *testing.T) {
		key := "list-test"

		// LPUSH
		err := client.LPush(ctx, key, "item1", "item2").Err()
		if err != nil {
			t.Errorf("LPush failed: %v", err)
		}

		// LLEN
		length, err := client.LLen(ctx, key).Result()
		if err != nil {
			t.Errorf("LLen failed: %v", err)
		}
		if length != 2 {
			t.Errorf("Expected list length 2, got %d", length)
		}

		// LRANGE
		items, err := client.LRange(ctx, key, 0, -1).Result()
		if err != nil {
			t.Errorf("LRange failed: %v", err)
		}
		if len(items) != 2 {
			t.Errorf("Expected 2 items, got %d", len(items))
		}
	})
}

// TestRedisErrorHandling tests error handling scenarios
func TestRedisErrorHandling(t *testing.T) {
	client := NewMiniRedis()
	ctx := context.Background()

	// Test GET on non-existent key
	_, err := client.Get(ctx, "non-existent-key").Result()
	if err != redis.Nil {
		t.Errorf("Expected redis.Nil for non-existent key, got: %v", err)
	}

	// Test HGET on non-existent hash
	_, err = client.HGet(ctx, "non-existent-hash", "field").Result()
	if err != redis.Nil {
		t.Errorf("Expected redis.Nil for non-existent hash field, got: %v", err)
	}
}

// TestRedisConcurrency tests concurrent operations
func TestRedisConcurrency(t *testing.T) {
	client := NewMiniRedis()
	ctx := context.Background()

	// Test concurrent SET operations
	done := make(chan bool, 10)

	for i := 0; i < 10; i++ {
		go func(i int) {
			key := fmt.Sprintf("concurrent-key-%d", i)
			value := fmt.Sprintf("value-%d", i)

			err := client.Set(ctx, key, value, 0).Err()
			if err != nil {
				t.Errorf("Concurrent Set failed for key %s: %v", key, err)
			}

			done <- true
		}(i)
	}

	// Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}

	// Verify all keys were set
	for i := 0; i < 10; i++ {
		key := fmt.Sprintf("concurrent-key-%d", i)
		expectedValue := fmt.Sprintf("value-%d", i)

		value, err := client.Get(ctx, key).Result()
		if err != nil {
			t.Errorf("Get failed for key %s: %v", key, err)
		}
		if value != expectedValue {
			t.Errorf("Expected %s for key %s, got %s", expectedValue, key, value)
		}
	}
}
