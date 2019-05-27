package match

import (
	"time"

	"github.com/jmartin82/mmock/pkg/mock"
	"github.com/twinj/uuid"
)

//Error contains the tested uri and the match error
type Error struct {
	URI    string `json:"uri"`
	Reason string `json:"reason"`
}

//Result contains the match result and the failing matches with different mocks and the reason or the fail.
type Result struct {
	Found  bool    `json:"match"`
	URI    string  `json:"uri"`
	Errors []Error `json:"errors"`
}

//Log contains the whole information about the request match. The http request, the final response received and the matching result.
type Log struct {
	ID       string         `json:"id"`
	Time     int64          `json:"time"`
	Request  *mock.Request  `json:"request"`
	Response *mock.Response `json:"response"`
	Result   *Result        `json:"result"`
}

func NewLog(request *mock.Request, response *mock.Response, result *Result) *Log {
	t := time.Now().Unix()
	uuid := uuid.NewV4().String()
	return &Log{ID: uuid, Time: t, Request: request, Response: response, Result: result}
}

type Spier interface {
	Find(mock.Request) []Log
	GetMatched() []Log
	GetUnMatched() []Log
	Storer
}

type Spy struct {
	store   Storer
	checker Matcher
}

func NewSpy(checker Matcher, matchStore Storer) *Spy {
	return &Spy{store: matchStore, checker: checker}
}

func (mc Spy) Find(r mock.Request) []Log {
	matches := mc.store.GetAll()
	result := []Log{}
	for _, match := range matches {
		if m, _ := mc.checker.Match(match.Request, &mock.Definition{Request: r}, false); m {
			result = append(result, match)
		}
	}
	return result

}

// ResetMatch ...
func (mc Spy) ResetMatch(r mock.Request) {
	mc.store.ResetMatch(r)
}

func (mc Spy) Save(match Log) {
	mc.store.Save(match)
}
func (mc Spy) Reset() {
	mc.store.Reset()
}

func (mc Spy) GetAll() []Log {
	return mc.store.GetAll()
}

func (mc Spy) Get(limit uint, offset uint) []Log {
	return mc.store.Get(limit, offset)
}

func (mc Spy) GetMatched() []Log {
	return mc.getMatchByResult(true)
}

func (mc Spy) GetUnMatched() []Log {
	return mc.getMatchByResult(false)
}

func (mc Spy) getMatchByResult(found bool) []Log {
	matches := mc.store.GetAll()
	result := []Log{}
	for _, match := range matches {
		if match.Result.Found == found {
			result = append(result, match)
		}
	}
	return result
}
