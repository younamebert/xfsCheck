package xfsmiddle

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"time"
	"xfsmiddle/common"

	"github.com/gin-gonic/gin"
	"github.com/go-resty/resty/v2"
	gateway "github.com/rpcxio/rpcx-gateway"
	"github.com/smallnest/rpcx/codec"
)

type rpcGatewayConfig struct {
	rpcServeAddr string
	gatesHost    string
	timeOut      string
	nodeAddr     string
}

type RpcGateway struct {
	Server     *gin.Engine
	Token      *TokenManage
	serviceMap map[string]*service
	config     rpcGatewayConfig
	msgpack    codec.MsgpackCodec
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

func NewRpcGateway(rpcaddr, gatewayhost, timeout, nodeaddr string, token *TokenManage, serviceMap map[string]*service) *RpcGateway {

	config := rpcGatewayConfig{
		rpcServeAddr: rpcaddr,
		timeOut:      timeout,
		gatesHost:    gatewayhost,
		nodeAddr:     nodeaddr,
	}

	return &RpcGateway{
		Server:     setupRouter(gatewayhost, 1),
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

		// Whitelist interface
		if common.IsHave(temp[0], whitelist) {

			reply, err := gates.sendApi(person)
			if err != nil {
				createError(err.Error(), person, c.Writer)
			}
			bs, err := json.Marshal(reply)
			if err != nil {
				createError(err.Error(), person, c.Writer)
			}
			createSuccess(bs, person, c.Writer)
			return
		}

		// Forwarding interface
		group, err := gates.Token.GetToken(token)
		if err != nil {
			createError(" token does not exist", person, c.Writer)
			return
		}
		if err := tokenCheck(string(group), person["method"].(string)); err != nil {
			createError(err.Error(), person, c.Writer)
			return
		}
		reply, err := gates.reqRepeater(person)
		if err != nil {
			createError(err.Error(), person, c.Writer)
			return
		}
		createSuccess(reply, person, c.Writer)
	})

	gates.Server.Run(gates.config.gatesHost)
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

func createSuccess(reply []byte, person map[string]interface{}, w gin.ResponseWriter) {
	out := make(map[string]interface{})
	out["jsonrpc"] = person["jsonrpc"]
	out["id"] = person["id"]
	out["result"] = json.RawMessage(reply)
	// fmt.Printf("result:%v\n", string(reply))
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

func SetStructFieldByJsonName(method *methodType, params interface{}) interface{} {

	var v reflect.Value
	if method.ArgType.Kind() == reflect.Ptr {
		v = reflect.New(method.ArgType.Elem())
	} else {
		v = reflect.New(method.ArgType)
	}
	v = v.Elem()

	if params == nil {
		return v.Interface()
	}
	fields := params.(map[string]interface{})

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

func (gates *RpcGateway) sendApi(person map[string]interface{}) (interface{}, error) {

	temp := strings.Split(person["method"].(string), ".")

	method := gates.serviceMap[temp[0]].method[temp[1]]
	var args interface{}
	if params, ok := person["params"].(map[string]interface{}); ok {
		args = SetStructFieldByJsonName(method, params)
	} else {
		args = SetStructFieldByJsonName(method, nil)
	}

	data, err := gates.msgpack.Encode(args)
	if err != nil {
		return nil, err
	}
	// temp := strings.Split(person["method"].(string), ".")
	req, err := http.NewRequest("POST", gates.config.rpcServeAddr, bytes.NewReader(data))
	if err != nil {
		return nil, err
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
		return nil, err
	}
	defer res.Body.Close()

	// handle http response
	replyData, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	// parse reply
	var reply interface{}
	err = gates.msgpack.Decode(replyData, &reply)
	if err != nil {
		return nil, err
	}
	return reply, nil
}

type jsonRPCReq struct {
	JsonRPC string      `json:"jsonrpc"`
	ID      int         `json:"id"`
	Method  string      `json:"method"`
	Params  interface{} `json:"params"`
}

type jsonRPCResp struct {
	JSONRPC string      `json:"jsonrpc"`
	Result  interface{} `json:"result"`
	Error   *RPCError   `json:"error"`
	ID      int         `json:"id"`
}

func (gates *RpcGateway) reqRepeater(person map[string]interface{}) ([]byte, error) {

	data := gates.repeaterStruct(person)
	if data == nil {
		return nil, fmt.Errorf("person to struct Error")
	}

	client := resty.New()

	timeDur, err := time.ParseDuration(gates.config.timeOut)
	if err != nil {
		return nil, err
	}
	client = client.SetTimeout(timeDur)

	var resp *jsonRPCResp = nil
	_, err = client.R().
		SetHeader("Content-Type", "application/json").
		SetBody(data).
		SetResult(&resp). // or SetResult(AuthSuccess{}).
		Post(gates.config.nodeAddr)
	if err != nil {
		return nil, err
	}
	if resp == nil {
		return nil, fmt.Errorf("resp null")
	}
	e := resp.Error
	if e != nil {
		return nil, e
	}
	bs, err := json.Marshal(resp.Result)
	if err != nil {
		return nil, err
	}
	return bs, nil
}

func (gates *RpcGateway) repeaterStruct(person map[string]interface{}) *jsonRPCReq {
	id := person["id"].(json.Number)
	id2int, err := id.Int64()
	if err != nil {
		return nil
	}

	req := &jsonRPCReq{
		JsonRPC: "2.0",
		ID:      int(id2int),
		Method:  person["method"].(string),
		Params:  person["params"],
	}
	return req
}
