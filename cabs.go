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
