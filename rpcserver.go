package xfsmiddle

import (
	"context"
	"errors"
	"flag"
	"reflect"
	"unicode"
	"unicode/utf8"
	"xfsmiddle/logs"

	"github.com/sirupsen/logrus"
	"github.com/smallnest/rpcx/server"
)

type Rpcserver struct {
	Serve      *server.Server
	serviceMap map[string]*service
	Logs       logs.ILogger
}

type methodType struct {
	method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type
	// numCalls   uint
}

// type functionType struct {
// 	fn        reflect.Value
// 	ArgType   reflect.Type
// 	ReplyType reflect.Type
// }

type service struct {
	name   string                 // name of service
	rcvr   reflect.Value          // receiver of methods for the service
	typ    reflect.Type           // type of the receiver
	method map[string]*methodType // registered methods
	// function map[string]*functionType // registered functions
}

func NewRpcServer() *Rpcserver {
	return &Rpcserver{
		Serve:      server.NewServer(),
		Logs:       logs.NewLogger("rpcserver"),
		serviceMap: make(map[string]*service),
	}
}

func (s *Rpcserver) RegisterName(name string, rcvr interface{}) error {
	return s.register(name, rcvr)
}

func (s *Rpcserver) Start(Apihost, Timeout string) error {
	addr := flag.String("addr", Apihost, "server address")
	if err := s.Serve.Serve("tcp", *addr); err != nil {
		return err
	}
	return nil
}

func (s *Rpcserver) ServiceMap() map[string]*service {
	return s.serviceMap
}

func (s *Rpcserver) register(name string, rcvr interface{}) error {
	s.Serve.RegisterName(name, rcvr, "")
	service := new(service)
	service.typ = reflect.TypeOf(rcvr)
	service.rcvr = reflect.ValueOf(rcvr)
	sname := reflect.Indirect(service.rcvr).Type().Name() // Type
	service.name = sname
	service.method = suitableMethods(service.typ, true)
	if len(service.method) == 0 {
		var errormsg string
		// To help the user, see if a pointer receiver would work.
		method := suitableMethods(reflect.PtrTo(service.typ), false)
		if len(method) != 0 {
			errormsg = "rpcx.Register: type " + sname + " has no exported methods of suitable type (hint: pass a pointer to value of that type)"
		} else {
			errormsg = "rpcx.Register: type " + sname + " has no exported methods of suitable type"
		}
		return errors.New(errormsg)
	}
	s.serviceMap[service.name] = service
	return nil
}

// suitableMethods returns suitable Rpc methods of typ, it will report
// error using logrus if reportErr is true.
func suitableMethods(typ reflect.Type, reportErr bool) map[string]*methodType {
	methods := make(map[string]*methodType)
	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name
		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}
		// Method needs four ins: receiver, context.Context, *args, *reply.
		if mtype.NumIn() != 4 {
			if reportErr {
				logrus.Debug("method ", mname, " has wrong number of ins:", mtype.NumIn())
			}
			continue
		}
		// First arg must be context.Context
		ctxType := mtype.In(1)
		if !ctxType.Implements(reflect.TypeOf((*context.Context)(nil)).Elem()) {
			if reportErr {
				logrus.Debug("method ", mname, " must use context.Context as the first parameter")
			}
			continue
		}

		// Second arg need not be a pointer.
		argType := mtype.In(2)
		if !isExportedOrBuiltinType(argType) {
			if reportErr {
				logrus.Info(mname, " parameter type not exported: ", argType)
			}
			continue
		}
		// Third arg must be a pointer.
		replyType := mtype.In(3)
		if replyType.Kind() != reflect.Ptr {
			if reportErr {
				logrus.Info("method", mname, " reply type not a pointer:", replyType)
			}
			continue
		}
		// Reply type must be exported.
		if !isExportedOrBuiltinType(replyType) {
			if reportErr {
				logrus.Info("method", mname, " reply type not exported:", replyType)
			}
			continue
		}
		// Method needs one out.
		if mtype.NumOut() != 1 {
			if reportErr {
				logrus.Info("method", mname, " has wrong number of outs:", mtype.NumOut())
			}
			continue
		}
		// The return type of the method must be error.
		if returnType := mtype.Out(0); returnType != reflect.TypeOf((*error)(nil)).Elem() {
			if reportErr {
				logrus.Info("method", mname, " returns ", returnType.String(), " not error")
			}
			continue
		}
		methods[mname] = &methodType{method: method, ArgType: argType, ReplyType: replyType}

		// init pool for reflect.Type of args and reply
		reflectTypePools.Init(argType)
		reflectTypePools.Init(replyType)
	}
	return methods
}

func isExported(name string) bool {
	rune, _ := utf8.DecodeRuneInString(name)
	return unicode.IsUpper(rune)
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	for t.Kind() == reflect.Ptr {
		t = t.Elem()
	}
	// PkgPath will be non-empty even for an exported type,
	// so we need to check the type name as well.
	return isExported(t.Name()) || t.PkgPath() == ""
}
