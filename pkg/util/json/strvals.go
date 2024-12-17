/*
Copyright 2021 The KodeRover Authors.

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

package json

import (
	"github.com/segmentio/encoding/json"

	"helm.sh/helm/v3/pkg/strvals"
)

// ToJSON takes a string of arguments(in this format: a=b,c.d=e) and converts to a JSON document.
func ToJSON(s string) ([]byte, error) {
	m, err := strvals.Parse(s)
	if err != nil {
		return nil, err
	}

	return json.Marshal(m)
}
