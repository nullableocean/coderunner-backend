package processpool

import (
	"codeproccesor/src/domain"
	"codeproccesor/src/pkg/semaphore"
	"codeproccesor/src/usecases"
	"log"
	"sync"
)

type ProcessPool struct {
	processor usecases.Processor
	sem       *semaphore.Semaphore
	wg        *sync.WaitGroup
	isStopped bool

	logger *log.Logger
}

func NewProcessPool(logger *log.Logger, processCount int, processor usecases.Processor) usecases.ProcessPool {
	return &ProcessPool{
		processor: processor,
		sem:       semaphore.NewSemaphore(processCount),
		wg:        &sync.WaitGroup{},
		isStopped: false,
		logger:    logger,
	}
}

func (p *ProcessPool) Push(t domain.Task) error {
	if p.isStopped {
		return ErrProcessPoolStopped
	}

	p.sem.Acquire()
	p.wg.Add(1)

	go p.process(t)

	return nil
}

func (p *ProcessPool) IsRunned() bool {
	return !p.isStopped
}

func (p *ProcessPool) Stop() {
	if !p.isStopped {
		p.isStopped = true
		p.wg.Wait()
	}
}

func (p *ProcessPool) process(t domain.Task) {
	defer p.wg.Done()
	defer p.sem.Release()

	err := p.processor.Process(t)

	if err != nil {
		p.handleProcessError(t, err)
	} else {
		p.logger.Printf("process completed. task %s", t.Uuid)
	}
}

func (p *ProcessPool) handleProcessError(t domain.Task, err error) {
	p.logger.Printf("[PROCESS ERROR] task_uuid: %s, error: %s", t.Uuid, err)
}
