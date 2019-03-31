// Go analogue of npm's content-addressable-blob-store
//
// Copyright 2019 Robert Charles Sheets
//
// See the LICENSE file for license terms.

// Package cabs implements a content-addressable blob store with an on-disk
// format that aims for compatibility with the on-disk format used by
// https://www.npmjs.com/package/content-addressable-blob-store so that blob
// stores created with that module can be used by Go projects.
package cabs

import (
	"crypto/sha256"
	"fmt"
	"io/ioutil"
	"os"
)

type blobstore interface {
	// Write takes a byte slice of blob data and writes it to the store,
	// returning its SHA256 sum as an array of bytes. If the write doesn't go
	// well, an error is returned instead.
	Write([]byte) ([]byte, error)

	// Read takes a 32-byte value and tries to find the corresponding blob.
	// If it finds one, it returns it. If not, it returns nil and an error.
	Read([]byte) ([]byte, error)
}

type filesystemBackedCABS struct {
	basepath string
}

func NewFilesystemBackedCABS(path string) (*filesystemBackedCABS, error) {
	var bs filesystemBackedCABS
	err := os.MkdirAll(path, 0777)
	if err != nil {
		return nil, err
	}
	bs.basepath = path
	return &bs, nil
}

func (c *filesystemBackedCABS) Write(blob []byte) ([]byte, error) {
	h := sha256.New()
	h.Write(blob)
	sum := h.Sum(nil)

	dirPath := fmt.Sprintf("%s/%x", c.basepath, sum[0:1])
	err := os.MkdirAll(dirPath, 0777)
	if err != nil {
		return []byte{}, err
	}

	filePath := fmt.Sprintf("%s/%x", dirPath, sum[1:])
	err = ioutil.WriteFile(filePath, blob, 0666)
	if err != nil {
		return []byte{}, err
	}

	return sum, nil
}

func (c *filesystemBackedCABS) Read(hash []byte) ([]byte, error) {
	path := fmt.Sprintf("%x/%x", hash[0:1], hash[1:])
	filePath := fmt.Sprintf("%s/%s", c.basepath, path)
	blob, err := ioutil.ReadFile(filePath)
	if err != nil {
		return []byte{}, err
	}

	return blob, nil
}
