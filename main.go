package main

import (
    "fmt"
    "bufio"
    "os"
)


func main(){
    reader := bufio.NewReader(os.Stdin)
    var text []string 
    for {
        fmt.Print("Pokemon >")
        line, err := reader.ReadString('\n')
        if err != nil{
            fmt.Printf("Error while Reading line: %v", err)
            break
        }
        line = line[:len(line)-1]
        text = cleanInput(line)
        fmt.Printf("Your Command was: %s\n", text[0])
    }
}
