// Copyright 2015 The tgbot Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package utils

import (
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"time"
)

const alnum = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Download downloads the given URL to the directory dir in a file with a random
// name and the extension ext and returns the path of the created file.
// If ext is empty string the file will be created with the same extension of
// the original file at the given url.
// If dir is the empty string, download uses the default directory for temporary
// files (see os.TempDir).
func Download(dir, ext, targetURL string) (filePath string, err error) {
	res, err := http.Get(targetURL)
	if err != nil {
		return "", err
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return "", fmt.Errorf("HTTP error: %v (%v)", res.Status, res.StatusCode)
	}

	if ext == "" {
		// Parse URL to get its extension
		u, err := url.Parse(targetURL)
		if err != nil {
			return "", err
		}
		ext = path.Ext(u.Path)
	}

	f, err := TempFile(dir, "", ext)
	if err != nil {
		return "", nil
	}
	defer f.Close()

	_, err = io.Copy(f, res.Body)
	if err != nil {
		return "", err
	}

	return f.Name(), nil
}

// TempFile creates a new temporary file in the directory dir with a name
// beginning with prefix and ending with suffix, opens the file for reading and
// writing, and returns the resulting *os.File.
// If dir is the empty string, tempFile uses the default directory for temporary
// files (see os.TempDir).
// The caller can use f.Name() to find the pathname of the file. It is the
// caller's responsibility to remove the file when no longer needed.
func TempFile(dir, prefix, suffix string) (*os.File, error) {
	if dir == "" {
		dir = os.TempDir()
	}

	// Seeded on init
	rnd := strconv.Itoa(rand.Int())
	name := filepath.Join(dir, prefix+rnd+suffix)

	f, err := os.OpenFile(name, os.O_RDWR|os.O_CREATE|os.O_EXCL, 0600)
	if os.IsExist(err) {
		return nil, err
	}

	return f, nil
}
