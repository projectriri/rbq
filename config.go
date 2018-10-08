package main

var config = Config{}

type Config struct {
	Addr string `toml:"addr"`
	UUID string `toml:"uuid"`
}
