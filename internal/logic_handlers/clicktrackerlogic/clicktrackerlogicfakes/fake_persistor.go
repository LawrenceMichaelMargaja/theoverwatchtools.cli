// Code generated by counterfeiter. DO NOT EDIT.
package clicktrackerlogicfakes

import (
	"context"
	"sync"

	"github.com/dembygenesis/local.tools/internal/model"
	"github.com/dembygenesis/local.tools/internal/persistence"
)

type FakePersistor struct {
	AddClickTrackerStub        func(context.Context, persistence.TransactionHandler, *model.CreateClickTracker) (*model.ClickTracker, error)
	addClickTrackerMutex       sync.RWMutex
	addClickTrackerArgsForCall []struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 *model.CreateClickTracker
	}
	addClickTrackerReturns struct {
		result1 *model.ClickTracker
		result2 error
	}
	addClickTrackerReturnsOnCall map[int]struct {
		result1 *model.ClickTracker
		result2 error
	}
	DeleteClickTrackerStub        func(context.Context, persistence.TransactionHandler, int) error
	deleteClickTrackerMutex       sync.RWMutex
	deleteClickTrackerArgsForCall []struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 int
	}
	deleteClickTrackerReturns struct {
		result1 error
	}
	deleteClickTrackerReturnsOnCall map[int]struct {
		result1 error
	}
	GetClickTrackerByNameStub        func(context.Context, persistence.TransactionHandler, string) (*model.ClickTracker, error)
	getClickTrackerByNameMutex       sync.RWMutex
	getClickTrackerByNameArgsForCall []struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 string
	}
	getClickTrackerByNameReturns struct {
		result1 *model.ClickTracker
		result2 error
	}
	getClickTrackerByNameReturnsOnCall map[int]struct {
		result1 *model.ClickTracker
		result2 error
	}
	GetClickTrackersStub        func(context.Context, persistence.TransactionHandler, *model.ClickTrackerFilters) (*model.PaginatedClickTrackers, error)
	getClickTrackersMutex       sync.RWMutex
	getClickTrackersArgsForCall []struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 *model.ClickTrackerFilters
	}
	getClickTrackersReturns struct {
		result1 *model.PaginatedClickTrackers
		result2 error
	}
	getClickTrackersReturnsOnCall map[int]struct {
		result1 *model.PaginatedClickTrackers
		result2 error
	}
	RestoreClickTrackerStub        func(context.Context, persistence.TransactionHandler, int) error
	restoreClickTrackerMutex       sync.RWMutex
	restoreClickTrackerArgsForCall []struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 int
	}
	restoreClickTrackerReturns struct {
		result1 error
	}
	restoreClickTrackerReturnsOnCall map[int]struct {
		result1 error
	}
	UpdateClickTrackerStub        func(context.Context, persistence.TransactionHandler, *model.UpdateClickTracker) (*model.ClickTracker, error)
	updateClickTrackerMutex       sync.RWMutex
	updateClickTrackerArgsForCall []struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 *model.UpdateClickTracker
	}
	updateClickTrackerReturns struct {
		result1 *model.ClickTracker
		result2 error
	}
	updateClickTrackerReturnsOnCall map[int]struct {
		result1 *model.ClickTracker
		result2 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *FakePersistor) AddClickTracker(arg1 context.Context, arg2 persistence.TransactionHandler, arg3 *model.CreateClickTracker) (*model.ClickTracker, error) {
	fake.addClickTrackerMutex.Lock()
	ret, specificReturn := fake.addClickTrackerReturnsOnCall[len(fake.addClickTrackerArgsForCall)]
	fake.addClickTrackerArgsForCall = append(fake.addClickTrackerArgsForCall, struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 *model.CreateClickTracker
	}{arg1, arg2, arg3})
	stub := fake.AddClickTrackerStub
	fakeReturns := fake.addClickTrackerReturns
	fake.recordInvocation("AddClickTracker", []interface{}{arg1, arg2, arg3})
	fake.addClickTrackerMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePersistor) AddClickTrackerCallCount() int {
	fake.addClickTrackerMutex.RLock()
	defer fake.addClickTrackerMutex.RUnlock()
	return len(fake.addClickTrackerArgsForCall)
}

