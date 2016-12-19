package main

import (
    "fmt"
    "regexp"
    "strings"
    "io/ioutil"
)

func main() {
    dat, _ := ioutil.ReadFile("./input.txt")
    // fmt.Println(string(dat))
    // re := regexp.MustCompile("ALERT [a-zA-Z]+\n\\s+IF.*\n\\s+FOR\\s[0-9]+[a-z]\n\\s+LABELS {.*}\n\\s+ANNOTATIONS {\n\\s+")
    re := regexp.MustCompile("query:.*took.*")
    arr := re.FindAllString(string(dat), -1)
    for _, m := range arr {
        fmt.Println(m[6:strings.Index(m, "took")])
    }
}