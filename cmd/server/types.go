package server

import "fmt"

// Response is embedded struct in every API's answer
type Response struct {
	Ok  bool   `json:"ok"`
	Err string `json:"error"`
}

func (r Response) Error() string {
	return r.Err
}

func ErrEmpty(fieldName string) Response {
	return Response{Err: fmt.Sprintf("%s is empty", fieldName)}
}

func ErrInvalid(fieldName string) Response {
	return Response{Err: fmt.Sprintf("%s is invalid", fieldName)}
}

var ErrInternal = Response{Err: "internal server error. PM me please https://t.me/armanokka"}
