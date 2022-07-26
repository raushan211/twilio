package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func main() {
	//creat gin app
	app := gin.Default()
	//set route
	app.GET("/", getHandler)
	app.GET("/ws", RegisterClient)
	//run app
	app.Run(":8080")
}

func getHandler(c *gin.Context) {
	res := gin.H{
		"message": "Hello World",
	}
	c.JSON(200, res)
	c.Header("Content-Type", "application/json")
	return
}

var wsupgraders = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func RegisterClient(c *gin.Context) {

	wsupgraders.CheckOrigin = func(r *http.Request) bool { return true }
	conn, err := wsupgraders.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		msg := fmt.Sprintf("Failed to set websocket upgrade: %+v", err)
		fmt.Println(msg)
		return
	}

	// for i := 0; i < 10; i++ {
	// 	mType, mByte, err := conn.ReadMessage()
	// 	fmt.Println("mByte: ", string(mByte))
	// 	fmt.Println("mType: ", mType)
	// 	fmt.Println("err: ", err)

	// 	if string(mByte) != "quote" {
	// 		link := getGIF(string(mByte))
	// 		fmt.Println("link: ", link)
	// 		image, err := downloadImage(link)
	// 		if err == nil {
	// 			conn.WriteMessage(websocket.BinaryMessage, image)
	// 		} else {
	// 			fmt.Println("error: ", err)
	// 			conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s", "unable to load image")))
	// 		}
	// 	} else {
	// 		conn.WriteMessage(websocket.TextMessage, []byte(fmt.Sprintf("%s", getQuote())))
	// 	}

	// }
	a := getQuote()
	sendSms(a)
	sendWhatsappMessage(a)
	conn.Close()

}

type Quote struct {
	Q string `json:"q"`
	A string `json:"a"`
	H string `json:"h"`
}

func getQuote() string {

	url := "https://zenquotes.io/api/random"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return "unable to generate quote"
	}
	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "unable to generate quote"
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "unable to generate quote"
	}
	data := []Quote{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
		return "unable to generate quote"
	}
	return data[0].Q
}

//download image from url
func downloadImage(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body, nil
}

func getGIF(query string) string {

	url := "https://g.tenor.com/v1/search?q=" + query + "&key=LIVDSRZULELA&limit=1"
	method := "GET"

	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
		return "unable to generate image"
	}
	req.Header.Add("apiKey", "0UTRbFtkMxAplrohufYco5IY74U8hOes")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return "unable to generate image"
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return "unable to generate image"
	}
	data := GIF{}
	err = json.Unmarshal(body, &data)
	if err != nil {
		fmt.Println(err)
		return "unable to generate image"
	}
	return data.Results[0].Media[0].Gif.URL
}

type GIF struct {
	Results []struct {
		Media []struct {
			Gif struct {
				URL string `json:"url"`
			} `json:"gif"`
		} `json:"media"`
	} `json:"results"`
}

func sendWhatsappMessage(quote string) {

	apiurl := "https://api.twilio.com/2010-04-01/Accounts/AC30a72f28f99b51a0b080b07716bc3f19/Messages.json"
	method := "POST"

	//prepare payload post req body
	data := url.Values{}
	data.Set("To", "whatsapp:+918340477211")
	data.Set("From", "whatsapp:+14155238886")
	data.Set("MessagingServiceSid", "MG92d25e86b99bd2d2ecdd4760d09806d2")
	data.Set("Body", quote)
	client := &http.Client{}
	req, err := http.NewRequest(method, apiurl, strings.NewReader(data.Encode()))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Basic QUMzMGE3MmYyOGY5OWI1MWEwYjA4MGIwNzcxNmJjM2YxOTpkZDllZmVhYTNhZTNhODI3NTczZDhlYjJiMjMwYjBmYg==")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}

func sendSms(quote string) {

	apiurl := "https://api.twilio.com/2010-04-01/Accounts/AC30a72f28f99b51a0b080b07716bc3f19/Messages.json"
	method := "POST"

	//payload := strings.NewReader("To=%2B918340477211&MessagingServiceSid=MG92d25e86b99bd2d2ecdd4760d09806d2&Body=dhoni")
	data := url.Values{}
	data.Set("To", "+918340477211")
	data.Set("MessagingServiceSid", "MG92d25e86b99bd2d2ecdd4760d09806d2")
	data.Set("Body", quote)

	client := &http.Client{}
	req, err := http.NewRequest(method, apiurl, strings.NewReader(data.Encode()))

	if err != nil {
		fmt.Println(err)
		return
	}
	req.Header.Add("Authorization", "Basic QUMzMGE3MmYyOGY5OWI1MWEwYjA4MGIwNzcxNmJjM2YxOTpkZDllZmVhYTNhZTNhODI3NTczZDhlYjJiMjMwYjBmYg==")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	res, err := client.Do(req)
	if err != nil {
		fmt.Println(err)
		return
	}
	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(string(body))
}
