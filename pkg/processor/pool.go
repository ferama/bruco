package processor

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
)

type EventCallback func(event *Response)

// Pool ...
type Pool struct {
	pythonMap   map[string]*Python
	wrapperPath string
	handlerPath string
	moduleName  string
	env         []EnvVar

	availableWorkers chan *Python
}

// GetPoolInstance ...
func NewPool(cfg *ProcessorConf, workingDir string) *Pool {
	data, _ := pythonWrapper.ReadFile("wrapper/wrapper.py")
	file, err := ioutil.TempFile(os.TempDir(), "bruco-python-wrapper-")
	if err != nil {
		log.Fatal(err)
	}
	if _, err = file.Write(data); err != nil {
		log.Fatal("[PROCESSOR] failed to write to temporary file", err)
	}
	if err := file.Close(); err != nil {
		log.Fatal(err)
	}
	log.Printf("[PROCESSOR] unpacked wrapper to %s", file.Name())

	pool := &Pool{
		pythonMap:        make(map[string]*Python),
		availableWorkers: make(chan *Python, cfg.Workers),
		wrapperPath:      file.Name(),
		handlerPath:      cfg.HandlerPath,
		env:              cfg.Env,
	}
	pool.moduleName = pool.resolveModuleName(cfg.ModuleName)
	for i := 0; i < cfg.Workers; i++ {
		name := fmt.Sprintf("worker-%d", i)
		pool.createPythonInstance(name, workingDir)

	}
	log.Printf("[PROCESSOR] allocated %d workers", cfg.Workers)
	return pool
}

func (p *Pool) resolveModuleName(name string) string {
	if name == "" {
		return "handler"
	}
	return name
}

func (p *Pool) HandleEventAsync(data []byte, callback EventCallback) {
	if len(data) == 0 {
		// log.Println("[PROCESSOR] WARNING: got 0 len event")
		res := Response{
			Data:  "",
			Error: "can't handle zero len event",
		}
		if callback != nil {
			callback(&res)
		}
		return
	}
	python := <-p.availableWorkers
	python.handleEvent(data)
	go func() {
		response := <-python.eventResponse
		if callback != nil {
			callback(&response)
		}
	}()
}

func (p *Pool) HandleEvent(data []byte) (*Response, error) {
	if len(data) == 0 {
		return nil, fmt.Errorf("can't handle zero len event")
	}
	python := <-p.availableWorkers
	err := python.handleEvent(data)
	if err != nil {
		return nil, err
	}
	response := <-python.eventResponse
	return &response, nil
}

// createPythonInstance ...
func (p *Pool) createPythonInstance(name string, workingDir string) *Python {
	python := NewPython(name,
		p.availableWorkers,
		p.wrapperPath,
		p.handlerPath,
		p.moduleName,
		p.env,
		workingDir,
	)
	p.pythonMap[name] = python

	return python
}

// Destroy ...
func (p *Pool) Destroy() {
	os.Remove(p.wrapperPath)
	for key, python := range p.pythonMap {
		python.kill()
		python = nil
		delete(p.pythonMap, key)
	}
}
