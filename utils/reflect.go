package utils

import "reflect"

func ReflectType(v interface{}) reflect.Type {
	tp := reflect.TypeOf(v)
	if tp.Kind() == reflect.Ptr {
		return tp.Elem()
	}
	return tp
}

func NewObject(src interface{}) interface{} {
	typ := reflect.TypeOf(src)
	if typ.Kind() == reflect.Ptr { //如果是指针类型
		typ = typ.Elem()               //获取源实际类型(否则为指针类型)
		dst := reflect.New(typ).Elem() //创建对象
		return dst.Addr().Interface()  //返回指针
	} else {
		dst := reflect.New(typ).Elem() //创建对象
		return dst.Interface()         //返回值
	}
}
