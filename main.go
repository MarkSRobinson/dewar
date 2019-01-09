package main

import (
	//"net/http"
	//"log"
	//"encoding/json"
	// "strings"
	//
	//"io/ioutil"
	//"./plaid"
)
import (
	"log"
	"net/http"
	"fmt"
	//"html/template"
	//"io/ioutil"
	"path/filepath"
	"html/template"
	"io/ioutil"
)

//

//
//var public_key string = "e63017e6f943f49da21151ff618a3a"
//var baseUrl string = "https://development.plaid.com"
//
//
//func createItem() {
//	createItem := CreateItemReq{}
//	createItem.Institution_id = "ins_9"
//	createItem.Public_key = public_key
//	createItem.Initial_products = []string{"transactions"}

//
//	bolB, _ := json.Marshal(createItem)
//
//	log.Print(string(bolB[:]))
//
//	resp, err := http.Post(baseUrl + "/link/item/create", "application/json", strings.NewReader(string(bolB)))
//
//	if(err == nil) {
//		log.Print("Success")
//		body, _ := ioutil.ReadAll(resp.Body)
//		s := string(body[:])
//		log.Print(s)
//	} else {
//		log.Print("Failure")
//		log.Print(err)
//	}
//}
//
//type ExchangeTokenReq struct {
//	Secret string        `json:"secret"`
//	Client_id string `json:"client_id"`
//	Public_token string         `json:"public_token"`
//}
//
//func convertToAccessToken(public_token string) {
//
//
//	bolB, _ := json.Marshal(createItem)
//
//	log.Print(string(bolB[:]))
//
//	if(err == nil) {
//		log.Print("Success")
//		body, _ := ioutil.ReadAll(resp.Body)
//		s := string(body[:])
//		log.Print(s)
//	} else {
//		log.Print("Failure")
//		log.Print(err)
//	}
//}
//
//type AccountReq struct {
//	Secret string        `json:"secret"`
//	Client_id string `json:"client_id"`
//	Access_token string         `json:"access_token"`
//}
//
//func getAccounts(access_token string) {
//	createItem := AccountReq{}
//
//	bolB, _ := json.Marshal(createItem)
//
//	log.Print(string(bolB[:]))
//
//	resp, err := http.Post(baseUrl + "/accounts/get", "application/json", strings.NewReader(string(bolB)))
//
//	if(err == nil) {
//		log.Print("Success")
//		body, _ := ioutil.ReadAll(resp.Body)
//		s := string(body[:])
//		log.Print(s)
//	} else {
//		log.Print("Failure")
//		log.Print(err)
//	}
//}
//
//type TransactionReq struct {
//	Secret string        `json:"secret"`
//	Client_id string `json:"client_id"`
//	Access_token string         `json:"access_token"`
//	StartDate string         `json:"start_date"`
//	EndDate string         `json:"end_date"`
//}
//
//type TransactionItem struct {
//
//}
//
//type TransactionResp struct {
//	EndDate []TransactionItem         `json:"transactions"`
//}
//
//
//func getTransactions(access_token string) {
//	createItem := TransactionReq{}
//	bolB, _ := json.Marshal(createItem)
//
//	log.Print(string(bolB[:]))
//
//	resp, err := http.Post(baseUrl + "/transactions/get", "application/json", strings.NewReader(string(bolB)))
//
//	if(err == nil) {
//		log.Print("Success")
//		body, _ := ioutil.ReadAll(resp.Body)
//		s := string(body[:])
//		log.Print(s)
//	} else {
//		log.Print("Failure")
//		log.Print(err)
//	}
//}
//
//

type Page struct {
	Title string
	Body  []byte
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

type Creds struct {
	Plaid_environment string
	Plaid_public_key string
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	if(r.URL.Path == "/") {
		t, err := template.ParseFiles("index.ejs")
		//http.ServeFile(w, r, "index.ejs")
		fmt.Print(err)

		c := Creds{"sandbox", "<Public Key>"}

		t.Execute(w, c)
	} else {
		http.ServeFile(w, r, filepath.Join(".", r.URL.Path))
	}
}

func processAccessToken(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)

	fmt.Print(string([]byte(data)))
	fmt.Print(err)
}

func main() {

	fmt.Printf("Output test")
	http.HandleFunc("/", viewHandler)
	http.HandleFunc("/get_access_token", processAccessToken)
	log.Fatal(http.ListenAndServe(":8080", nil))
}
