package main

import(
    "strings"
    "fmt"
    "os"
    "net/http"
    "encoding/json"
    "time"
    "math/rand"
)

var capturedMap map[string]Pokemon

type Result struct{
    Name string `json:"name"`
}

type Location struct{
    Results []Result `json:"results"`
    Next string `json:"next"`
    Prev string `json:"previous"`
}

type PokemonEncounter struct{
    Pokemon Pokemon `json:"pokemon"`
}

type Type struct{
    Type Type_sub `json:"type"`
}

type Type_sub struct{
    Name string `json:"name"`
}

type Stat struct{
    Basestat int `json:"base_stat"`
    Stat Stat_sub `json:"stat"`
}

type Stat_sub struct{
    Name string `json:"name"`
}

type Pokemon struct{
    Name string `json:"name"`
    Height int `json:"height"`
    Weight int `json:"weight"`
    Stats []Stat `json:"stats"`
    Types []Type `json:"types"`
    CaptureRate int `json:"capture_rate"`
}

type LocationArea struct{
    Pokemon_encounters []PokemonEncounter `json:"pokemon_encounters"`
}


func commandExit() error{
    fmt.Printf("Closing the Pokedex... Goodbye!")
    os.Exit(0)
    return nil
}

func commandHelp() error{
    fmt.Printf("Welcome to the Pokedex!\nUsage:\n\n")
    for _, value := range commandDict{
        fmt.Printf("%s: %s\n", value.name, value.desc)
    }
    return nil
}

func apiResp[T any](url string, parsed_resp *T) error{

    cacheval, exists := cache.Get(url)
      
    if !exists{
        req, err := http.NewRequest("GET", url, nil)
        if err != nil{
            return err
        }

        client := http.Client{}
        res, err := client.Do(req)
        if err != nil{
            return err
        }

        decoder := json.NewDecoder(res.Body)
        err = decoder.Decode(&parsed_resp)

        if err != nil{
            return err
        }
        cacheMarshVal, err := json.Marshal(parsed_resp)

        if err != nil{
            return err
        }
        cache.Add(url, cacheMarshVal)
        return err 

    }else{
        err := json.Unmarshal(cacheval, &parsed_resp)
        return err
    }

}

func commandMap() error{
    pagination, exists := apiMap["map"]
    if !exists{
        return fmt.Errorf("Error retrieving url")
    }
    fullPath := pagination.next
    if fullPath == ""{
        return fmt.Errorf("Error getting the next page")
    }

    var locations Location
     
        if err := apiResp(fullPath, &locations); err != nil{
            return err
        }


    for _,location := range locations.Results{
        fmt.Printf("%s\n", location.Name)
    }
    if locations.Next == ""{
        pagination.next = fullPath
    } else{
        pagination.next = locations.Next
    }
    pagination.prev = locations.Prev
   return nil 
}

func commandMapb() error{
    pagination, exists := apiMap["map"]
    if !exists{
        return fmt.Errorf("Error retrieving url")
    }
    fullPath := apiMap["map"].prev
    if fullPath == ""{
        fmt.Printf("you're on the first page\n")
    }else{

        var locations Location

        if err := apiResp(fullPath, &locations); err != nil{
            return err
        }

        for _, location := range locations.Results{
            fmt.Printf("%s\n", location.Name)
        }
        pagination.next = locations.Next
        pagination.prev = locations.Prev
        
    }
    return nil
}
func commandExplore() error{
    
    area, _ := secondParamMap["explore"] 

    fmt.Printf("Exploring %s\n", area) 

    getLocaitonApiBase := "https://pokeapi.co/api/v2/location-area/"
    
    fullPath := getLocaitonApiBase + area + "/"
    
    var pokemonByLocation LocationArea
   
    if err := apiResp(fullPath, &pokemonByLocation); err != nil{
            return err
        }
    fmt.Printf("Found Pokemon!\n")
    for _, pokemon_encounter := range pokemonByLocation.Pokemon_encounters{
        fmt.Printf("- %s\n", pokemon_encounter.Pokemon.Name)
    }

    return nil
}

