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
package main


import (
	"github.com/spf13/cobra"

	"github.com/pingcap/go-ycsb/pkg/prop"
)

func newWorkerCommand() *cobra.Command {
	m := &cobra.Command{
		Use:   "shell db",
		Short: "YCSB Command Line Client",
		Args:  cobra.MinimumNArgs(1),
		Run:   runWorkerCommmandFunc,
	}
	m.Flags().StringSliceVarP(&propertyFiles, "property_file", "P", nil, "Spefify a property file")
	m.Flags().StringSliceVarP(&propertyValues, "prop", "p", nil, "Specify a property value with name=value")
	m.Flags().StringVar(&tableName, "table", "", "Use the table name instead of the default \""+prop.TableNameDefault+"\"")
	return m
}

func runWorkerCommmandFunc(cmd *cobra.Command, args []string) {
	dbName := args[0]
	initialGlobal(dbName, nil)

	shellContext = globalWorkload.InitThread(globalContext, 0, 1)
	shellContext = globalDB.InitThread(shellContext, 0, 1)

	shellLoop()

	globalDB.CleanupThread(shellContext)
	globalWorkload.CleanupThread(shellContext)
}
