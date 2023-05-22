/*
Copyright 2023.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package version contains the operator's version, as well as versions of underlying components.
package version

import (
	"fmt"
	"runtime"
)

var (
	version                   string
	buildDate                 string
	autoInstrumentationJava   string
	autoInstrumentationNodeJS string
	autoInstrumentationPython string
	autoInstrumentationDotNet string
)

// Version holds this Operator's version as well as the version of some of the components it uses.
type Version struct {
	Operator                  string `json:"newrelic-agent-operator"`
	BuildDate                 string `json:"build-date"`
	Go                        string `json:"go-version"`
	AutoInstrumentationJava   string `json:"newrelic-instrumentation-java"`
	AutoInstrumentationNodeJS string `json:"newrelic-instrumentation-nodejs"`
	AutoInstrumentationPython string `json:"newrelic-instrumentation-python"`
	AutoInstrumentationDotNet string `json:"newrelic-instrumentation-dotnet"`
}

// Get returns the Version object with the relevant information.
func Get() Version {
	return Version{
		Operator:                  version,
		BuildDate:                 buildDate,
		Go:                        runtime.Version(),
		AutoInstrumentationJava:   AutoInstrumentationJava(),
		AutoInstrumentationNodeJS: AutoInstrumentationNodeJS(),
		AutoInstrumentationPython: AutoInstrumentationPython(),
		AutoInstrumentationDotNet: AutoInstrumentationDotNet(),
	}
}

func (v Version) String() string {
	return fmt.Sprintf(
		"Version(Operator='%v', BuildDate='%v', Go='%v', AutoInstrumentationJava='%v', AutoInstrumentationNodeJS='%v', AutoInstrumentationPython='%v', AutoInstrumentationDotNet='%v')",
		v.Operator,
		v.BuildDate,
		v.Go,
		v.AutoInstrumentationJava,
		v.AutoInstrumentationNodeJS,
		v.AutoInstrumentationPython,
		v.AutoInstrumentationDotNet,
	)
}

func AutoInstrumentationJava() string {
	if len(autoInstrumentationJava) > 0 {
		return autoInstrumentationJava
	}
	return "0.0.0"
}

func AutoInstrumentationNodeJS() string {
	if len(autoInstrumentationNodeJS) > 0 {
		return autoInstrumentationNodeJS
	}
	return "0.0.0"
}

func AutoInstrumentationPython() string {
	if len(autoInstrumentationPython) > 0 {
		return autoInstrumentationPython
	}
	return "0.0.0"
}

func AutoInstrumentationDotNet() string {
	if len(autoInstrumentationDotNet) > 0 {
		return autoInstrumentationDotNet
	}
	return "0.0.0"
}
