package app

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/mohae/contour"
)

// Archive holds information about an archive.
type Archive struct {
	// Name is the name of the archive, w/o extensions
	Name string
	// Path to the target directory for the archive output.
	OutDir string
	// Compression type to be used.
	Type string
	// List of files to add to the archive.
	directory
}

func NewArchive(s string) *Archive {
	return &Archive{Name: s}
}
func (a *Archive) addFile(tW *tar.Writer, filename string) error {
	// Add the passed file, if it exists, to the archive, otherwise error.
	// This preserves mode and modification.
	// TODO check ownership/permissions
	file, err := os.Open(filename)
	if err != nil {
		return archivePriorBuildErr(err)
	}
	defer file.Close()
	var fileStat os.FileInfo
	fileStat, err = file.Stat()
	if err != nil {
		return archivePriorBuildErr(err)
	}
	// Don't add directories--they result in tar header errors.
	fileMode := fileStat.Mode()
	if fileMode.IsDir() {
		return nil
	}
	// Create the tar header stuff.
	tH := new(tar.Header)
	tH.Name = filename
	tH.Size = fileStat.Size()
	tH.Mode = int64(fileStat.Mode())
	tH.ModTime = fileStat.ModTime()
	// Write the file header to the tarball.
	err = tW.WriteHeader(tH)
	if err != nil {
		return archivePriorBuildErr(err)
	}
	// Add the file to the tarball.
	_, err = io.Copy(tW, file)
	if err != nil {
		return archivePriorBuildErr(err)
	}
	return nil
}

// priorBuild handles archiving prior build artifacts, if it exists, and then
// deleting those artifacts. This prevents any stale elements from persisting
// to the new build.
func (a *Archive) priorBuild(p string) error {
	if !contour.GetBool(ArchivePriorBuild) {
		return nil
	}
	// See if src exists, if it doesn't then don't do anything
	_, err := os.Stat(p)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return archivePriorBuildErr(err)
	}
	// Archive the old artifacts.
	err = a.create(p)
	if err != nil {
		return archivePriorBuildErr(err)
	}
	// Delete the old artifacts.
	err = a.deletePriorBuild(p)
	if err != nil {
		return err
	}
	return nil
}

func (a *Archive) create(p string) error {
	// examples don't get archived
	if contour.GetBool(Example) {
		return nil
	}
	// Get a list of directory contents
	err := a.DirWalk(p)
	if err != nil {
		return err
	}
	if len(a.Files) <= 1 {
		// This isn't a real error, just log it and return a non-error state.
		return nil
	}
	// Get the relative path so that it can be added to the tarball name.
	relPath := filepath.Dir(filepath.Clean(p))
	// The tarball's name is the directory name + extension.  If there is a collision
	// on the resulting name, a unique name will be generated and returned.
	tBName, err := archiveFilename(relPath, a.Name)
	if err != nil {
		return err
	}
	// Create the new archive file.
	tBall, err := os.Create(tBName)
	if err != nil {
		return err
	}
	// Close the file with error handling
	defer func() {
		cerr := tBall.Close()
		if cerr != nil && err == nil {
			err = cerr
		}
	}()
	// The tarball gets compressed with gzip
	gw := gzip.NewWriter(tBall)
	defer func() {
		cerr := gw.Close()
		if cerr != nil && err == nil {
			err = cerr
		}
	}()
	// Create the tar writer.
	tW := tar.NewWriter(gw)
	defer func() {
		cerr := tW.Close()
		if cerr != nil && err == nil {
			err = cerr
		}
	}()
	// Go through each file in the path and add it to the archive
	for _, f := range a.Files {
		err := a.addFile(tW, filepath.Join(p, f.p))
		if err != nil {
			return err
		}
	}
	return nil
}

func (a *Archive) deletePriorBuild(p string) error {
	//delete the contents of the passed directory
	err := deleteDir(p)
	if err != nil {
		return fmt.Errorf("deletePriorBuild failed: %s", err.Error())
	}
	return nil
}

// directory is a container for files to add to an archive.
type directory struct {
	// A slice of file structs.
	Files []file
}

// file contains information about a file
type file struct {
	// The file's path
	p string
	// The file's FileInfo
	info os.FileInfo
}

// DirWalk walks the passed path, making a list of all the files that are
// children of the path.
func (d *directory) DirWalk(dirPath string) error {
	// If the directory exists, create a list of its contents.
	if dirPath == "" {
		// If nothing was passed, do nothing. This is not an error.
		// However archive.Files will be nil
		return nil
	}
	// See if the path exists
	exists, err := pathExists(dirPath)
	if err != nil {
		return archivePriorBuildErr(err)
	}
	if !exists {
		return archivePriorBuildErr(fmt.Errorf("%s does not exist", dirPath))
	}
	fullPath, err := filepath.Abs(dirPath)
	if err != nil {
		return archivePriorBuildErr(err)
	}
	// Set up the call back function.
	callback := func(p string, fi os.FileInfo, err error) error {
		return d.addFilename(fullPath, p, fi, err)
	}
	// Walk the tree.
	return filepath.Walk(fullPath, callback)
}

// Add the current file information to the file slice.
func (d *directory) addFilename(root, p string, fi os.FileInfo, err error) error {
	// Add a file to the slice of files for which an archive will be created.
	// See if the path exists
	var exists bool
	exists, err = pathExists(p)
	if err != nil {
		return err
	}
	if !exists {
		return archivePriorBuildErr(fmt.Errorf("%s does not exist", p))
	}
	// Get the relative information.
	rel, err := filepath.Rel(root, p)
	if err != nil {
		return archivePriorBuildErr(err)
	}
	if rel == "." {
		return nil
	}
	// Add the file information.
	d.Files = append(d.Files, file{p: rel, info: fi})
	return nil
}

// deleteDir deletes the contents of a directory.
func deleteDir(dir string) error {
	var dirs []string
	// see if the directory exists first, actually any error results in the
	// same handling so just return on any error instead of doing an
	// os.IsNotExist(err)
	_, err := os.Stat(dir)
	if err != nil {
		if os.IsNotExist(err) {
			return fmt.Errorf("deleteDir: %s", err.Error())
		}
	}
	dirInf := directory{}
	dirInf.DirWalk(dir)
	dir = appendSlash(dir)
	for _, file := range dirInf.Files {
		if file.info.IsDir() {
			dirs = append(dirs, dir+file.p)
			continue
		}
		err := os.Remove(dir + file.p)
		if err != nil {
			return fmt.Errorf("deleteDir: %s", err)
		}
	}
	// all the files should now be deleted so its safe to delete the directories
	// do this in reverse order
	for i := len(dirs) - 1; i >= 0; i-- {
		err = os.Remove(dirs[i])
		if err != nil {
			return fmt.Errorf("deleteDir: %s", err.Error())
		}
	}
	return nil
}

// archiveFilename returns the name of the archive to be created
func archiveFilename(p, name string) (string, error) {
	name = fmt.Sprintf("%s.tar.gz", filepath.Join(appendSlash(p), name))
	// ensure the archive name is unique
	return getUniqueFilename(name, "2006-01-02")
}
