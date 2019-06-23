// Copyright 2018 PingCAP, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.
package worker

import (
	"encoding/json"
	"fmt"
	"net/http"
	"path"
	"sync"
)

type ErrorWithCode struct {
	Err  error
	Code int
}

type Server struct {
	storeDir         string
	jobExecutorToken chan struct{}
	jobs             map[string]*Job
	jobsMu           sync.Mutex
}

func NewServer(storeDir string) *Server {
	return &Server{
		storeDir:         storeDir,
		jobExecutorToken: make(chan struct{}, 0),
		jobs:             make(map[string]*Job),
	}
}

func (s *Server) Start(addr string) error {
	mux := http.NewServeMux()
	mux.HandleFunc("/job/list", func(w http.ResponseWriter, r *http.Request) {

	})

	mux.HandleFunc("/job/status/", func(w http.ResponseWriter, r *http.Request) {
		session, err := s.getJob(r)
		if err != nil {
			s.renderError(w, err.Err, err.Code)
			return
		}
		if err := json.NewEncoder(w).Encode(session); err != nil {
			s.renderError(w, err, http.StatusInternalServerError)
			return
		}
	})

	mux.HandleFunc("/job/run/", func(w http.ResponseWriter, r *http.Request) {
		session, err := s.createJob(r)
		if err != nil {
			s.renderError(w, err.Err, err.Code)
			return
		}

		if err := s.acquireJobLock(); err != nil {
			s.renderError(w, err, http.StatusLocked)
			return
		}
		defer s.releaseJobLock()

		if err := session.Run(); err != nil {
			s.renderError(w, err.Err, err.Code)
		}
	})

	mux.HandleFunc("/job/download/", func(w http.ResponseWriter, r *http.Request) {
		session, err := s.getJob(r)
		if err != nil {
			s.renderError(w, err.Err, err.Code)
			return
		}
		if session.State != StateFinished {
			s.renderError(w, fmt.Errorf("job %s is not ready for start (current: %s)", session.ID, session.State), http.StatusBadRequest)
			return
		}
		key := r.URL.Query().Get("key")
		http.ServeFile(w, r, path.Join(s.storeDir, session.ID, key))
	})
	return http.ListenAndServe(addr, mux)
}

func (s *Server) acquireJobLock() error {
	select {
	case <-s.jobExecutorToken:
		return nil
	default:
		return fmt.Errorf("job in progress")
	}
}

func (s *Server) releaseJobLock() {
	s.jobExecutorToken <- struct{}{}
}

func (s *Server) createJob(r *http.Request) (*Job, *ErrorWithCode) {
	_, sessionId := path.Split(r.URL.Path)
	if len(sessionId) == 0 {
		return nil, &ErrorWithCode{
			Err:  fmt.Errorf("sessionId not specified"),
			Code: http.StatusBadRequest,
		}
	}
	s.jobsMu.Lock()
	defer s.jobsMu.Unlock()
	if _, ok := s.jobs[sessionId]; ok {
		return nil, &ErrorWithCode{
			Err:  fmt.Errorf("sessionId %s already exists", sessionId),
			Code: http.StatusBadRequest,
		}
	}
	sess, err := CreateJobFromRequest(sessionId, r)
	if err != nil {
		return nil, err
	}
	s.jobs[sessionId] = sess
	return sess, nil
}

func (s *Server) getJob(r *http.Request) (*Job, *ErrorWithCode) {
	_, sessionId := path.Split(r.URL.Path)
	if len(sessionId) == 0 {
		return nil, &ErrorWithCode{
			Err:  fmt.Errorf("sessionId not specified"),
			Code: http.StatusBadRequest,
		}
	}
	s.jobsMu.Lock()
	defer s.jobsMu.Unlock()
	if sess, ok := s.jobs[sessionId]; ok {
		return sess, nil
	}
	sess, err := OpenJob(sessionId)
	if err != nil {
		return nil, err
	}
	s.jobs[sessionId] = sess
	return sess, nil
}

func (s *Server) renderError(w http.ResponseWriter, err error, code int) {
	json.NewEncoder(w).Encode(errorResponse{Message: err.Error()})
	w.WriteHeader(code)
}

type errorResponse struct {
	Message string `json:"message"`
}
