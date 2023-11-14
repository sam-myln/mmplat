package fs

import (
	"golang.org/x/exp/slices"
	"io"
	"io/fs"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// IItem basically, if remote mechanism is used to work with files,
// then methods to offset data must be used (or local cache, witch can blow,
// so maybe memlimit of file size limit ensured
type IItem interface {
	// Path fullpath
	Path() string
	// Metadata return file format in fasthttp metadata fashion
	Metadata() string
	// Name basename
	Name() string
	Size() int64
	IsDir() bool
	Reader() (io.ReadSeekCloser, error)
}

type Item struct {
	path     string
	fileInfo os.FileInfo
	metadata string
}

func (item *Item) Path() string {
	return item.path
}

func (item *Item) Metadata() string {
	return item.metadata
}

// Name basename
func (item *Item) Name() string {
	return item.fileInfo.Name()
}

func (item *Item) Size() int64 {
	return item.fileInfo.Size()
}

func (item *Item) IsDir() bool {
	return item.fileInfo.IsDir()
}

// Reader closed by caller
func (item *Item) Reader() (io.ReadSeekCloser, error) {
	return os.Open(item.path)
}

// FSWorker merges multiple dirs under the `root` /
type FSWorker struct {
	recursive bool
	fmts      []string
	dirs      []string
	tree      *Tree
}

func NewFSWorker(recurse bool,
	dirs []string,
	fmts []string,
) *FSWorker {
	var afmts []string
	var arecurse bool
	var adirs []string
	afmts = fmts
	if len(fmts) == 0 {
		afmts = []string{"*"}
	}
	arecurse = recurse
	adirs = dirs
	if len(dirs) == 0 {
		adirs = []string{"./"}
	}
	return &FSWorker{
		arecurse,
		afmts,
		adirs,
		NewTree(),
	}
}

func (fs *FSWorker) Withfmts(fmts ...string) *FSWorker {
	fs.fmts = append(fs.fmts, fmts...)
	return fs
}

func (fs *FSWorker) WithDirs(dirs ...string) *FSWorker {
	fs.dirs = append(fs.dirs, dirs...)
	return fs
}

func (fs *FSWorker) WithRecurse(r bool) *FSWorker {
	fs.recursive = r
	return fs
}

// BuildTree dir - directory, ext - permitted extensions
//
//	do NOT follow symlnix or go lower than root
//
// use stat + base
// root item has Null IItem
func (lfs *FSWorker) BuildTree() error {
	// top level subroot of the actual root, corresponds to the dir being iterated
	var root *Node
	// subrtoot corresponds to the nested dir being traversed
	var subRoot *Node
	var walkDirFunc func(path string, de fs.DirEntry, err error) error
	for _, dir := range lfs.dirs {
		root = NewNode(Dir, dir, lfs.Tree().Root())
		subRoot = root
		lfs.Tree().Root().Add(root)
		// path is full path /W name, de.name is filename
		walkDirFunc = func(path string, de fs.DirEntry, err error) error {
			if err != nil {
				return err
			}
			switch de.IsDir() {
			case true:
				str, found := strings.CutSuffix(path, "/"+de.Name())
				if !lfs.recursive && found && subRoot.Item().Path() == str {
					return filepath.SkipAll
				}
				if path != subRoot.Item().Path() {
					dupl := lfs.Tree().FindBy(func(needle *Node) *Node {
						if needle.Item().Path() == path {
							return needle
						}
						return nil
					})
					if dupl != nil {
						return filepath.SkipDir
					}
					t := NewNode(Dir, path, subRoot)
					subRoot.Add(t)
					subRoot = t
					filepath.WalkDir(path, walkDirFunc)
					// reset subroot up on lvl
					if subRoot.Parent() != lfs.Tree().Root() {
						subRoot = subRoot.Parent()
					}
				}
			case false:
				dupl := lfs.Tree().FindBy(func(needle *Node) *Node {
					if needle.Item().Path() == path {
						return needle
					}
					return nil
				})
				if dupl != nil {
					return filepath.SkipDir
				}
				if lfs.FilterFormat(path) {
					subRoot.Add(NewNode(File, path, subRoot))
				}
			}

			return nil
		}
		err := filepath.WalkDir(dir, walkDirFunc)
		if err != nil && err != filepath.SkipDir && err != filepath.SkipAll {
			return err
		}
	}
	return nil
}

func CreateItem(path string) (*Item, error) {
	var e error

	item := &Item{}
	item.path = path
	var err error
	item.fileInfo, err = os.Stat(path)
	if err != nil {
		return nil, err
	}
	f, e := os.Open(path)
	if e == nil {
		defer f.Close()
		buf := make([]byte, 512)
		_, e = f.Read(buf)
		if e != nil && !item.fileInfo.IsDir() {
			item.metadata = ""
		} else if item.fileInfo.IsDir() {
			item.metadata = "inode/directory"
		} else {
			item.metadata = http.DetectContentType(buf)
		}
	}

	return item, e
}

func (fs *FSWorker) Tree() *Tree {
	return fs.tree
}

//FilterFormat if preceded by - excludes, includes otherwise
func (fs *FSWorker) FilterFormat(path string) bool {
	if len(fs.fmts) == 1 && fs.fmts[0] == `*` {
		return true
	}
	ext := filepath.Ext(path)
	if ext != "" {
		ext = ext[1:]
	}
	idx := slices.Index(fs.fmts, ext)
	return idx != -1 &&	!strings.Contains(fs.fmts[idx], "-")
}
