package main

import "fmt"

func main() {
    fmt.Println("Hello, World!")

    ip := "172.31.233.128"

    err := SendFile(ip, "ubuntu", "/root/VCCSRVADM/Code/GO/SendFile/test.txt")
    if err != nil {
        fmt.Printf("Error sending file: %v\n", err)
    }
}