// Code generated by mockery v1.0.1. DO NOT EDIT.

package mocks

import (
	core "github.com/flyteorg/flyteplugins/go/tasks/pluginmachinery/core"
	io "github.com/flyteorg/flyteplugins/go/tasks/pluginmachinery/io"

	mock "github.com/stretchr/testify/mock"

	storage "github.com/flyteorg/flytestdlib/storage"
)

// StatusContext is an autogenerated mock type for the StatusContext type
type StatusContext struct {
	mock.Mock
}

type StatusContext_DataStore struct {
	*mock.Call
}

func (_m StatusContext_DataStore) Return(_a0 *storage.DataStore) *StatusContext_DataStore {
	return &StatusContext_DataStore{Call: _m.Call.Return(_a0)}
}

func (_m *StatusContext) OnDataStore() *StatusContext_DataStore {
	c := _m.On("DataStore")
	return &StatusContext_DataStore{Call: c}
}

func (_m *StatusContext) OnDataStoreMatch(matchers ...interface{}) *StatusContext_DataStore {
	c := _m.On("DataStore", matchers...)
	return &StatusContext_DataStore{Call: c}
}

// DataStore provides a mock function with given fields:
func (_m *StatusContext) DataStore() *storage.DataStore {
	ret := _m.Called()

	var r0 *storage.DataStore
	if rf, ok := ret.Get(0).(func() *storage.DataStore); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(*storage.DataStore)
		}
	}

	return r0
}

type StatusContext_InputReader struct {
	*mock.Call
}

func (_m StatusContext_InputReader) Return(_a0 io.InputReader) *StatusContext_InputReader {
	return &StatusContext_InputReader{Call: _m.Call.Return(_a0)}
}

func (_m *StatusContext) OnInputReader() *StatusContext_InputReader {
	c := _m.On("InputReader")
	return &StatusContext_InputReader{Call: c}
}

func (_m *StatusContext) OnInputReaderMatch(matchers ...interface{}) *StatusContext_InputReader {
	c := _m.On("InputReader", matchers...)
	return &StatusContext_InputReader{Call: c}
}

// InputReader provides a mock function with given fields:
func (_m *StatusContext) InputReader() io.InputReader {
	ret := _m.Called()

	var r0 io.InputReader
	if rf, ok := ret.Get(0).(func() io.InputReader); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.InputReader)
		}
	}

	return r0
}

type StatusContext_MaxDatasetSizeBytes struct {
	*mock.Call
}

func (_m StatusContext_MaxDatasetSizeBytes) Return(_a0 int64) *StatusContext_MaxDatasetSizeBytes {
	return &StatusContext_MaxDatasetSizeBytes{Call: _m.Call.Return(_a0)}
}

func (_m *StatusContext) OnMaxDatasetSizeBytes() *StatusContext_MaxDatasetSizeBytes {
	c := _m.On("MaxDatasetSizeBytes")
	return &StatusContext_MaxDatasetSizeBytes{Call: c}
}

func (_m *StatusContext) OnMaxDatasetSizeBytesMatch(matchers ...interface{}) *StatusContext_MaxDatasetSizeBytes {
	c := _m.On("MaxDatasetSizeBytes", matchers...)
	return &StatusContext_MaxDatasetSizeBytes{Call: c}
}

// MaxDatasetSizeBytes provides a mock function with given fields:
func (_m *StatusContext) MaxDatasetSizeBytes() int64 {
	ret := _m.Called()

	var r0 int64
	if rf, ok := ret.Get(0).(func() int64); ok {
		r0 = rf()
	} else {
		r0 = ret.Get(0).(int64)
	}

	return r0
}

type StatusContext_OutputWriter struct {
	*mock.Call
}

func (_m StatusContext_OutputWriter) Return(_a0 io.OutputWriter) *StatusContext_OutputWriter {
	return &StatusContext_OutputWriter{Call: _m.Call.Return(_a0)}
}

func (_m *StatusContext) OnOutputWriter() *StatusContext_OutputWriter {
	c := _m.On("OutputWriter")
	return &StatusContext_OutputWriter{Call: c}
}

func (_m *StatusContext) OnOutputWriterMatch(matchers ...interface{}) *StatusContext_OutputWriter {
	c := _m.On("OutputWriter", matchers...)
	return &StatusContext_OutputWriter{Call: c}
}

// OutputWriter provides a mock function with given fields:
func (_m *StatusContext) OutputWriter() io.OutputWriter {
	ret := _m.Called()

	var r0 io.OutputWriter
	if rf, ok := ret.Get(0).(func() io.OutputWriter); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(io.OutputWriter)
		}
	}

	return r0
}

type StatusContext_Resource struct {
	*mock.Call
}

func (_m StatusContext_Resource) Return(_a0 interface{}) *StatusContext_Resource {
	return &StatusContext_Resource{Call: _m.Call.Return(_a0)}
}

