package kash

import "time"

type EvictionPolicy uint8

const (
	Random EvictionPolicy = 0
	LRU    EvictionPolicy = 1
)

type Config struct {
	evictionPolicy EvictionPolicy
	defaultTimeout time.Duration
}