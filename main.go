package main

import (
	"fmt"
	"log"
	"net/url"
	"os"
	"regexp"

	"flag"

	"github.com/jackdanger/collectlinks"

	"net/http"
)


var (
	ur1 string
	txtname string
	domain string
)


func crawl(uri string) []string {
	defer func(){
		if err := recover(); err != nil{
			return
		}
	}()
	sublists := make([]string, 0)
	resp, err1 := http.Get(uri)
	if err1 != nil {
		log.Print(err1)
	}

	links := collectlinks.All(resp.Body)
	resp.Body.Close()


	for _, link := range links {
		parlink:=urlparse(link,uri)
		if uri!=" "{
			write(parlink,txtname)
		}
		urls,err:=url.Parse(link)
		if err!=nil {
			log.Println(err)
			continue
		}
		if isSubdomain(urls.Host,domain){
			sublists=append(sublists,"http://"+urls.Host)
		}
	}
	sublists=remove(sublists)

	return sublists
}



func main() {
	worklists := make(chan []string)
	sublists := make([]string, 0)
	seen := make(map[string]bool)

	flag.StringVar(&ur1, "u","","待爬url,如http://www.jd.com")
	flag.StringVar(&domain, "s","","域名,如jd.com")
	flag.Parse()

	sublists = append(sublists,ur1)

	go func() { worklists <- sublists }()

	txtname=domain+".txt"

	f,err:=os.Create(txtname)
	if err!=nil{
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
	go func() { worklists <- sublists }()

	for list := range worklists {
		for _, link := range list {
			if !seen[link] {
				seen[link] = true
				go func(link string) {
					worklists <- crawl(link)
				}(link)
			}
		}
	}
}



func urlparse(href, base string) string {
	uri, err := url.Parse(href)
	if err != nil {
		return " "
	}
	baseUrl, err := url.Parse(base)
	if err != nil {
		return " "
	}
	return baseUrl.ResolveReference(uri).String()
}


func isSubdomain(rawURL, domain string) bool {
	reg,err:=regexp.MatchString(domain,rawURL)
	if err!=nil{
		log.Println(err)
	}
	if reg{
		return true
	}
	return false
}


func remove(languages []string) []string {
	result := make([]string, 0, len(languages))
	temp := map[string]struct{}{}
	for _, item := range languages {
		if _, ok := temp[item]; !ok {
			temp[item] = struct{}{}
			result = append(result, item)
		}
	}
	return result
}



func write(link string,textname string) {
	f, err := os.OpenFile(textname, os.O_APPEND|os.O_WRONLY, 0644)
	if err != nil {
		fmt.Println(err)
		return
	}
	_, err = fmt.Fprintln(f, link)
	if err != nil {
		fmt.Println(err)
		f.Close()
		return
	}
	err = f.Close()
	if err != nil {
		fmt.Println(err)
		return
	}
}
