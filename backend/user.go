package main

type User struct {
    UUID     string   
    Bindings []string 
}

func NewUser(uuid string, bindings []string) *User {
    return &User{
        UUID:     uuid,
        Bindings: bindings,
    }
}