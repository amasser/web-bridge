package rest

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/duality-solutions/web-bridge/rpc/dynamic"
	"github.com/gin-gonic/gin"
)

type jsonRPC struct {
	JSONRPC interface{}   `json:"jsonrpc"`
	Method  string        `json:"method"`
	Params  []interface{} `json:"params"`
	ID      interface{}   `json:"id"`
}

func (w *WebBridgeRunner) handleJSONRPC(c *gin.Context) {
	reqInput := jsonRPC{}
	err := json.NewDecoder(c.Request.Body).Decode(&reqInput)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err})
		return
	}
	strRequest := "dynamic-cli " + reqInput.Method
	for _, param := range reqInput.Params {
		switch param.(type) {
		case int:
			val, ok := param.(int)
			if ok {
				strRequest += ` ` + fmt.Sprintf("%v", val)
			}
		case float64:
			val, ok := param.(float64)
			if ok {
				strRequest += ` ` + fmt.Sprintf("%v", val)
			}
		case bool:
			val, ok := param.(bool)
			if ok {
				strRequest += ` ` + fmt.Sprintf("%v", val)
			}
		case string:
			strRequest += ` "` + fmt.Sprintf("%v", param) + `"`
		}
	}
	reqOutput, _ := dynamic.NewRequest(strRequest)
	response, _ := <-w.dynamicd.ExecCmdRequest(reqOutput)
	c.JSON(http.StatusOK, gin.H{"result": response})
}
