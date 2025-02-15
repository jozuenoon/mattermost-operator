// Copyright © 2019 Banzai Cloud
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package objectmatch

import (
	"encoding/json"

	"github.com/goph/emperror"
	appsv1 "k8s.io/api/apps/v1"
	v1 "k8s.io/kubernetes/pkg/apis/apps/v1"
)

type daemonSetMatcher struct {
	objectMatcher ObjectMatcher
}

func NewDaemonSetMatcher(objectMatcher ObjectMatcher) *daemonSetMatcher {
	return &daemonSetMatcher{
		objectMatcher: objectMatcher,
	}
}

// Match compares two appsv1.ClusterDaemonSet objects
func (m daemonSetMatcher) Match(oldOrig, newOrig *appsv1.DaemonSet) (bool, error) {
	old := oldOrig.DeepCopy()
	new := newOrig.DeepCopy()

	v1.SetObjectDefaults_DaemonSet(new)

	type DaemonSet struct {
		ObjectMeta
		Spec appsv1.DaemonSetSpec
	}

	oldData, err := json.Marshal(DaemonSet{
		ObjectMeta: m.objectMatcher.GetObjectMeta(old.ObjectMeta),
		Spec:       old.Spec,
	})
	if err != nil {
		return false, emperror.WrapWith(err, "could not marshal old object", "name", old.Name)
	}
	newObject := DaemonSet{
		ObjectMeta: m.objectMatcher.GetObjectMeta(new.ObjectMeta),
		Spec:       new.Spec,
	}
	newData, err := json.Marshal(newObject)
	if err != nil {
		return false, emperror.WrapWith(err, "could not marshal new object", "name", new.Name)
	}

	matched, err := m.objectMatcher.MatchJSON(oldData, newData, newObject)
	if err != nil {
		return false, emperror.WrapWith(err, "could not match objects", "name", new.Name)
	}

	return matched, nil
}