func (fake *FakePersistor) AddClickTrackerCalls(stub func(context.Context, persistence.TransactionHandler, *model.CreateClickTracker) (*model.ClickTracker, error)) {
	fake.addClickTrackerMutex.Lock()
	defer fake.addClickTrackerMutex.Unlock()
	fake.AddClickTrackerStub = stub
}

func (fake *FakePersistor) AddClickTrackerArgsForCall(i int) (context.Context, persistence.TransactionHandler, *model.CreateClickTracker) {
	fake.addClickTrackerMutex.RLock()
	defer fake.addClickTrackerMutex.RUnlock()
	argsForCall := fake.addClickTrackerArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakePersistor) AddClickTrackerReturns(result1 *model.ClickTracker, result2 error) {
	fake.addClickTrackerMutex.Lock()
	defer fake.addClickTrackerMutex.Unlock()
	fake.AddClickTrackerStub = nil
	fake.addClickTrackerReturns = struct {
		result1 *model.ClickTracker
		result2 error
	}{result1, result2}
}

func (fake *FakePersistor) AddClickTrackerReturnsOnCall(i int, result1 *model.ClickTracker, result2 error) {
	fake.addClickTrackerMutex.Lock()
	defer fake.addClickTrackerMutex.Unlock()
	fake.AddClickTrackerStub = nil
	if fake.addClickTrackerReturnsOnCall == nil {
		fake.addClickTrackerReturnsOnCall = make(map[int]struct {
			result1 *model.ClickTracker
			result2 error
		})
	}
	fake.addClickTrackerReturnsOnCall[i] = struct {
		result1 *model.ClickTracker
		result2 error
	}{result1, result2}
}

func (fake *FakePersistor) DeleteClickTracker(arg1 context.Context, arg2 persistence.TransactionHandler, arg3 int) error {
	fake.deleteClickTrackerMutex.Lock()
	ret, specificReturn := fake.deleteClickTrackerReturnsOnCall[len(fake.deleteClickTrackerArgsForCall)]
	fake.deleteClickTrackerArgsForCall = append(fake.deleteClickTrackerArgsForCall, struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 int
	}{arg1, arg2, arg3})
	stub := fake.DeleteClickTrackerStub
	fakeReturns := fake.deleteClickTrackerReturns
	fake.recordInvocation("DeleteClickTracker", []interface{}{arg1, arg2, arg3})
	fake.deleteClickTrackerMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakePersistor) DeleteClickTrackerCallCount() int {
	fake.deleteClickTrackerMutex.RLock()
	defer fake.deleteClickTrackerMutex.RUnlock()
	return len(fake.deleteClickTrackerArgsForCall)
}

func (fake *FakePersistor) DeleteClickTrackerCalls(stub func(context.Context, persistence.TransactionHandler, int) error) {
	fake.deleteClickTrackerMutex.Lock()
	defer fake.deleteClickTrackerMutex.Unlock()
	fake.DeleteClickTrackerStub = stub
}

func (fake *FakePersistor) DeleteClickTrackerArgsForCall(i int) (context.Context, persistence.TransactionHandler, int) {
	fake.deleteClickTrackerMutex.RLock()
	defer fake.deleteClickTrackerMutex.RUnlock()
	argsForCall := fake.deleteClickTrackerArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakePersistor) DeleteClickTrackerReturns(result1 error) {
	fake.deleteClickTrackerMutex.Lock()
	defer fake.deleteClickTrackerMutex.Unlock()
	fake.DeleteClickTrackerStub = nil
	fake.deleteClickTrackerReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakePersistor) DeleteClickTrackerReturnsOnCall(i int, result1 error) {
	fake.deleteClickTrackerMutex.Lock()
	defer fake.deleteClickTrackerMutex.Unlock()
	fake.DeleteClickTrackerStub = nil
	if fake.deleteClickTrackerReturnsOnCall == nil {
		fake.deleteClickTrackerReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.deleteClickTrackerReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakePersistor) GetClickTrackerByName(arg1 context.Context, arg2 persistence.TransactionHandler, arg3 string) (*model.ClickTracker, error) {
	fake.getClickTrackerByNameMutex.Lock()
	ret, specificReturn := fake.getClickTrackerByNameReturnsOnCall[len(fake.getClickTrackerByNameArgsForCall)]
	fake.getClickTrackerByNameArgsForCall = append(fake.getClickTrackerByNameArgsForCall, struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 string
	}{arg1, arg2, arg3})
	stub := fake.GetClickTrackerByNameStub
	fakeReturns := fake.getClickTrackerByNameReturns
	fake.recordInvocation("GetClickTrackerByName", []interface{}{arg1, arg2, arg3})
	fake.getClickTrackerByNameMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePersistor) GetClickTrackerByNameCallCount() int {
	fake.getClickTrackerByNameMutex.RLock()
	defer fake.getClickTrackerByNameMutex.RUnlock()
	return len(fake.getClickTrackerByNameArgsForCall)
}

