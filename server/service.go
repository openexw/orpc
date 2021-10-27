package server

import (
	"fmt"
	"go/ast"
	"log"
	"reflect"
	"runtime"
)

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

type methodType struct {
	method    reflect.Method // 方法名
	ArgType   reflect.Type   // 方法的参数
	ReplyType reflect.Type   // 返回值
}

func (m *methodType) newArg() (argv reflect.Value) {
	if m.ArgType.Kind() == reflect.Ptr {
		argv = reflect.New(m.ArgType.Elem())
	} else {
		argv = reflect.New(m.ArgType)
	}
	return
}

func (m *methodType) newReply() (replyv reflect.Value) {
	replyv = reflect.New(m.ReplyType.Elem())
	switch m.ReplyType.Elem().Kind() {
	case reflect.Map:
		replyv.Elem().Set(reflect.MakeMap(m.ReplyType.Elem()))
	case reflect.Slice:
		replyv.Elem().Set(reflect.MakeSlice(m.ReplyType.Elem(), 0, 0))
	}
	return
}

type service struct {
	name   string                 // service 名称
	method map[string]*methodType // service 的方法列表
	rcvr   reflect.Value          // service 对象的 value
	typ    reflect.Type           // service 的类型
}

func newService(rcvr interface{}) *service {
	s := new(service)
	s.typ = reflect.TypeOf(rcvr)
	s.rcvr = reflect.ValueOf(rcvr)

	// 获取 server name
	s.name = reflect.Indirect(s.rcvr).Type().Name()

	// 获取 method
	//s.method =
	s.method = s.suitableMethods(s.typ)
	return s
}

// call 调用 service 中的函数
func (srv *service) call(mType *methodType, argv, replyv reflect.Value) error {
	// 异常处理
	defer func() {
		if r := recover(); r != nil {
			buf := make([]byte, 4096)
			n := runtime.Stack(buf, false)
			buf = buf[:n]
			err := fmt.Errorf("[service internal error]: %v, method: %s, argv: %+v, stack: %s",
				r, mType.method.Name, argv.Interface(), buf)
			log.Println(err)
		}
	}()

	fn := mType.method.Func

	returnVals := fn.Call([]reflect.Value{srv.rcvr, argv, replyv})
	if errInter := returnVals[0].Interface(); errInter != nil {
		return errInter.(error)
	}
	log.Printf("internal data:%v", replyv.Interface())
	return nil
}

// suitableMethods 从反射中获取 Service 的 Method 列表
// Method 必须满足一下条件：
// 1. 方法可导出
// 2. 必须包含两个参数，一个参数、一个返回值，并且可被导出
// 3. 参数和返回值都必须是指针
// 4. 必须包含一个返回值且是 error 类型
func (srv *service) suitableMethods(typ reflect.Type) map[string]*methodType {
	methods := make(map[string]*methodType)
	for i := 0; i < typ.NumMethod(); i++ {
		m := typ.Method(i)
		mType := m.Type
		mName := m.Name

		// 判断是否可导出
		if !m.IsExported() {
			continue
		}
		// 1. 方法必须包含 receiver, *args, *reply.
		// 2. 必须有一个返回值
		if mType.NumIn() != 3 || mType.NumOut() != 1 {
			continue
		}
		argType := mType.In(1)
		replyType := mType.In(2)
		// *args、*reply 必须是一个指针
		if argType.Kind() != reflect.Ptr || replyType.Kind() != reflect.Ptr {
			continue
		}
		// *reply 必须是可以被导出的
		if !isExportedOrBuiltinType(replyType) {
			continue
		}
		// 判断返回值的类型
		if returnType := mType.Out(0); returnType != typeOfError {
			continue
		}
		methods[mName] = &methodType{
			method:    m,
			ArgType:   argType,
			ReplyType: replyType,
		}
	}
	return methods
}

func isExportedOrBuiltinType(t reflect.Type) bool {
	return ast.IsExported(t.Name()) || t.PkgPath() == ""
}
