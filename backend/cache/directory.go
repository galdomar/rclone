// +build !plan9,go1.7

package cache

import (
	"time"

	"path"

	"github.com/ncw/rclone/fs"
)

// Directory is a generic dir that stores basic information about it
type Directory struct {
	fs.Directory `json:"-"`

	CacheFs      *Fs    `json:"-"`       // cache fs
	Name         string `json:"name"`    // name of the directory
	Dir          string `json:"dir"`     // abs path of the directory
	CacheModTime int64  `json:"modTime"` // modification or creation time - IsZero for unknown
	CacheSize    int64  `json:"size"`    // size of directory and contents or -1 if unknown

	CacheItems int64      `json:"items"`     // number of objects or -1 for unknown
	CacheType  string     `json:"cacheType"` // object type
	CacheTs    *time.Time `json:",omitempty"`
}

// NewDirectory builds an empty dir which will be used to unmarshal data in it
func NewDirectory(f *Fs, remote string) *Directory {
	cd := ShallowDirectory(f, remote)
	t := time.Now()
	cd.CacheTs = &t

	return cd
}

// ShallowDirectory builds an empty dir which will be used to unmarshal data in it
func ShallowDirectory(f *Fs, remote string) *Directory {
	var cd *Directory
	fullRemote := cleanPath(path.Join(f.Root(), remote))

	// build a new one
	dir := cleanPath(path.Dir(fullRemote))
	name := cleanPath(path.Base(fullRemote))
	cd = &Directory{
		CacheFs:      f,
		Name:         name,
		Dir:          dir,
		CacheModTime: time.Now().UnixNano(),
		CacheSize:    0,
		CacheItems:   0,
		CacheType:    "Directory",
	}

	return cd
}

// DirectoryFromOriginal builds one from a generic fs.Directory
func DirectoryFromOriginal(f *Fs, d fs.Directory) *Directory {
	var cd *Directory
	fullRemote := path.Join(f.Root(), d.Remote())

	dir := cleanPath(path.Dir(fullRemote))
	name := cleanPath(path.Base(fullRemote))
	t := time.Now()
	cd = &Directory{
		Directory:    d,
		CacheFs:      f,
		Name:         name,
		Dir:          dir,
		CacheModTime: d.ModTime().UnixNano(),
		CacheSize:    d.Size(),
		CacheItems:   d.Items(),
		CacheType:    "Directory",
		CacheTs:      &t,
	}

	return cd
}

// Fs returns its FS info
func (d *Directory) Fs() fs.Info {
	return d.CacheFs
}

// String returns a human friendly name for this object
func (d *Directory) String() string {
	if d == nil {
		return "<nil>"
	}
	return d.Remote()
}

// Remote returns the remote path
func (d *Directory) Remote() string {
	return d.CacheFs.cleanRootFromPath(d.abs())
}

// abs returns the absolute path to the dir
func (d *Directory) abs() string {
	return cleanPath(path.Join(d.Dir, d.Name))
}

// parentRemote returns the absolute path parent remote
func (d *Directory) parentRemote() string {
	absPath := d.abs()
	if absPath == "" {
		return ""
	}
	return cleanPath(path.Dir(absPath))
}

// ModTime returns the cached ModTime
func (d *Directory) ModTime() time.Time {
	return time.Unix(0, d.CacheModTime)
}

// Size returns the cached Size
func (d *Directory) Size() int64 {
	return d.CacheSize
}

// Items returns the cached Items
func (d *Directory) Items() int64 {
	return d.CacheItems
}

var (
	_ fs.Directory = (*Directory)(nil)
)
