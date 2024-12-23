package main

import(
    "unicode"
    "strings"
)

func cleanInput(text string) []string{
    
    result := make([]string, 0)
    result = removeWhitespace(strings.ToLower(text))
    return result
}

func removeWhitespace(text string) []string{
    
    var newstring []string
    var char rune
    addflag := false
    var tempstring string
    for _, char = range text{
            
        if !unicode.IsSpace(char){
            tempstring += string(char)
            addflag = true
            continue
        }

        if addflag{
            newstring = append(newstring, tempstring)
            tempstring = ""
            addflag = false
        }
    }
    return newstring
}

func removeLeadingWhitespace(text string) string{
    var newstring string
    var char rune
    var idx_start int
    for idx_start, char = range text{
        if unicode.IsSpace(char) {
            continue
        }
        break
    }
    original_len := len(text)
    var last_char_pos int
    runes := []rune(text)
    for idx_end := idx_start;idx_end < original_len;idx_end++{
        if unicode.IsSpace(runes[idx_end]) {
            continue
        } else{
            last_char_pos= idx_end
        }
    }
    newstring = text[idx_start:last_char_pos+1]
    return newstring
}
        
