package main

import (
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/L-oris/go-heroku-continuous_deployment/controller"
	"github.com/julienschmidt/httprouter"
)

const defaultPort = "8080"

func main() {
	router := httprouter.New()
	controller := controller.NewController()

	router.GET("/", controller.PrintMessage)
	router.POST("/", controller.AddToMessage)
	router.DELETE("/", controller.ResetMessage)

	port := determineEnvPort()
	if err := http.ListenAndServe(port, router); err != nil {
		fmt.Println("error listening to port", port, ":", err)
		return
	}

	c := make(chan string)
	links := []string{
		"https://stackoverflow.com",
	}

	fmt.Println("checking links..")
	for _, link := range links {
		go checkLink(link, c)
	}

	for link := range c {
		fmt.Println("ADDRESS MAIN", &link)
		go func() {
			checkLink(link, c)
		}()
	}
}

func determineEnvPort() string {
	port := os.Getenv("PORT")
	log.Println("'PORT' received from env variables: ", port)

	if port == "" {
		log.Println("error 'PORT' not found in env variables, using default ", defaultPort)
		port = defaultPort
	}

	log.Println("listen on port:", port)
	return ":" + port
}

func checkLink(link string, c chan string) {
	fmt.Println("ADDRESS CHECKLINK", &link)
	_, err := http.Get(link)
	if err != nil {
		// fmt.Println(link + " may be down")
		c <- link
		return
	}

	// fmt.Println(link + " is up!")
	c <- link
}
