// Code generated by counterfeiter. DO NOT EDIT.
package fakes

import (
	"sync"

	"github.com/pivotal-cf/kiln/commands"
	"github.com/pivotal-cf/kiln/release"
)

type ReleaseDownloader struct {
	DownloadReleaseStub        func(string, release.ReleaseRequirement) (release.LocalRelease, string, string, error)
	downloadReleaseMutex       sync.RWMutex
	downloadReleaseArgsForCall []struct {
		arg1 string
		arg2 release.ReleaseRequirement
	}
	downloadReleaseReturns struct {
		result1 release.LocalRelease
		result2 string
		result3 string
		result4 error
	}
	downloadReleaseReturnsOnCall map[int]struct {
		result1 release.LocalRelease
		result2 string
		result3 string
		result4 error
	}
	invocations      map[string][][]interface{}
	invocationsMutex sync.RWMutex
}

func (fake *ReleaseDownloader) DownloadRelease(arg1 string, arg2 release.ReleaseRequirement) (release.LocalRelease, string, string, error) {
	fake.downloadReleaseMutex.Lock()
	ret, specificReturn := fake.downloadReleaseReturnsOnCall[len(fake.downloadReleaseArgsForCall)]
	fake.downloadReleaseArgsForCall = append(fake.downloadReleaseArgsForCall, struct {
		arg1 string
		arg2 release.ReleaseRequirement
	}{arg1, arg2})
	fake.recordInvocation("DownloadRelease", []interface{}{arg1, arg2})
	fake.downloadReleaseMutex.Unlock()
	if fake.DownloadReleaseStub != nil {
		return fake.DownloadReleaseStub(arg1, arg2)
	}
	if specificReturn {
		return ret.result1, ret.result2, ret.result3, ret.result4
	}
	fakeReturns := fake.downloadReleaseReturns
	return fakeReturns.result1, fakeReturns.result2, fakeReturns.result3, fakeReturns.result4
}

func (fake *ReleaseDownloader) DownloadReleaseCallCount() int {
	fake.downloadReleaseMutex.RLock()
	defer fake.downloadReleaseMutex.RUnlock()
	return len(fake.downloadReleaseArgsForCall)
}

func (fake *ReleaseDownloader) DownloadReleaseCalls(stub func(string, release.ReleaseRequirement) (release.LocalRelease, string, string, error)) {
	fake.downloadReleaseMutex.Lock()
	defer fake.downloadReleaseMutex.Unlock()
	fake.DownloadReleaseStub = stub
}

func (fake *ReleaseDownloader) DownloadReleaseArgsForCall(i int) (string, release.ReleaseRequirement) {
	fake.downloadReleaseMutex.RLock()
	defer fake.downloadReleaseMutex.RUnlock()
	argsForCall := fake.downloadReleaseArgsForCall[i]
	return argsForCall.arg1, argsForCall.arg2
}

func (fake *ReleaseDownloader) DownloadReleaseReturns(result1 release.LocalRelease, result2 string, result3 string, result4 error) {
	fake.downloadReleaseMutex.Lock()
	defer fake.downloadReleaseMutex.Unlock()
	fake.DownloadReleaseStub = nil
	fake.downloadReleaseReturns = struct {
		result1 release.LocalRelease
		result2 string
		result3 string
		result4 error
	}{result1, result2, result3, result4}
}

func (fake *ReleaseDownloader) DownloadReleaseReturnsOnCall(i int, result1 release.LocalRelease, result2 string, result3 string, result4 error) {
	fake.downloadReleaseMutex.Lock()
	defer fake.downloadReleaseMutex.Unlock()
	fake.DownloadReleaseStub = nil
	if fake.downloadReleaseReturnsOnCall == nil {
		fake.downloadReleaseReturnsOnCall = make(map[int]struct {
			result1 release.LocalRelease
			result2 string
			result3 string
			result4 error
		})
	}
	fake.downloadReleaseReturnsOnCall[i] = struct {
		result1 release.LocalRelease
		result2 string
		result3 string
		result4 error
	}{result1, result2, result3, result4}
}

func (fake *ReleaseDownloader) Invocations() map[string][][]interface{} {
	fake.invocationsMutex.RLock()
	defer fake.invocationsMutex.RUnlock()
	fake.downloadReleaseMutex.RLock()
	defer fake.downloadReleaseMutex.RUnlock()
	copiedInvocations := map[string][][]interface{}{}
	for key, value := range fake.invocations {
		copiedInvocations[key] = value
	}
	return copiedInvocations
}

func (fake *ReleaseDownloader) recordInvocation(key string, args []interface{}) {
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

var _ commands.ReleaseDownloader = new(ReleaseDownloader)
