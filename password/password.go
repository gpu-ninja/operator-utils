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

package password

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789!@#$%^&*()-_=+<>?[]{}"

// GeneratePassword generates a secure random password of length n.
func GeneratePassword(n int) (string, error) {
	if n <= 0 {
		return "", fmt.Errorf("invalid password length: %d", n)
	}

	passwordBytes := make([]byte, n)
	for i := 0; i < n; i++ {
		index, err := randInt(0, len(charset)-1)
		if err != nil {
			return "", fmt.Errorf("failed to generate random integer: %w", err)
		}

		passwordBytes[i] = charset[index]
	}

	return string(passwordBytes), nil
}

// randInt generates a random integer between min and max (inclusive).
func randInt(min int, max int) (int, error) {
	diff := max - min + 1
	bigInt, err := rand.Int(rand.Reader, big.NewInt(int64(diff)))
	if err != nil {
		return 0, err
	}

	return int(bigInt.Int64()) + min, nil
}
