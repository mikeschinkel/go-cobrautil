package testutil

import (
	"bytes"
	"net/http/httptest"
)

//goland:noinspection GoUnusedGlobalVariable
var TestLogOutput bytes.Buffer

type Test interface {
	Common() *TestCommon
}

var initializers []InitializerFunc

type InitializerFunc func(ContextForTests)

//goland:noinspection GoUnusedExportedFunction
func AddInitializer(initFunc InitializerFunc) {
	initializers = append(initializers, initFunc)
}

//goland:noinspection GoUnusedExportedFunction
func InitializeTest(c4t ContextForTests, test Test) {
	c4t.Helper()
	TestLogOutput = bytes.Buffer{}
	c4t.Update(ContextForTestArgs{
		Test:     test,
		Recorder: httptest.NewRecorder(),
	})
	for _, initFunc := range initializers {
		initFunc(c4t)
	}
}
