package xfsmiddle

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
	"xfsmiddle/common"

	"github.com/gin-gonic/gin"
)

type RpcGatewayConfig struct {
	RpcServerHost     string
	GatewayServerHost string
}

type RpcGateway struct {
	Server *gin.Engine
	Token  *TokenManage
	config RpcGatewayConfig
}

var whitelist = []string{"Token"}

func timeoutMiddleware(timeout time.Duration) func(c *gin.Context) {
	return func(c *gin.Context) {

		// wrap the request context with a timeout
		ctx, cancel := context.WithTimeout(c.Request.Context(), timeout)

		defer func() {
			// check if context timeout was reached
			if ctx.Err() == context.DeadlineExceeded {

				// write response and abort the request
				c.Writer.WriteHeader(http.StatusGatewayTimeout)
				c.Abort()
			}

			//cancel to clear resources after finished
			cancel()
		}()

		// replace request with context wrapped request
		c.Request = c.Request.WithContext(ctx)
		c.Next()
	}
}

func NewRpcGateway(rpchost, gatewayhost string, timeout int, token *TokenManage) *RpcGateway {

	config := RpcGatewayConfig{
		RpcServerHost:     rpchost,
		GatewayServerHost: gatewayhost,
	}

	return &RpcGateway{
		Server: setupRouter(gatewayhost, timeout),
		Token:  token,
		config: config,
	}
}

func setupRouter(gatewayhost string, timeout int) *gin.Engine {
	engine := gin.New()
	engine.Use(timeoutMiddleware(time.Second * time.Duration(timeout)))
	return engine
}

func (gates *RpcGateway) Start() {
	gates.Server.Any("/", func(c *gin.Context) {
		person, err := formatRule(c.Request)
		if err != nil {
			createError(err.Error(), person, c.Writer)
			return
		}

		token := c.Query("token")
		temp := strings.Split(person["method"].(string), ".") // get methods
		if !common.IsHave(temp[0], whitelist) {
			sendApi(person)
			return
		}

		group, err := gates.Token.GetToken(token)
		if err != nil {
			createError(err.Error(), person, c.Writer)
			return
		}

		if err := tokenCheck(string(group), person["method"].(string)); err != nil {
			createError(err.Error(), person, c.Writer)
			return
		}
		sendApi(person)
	})
}

func formatRule(req *http.Request) (map[string]interface{}, error) {

	if req.Method != "POST" {
		return nil, errors.New("request mode error")
	}
	person := make(map[string]interface{})

	body, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return nil, err
	}

	decoder := json.NewDecoder(bytes.NewReader(body))
	decoder.UseNumber()
	decoder.Decode(&person)

	if _, ok := person["method"]; !ok {
		return nil, errors.New("no rule no match")
	}
	if _, ok := person["params"]; !ok {
		return nil, errors.New("no rule no match")
	}
	if _, ok := person["id"]; !ok {
		return nil, errors.New("no rule no match")

	}
	if _, ok := person["jsonrpc"]; !ok {
		return nil, errors.New("no rule no match")
	}

	return person, nil
}

func createError(errmsg string, person map[string]interface{}, w gin.ResponseWriter) {
	out := make(map[string]interface{})
	out["jsonrpc"] = person["jsonrpc"]
	out["id"] = person["id"]
	out["error"] = errmsg
	bs, _ := json.Marshal(out)
	_, _ = w.Write(bs)
}

func tokenCheck(group, methods string) error {
	want := strings.Split(string(group), ",") // user group all
	got := strings.Split(methods, ".")        // get methods

	if !common.IsHave(got[0], want) {
		if !common.IsHave(methods, want) {
			return errors.New("insufficient token permissions")
		}
		return nil
	}
	return nil

}

func sendApi(person map[string]interface{}) {

}

// func (gates *RpcGateway) specialHandling() {

// }
