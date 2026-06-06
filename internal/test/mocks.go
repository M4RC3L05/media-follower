package test

import (
	"reflect"
	"slices"
)

type CallInfo struct {
	Args []any
}

type InterfaceMock struct {
	Calls       map[string][]CallInfo
	MockReturns map[string][][]any
}

func RecordCall(im InterfaceMock, mName string, args ...any) {
	_, ok := im.Calls[mName]

	if !ok {
		im.Calls[mName] = []CallInfo{}
	}

	im.Calls[mName] = append(
		im.Calls[mName],
		CallInfo{Args: args},
	)
}

func MockReturn(im InterfaceMock, mName string, r ...any) {
	returns, ok := im.MockReturns[mName]

	if ok {
		if len(returns) <= 0 {
			panic("no more return mocks")
		}

		retValues := returns[0]
		im.MockReturns[mName] = slices.Delete(returns, 0, 1)

		if len(r) != len(retValues) {
			panic("return value number missmatch")
		}

		for i, h := range r {
			if retValues[i] == nil {
				v := reflect.ValueOf(h)
				v.Elem().Set(reflect.Zero(v.Elem().Type()))
				continue
			}

			reflect.ValueOf(h).Elem().Set(reflect.ValueOf(retValues[i]))
		}
	} else {
		panic("must mock")
	}
}
