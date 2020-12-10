package main

type Config struct {
	Url      string `yaml:"url"`
	Sql      string `yaml:"sql"`
	Interval int    `yaml:"interval"`
}
