package main

import (
	"context"
	"github.com/redis/go-redis/v9"
	"strconv"
	"strings"
	"sync"
)

func DecodeLoginPassword(decoded string) (string, string) {
	output := strings.Split(decoded, ",")
	return output[0], output[1]
}

func EncodeLoginPassword(login, password string) string {
	return login + "," + password
}

func EncodeService(chatId int64, service string) string {
	return strconv.FormatInt(chatId, 10) + "," + service
}

func DecodeService(decoded string) (int64, string) {
	output := strings.Split(decoded, ",")
	val, _ := strconv.ParseInt(output[0], 10, 64)
	return val, output[1]
}

type Server struct {
	rb    *redis.Client
	mutex sync.Mutex
}

func NewServer() *Server {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	server := Server{rb: rdb, mutex: sync.Mutex{}}
	return &server
}

func (s *Server) Get(chatId int64, service string) (string, string, error) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	key := EncodeService(chatId, service)
	val, err := s.rb.Get(context.Background(), key).Result()
	if err != nil {
		return "", "", err
	}
	login, password := DecodeLoginPassword(val)
	return login, password, err
}

func (s *Server) Set(chatId int64, service, login, password string) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	key := EncodeService(chatId, service)
	value := EncodeLoginPassword(login, password)
	err := s.rb.Set(context.Background(), key, value, 0).Err()
	return err
}

func (s *Server) Delete(chatId int64, service string) {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	key := EncodeService(chatId, service)
	s.rb.Del(context.Background(), key)
}

func (s *Server) GetKeys() []string {
	s.mutex.Lock()
	defer s.mutex.Unlock()

	keys := make([]string, 0)
	iter := s.rb.Scan(context.Background(), 0, "", 0).Iterator()
	for iter.Next(context.Background()) {
		keys = append(keys, iter.Val())
	}
	return keys
}
