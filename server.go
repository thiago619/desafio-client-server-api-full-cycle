package main

import(
	"net/http"
	"io"
	"encoding/json"
	_ "github.com/ncruces/go-sqlite3/driver"
    "database/sql"
    "context"
    "time"
    "log"
    "errors"
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

type HttpError struct{
	Erro string `json:"erro"`
}


func main(){
	log.Println("[info] iniciando o servidor")
	log.Println("[info] iniciando o banco de dados")
	err := criarDatabase()

	if err != nil{
		log.Panic(err)
	}
	log.Println("[info] banco de dados iniciado com sucesso")

	mux := http.NewServeMux()

	mux.HandleFunc("GET /cotacao",cotacao)

	log.Println("[info] servidor iniciado com sucesso")
	err = http.ListenAndServe(":8080",mux)

	if(err != nil){
		log.Panic(err)
	}
}

func cotacao(w http.ResponseWriter, r *http.Request){
	w.Header().Set("Content-Type", "application/json; charset=utf-8")
	client := &http.Client{}
	url := "https://economia.awesomeapi.com.br/json/last/USD-BRL"

	ctxReq, cancelCtxReq := context.WithTimeout(r.Context(), 300 * time.Millisecond)
	defer cancelCtxReq()

	req, err := http.NewRequestWithContext(ctxReq, "GET",url,nil)
	if err != nil{
		LogHttpError(w,err)
		return
	}
	resp, err := client.Do(req)
	if err != nil{
		LogHttpError(w,err)
		return
	}
	defer resp.Body.Close()

	data, err := io.ReadAll(resp.Body)
	if err != nil{
		LogHttpError(w,err)
		return
	}

	var cotacao AAResponse

	err = json.Unmarshal(data, &cotacao)
	if err != nil{
		LogHttpError(w,err)
		return
	}

	dadosNovos := MinhaResposta{
		Moeda: "USDBRL",
		Cotacao: cotacao.USDBRL.Cotacao,
	}
	ctxDb,ctxDbCancel := context.WithTimeout(r.Context(),10*time.Millisecond)
	defer ctxDbCancel()
	err = RegistrarCotacao(ctxDb,dadosNovos)
	if err != nil{
		LogHttpError(w,err)
		return
	}

	err = json.NewEncoder(w).Encode(dadosNovos)

	if err != nil{
		LogHttpError(w,err)
		return
	}

}

func RegistrarCotacao(ctx context.Context, cotacao MinhaResposta) error{
	db, err := sql.Open("sqlite3","file:database.sqlite")
	if err != nil{
		return err
	}
	defer db.Close()
	stmt,err := db.PrepareContext(ctx,"INSERT INTO cotacoes(moeda,valor) VALUES(?,?)")
	if err != nil{
		return err
	}
	defer stmt.Close()
	_,err = stmt.ExecContext(ctx,cotacao.Moeda,cotacao.Cotacao)
	if err != nil{
		return err
	}
	return nil
}

func criarDatabase() error{
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
	return nil
}

func LogHttpError(w http.ResponseWriter, err error){
	w.WriteHeader(http.StatusInternalServerError)
	log.Printf("[erro] %v",err)
	json.NewEncoder(w).Encode(HttpError{Erro: err.Error()})
}