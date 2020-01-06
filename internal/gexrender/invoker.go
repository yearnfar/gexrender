package gexrender

import (
	log "github.com/sirupsen/logrus"
)

// Invoker 调度程序
type Invoker interface {
	Invoke(job *Job, data []byte) error
}

var invokers = make(map[string]Invoker)

// Register 注册invoker
func Register(name string, invoker Invoker) {
	_, ok := invokers[name]
	if ok {
		log.Fatalf("已注册调度程序%s", name)
	}

	invokers[name] = invoker
}

// Invoke 调用invoke
func Invoke(name string, job *Job, data []byte) error {
	invoker, ok := invokers[name]
	if !ok {
		log.Fatalf("调度程序%s不存在", name)
	}

	log.Infof("[%s] invoke %s, parameter %s", job.Uid, name, string(data))
	return invoker.Invoke(job, data)
}
