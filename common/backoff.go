// Copyright 2017. box.la authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package common

import (
	"math/rand"
	"time"
)

var DefaultBackoff = Backoff{
	MaxDelay:  120 * time.Second,
	baseDelay: 1.0 * time.Second,
	factor:    1.6,
	jitter:    0.2,
}

type Backoff struct {
	MaxDelay  time.Duration
	baseDelay time.Duration
	factor    float64
	jitter    float64
}

func (bc Backoff) Duration(retries int) time.Duration {
	if retries == 0 {
		return bc.baseDelay
	}

	backoff, max := float64(bc.baseDelay), float64(bc.MaxDelay)
	for backoff < max && retries > 0 {
		backoff *= bc.factor
		retries--
	}

	if backoff > max {
		backoff = max
	}

	backoff *= 1 + bc.jitter*(rand.Float64()*2-1)
	if backoff < 0 {
		return 0
	}

	return time.Duration(backoff)
}
