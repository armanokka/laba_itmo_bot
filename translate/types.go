package translate

import "fmt"

type HTTPError struct {
    Code int
    Description string
}

func (c HTTPError) Error() string {
    return fmt.Sprintf("HTTP Error [code:%v]:%s", c.Code, c.Description)
}


type GoogleAPIResponse struct {
    Text, FromLang string
}

type Player struct {
    Lang string // Language code into we have to translate (ISO-6391)
}