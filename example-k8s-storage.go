package main

import (
  "os"
  "fmt"
  "log"
  "time"
  ioutil "io/ioutil"
  http "net/http"
)

func statusServer() {
    http.HandleFunc("/pvDataReturn", pvDataReturn)
    http.HandleFunc("/pvDataSet", pvDataSet)

    //  create server that doesn't leave things open forever
    s := &http.Server{
            Addr:           ":8080",
            ReadTimeout:    10 * time.Second,
            WriteTimeout:   10 * time.Second,
        }
    s.ListenAndServe()
}

func pvDataReturn(w http.ResponseWriter, r *http.Request){
    if r.Method == "GET" {
        pvData, err := ioutil.ReadFile("/mnt/pv/pvData")
        if err != nil {
            fmt.Fprintf(w, "PV not present\n")
        } else {
            fmt.Fprintf(w, string(pvData[:]) + "\n")
        }
    }
}

func pvDataSet(w http.ResponseWriter, r *http.Request){
    if r.Method == "PUT" {
        bytes := make([]byte, 100)
        _, err := r.Body.Read(bytes)
        if err != nil && err.Error() != "EOF" {
            log.Println(err)
            w.WriteHeader(http.StatusInternalServerError)            
            fmt.Fprintf(w, "Failed to read PUT data\n")
        }
        
        err = ioutil.WriteFile("/mnt/pv/pvData", bytes, os.ModePerm)
        if err != nil {
            w.WriteHeader(http.StatusInternalServerError)
            fmt.Fprintf(w, "PV not present\n")
        } else {
            fmt.Fprintf(w, "Data saved\n")
        }
    }
}

func main() {
    statusServer()
}