func (fake *FakePersistor) GetClickTrackerByNameCalls(stub func(context.Context, persistence.TransactionHandler, string) (*model.ClickTracker, error)) {
	fake.getClickTrackerByNameMutex.Lock()
	defer fake.getClickTrackerByNameMutex.Unlock()
	fake.GetClickTrackerByNameStub = stub
}

func (fake *FakePersistor) GetClickTrackerByNameArgsForCall(i int) (context.Context, persistence.TransactionHandler, string) {
	fake.getClickTrackerByNameMutex.RLock()
	defer fake.getClickTrackerByNameMutex.RUnlock()
	argsForCall := fake.getClickTrackerByNameArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakePersistor) GetClickTrackerByNameReturns(result1 *model.ClickTracker, result2 error) {
	fake.getClickTrackerByNameMutex.Lock()
	defer fake.getClickTrackerByNameMutex.Unlock()
	fake.GetClickTrackerByNameStub = nil
	fake.getClickTrackerByNameReturns = struct {
		result1 *model.ClickTracker
		result2 error
	}{result1, result2}
}

func (fake *FakePersistor) GetClickTrackerByNameReturnsOnCall(i int, result1 *model.ClickTracker, result2 error) {
	fake.getClickTrackerByNameMutex.Lock()
	defer fake.getClickTrackerByNameMutex.Unlock()
	fake.GetClickTrackerByNameStub = nil
	if fake.getClickTrackerByNameReturnsOnCall == nil {
		fake.getClickTrackerByNameReturnsOnCall = make(map[int]struct {
			result1 *model.ClickTracker
			result2 error
		})
	}
	fake.getClickTrackerByNameReturnsOnCall[i] = struct {
		result1 *model.ClickTracker
		result2 error
	}{result1, result2}
}

func (fake *FakePersistor) GetClickTrackers(arg1 context.Context, arg2 persistence.TransactionHandler, arg3 *model.ClickTrackerFilters) (*model.PaginatedClickTrackers, error) {
	fake.getClickTrackersMutex.Lock()
	ret, specificReturn := fake.getClickTrackersReturnsOnCall[len(fake.getClickTrackersArgsForCall)]
	fake.getClickTrackersArgsForCall = append(fake.getClickTrackersArgsForCall, struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 *model.ClickTrackerFilters
	}{arg1, arg2, arg3})
	stub := fake.GetClickTrackersStub
	fakeReturns := fake.getClickTrackersReturns
	fake.recordInvocation("GetClickTrackers", []interface{}{arg1, arg2, arg3})
	fake.getClickTrackersMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePersistor) GetClickTrackersCallCount() int {
	fake.getClickTrackersMutex.RLock()
	defer fake.getClickTrackersMutex.RUnlock()
	return len(fake.getClickTrackersArgsForCall)
}

func (fake *FakePersistor) GetClickTrackersCalls(stub func(context.Context, persistence.TransactionHandler, *model.ClickTrackerFilters) (*model.PaginatedClickTrackers, error)) {
	fake.getClickTrackersMutex.Lock()
	defer fake.getClickTrackersMutex.Unlock()
	fake.GetClickTrackersStub = stub
}

