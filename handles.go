package main

import (
	"fmt"
	"strconv"
	"sync"
	"time"
)

// goroutine (concurrency) safe key-value store
var db = sync.Map{}

// https://redis.io/commands/ping/
func handlePING(contents []string) (string, error) {
	if len(contents) == 1 {
		return "PONG", nil
	} else if len(contents) == 2 {
		return contents[1], nil
	} else {
		return "", fmt.Errorf("wrong number of arguments for 'ping' command")
	}
}

// https://redis.io/commands/echo/
func handleECHO(contents []string) (string, error) {
	if len(contents) == 2 {
		return contents[1], nil
	} else {
		return "", fmt.Errorf("wrong number of arguments for 'echo' command")
	}
}

// https://redis.io/commands/set/
func handleSET(contents []string) (string, error) {
	if len(contents) == 3 {
		key := contents[1]
		value := contents[2]
		db.Store(key, value)
		return "OK", nil
	} else if len(contents) == 5 {
		key := contents[1]
		value := contents[2]
		switch contents[3] {
		case "EX":
			delta, err := strconv.Atoi(contents[4])
			if err != nil {
				return "", fmt.Errorf("value is not an integer or out of range")
			}
			if delta <= 0 {
				return "", fmt.Errorf("invalid expire time in 'set' command")
			}

			db.Store(key, value)
			time.AfterFunc(time.Duration(delta)*time.Second, func() { db.Delete(key) })
			return "OK", nil

		case "PX":
			delta, err := strconv.Atoi(contents[4])
			if err != nil {
				return "", fmt.Errorf("value is not an integer or out of range")
			}
			if delta <= 0 {
				return "", fmt.Errorf("invalid expire time in 'set' command")
			}

			db.Store(key, value)
			time.AfterFunc(time.Duration(delta)*time.Millisecond, func() { db.Delete(key) })
			return "OK", nil

		case "EXAT":
			timestamp, err := strconv.ParseInt(contents[4], 10, 64)
			if err != nil {
				return "", fmt.Errorf("value is not an integer or out of range")
			}
			if timestamp <= 0 {
				return "", fmt.Errorf("invalid expire time in 'set' command")
			}

			db.Store(key, value)

			delta := timestamp - time.Now().Unix()
			if delta <= 0 {
				return "", fmt.Errorf("invalid expire time in 'set' command")
			}

			time.AfterFunc(time.Duration(delta)*time.Second, func() { db.Delete(key) })
			return "OK", nil

		case "PXAT":
			timestamp, err := strconv.ParseInt(contents[4], 10, 64)
			if err != nil {
				return "", fmt.Errorf("value is not an integer or out of range")
			}
			if timestamp <= 0 {
				return "", fmt.Errorf("invalid expire time in 'set' command")
			}

			db.Store(key, value)

			delta := timestamp - time.Now().UnixMilli()
			if delta <= 0 {
				return "", fmt.Errorf("invalid expire time in 'set' command")
			}

			time.AfterFunc(time.Duration(delta)*time.Millisecond, func() { db.Delete(key) })
			return "OK", nil

		default:
			return "", fmt.Errorf("syntax error")
		}
	}
	return "", fmt.Errorf("wrong number of arguments for 'set' command")
}

// https://redis.io/commands/get/
func handleGET(contents []string) (string, error) {
	if len(contents) == 2 {
		key := contents[1]
		value, ok := db.Load(key)
		if !ok {
			return "", fmt.Errorf("NULL")
		}
		return value.(string), nil
	}
	return "", fmt.Errorf("wrong number of arguments for 'get' command")
}
