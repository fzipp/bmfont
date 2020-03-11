// Copyright 2020 Frederik Zipp. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package bmfont

type errorList []error

func (list errorList) Err() error {
	if len(list) == 0 {
		return nil
	}
	return list
}

func (list errorList) Error() string {
	if len(list) == 0 {
		return "no errors"
	}
	return list[0].Error()
}