func (fake *FakePersistor) GetClickTrackersArgsForCall(i int) (context.Context, persistence.TransactionHandler, *model.ClickTrackerFilters) {
	fake.getClickTrackersMutex.RLock()
	defer fake.getClickTrackersMutex.RUnlock()
	argsForCall := fake.getClickTrackersArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakePersistor) GetClickTrackersReturns(result1 *model.PaginatedClickTrackers, result2 error) {
	fake.getClickTrackersMutex.Lock()
	defer fake.getClickTrackersMutex.Unlock()
	fake.GetClickTrackersStub = nil
	fake.getClickTrackersReturns = struct {
		result1 *model.PaginatedClickTrackers
		result2 error
	}{result1, result2}
}

func (fake *FakePersistor) GetClickTrackersReturnsOnCall(i int, result1 *model.PaginatedClickTrackers, result2 error) {
	fake.getClickTrackersMutex.Lock()
	defer fake.getClickTrackersMutex.Unlock()
	fake.GetClickTrackersStub = nil
	if fake.getClickTrackersReturnsOnCall == nil {
		fake.getClickTrackersReturnsOnCall = make(map[int]struct {
			result1 *model.PaginatedClickTrackers
			result2 error
		})
	}
	fake.getClickTrackersReturnsOnCall[i] = struct {
		result1 *model.PaginatedClickTrackers
		result2 error
	}{result1, result2}
}

func (fake *FakePersistor) RestoreClickTracker(arg1 context.Context, arg2 persistence.TransactionHandler, arg3 int) error {
	fake.restoreClickTrackerMutex.Lock()
	ret, specificReturn := fake.restoreClickTrackerReturnsOnCall[len(fake.restoreClickTrackerArgsForCall)]
	fake.restoreClickTrackerArgsForCall = append(fake.restoreClickTrackerArgsForCall, struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 int
	}{arg1, arg2, arg3})
	stub := fake.RestoreClickTrackerStub
	fakeReturns := fake.restoreClickTrackerReturns
	fake.recordInvocation("RestoreClickTracker", []interface{}{arg1, arg2, arg3})
	fake.restoreClickTrackerMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1
	}
	return fakeReturns.result1
}

func (fake *FakePersistor) RestoreClickTrackerCallCount() int {
	fake.restoreClickTrackerMutex.RLock()
	defer fake.restoreClickTrackerMutex.RUnlock()
	return len(fake.restoreClickTrackerArgsForCall)
}

func (fake *FakePersistor) RestoreClickTrackerCalls(stub func(context.Context, persistence.TransactionHandler, int) error) {
	fake.restoreClickTrackerMutex.Lock()
	defer fake.restoreClickTrackerMutex.Unlock()
	fake.RestoreClickTrackerStub = stub
}

func (fake *FakePersistor) RestoreClickTrackerArgsForCall(i int) (context.Context, persistence.TransactionHandler, int) {
	fake.restoreClickTrackerMutex.RLock()
	defer fake.restoreClickTrackerMutex.RUnlock()
	argsForCall := fake.restoreClickTrackerArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakePersistor) RestoreClickTrackerReturns(result1 error) {
	fake.restoreClickTrackerMutex.Lock()
	defer fake.restoreClickTrackerMutex.Unlock()
	fake.RestoreClickTrackerStub = nil
	fake.restoreClickTrackerReturns = struct {
		result1 error
	}{result1}
}

func (fake *FakePersistor) RestoreClickTrackerReturnsOnCall(i int, result1 error) {
	fake.restoreClickTrackerMutex.Lock()
	defer fake.restoreClickTrackerMutex.Unlock()
	fake.RestoreClickTrackerStub = nil
	if fake.restoreClickTrackerReturnsOnCall == nil {
		fake.restoreClickTrackerReturnsOnCall = make(map[int]struct {
			result1 error
		})
	}
	fake.restoreClickTrackerReturnsOnCall[i] = struct {
		result1 error
	}{result1}
}