func (_m *StatusContext) OnResource() *StatusContext_Resource {
	c := _m.On("Resource")
	return &StatusContext_Resource{Call: c}
}

func (_m *StatusContext) OnResourceMatch(matchers ...interface{}) *StatusContext_Resource {
	c := _m.On("Resource", matchers...)
	return &StatusContext_Resource{Call: c}
}

// Resource provides a mock function with given fields:
func (_m *StatusContext) Resource() interface{} {
	ret := _m.Called()

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

type StatusContext_ResourceMeta struct {
	*mock.Call
}

func (_m StatusContext_ResourceMeta) Return(_a0 interface{}) *StatusContext_ResourceMeta {
	return &StatusContext_ResourceMeta{Call: _m.Call.Return(_a0)}
}

func (_m *StatusContext) OnResourceMeta() *StatusContext_ResourceMeta {
	c := _m.On("ResourceMeta")
	return &StatusContext_ResourceMeta{Call: c}
}

func (_m *StatusContext) OnResourceMetaMatch(matchers ...interface{}) *StatusContext_ResourceMeta {
	c := _m.On("ResourceMeta", matchers...)
	return &StatusContext_ResourceMeta{Call: c}
}

// ResourceMeta provides a mock function with given fields:
func (_m *StatusContext) ResourceMeta() interface{} {
	ret := _m.Called()

	var r0 interface{}
	if rf, ok := ret.Get(0).(func() interface{}); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(interface{})
		}
	}

	return r0
}

type StatusContext_SecretManager struct {
	*mock.Call
}

func (_m StatusContext_SecretManager) Return(_a0 core.SecretManager) *StatusContext_SecretManager {
	return &StatusContext_SecretManager{Call: _m.Call.Return(_a0)}
}

func (_m *StatusContext) OnSecretManager() *StatusContext_SecretManager {
	c := _m.On("SecretManager")
	return &StatusContext_SecretManager{Call: c}
}

func (_m *StatusContext) OnSecretManagerMatch(matchers ...interface{}) *StatusContext_SecretManager {
	c := _m.On("SecretManager", matchers...)
	return &StatusContext_SecretManager{Call: c}
}

// SecretManager provides a mock function with given fields:
func (_m *StatusContext) SecretManager() core.SecretManager {
	ret := _m.Called()

	var r0 core.SecretManager
	if rf, ok := ret.Get(0).(func() core.SecretManager); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.SecretManager)
		}
	}

	return r0
}

type StatusContext_TaskExecutionMetadata struct {
	*mock.Call
}

func (_m StatusContext_TaskExecutionMetadata) Return(_a0 core.TaskExecutionMetadata) *StatusContext_TaskExecutionMetadata {
	return &StatusContext_TaskExecutionMetadata{Call: _m.Call.Return(_a0)}
}

func (_m *StatusContext) OnTaskExecutionMetadata() *StatusContext_TaskExecutionMetadata {
	c := _m.On("TaskExecutionMetadata")
	return &StatusContext_TaskExecutionMetadata{Call: c}
}

func (_m *StatusContext) OnTaskExecutionMetadataMatch(matchers ...interface{}) *StatusContext_TaskExecutionMetadata {
	c := _m.On("TaskExecutionMetadata", matchers...)
	return &StatusContext_TaskExecutionMetadata{Call: c}
}

// TaskExecutionMetadata provides a mock function with given fields:
func (_m *StatusContext) TaskExecutionMetadata() core.TaskExecutionMetadata {
	ret := _m.Called()

	var r0 core.TaskExecutionMetadata
	if rf, ok := ret.Get(0).(func() core.TaskExecutionMetadata); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.TaskExecutionMetadata)
		}
	}

	return r0
}

type StatusContext_TaskReader struct {
	*mock.Call
}

func (_m StatusContext_TaskReader) Return(_a0 core.TaskReader) *StatusContext_TaskReader {
	return &StatusContext_TaskReader{Call: _m.Call.Return(_a0)}
}

func (_m *StatusContext) OnTaskReader() *StatusContext_TaskReader {
	c := _m.On("TaskReader")
	return &StatusContext_TaskReader{Call: c}
}

func (_m *StatusContext) OnTaskReaderMatch(matchers ...interface{}) *StatusContext_TaskReader {
	c := _m.On("TaskReader", matchers...)
	return &StatusContext_TaskReader{Call: c}
}

// TaskReader provides a mock function with given fields:
func (_m *StatusContext) TaskReader() core.TaskReader {
	ret := _m.Called()

	var r0 core.TaskReader
	if rf, ok := ret.Get(0).(func() core.TaskReader); ok {
		r0 = rf()
	} else {
		if ret.Get(0) != nil {
			r0 = ret.Get(0).(core.TaskReader)
		}
	}

	return r0
}
