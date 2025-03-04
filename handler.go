package main

import (
	"sync"
)

type handler func(args []Value) Value

var Handlers = map[string]handler{
	"PING":    ping,
	"COMMAND": command,
	"SET":     set,
	"GET":     get,
	"HSET":    hset,
	"HGET":    hget,
	"HGETALL": hgetAll,
}

func ping(_ []Value) Value {
	return Value{typ: "string", str: "PONG"}
}

func command(_ []Value) Value {
	return Value{typ: "string", str: ""}
}

type hashmap[T any] map[string]T

var setMap = hashmap[string]{}
var setMapMu = sync.RWMutex{}

func set(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'set' command"}
	}

	key := args[0].bulk
	value := args[1].bulk

	setMapMu.Lock()
	setMap[key] = value
	setMapMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func get(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'get' command"}
	}

	key := args[0].bulk

	setMapMu.RLock()
	value, ok := setMap[key]
	setMapMu.RUnlock()
	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

var hsMap = hashmap[hashmap[string]]{}
var hsMapMu = sync.RWMutex{}

func hset(args []Value) Value {
	if len(args) != 3 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hset' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk
	value := args[2].bulk

	hsMapMu.Lock()
	if _, ok := hsMap[hash]; !ok {
		hsMap[hash] = hashmap[string]{}
	}

	hsMap[hash][key] = value
	hsMapMu.Unlock()

	return Value{typ: "string", str: "OK"}
}

func hget(args []Value) Value {
	if len(args) != 2 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hget' command"}
	}

	hash := args[0].bulk
	key := args[1].bulk

	hsMapMu.RLock()
	value, ok := hsMap[hash][key]
	hsMapMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	return Value{typ: "bulk", bulk: value}
}

func hgetAll(args []Value) Value {
	if len(args) != 1 {
		return Value{typ: "error", str: "ERR wrong number of arguments for 'hgetall' command"}
	}

	hash := args[0].bulk

	hsMapMu.RLock()
	hm, ok := hsMap[hash]
	hsMapMu.RUnlock()

	if !ok {
		return Value{typ: "null"}
	}

	val := []Value{}
	for k, v := range hm {
		val = append(val, Value{typ: "bulk", bulk: k})
		val = append(val, Value{typ: "bulk", bulk: v})
	}

	return Value{typ: "array", array: val}
}