func (fake *FakePersistor) UpdateClickTracker(arg1 context.Context, arg2 persistence.TransactionHandler, arg3 *model.UpdateClickTracker) (*model.ClickTracker, error) {
	fake.updateClickTrackerMutex.Lock()
	ret, specificReturn := fake.updateClickTrackerReturnsOnCall[len(fake.updateClickTrackerArgsForCall)]
	fake.updateClickTrackerArgsForCall = append(fake.updateClickTrackerArgsForCall, struct {
		arg1 context.Context
		arg2 persistence.TransactionHandler
		arg3 *model.UpdateClickTracker
	}{arg1, arg2, arg3})
	stub := fake.UpdateClickTrackerStub
	fakeReturns := fake.updateClickTrackerReturns
	fake.recordInvocation("UpdateClickTracker", []interface{}{arg1, arg2, arg3})
	fake.updateClickTrackerMutex.Unlock()
	if stub != nil {
		return stub(arg1, arg2, arg3)
	}
	if specificReturn {
		return ret.result1, ret.result2
	}
	return fakeReturns.result1, fakeReturns.result2
}

func (fake *FakePersistor) UpdateClickTrackerCallCount() int {
	fake.updateClickTrackerMutex.RLock()
	defer fake.updateClickTrackerMutex.RUnlock()
	return len(fake.updateClickTrackerArgsForCall)
}

func (fake *FakePersistor) UpdateClickTrackerCalls(stub func(context.Context, persistence.TransactionHandler, *model.UpdateClickTracker) (*model.ClickTracker, error)) {
	fake.updateClickTrackerMutex.Lock()
	defer fake.updateClickTrackerMutex.Unlock()
	fake.UpdateClickTrackerStub = stub
}

func (fake *FakePersistor) UpdateClickTrackerArgsForCall(i int) (context.Context, persistence.TransactionHandler, *model.UpdateClickTracker) {
	fake.updateClickTrackerMutex.RLock()
	defer fake.updateClickTrackerMutex.RUnlock()
	argsForCall := fake.updateClickTrackerArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2, argsForCall.arg3
}

func (fake *FakePersistor) UpdateClickTrackerReturns(result1 *model.ClickTracker, result2 error) {
	fake.updateClickTrackerMutex.Lock()
	defer fake.updateClickTrackerMutex.Unlock()
	fake.UpdateClickTrackerStub = nil
	fake.updateClickTrackerReturns = struct {
		result1 *model.ClickTracker
		result2 error
	}{result1, result2}
}

func (fake *FakePersistor) UpdateClickTrackerReturnsOnCall(i int, result1 *model.ClickTracker, result2 error) {
	fake.updateClickTrackerMutex.Lock()
	defer fake.updateClickTrackerMutex.Unlock()
	fake.UpdateClickTrackerStub = nil
	if fake.updateClickTrackerReturnsOnCall == nil {
		fake.updateClickTrackerReturnsOnCall = make(map[int]struct {
			result1 *model.ClickTracker
			result2 error
		})
	}
	fake.updateClickTrackerReturnsOnCall[i] = struct {
		result1 *model.ClickTracker
		result2 error
	}{result1, result2}
}

func (fake *FakePersistor) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.addClickTrackerMutex.RLock()
	defer fake.addClickTrackerMutex.RUnlock()
	fake.deleteClickTrackerMutex.RLock()
	defer fake.deleteClickTrackerMutex.RUnlock()
	fake.getClickTrackerByNameMutex.RLock()
	defer fake.getClickTrackerByNameMutex.RUnlock()
	fake.getClickTrackersMutex.RLock()
	defer fake.getClickTrackersMutex.RUnlock()
	fake.restoreClickTrackerMutex.RLock()
	defer fake.restoreClickTrackerMutex.RUnlock()
	fake.updateClickTrackerMutex.RLock()
	defer fake.updateClickTrackerMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *FakePersistor) recordInvocation(key string, args []interface{}) {
	fake.invocationsMutex.Lock()
	defer fake.invocationsMutex.Unlock()
	if fake.invocations == nil {
		fake.invocations = map[string][][]interface{}{}
	}
	if fake.invocations[key] == nil {
		fake.invocations[key] = [][]interface{}{}
	}
	fake.invocations[key] = append(fake.invocations[key], args)
}
