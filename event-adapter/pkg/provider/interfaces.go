/*
Copyright 2017 The Kubernetes Authors.
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

package provider

import (
	"k8s.io/apimachinery/pkg/runtime/schema"
	"github.com/GoogleCloudPlatform/k8s-stackdriver/event-adapter/pkg/types"
)

// Info relevant for an Event
type EventInfo struct {
	GroupResource          schema.GroupResource
	Namespaced             bool
	Event                  string
}

// Interfaces that contains the methods that will provide info for the given events
type EventsProvider interface {
	GetNamespacedEventsByName( namespace, eventName string) (*types.EventValue, error)
}