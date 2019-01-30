package main

import (
	//"blockchain"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"proof_of_stack"
	//"util"
	//"tsl"
)

var HTML_PATH string = "src/main/html/"

type Page struct {
	Title string
	Body  []byte
}

func (p *Page) save() error {
	filename := p.Title + ".txt"
	return ioutil.WriteFile(filename, p.Body, 0600)
}

func loadPage(title string) (*Page, error) {
	filename := title + ".txt"
	body, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	return &Page{Title: title, Body: body}, nil
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi there, I love %s!", r.URL.Path[1:])
}

func loadHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hi load there, I love %s!", r.URL.Path[1:])
}

func HelloServer(w http.ResponseWriter, req *http.Request) {
	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("This is an example server.\n"))
	// fmt.Fprintf(w, "This is an example server.\n")
	// io.WriteString(w, "This is an example server.\n")
}

func startHttpServer(port string) {
	fmt.Println("start http server port " + port)
	p1 := &Page{Title: "TestPage", Body: []byte("This is a sample Page.")}
	p1.save()
	p2, _ := loadPage("TestPage")
	fmt.Println(string(p2.Body))

	http.HandleFunc("/view", viewHandler)
	http.HandleFunc("/edit/", editHandler)
	http.HandleFunc("/save/", saveHandler)
	log.Fatal(http.ListenAndServe(port, nil))
}

func main() {

	//log.Fatal(http.ListenAndServe(":8080", nil))
	/*
		http.HandleFunc("/hello", HelloServer)
		err := http.ListenAndServeTLS(":443", "server.crt", "server.key", nil)
		if err != nil {
			log.Fatal("ListenAndServe: ", err)
		}*/
	//go startHttpServer(":8080")
	//tsl.StartTLSServer()
	//fmt.Scanln()
	//fmt.Println("done")
	//blockchain.StartBlockChainServer()
	//blockchain.GenerateBlock()
	//httpMethod.StartDefaultServer()
	proof_of_stack.StartProofOfStack()
	//util.TestPublicPrivateKey()
}

func viewHandler(w http.ResponseWriter, r *http.Request) {
	//title := r.URL.Path[len("/view/"):]
	p, _ := loadPage("TestPage")
	fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", p.Title, p.Body)
}
func saveHandler(w http.ResponseWriter, r *http.Request) {
	t, _ := template.ParseFiles(HTML_PATH + "edit.html")

	details := r.FormValue("body")
	p := &Page{Title: "Receive 1", Body: []byte(details)}

	t.Execute(w, p)
	//fmt.Fprintf(w, "<h1>%s</h1><div>%s</div>", "Received", details)

}
func editHandler(w http.ResponseWriter, r *http.Request) {
	title := "TestPage" // r.URL.Path[len("/edit/"):]
	p, err := loadPage(title)
	if err != nil {
		p = &Page{Title: title}
	}
	t, _ := template.ParseFiles(HTML_PATH + "edit.html")
	t.Execute(w, p)
}
