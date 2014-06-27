package main

type Hook struct {
	name          string
	commands      []string
	allowFailures bool
}

var hookNames = [...]string{
	"before_deploy",
	"before_setup",
	"after_setup",
	"before_code_update",
	"after_code_update",
	"before_link_release",
	"after_link_release",
	"after_deploy",
}
