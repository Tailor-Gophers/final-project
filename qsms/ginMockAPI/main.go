package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.POST("/send_sms", func(c *gin.Context) {
		data, err := ioutil.ReadAll(c.Request.Body)
		if err != nil {
			// Handle error
		}
		jsonData := make(map[string]string)
		if e := json.Unmarshal(data, &jsonData); e != nil {
			c.JSON(http.StatusBadRequest, gin.H{"msg": e.Error()})
			return
		}
		if len(jsonData["phone"]) < 11 {
			c.JSON(http.StatusBadRequest, gin.H{"msg": "invalid phone number"})
			return
		}
		fmt.Println(jsonData)
		c.JSON(200, gin.H{
			"result": fmt.Sprintf("Message sent to: %s", jsonData["phone"]),
		})
	})
	r.Run() // listen and serve on 0.0.0.0:8080
}


// func main() {
// //Encode the data
//    postBody, _ := json.Marshal(map[string]string{
//       "name":  "Toby",
//       "email": "Toby@example.com",
//    })
//    responseBody := bytes.NewBuffer(postBody)
// //Leverage Go's HTTP Post function to make request
//    resp, err := http.Post("https://postman-echo.com/post", "application/json", responseBody)
// //Handle Error
//    if err != nil {
//       log.Fatalf("An Error Occured %v", err)
//    }
//    defer resp.Body.Close()
// //Read the response body
//    body, err := ioutil.ReadAll(resp.Body)
//    if err != nil {
//       log.Fatalln(err)
//    }
//    sb := string(body)
//    log.Printf(sb)
// }