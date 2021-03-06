/*
 * Copyright 2019 The CovenantSQL Authors.
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

package naconn

import "github.com/SQLess/SQLess/proto"

// Resolver defines the node ID resolver interface for node-oriented connection.
type Resolver interface {
	Resolve(id *proto.RawNodeID) (string, error)
	ResolveEx(id *proto.RawNodeID) (*proto.Node, error)
}

var (
	defaultResolver Resolver
)

// RegisterResolver registers the default resolver.
func RegisterResolver(resolver Resolver) {
	defaultResolver = resolver
}
