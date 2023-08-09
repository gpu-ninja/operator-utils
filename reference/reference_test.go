/* SPDX-License-Identifier: Apache-2.0
 *
 * Copyright 2023 Damian Peckett <damian@pecke.tt>.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package reference_test

import (
	"context"
	"testing"

	"github.com/gpu-ninja/operator-utils/reference"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"sigs.k8s.io/controller-runtime/pkg/client/fake"
)

func TestResolveReference(t *testing.T) {
	clientScheme := runtime.NewScheme()
	clientScheme.AddKnownTypes(testGV, &MyObject{})
	_ = corev1.AddToScheme(clientScheme)

	reader := fake.NewClientBuilder().WithScheme(clientScheme).WithObjects(&MyObject{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "first",
			Namespace: "default",
		},
	}, &corev1.Secret{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "demo",
			Namespace: "default",
		},
		Data: map[string][]byte{
			"secret": []byte("change-me"),
		},
	}).Build()

	ctx := context.Background()

	// Intentionally don't register the secret type.
	scheme := runtime.NewScheme()
	scheme.AddKnownTypes(testGV, &MyObject{})

	parent := &MyObject{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "second",
			Namespace: "default",
		},
	}

	t.Run("Registered Type", func(t *testing.T) {
		ref := reference.ObjectReference{
			Name:       "first",
			APIVersion: "example.com/v1",
			Kind:       "MyObject",
		}

		obj, err := ref.Resolve(ctx, reader, scheme, parent)
		require.NoError(t, err)

		assert.IsType(t, &MyObject{}, obj)
	})

	t.Run("Unregistered Type", func(t *testing.T) {
		ref := reference.ObjectReference{
			Name:       "demo",
			APIVersion: "v1",
			Kind:       "Secret",
		}

		obj, err := ref.Resolve(ctx, reader, scheme, parent)
		require.NoError(t, err)

		assert.IsType(t, &unstructured.Unstructured{}, obj)
	})

	t.Run("Same APIVersion", func(t *testing.T) {
		ref := reference.ObjectReference{
			Name: "first",
			Kind: "MyObject",
		}

		obj, err := ref.Resolve(ctx, reader, scheme, parent)
		require.NoError(t, err)

		assert.IsType(t, &MyObject{}, obj)
	})
}

var testGV = schema.GroupVersion{
	Group:   "example.com",
	Version: "v1",
}

type MyObject struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata"`
}

func (in *MyObject) DeepCopyObject() runtime.Object {
	out := MyObject{}
	in.DeepCopyInto(&out)

	return &out
}

func (in *MyObject) DeepCopyInto(out *MyObject) {
	out.TypeMeta = in.TypeMeta
	out.ObjectMeta = in.ObjectMeta
}
