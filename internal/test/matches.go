package test

import (
	"fmt"
	"reflect"

	"github.com/onsi/gomega/gcustom"
	"github.com/onsi/gomega/types"
)

func HaveInterfaceBeenCalledTimes(times int) types.GomegaMatcher {
	return gcustom.MakeMatcher(func(m struct{ InterfaceMock }) (bool, error) {
		return len(m.InterfaceMock.Calls) == times, nil
	}).WithTemplate(fmt.Sprintf("Expected interface {{ .To }} have been called %d time(s), but was called {{ len .Actual.Calls }} time(s)", times))
}

func HaveMethodBeenCalledTimes(mName string, times int) types.GomegaMatcher {
	return gcustom.MakeMatcher(func(m struct{ InterfaceMock }) (bool, error) {
		return len(m.InterfaceMock.Calls[mName]) == times, nil
	}).WithTemplate(fmt.Sprintf("Expected method \"%s\" {{ .To }} have been called %d time(s), but was called {{ len (index .Actual.Calls \"%s\") }} time(s)", mName, times, mName))
}

func HaveMethodBeenNthCalledWith(mName string, nth int, args ...any) types.GomegaMatcher {
	return gcustom.MakeMatcher(func(m struct{ InterfaceMock }) (bool, error) {
		callInfo := m.InterfaceMock.Calls[mName][nth]

		if len(callInfo.Args) != len(args) {
			return false, fmt.Errorf(
				"refusing to compare args with different lengths.\n provided length %d, atual length %d",
				len(args),
				len(callInfo.Args),
			)
		}

		return reflect.DeepEqual(args, callInfo.Args), nil
	}).
		WithTemplate(
			fmt.Sprintf(
				"Expected method \"%s\" {{ .To }} have been called with {{ format .Data }}, but was called with {{ format (index (index .Actual.Calls \"%s\") %d).Args }}",
				mName,
				mName,
				nth,
			),
		).
		WithTemplateData(args)
}
