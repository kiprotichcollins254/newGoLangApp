package main

import "fmt"

func greetings(name string) string{
	message := fmt.Sprintf("Hi, %v. Welcome!", name)
    return message
}