func commandCatch() error{

    pokemon,_ := secondParamMap["catch"]

    getPokemonDataApiBase := "https://pokeapi.co/api/v2/pokemon/"
    
    fullPath := getPokemonDataApiBase + pokemon + "/"
    
    var pokemonData Pokemon
    
    if err:= apiResp(fullPath, &pokemonData); err != nil{
        return err
    }

    captured := false

    var seed int64 = time.Now().UnixNano()

    rand.Seed(seed)
    
    captureRateApiBase := "https://pokeapi.co/api/v2/pokemon-species/"

    captureRateApiPath := captureRateApiBase + pokemon + "/"

    var tempCaptureRate struct{
        CaptureRate int `json:"capture_rate"`
    }

    if err:= apiResp(captureRateApiPath, &tempCaptureRate); err != nil{
        return err
    }
    pokemonData.CaptureRate = tempCaptureRate.CaptureRate

    for ;captured == false; {
        fmt.Printf("Throwing a Pokeball at %s...\n", pokemonData.Name)
        randomValue := rand.Intn(100) + 1 //generate a random number from 1.., 100

        if randomValue <= pokemonData.CaptureRate{
            captured = true
            var i int
            for ;i<3;i++{
                dots := strings.Repeat(".", i+1) 
                fmt.Printf("%s\n", dots)
                time.Sleep(time.Second)
            }
            if _, exists := capturedMap[pokemonData.Name]; !exists{
                capturedMap[pokemonData.Name] = pokemonData // disk persistence maybe if login implemented?
                fmt.Printf("%s has been caught! Data has been added to the Pokedex\n", pokemonData.Name)
            }else{
                fmt.Printf("%s has been caught!\n", pokemonData.Name)
            }
        } else{
            
            
            blockSize := (100 - pokemonData.CaptureRate) / 3
            if blockSize == 0 {
            blockSize = 1 // Avoid division by zero
            }
            blockPos := (randomValue - pokemonData.CaptureRate) / blockSize

            var i int
            for ;i < blockPos;i++ {
                dots := strings.Repeat(".", i+1) 
                fmt.Printf("%s\n", dots)
                time.Sleep(time.Second)
            }
            if i == 2{
                fmt.Printf("So close ... Maybe the next one will get it!\nWould you like to try again (y/n)?\n")
            }else{
                fmt.Printf("It broke free. Would you like to try again (y/n)?\n")
            }

            var input string
            for {
                fmt.Scanln(&input) // Reads a single input until newline
                if input == "y"{
                   break 
                }
                if input == "n"{
                    captured = true
                    break
                }else{
                    fmt.Println("Invalid choice type 'y' or 'n' in lower or uppercase only.\n")
                }
            }
            
        }
    }
    return nil
}

func commandInspect() error{
    var pokemon Pokemon
    var exists bool
    if pokemon, exists = capturedMap[secondParamMap["inspect"]]; !exists{
        fmt.Printf("You havent caught this pokemon yet.\n")
        return nil
    }
    fmt.Printf("Name: %s\nHeight: %d\nWeight %d\n", pokemon.Name, pokemon.Height, pokemon.Weight)

    fmt.Println("Stats:")
    for _, stat := range pokemon.Stats{
        fmt.Printf(" -%s: %d\n", stat.Stat.Name, stat.Basestat)
    }

    fmt.Println("Types:")
    for _, poketype := range pokemon.Types{
        fmt.Printf(" -%s\n", poketype.Type.Name)
    }
    return nil
}

func commandPokedex() error{
    fmt.Printf("Your Pokedex:")
    for pokemonName,_ := range capturedMap{
        fmt.Printf(" -%s\n", pokemonName)
    }
    return nil
}

