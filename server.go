package main

import(
	"fmt"
	"net/http"
	"io"
	"encoding/json"
)

type AAResponse struct{
	USDBRL Moeda `json:"USDBRL"`
}

type Moeda struct{
	Cotacao string `json:"bid"`
}

type MinhaResposta struct{
	Moeda string `json:"moeda"`
	Cotacao string `json:"cotacao"`
}


func main(){
	fmt.Println("[SERVIDOR] - INICIALIZANDO...")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /cotacao",cotacao)

	err := http.ListenAndServe(":8080",mux)

	if(err != nil){
		fmt.Println(err)
	}
}

func cotacao(w http.ResponseWriter, r *http.Request){
	client := &http.Client{}
	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"
	req, err := http.NewRequest("GET",url,nil)
	if err != nil{
		panic(err)
	}
	resp, err := client.Do(req)
	if err != nil{
		panic(err)
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil{
		panic(err)
	}

	var cotacao AAResponse

	err = json.Unmarshal(data, &cotacao)
	if err != nil{
		panic(err)
	}

	dadosNovos := MinhaResposta{
		Moeda: "USDBRL",
		Cotacao: cotacao.USDBRL.Cotacao,
	}

	err = json.NewEncoder(w).Encode(dadosNovos)

	if err != nil{
		panic(err)
	}

}

