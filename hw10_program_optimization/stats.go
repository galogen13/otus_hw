package hw10_program_optimization //nolint:golint,stylecheck

import (
	"bufio"
	"errors"
	"io"
	"strings"
	"sync"
)

// User .
//easyjson:json
type User struct {
	Email string
}

// DomainStat .
type DomainStat map[string]int

var ErrEmptyDomain = errors.New("empty domain")

var mutex = &sync.Mutex{}

var result *DomainStat

var linePool = sync.Pool{
	New: func() interface{} { return []byte{} },
}

// GetDomainStat .
func GetDomainStat(r io.Reader, domain string) (DomainStat, error) {
	result = &DomainStat{}

	if domain == "" {
		return nil, ErrEmptyDomain
	}

	reader := bufio.NewReader(r)

	waitCh := make(chan struct{})
	errCh := make(chan error)
	doneCh := make(chan struct{})

	go func() {
		wg := &sync.WaitGroup{}

		for {
			line := linePool.Get().([]byte)
			line, err := reader.ReadBytes('\n')
			if err == io.EOF {
				if len(line) == 0 {
					break
				}
				wg.Add(1)
				go getDomainStatInLine(line, domain, errCh, doneCh, wg)
				break
			}
			if err != nil {
				errCh <- err
				return
			}

			wg.Add(1)
			go getDomainStatInLine(line, domain, errCh, doneCh, wg)
		}
		wg.Wait()
		close(waitCh)
	}()

	select {
	case <-waitCh:
		close(errCh)
		close(doneCh)
		break
	case err := <-errCh:
		close(doneCh)
		return nil, err
	}
	return *result, nil
}

func getDomainStatInLine(line []byte, domain string, errCh chan error, done <-chan struct{}, wg *sync.WaitGroup) {
	defer wg.Done()

	select {
	case <-done:
		return
	default:
	}

	user := User{}
	err := user.UnmarshalJSON(line)
	if err != nil {
		select {
		case <-done:
		default:
			errCh <- err
		}
		return
	}

	if user.Email == "" {
		return
	}

	if contain := strings.Contains(user.Email, "."+domain); contain {
		domainName := strings.ToLower(strings.SplitN(user.Email, "@", 2)[1])
		mutex.Lock()
		(*result)[domainName]++
		mutex.Unlock()
	}
}
