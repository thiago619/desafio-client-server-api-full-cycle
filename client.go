package main

import(
	"net/http"
	"encoding/json"
	"context"
	"time"
	"io"
	"log"
	"os"
	"fmt"
)

type MinhaResposta struct{
	Moeda string `json:"moeda"`
	Cotacao string `json:"cotacao"`
}

func main(){
	client := &http.Client{}
	url := "http://localhost:8080/cotacao"

	ctxReq, cancelCtxReq := context.WithTimeout(context.Background(), 300 * time.Millisecond)
	defer cancelCtxReq()

	req, err := http.NewRequestWithContext(ctxReq, "GET",url,nil)
	if err != nil{
		log.Println(err)
		return
	}
	resp, err := client.Do(req)
	if err != nil{
		log.Println(err)
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil{
		log.Println(err)
		return
	}

	var cotacao MinhaResposta

	err = json.Unmarshal(data, &cotacao)
	if err != nil{
		log.Println(err)
		return
	}
	err = RegistrarCotacao(cotacao)
	if err != nil{
		log.Println(err)
		return
	}

}

func RegistrarCotacao(cotacao MinhaResposta) error{
	f, err := os.OpenFile("cotacao.txt",os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil{
		return err
	}
	defer f.Close()
	linha := fmt.Sprintf("Dólar: %s\n", cotacao.Cotacao)
	_, err = f.WriteString(linha)
	return err
}