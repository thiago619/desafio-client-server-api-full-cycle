package main

import(
	"fmt"
	"net/http"
)

func main(){
	fmt.Println("[SERVIDOR] - INICIALIZANDO...")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /cotacao",func(w http.ResponseWriter, r *http.Request){
		w.Write([]byte("funcionou"))
	})

	err := http.ListenAndServe(":8080",mux)

	if(err != nil){
		fmt.Println(err)
	}



}
