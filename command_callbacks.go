package main

import(
    "fmt"
    "os"
    "net/http"
    "encoding/json"
)


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

type Pokemon struct{
    Name string `json:"name"`
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
    
    area := secondParamMap["explore"]["area"] 

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



    

