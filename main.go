package main

import (
    "fmt"
    "bufio"
    "os"
)

type Commands struct{
    name string
    desc string
    callback func() error
}
var commandMap map[string]Commands
func init(){
 commandMap = map[string]Commands {
    "help":{ name: "help",
             desc: "Displays a help message",
             callback: commandHelp,
         },
     "exit":{name: "exit",
             desc: "Exit the Pokedex",
             callback: commandExit,
         },
     }
 }

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
        if cmd, exists := commandMap[text[0]]; exists {
            if err := cmd.callback(); err != nil{
                fmt.Printf("Error in command callback, %s \n %v ", text[0], err)
            }
         } //this calls the command as per the first word in the user input
    }
}
