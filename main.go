package main

import (
    "fmt"
    "bufio"
    "os"
    "time"
    "github.com/shreyasganesh0/caching"
)

type Commands struct{
    name string
    desc string
    callback func() error
}

type Pagination struct{
    next string
    prev string
}

var commandDict map[string]Commands
var apiMap map[string]*Pagination
var cache *caching.Cache
func init(){
 commandDict = map[string]Commands {
    "help":{ name: "help",
             desc: "Displays a help message",
             callback: commandHelp,
         },
     "exit":{name: "exit",
             desc: "Exit the Pokedex",
             callback: commandExit,
         },
     "mapb":{name: "mapb",
             desc: "List previous 20 locations listed in map command.",
             callback: commandMapb,
         },
     "map":{name: "map",
             desc: "List 20 locations - will keep track of the locations displayed,\nnext call will display the next list of 20 locations.\nRead mapb for previous 20 locations",
             callback: commandMap,
         },
     }
 apiMap = map[string]*Pagination {
    "map":{ next: "https://pokeapi.co/api/v2/location-area",
            prev: "",
        },
    }
 cache = caching.CreateCache(5* time.Second) 
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
        if cmd, exists := commandDict[text[0]]; exists {
            if err := cmd.callback(); err != nil{
                fmt.Printf("Error in command callback, %s \n %v ", text[0], err)
            }
         } //this calls the command as per the first word in the user input
    }
}
