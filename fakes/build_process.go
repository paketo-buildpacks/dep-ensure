package fakes

import "sync"

type BuildProcess struct {
	ExecuteCall struct {
		sync.Mutex
		CallCount int
		Receives  struct {
			Workspace  string
			GoPath     string
			Gocachedir string
		}
		Returns struct {
			Err error
		}
		Stub func(string, string, string) error
	}
}

func (f *BuildProcess) Execute(param1 string, param2 string, param3 string) error {
	f.ExecuteCall.Lock()
	defer f.ExecuteCall.Unlock()
	f.ExecuteCall.CallCount++
	f.ExecuteCall.Receives.Workspace = param1
	f.ExecuteCall.Receives.GoPath = param2
	f.ExecuteCall.Receives.Gocachedir = param3
	if f.ExecuteCall.Stub != nil {
		return f.ExecuteCall.Stub(param1, param2, param3)
	}
	return f.ExecuteCall.Returns.Err
}
