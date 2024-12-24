package main

import(
    "fmt"
    "os"
)


func commandExit() error{
    fmt.Printf("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}

func commandHelp() error{
    fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
    for _, value := range commandMap{
        fmt.Printf("%s: %s\n", value.name, value.desc)
    }
    return nil
}


