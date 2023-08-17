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

// Package updater provides a way to update Kubernetes resources.
// Inspired by controllerutil and:
// https://github.com/gardener/gardener/blob/master/docs/development/kubernetes-clients.md
package updater

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"sigs.k8s.io/controller-runtime/pkg/client"
)

type MutateFunc func() error

func CreateOrUpdate(ctx context.Context, c client.Client, key client.ObjectKey, obj client.Object, f MutateFunc) error {
	if err := c.Get(ctx, key, obj); err != nil {
		if !apierrors.IsNotFound(err) {
			return fmt.Errorf("failed to get object: %w", err)
		}

		obj.SetName(key.Name)
		obj.SetNamespace(key.Namespace)

		if f != nil {
			if err := f(); err != nil {
				return fmt.Errorf("failed to mutate object: %w", err)
			}
		}

		if err := c.Create(ctx, obj); err != nil {
			return fmt.Errorf("failed to create object: %w", err)
		}

		return nil
	}

	if f != nil {
		if err := f(); err != nil {
			return fmt.Errorf("failed to mutate object: %w", err)
		}
	}

	if err := c.Update(ctx, obj); err != nil {
		return fmt.Errorf("failed to update object: %w", err)
	}

	// Get the latest version of the object before returning.
	if err := c.Get(ctx, key, obj); err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("failed to get updated object: %w", err)
	}

	return nil
}

func UpdateStatus(ctx context.Context, c client.Client, key client.ObjectKey, obj client.Object, f MutateFunc) error {
	if err := c.Get(ctx, key, obj); err != nil {
		return fmt.Errorf("failed to get object: %w", err)
	}

	if f != nil {
		if err := f(); err != nil {
			return fmt.Errorf("failed to mutate object: %w", err)
		}
	}

	if err := c.Status().Update(ctx, obj); err != nil {
		return fmt.Errorf("failed to update object: %w", err)
	}

	// Get the latest version of the object before returning.
	if err := c.Get(ctx, key, obj); err != nil && !apierrors.IsNotFound(err) {
		return fmt.Errorf("failed to get updated object: %w", err)
	}

	return nil
}
