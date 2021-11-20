package xfsmiddle

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
	"xfsmiddle/common"

	"github.com/gin-gonic/gin"
	gateway "github.com/rpcxio/rpcx-gateway"
	"github.com/smallnest/rpcx/codec"
)

type RpcGatewayConfig struct {
	RpcServerHost     string
	GatewayServerHost string
}

type RpcGateway struct {
	Server     *gin.Engine
	Token      *TokenManage
	serviceMap map[string]*service
	config     RpcGatewayConfig
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

func NewRpcGateway(rpchost, gatewayhost string, timeout int, token *TokenManage, serviceMap map[string]*service) *RpcGateway {

	config := RpcGatewayConfig{
		RpcServerHost:     rpchost,
		GatewayServerHost: gatewayhost,
	}

	return &RpcGateway{
		Server:     setupRouter(gatewayhost, timeout),
		Token:      token,
		serviceMap: serviceMap,
		config:     config,
	}
}

func setupRouter(gatewayhost string, timeout int) *gin.Engine {
	// if setting.Conf.Asc.Release {
	// 	gin.SetMode(gin.ReleaseMode)
	// }
	engine := gin.Default()
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

		// 白名单接口
		if common.IsHave(temp[0], whitelist) {
			if err := gates.sendApi(person, gates.config.GatewayServerHost); err != nil {
				createError(err.Error(), person, c.Writer)
			}
			return
		}

		// 转发接口
		group, err := gates.Token.GetToken(token)
		if err != nil {
			createError(" token does not exist", person, c.Writer)
			return
		}
		if err := tokenCheck(string(group), person["method"].(string)); err != nil {
			createError(err.Error(), person, c.Writer)
			return
		}
		if err := gates.sendApi(person, gates.config.GatewayServerHost); err != nil {
			createError(err.Error(), person, c.Writer)
		}
	})

	gates.Server.Run(":9004")
	// server.logger.Infof("RPC Service listen on: %s", ln.Addr())
	// return server.ginEngine.RunListener(ln)
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

// func createSuccess(reply interface{}, person map[string]interface{}, w gin.ResponseWriter) {
// 	out := make(map[string]interface{})
// 	out["jsonrpc"] = person["jsonrpc"]
// 	out["id"] = person["id"]
// 	out["result"] = reply
// 	bs, _ := json.Marshal(out)
// 	_, _ = w.Write(bs)
// }

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

func SetStructFieldByJsonName(method *methodType, fields map[string]interface{}) interface{} {
	var v reflect.Value
	if method.ArgType.Kind() == reflect.Ptr {
		v = reflect.New(method.ArgType.Elem())
	} else {
		v = reflect.New(method.ArgType)
	}
	v = v.Elem()

	for i := 0; i < v.NumField(); i++ {

		fieldInfo := v.Type().Field(i) // a reflect.StructField
		tag := fieldInfo.Tag           // a reflect.StructTag
		name := tag.Get("json")

		if name == "" {
			name = strings.ToLower(fieldInfo.Name)
		}
		name = strings.Split(name, ",")[0]

		if value, ok := fields[name]; ok {
			if reflect.ValueOf(value).Type() == v.FieldByName(fieldInfo.Name).Type() {
				v.FieldByName(fieldInfo.Name).Set(reflect.ValueOf(value))
			}

		}
	}
	return v.Interface()
}

// v := reflect.ValueOf(ptr).Elem() // the struct variable
func (gates *RpcGateway) sendApi(person map[string]interface{}, GatewayServerHost string) error {
	cc := &codec.MsgpackCodec{}
	temp := strings.Split(person["method"].(string), ".")

	method := gates.serviceMap[temp[0]].method[temp[1]]

	args := SetStructFieldByJsonName(method, person["params"].(map[string]interface{}))
	// fmt.Printf("argv:%v\n", argv)
	data, err := cc.Encode(args)
	if err != nil {
		return err
	}

	req, err := http.NewRequest("POST", "http://127.0.0.1:9002/", bytes.NewReader(data))
	if err != nil {
		return err
	}

	// set extra headers
	h := req.Header
	h.Set(gateway.XMessageID, "10000")
	h.Set(gateway.XMessageType, "0")
	h.Set(gateway.XSerializeType, "3")
	h.Set(gateway.XServicePath, temp[0])
	h.Set(gateway.XServiceMethod, temp[1])

	// send to gateway
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		log.Fatal("failed to call: ", err)
	}
	defer res.Body.Close()

	// handle http response
	replyData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatal("failed to read response: ", err)
	}

	// parse reply
	var reply string
	err = cc.Decode(replyData, &reply)
	if err != nil {
		fmt.Printf("eerr;%v\n", err)
		return err
	}

	// cli := NewClient(GatewayServerHost, "10s")
	// result := make(map[string])
	// return nil
	// cli.Call(person["method"].(string),)
	return nil
}
