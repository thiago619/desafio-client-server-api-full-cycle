package main

import(
	"fmt"
	"net/http"
	"io"
	"encoding/json"
	_ "github.com/ncruces/go-sqlite3/driver"
    _ "github.com/ncruces/go-sqlite3/embed"
    "database/sql"
    "context"
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

	err = RegistrarCotacao(r.Context(),dadosNovos)
	if err != nil{
		panic(err)
	}

	err = json.NewEncoder(w).Encode(dadosNovos)

	if err != nil{
		panic(err)
	}

}

func RegistrarCotacao(ctx context.Context, cotacao MinhaResposta) error{
	select{
	case <- ctx.Done():
		return ctx.Err()
	default:
		db, err := sql.Open("sqlite3","file:database.sqlite")
		if err != nil{
			return err
		}
		query := "CREATE TABLE IF NOT EXISTS cotacoes(id INTEGER PRIMARY KEY, moeda TEXT NOT NULL,valor NUMERIC NOT NULL,data_hora TEXT NOT NULL DEFAULT CURRENT_TIMESTAMP)"
		defer db.Close()
		_, err = db.Exec(query)
		if err != nil{
			return err
		}
		stmt,err := db.Prepare("INSERT INTO cotacoes(moeda,valor) VALUES(?,?)")
		if err != nil{
			return err
		}
		defer stmt.Close()
		_,err = stmt.Exec(cotacao.Moeda,cotacao.Cotacao)
		if err != nil{
			return err
		}
		return nil
	}
}

