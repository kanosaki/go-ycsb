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
package epsilon

import (
	"context"

	"github.com/magiconair/properties"

	"github.com/pingcap/go-ycsb/pkg/ycsb"
)

type epsilonCreator struct {
}

type epsilonDB struct {
}

func (epsilonCreator) Create(p *properties.Properties) (ycsb.DB, error) {
	return &epsilonDB{}, nil
}

func (*epsilonDB) Close() error {
	return nil
}

func (*epsilonDB) InitThread(ctx context.Context, threadID int, threadCount int) context.Context {
	return ctx
}

func (*epsilonDB) CleanupThread(ctx context.Context) {
}

func (*epsilonDB) Read(ctx context.Context, table string, key string, fields []string) (map[string][]byte, error) {
	return nil, nil
}

func (*epsilonDB) Scan(ctx context.Context, table string, startKey string, count int, fields []string) ([]map[string][]byte, error) {
	return nil, nil
}

func (*epsilonDB) Update(ctx context.Context, table string, key string, values map[string][]byte) error {
	return nil
}

func (*epsilonDB) Insert(ctx context.Context, table string, key string, values map[string][]byte) error {
	return nil
}

func (*epsilonDB) Delete(ctx context.Context, table string, key string) error {
	return nil
}

func init() {
	ycsb.RegisterDBCreator("epsilon", epsilonCreator{})
}
