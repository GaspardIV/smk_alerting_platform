package pkg

import "sync"

type SharedData struct {
        Wg            sync.WaitGroup
        mutex         sync.Mutex
        ErrorMessages string
}

func (s *SharedData) AppendErrorMessage(message string) {
        s.mutex.Lock()
        s.ErrorMessages += message
        s.ErrorMessages += "\n"
        s.mutex.Unlock()
}