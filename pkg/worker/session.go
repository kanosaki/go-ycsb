package worker

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/magiconair/properties"

	"github.com/pingcap/go-ycsb/pkg/client"
	"github.com/pingcap/go-ycsb/pkg/ycsb"
)

type Job struct {
	ID           string       `json:"id"`
	StartedAt    time.Time    `json:"started_at"`
	FinishedAt   time.Time    `json:"finished_at"`
	State        SessionState `json:"state"`
	Workload     string       `json:"workload"`
	DatabaseName string       `json:"database"`

	props      *properties.Properties `json:"props"`
	dir        string
	ycsbClient *client.Client
	mu         sync.Mutex
}

func CreateJobFromRequest(sessionId string, r *http.Request) (*Job, *ErrorWithCode) {
	queries := r.URL.Query()
	fp, fh, err := r.FormFile("properties")
	if err != nil {
		if err == http.ErrMissingFile {
			return nil, &ErrorWithCode{
				Err:  fmt.Errorf("form file 'properties' is required"),
				Code: http.StatusBadRequest,
			}
		} else {
			return nil, &ErrorWithCode{
				Err:  err,
				Code: http.StatusInternalServerError,
			}
		}
	}
	buf := bytes.NewBuffer(make([]byte, 0, fh.Size))
	if _, err := io.Copy(buf, fp); err != nil {
		return nil, &ErrorWithCode{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	props, err := properties.Load(buf.Bytes(), properties.UTF8)
	if err != nil {
		return nil, &ErrorWithCode{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	dbName := queries.Get("database")
	if len(dbName) == 0 {
		return nil, &ErrorWithCode{
			Err:  fmt.Errorf("query parameter 'database' is required"),
			Code: http.StatusBadRequest,
		}
	}
	workload := queries.Get("workload")
	if len(workload) == 0 {
		return nil, &ErrorWithCode{
			Err:  fmt.Errorf("query parameter 'workload' is required"),
			Code: http.StatusBadRequest,
		}
	}
	sess, err := CreateJob(sessionId, props, workload, dbName)
	if err != nil {
		return nil, &ErrorWithCode{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return sess, nil
}

func CreateJob(sessionId string, props *properties.Properties, workload, dbName string) (*Job, error) {
	sess := &Job{
		ID:           sessionId,
		props:        props,
		Workload:     workload,
		DatabaseName: dbName,
	}
	if err := sess.init(); err != nil {
		return nil, err
	}
	return sess, nil
}

func OpenJob(sessionId, dir string) (*Job, *ErrorWithCode) {
	sess := &Job{
		ID:  sessionId,
		dir: dir,
	}
	if err := sess.load(); err != nil {
		if err == os.ErrNotExist {
			return nil, &ErrorWithCode{
				Err:  err,
				Code: http.StatusNotFound,
			}
		} else {
			return nil, &ErrorWithCode{
				Err:  err,
				Code: http.StatusInternalServerError,
			}
		}
	}
	if err := sess.init(); err != nil {
		return nil, &ErrorWithCode{
			Err:  err,
			Code: http.StatusInternalServerError,
		}
	}
	return sess, nil
}

func (s *Job) Save() error {
	s.mu.Lock()
	defer s.mu.Unlock()
}

func (s *Job) load() error {
	s.mu.Lock()
	defer s.mu.Unlock()
}

func (s *Job) init() error {
	s.mu.Lock()
	defer s.mu.Unlock()
	wc := ycsb.GetWorkloadCreator(s.Workload)
	if wc == nil {
		return fmt.Errorf("no such workload: %s", s.Workload)
	}
	w, err := wc.Create(s.props)
	if err != nil {
		return fmt.Errorf("failed to initialize workload: %v", err)
	}
	dbc := ycsb.GetDBCreator(s.DatabaseName)
	db, err := dbc.Create(s.props)
	if err != nil {
		return fmt.Errorf("failed to initialize database: %v", err)
	}
	s.ycsbClient = client.NewClient(s.props, w, db)
	return nil
}

func (s *Job) Run(ctx context.Context) *ErrorWithCode {
	if s.State != StateInitialized {
		return &ErrorWithCode{
			Err:  fmt.Errorf("job %s is not ready for prepare (current: %s)", s.ID, s.State),
			Code: http.StatusBadRequest,
		}
	}
	s.ycsbClient.Run(ctx)
	return nil
}

func (s *Job) Download(key string, r *http.Request, w http.ResponseWriter) {
	http.ServeFile(w, r, key)
}
