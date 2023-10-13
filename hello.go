package main

import "fmt" // packaeg for printing to stdio

import "rsc.io/quote"

// entry point
func main() {
    fmt.Println("Gorangu!")

    fmt.Println("Quote:")
    fmt.Println(quote.Go())
}
