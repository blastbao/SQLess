/*
 * Copyright 2018 The CovenantSQL Authors.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package types

import "github.com/SQLess/SQLess/proto"

// ApplyRequest defines the apply request entity.
type ApplyRequest struct {
	proto.Envelope
	Instance string
	Log      *Log
}

// FetchRequest defines the fetch request entity.
type FetchRequest struct {
	proto.Envelope
	Instance string
	Index    uint64
}

// FetchResponse defines the fetch response entity.
type FetchResponse struct {
	proto.Envelope
	Instance string
	Log      *Log
}
