package main

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
	

	"github.com/PuerkitoBio/goquery"
)

var tpl *template.Template

func init() {
	tpl = template.Must(template.ParseGlob("index.html"))
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
}

func main() {
	http.HandleFunc("/", index)

	http.ListenAndServe(":4200", nil)

}

func index(w http.ResponseWriter, r *http.Request) {

	tpl.ExecuteTemplate(w, "index.html", nil)
	if r.Method != "POST" {
		fmt.Printf("The Value is not available")
		return
	}
	address := r.PostFormValue("address")

	response, err := http.Get(address)

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		log.Fatal(err)
	}
	
	u, err := url.Parse(address)
	//log.Println("u------------:", u)
	if err != nil {
		log.Fatal(err)
	}

	parts := strings.Split(u.Hostname(), ".")
	//log.Println("parts------------:", parts)

	domain := parts[len(parts)-2]
	//log.Println("domain------------:", domain)

	doc2, err := goquery.NewDocument(address)

	if err != nil {
		log.Fatal(err)
	}

	if err != nil {
		log.Fatal("http get error:", err)
	}
	//log.Println("response body -----------", response.Body)
	//log.Println("response:", body)
	if err != nil {
		log.Fatal("http get error:", err)
	}

	//var links [100]string
	var i int = 0
	var j int = 0
	var k int = 0
	iLinks := make([]string, 100)
	eLinks := make([]string, 100)
	bLinks := make([]string, 100)
	var headerCount int = 0
	var intlinks int = 0
	var extlinks int = 0
	var blinkscount int = 0
	var loginPage = false
	var dname = strings.Split(address, ".")[1]

	log.Println("dname------------:", dname)

	doc2.Find("h1,h2,h3,h4,h5,h6").Each(func(index int, item *goquery.Selection) {
		headerCount++
	})
	doc2.Find("a[href]").Each(func(index int, item *goquery.Selection) {
		href, _ := item.Attr("href")

		if strings.Contains(item.Text(), "Passwort") || strings.Contains(item.Text(), "Password") {
			loginPage = true
		}

		if strings.Contains(href, "http") {

			_, err := http.Get(href)

			if err != nil {
				bLinks = append(bLinks, href)
				k++
				log.Println("broken links from anchor tag------------:", href)
			} else {

				if strings.Contains(href, dname) {
					iLinks = append(iLinks, href)
					intlinks++

					i++
				} else {
					eLinks = append(eLinks, href)
					extlinks++

					j++
				}
			}
		} else {
			bLinks = append(bLinks, href)
			blinkscount++

			k++
		}
	})

	doc2.Find("link[href]").Each(func(index int, item *goquery.Selection) {
		href, _ := item.Attr("href")

		if strings.Contains(href, "http") {

			_, err := http.Get(href)

			if err != nil {
				bLinks = append(bLinks, href)
				k++
				log.Println("broken links from anchor tag------------:", href)
			} else {

				if strings.Contains(href, dname) {
					iLinks = append(iLinks, href)

					i++
				} else {
					eLinks = append(eLinks, href)

					j++
				}
			}
		} else {
			bLinks = append(bLinks, href)

			k++
		}
	})

	title := doc2.Find("title").Text()

	/*login page finder code*/
	doc2.Find("a[href]").Each(func(index int, item1 *goquery.Selection) {
		login1, _ := item1.Attr("forgot password")
		fmt.Printf(login1, item1.Text())

	})

	tpl.Execute(w, struct {
		Details     string
		
		Success     bool
		DomainValue string
		Intlinkssount int
		ILinks      []string
		Extlinkscount int
		ELinks      []string
		Title       string
		Blinkscount  int
		BLinks      []string
		LoginPage   bool
		HeaderCount int
		
	}{string(body), true, domain,intlinks, iLinks,extlinks, eLinks, title,blinkscount, bLinks, loginPage, headerCount})

	if err != nil {
		log.Fatal(err)
	}

}
