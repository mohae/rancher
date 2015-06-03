package app

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/mohae/contour"
	jww "github.com/spf13/jwalterweatherman"
)

// rawTemplate holds all the information for a Rancher template. This is used
// to generate the Packer Build.
type rawTemplate struct {
	PackerInf
	IODirInf
	BuildInf
	// holds release information
	releaseISO releaser
	// the builder specific string for the template's OS and Arch
	osType string
	// Current date in ISO 8601
	date string
	// The character(s) used to identify variables for Rancher. By default
	// this is a colon, :. Currently only a starting delimeter is supported.
	delim string
	// The distro that this template targets. The type must be a supported
	// type, i.e. defined in supported.toml. The values for type are
	// consistent with Packer values.
	Distro string
	// The architecture for the ISO image, this is either 32bit or 64bit,
	// with the actual values being dependent on the operating system and
	// the target builder.
	Arch string
	// The image for the ISO. This is distro dependent.
	Image string
	// The release, or version, for the ISO. Usage and values are distro
	// dependent, however only version currently supported images that are
	// available on the distro's download site are supported.
	Release string
	// varVals is a variable replacement map used in finalizing the value of strings for
	// which variable replacement is supported.
	varVals map[string]string
	// Contains all the build information needed to create the target Packer
	// template and its associated artifacts.
	build
	// files maps destination files to their sources. These are the actual file locations
	// after they have been resolved. The destination file is the key, the source file
	// is the value
	files map[string]string
	// dirs maps destination directories to their source directories. Everything within
	// the directory will be copied. The same resolution rules apply for dirs as for
	// filies. The destination directory is the key, the source directory is the value
	dirs map[string]string
}

// mewRawTemplate returns a rawTemplate with current date in ISO 8601 format.
// This should be called when a rawTemplate with the current date is desired.
func newRawTemplate() *rawTemplate {
	// Set the date, formatted to ISO 8601
	date := time.Now()
	splitDate := strings.Split(date.String(), " ")
	return &rawTemplate{date: splitDate[0], delim: contour.GetString(ParamDelimStart), files: make(map[string]string), dirs: make(map[string]string)}
}

// r.createPackerTemplate creates a Packer template from the rawTemplate that
// can be marshalled to JSON.
func (r *rawTemplate) createPackerTemplate() (packerTemplate, error) {
	var err error
	// Resolve the Rancher variables to their final values.
	r.mergeVariables()
	// General Packer Stuff
	p := packerTemplate{}
	p.MinPackerVersion = r.MinPackerVersion
	p.Description = r.Description
	// Builders
	p.Builders, err = r.createBuilders()
	if err != nil {
		jww.ERROR.Println(err)
		return p, err
	}
	// Post-Processors
	p.PostProcessors, err = r.createPostProcessors()
	if err != nil {
		jww.ERROR.Println(err)
		return p, err
	}
	// Provisioners
	p.Provisioners, err = r.createProvisioners()
	if err != nil {
		jww.ERROR.Println(err)
		return p, err
	}
	// Now we can create the Variable Section
	// TODO: currently not implemented/supported
	// Return the generated Packer Template
	return p, nil
}

// replaceVariables checks incoming string for variables and replaces them with
// their values.
func (r *rawTemplate) replaceVariables(s string) string {
	//see if the delim is in the string, if not, nothing to replace
	if strings.Index(s, r.delim) < 0 {
		return s
	}
	// Go through each variable and replace as applicable.
	for vName, vVal := range r.varVals {
		s = strings.Replace(s, vName, vVal, -1)
	}
	return s
}

// r.setDefaults takes the incoming distro settings and merges them with its
// existing settings, which are set to rancher's defaults, to create the
// default template.
func (r *rawTemplate) setDefaults(d *distro) {
	// merges Settings between an old and new template.
	// Note: Arch, Image, and Release are not updated here as how these fields
	// are updated depends on whether this is a build from a distribution's
	// default template or from a defined build template.
	r.IODirInf.update(d.IODirInf)
	r.PackerInf.update(d.PackerInf)
	r.BuildInf.update(d.BuildInf)
	// If defined, BuilderTypes override any prior BuilderTypes Settings
	if d.BuilderTypes != nil && len(d.BuilderTypes) > 0 {
		r.BuilderTypes = d.BuilderTypes
	}
	// If defined, PostProcessorTypes override any prior PostProcessorTypes Settings
	if d.PostProcessorTypes != nil && len(d.PostProcessorTypes) > 0 {
		r.PostProcessorTypes = d.PostProcessorTypes
	}
	// If defined, ProvisionerTypes override any prior ProvisionerTypes Settings
	if d.ProvisionerTypes != nil && len(d.ProvisionerTypes) > 0 {
		r.ProvisionerTypes = d.ProvisionerTypes
	}
	// merge the build portions.
	r.updateBuilders(d.Builders)
	r.updatePostProcessors(d.PostProcessors)
	r.updateProvisioners(d.Provisioners)
	return
}

// r.updateBuildSettings merges Settings between an old and new template.
// Note:  Arch, Image, and Release are not updated here as how these fields are
// updated depends on whether this is a build from a distribution's default
// template or from a defined build template.
func (r *rawTemplate) updateBuildSettings(bld *rawTemplate) {
	r.IODirInf.update(bld.IODirInf)
	r.PackerInf.update(bld.PackerInf)
	r.BuildInf.update(bld.BuildInf)
	// If defined, Builders override any prior builder Settings.
	if bld.BuilderTypes != nil && len(bld.BuilderTypes) > 0 {
		r.BuilderTypes = bld.BuilderTypes
	}
	// If defined, PostProcessorTypes override any prior PostProcessorTypes Settings
	if bld.PostProcessorTypes != nil && len(bld.PostProcessorTypes) > 0 {
		r.PostProcessorTypes = bld.PostProcessorTypes
	}
	// If defined, ProvisionerTypes override any prior ProvisionerTypes Settings
	if bld.ProvisionerTypes != nil && len(bld.ProvisionerTypes) > 0 {
		r.ProvisionerTypes = bld.ProvisionerTypes
	}
	// merge the build portions.
	r.updateBuilders(bld.Builders)
	r.updatePostProcessors(bld.PostProcessors)
	r.updateProvisioners(bld.Provisioners)
}

// mergeVariables goes through the template variables and finalizes the values
// of any :vars found within the strings.
//
// Supported:
//  distro                   the name of the distro
//  release                  the release version being used
//  arch                     the target architecture for the build
//  image                    the image used, e.g. server
//  date                     the current datetime, time.Now()
//  build_name               the name of the build template
//  out_dir                  the directory to write the build output to
//  src_dir                  the directory of any source files used in the build*
//
// Note: src_dir must be set. Rancher searches for referenced files and uses
// src_dir/distro as the last search directory. This directory is also used as
// the base directory for any specified src directories.
//
// TODO should there be a flag to not prefix src paths with src_dir to allow for
// specification of files that are not in src? If the flag is set to not prepend
// src_dir, src_dir could still be used by adding it to the specific variable.
func (r *rawTemplate) mergeVariables() {
	// Get the delim and set the replacement map, resolve name information
	r.setBaseVarVals()
	// get final value for name first
	r.Name = r.replaceVariables(r.Name)
	r.varVals[r.delim+"name"] = r.Name
	// then merge the sourc and out dirs and set them
	r.mergeSrcDir()
	r.mergeOutDir()
	r.varVals[r.delim+"out_dir"] = r.OutDir
	r.varVals[r.delim+"src_dir"] = r.SrcDir
}

// setBaseVarVals sets the varVals for the base variables
func (r *rawTemplate) setBaseVarVals() {
	r.varVals = map[string]string{
		r.delim + "distro":     r.Distro,
		r.delim + "release":    r.Release,
		r.delim + "arch":       r.Arch,
		r.delim + "image":      r.Image,
		r.delim + "date":       r.date,
		r.delim + "build_name": r.BuildName,
	}
}

// mergeVariable does a variable replacement on the passed string and returns
// the finalized value. If the passed string is empty, the default value, d, is
// returned
func (r *rawTemplate) mergeString(s, d string) string {
	if s == "" {
		return d
	}
	return strings.TrimSuffix(r.replaceVariables(s), "/")
}

// mergeSrcDir sets whether or not a custom source directory was used, does any
// necessary variable replacement, and normalizes the string to not end in /
func (r *rawTemplate) mergeSrcDir() {
	// variable replacement is only necessary if the SrcDir has the variable delims
	if !strings.Contains(r.SrcDir, r.delim) {
		// normalize to no ending /
		r.SrcDir = strings.TrimSuffix(r.replaceVariables(r.SrcDir), "/")
		return
	}
	// normalize to no ending /
	r.SrcDir = strings.TrimSuffix(r.replaceVariables(r.SrcDir), "/")
}

// mergeOutDir resolves the out_dir for this template.
func (r *rawTemplate) mergeOutDir() {
	// variable replacement is only necessary if the SrcDir has the variable delims
	if !strings.Contains(r.OutDir, r.delim) {
		// normalize to no ending /
		r.OutDir = strings.TrimSuffix(r.replaceVariables(r.OutDir), "/")
		return
	}
	// normalize to no ending /
	r.OutDir = strings.TrimSuffix(r.replaceVariables(r.OutDir), "/")
}

// ISOInfo sets the ISO info for the template's supported distro type. This
// also sets the builder specific string, when applicable.
// TODO: these should use new functions in release.go. instead of creating the
// structs here
func (r *rawTemplate) ISOInfo(builderType Builder, settings []string) error {
	var k, v, checksumType string
	var err error
	// Only the iso_checksum_type is needed for this.
	for _, s := range settings {
		k, v = parseVar(s)
		switch k {
		case "iso_checksum_type":
			checksumType = v
		}
	}
	switch r.Distro {
	case CentOS.String():
		r.releaseISO = &centOS{
			release: release{
				iso: iso{
					BaseURL:      r.BaseURL,
					ChecksumType: checksumType,
				},
				Arch:    r.Arch,
				Distro:  r.Distro,
				Image:   r.Image,
				Release: r.Release,
			},
		}
		r.releaseISO.SetISOInfo()
		r.osType, err = r.releaseISO.(*centOS).getOSType(builderType.String())
		if err != nil {
			jww.ERROR.Println(err)
			return err
		}
	case Debian.String():
		r.releaseISO = &debian{
			release: release{
				iso: iso{
					BaseURL:      r.BaseURL,
					ChecksumType: checksumType,
				},
				Arch:    r.Arch,
				Distro:  r.Distro,
				Image:   r.Image,
				Release: r.Release,
			},
		}
		r.releaseISO.SetISOInfo()
		r.osType, err = r.releaseISO.(*debian).getOSType(builderType.String())
		if err != nil {
			jww.ERROR.Println(err)
			return err
		}
	case Ubuntu.String():
		r.releaseISO = &ubuntu{
			release: release{
				iso: iso{
					BaseURL:      r.BaseURL,
					ChecksumType: checksumType,
				},
				Arch:    r.Arch,
				Distro:  r.Distro,
				Image:   r.Image,
				Release: r.Release,
			},
		}
		r.releaseISO.SetISOInfo()
		r.osType, err = r.releaseISO.(*ubuntu).getOSType(builderType.String())
		if err != nil {
			jww.ERROR.Println(err)
			return err
		}
	default:
		err := fmt.Errorf("unable to set ISO related information for the unsupported distro: %q", r.Distro)
		jww.ERROR.Println(err)
		return err
	}
	return nil
}

// commandsFromFile returns the commands within the requested file, if it can
// be found. No validation of the contents is done.
func (r *rawTemplate) commandsFromFile(component, name string) (commands []string, err error) {
	// find the file
	src, err := r.findCommandFile(component, name)
	if err != nil {
		return nil, err
	}
	f, err := os.Open(src)
	if err != nil {
		jww.ERROR.Println(err)
		return nil, err
	}
	// always close what's been opened and check returned error
	defer func() {
		cerr := f.Close()
		if cerr != nil && err == nil {
			jww.WARN.Println(cerr)
			err = cerr
		}
	}()
	//New Reader for the string
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		commands = append(commands, scanner.Text())
	}
	err = scanner.Err()
	if err != nil {
		jww.WARN.Println(err)
		return nil, err
	}
	return commands, nil
}

// findCommandFile locates the requested command file. If a match cannot be
// found, an os.ErrNotExist is returned. Any other errors will result in a
// termination of the search.
//
// The request string is build with the following order:
//    commands/{name}
//    {name}
//
// findComponentSource is called to handle the actual location of the file. If
// no match is found an os.ErrNotExist will be returned.
func (r *rawTemplate) findCommandFile(component, name string) (string, error) {
	if name == "" {
		err := fmt.Errorf("the passed command filename was empty")
		jww.ERROR.Println(err)
		return "", err
	}
	findPath := filepath.Join("commands", name)
	src, err := r.findComponentSource(component, findPath)
	// return the error for any error other than ErrNotExist
	if err != nil && err != os.ErrNotExist {
		return "", err
	}
	// if err is nil, the source was found
	if err == nil {
		return src, nil
	}
	return r.findComponentSource(component, name)
}

// findComponentSource attempts to locate the source file or directory referred
// to in p for the requested component and return it's actual location within
// the src_dir.  If the component is not empty, it is added to the path to see
// if there are any component specific files that match.  If none are found,
// just the path is used.  Any match is returned, otherwise an os.ErrNotFound
// error is returned.  Any other error encountered will also be returned.
//
// The search path is built, in order of precedence:
//    component/path
//    component-base/path
//    path
//
// Component is the name of the packer component that this path belongs to,
// e.g. vagrant, chef-client, shell, etc.  The component-base is the base name
// of the packer component that this path belongs to, if applicable, e.g.
// chef-client's base would be chef as would chef-solo's.
func (r *rawTemplate) findComponentSource(component, p string) (string, error) {
	var tmpPath string
	var err error
	// if len(cParts) > 1, there was a - and component-base processing should be done
	if component != "" {
		tmpPath = filepath.Join(component, p)
		tmpPath, err = r.findSource(tmpPath)
		if err != nil && err != os.ErrNotExist {
			return "", err
		}
		if err == nil {
			return tmpPath, nil
		}
		cParts := strings.Split(component, "-")
		if len(cParts) > 1 {
			// first element is the base
			tmpPath = filepath.Join(cParts[0], p)
			tmpPath, err = r.findSource(tmpPath)
			if err != nil && err != os.ErrNotExist {
				return "", err
			}
			if err == nil {
				return tmpPath, nil
			}
		}
	}
	// look for the source as using just the passed path
	tmpPath, err = r.findSource(p)
	if err == nil {
		return tmpPath, nil
	}
	return "", err
}

// findSource searches for the specified sub-path using Rancher's algorithm for
// finding the correct location.  Passed names may include relative path
// information and may be either a filename or a directory.  Releases may have
// "."'s in them.  In addition to searching for the requested source within the
// point release, the "." are stripped out and the resulting value is searched:
// e.g. 14.04 becomes 1404.  The base release number is also checked: e.g. 14 is
// searched for 14.04.
// Search order:
//   src_dir/distro/release/build_name/
//   src_dir/distro/releaseBase/build_name/
//   src_dir/distro/build_name/
//   src_dir/build_name/
//   src_dir/distro/release/arch/
//   src_dir/distro/releaseBase/arch/
//   src_dir/distro/release/
//   src_dir/distro/releaseBase/
//   src_dir/distro/arch
//   src_dir/distro/
//   src_dir/
//
// If the passed path is not found, an os.ErrNotExist will be returned
func (r *rawTemplate) findSource(p string) (string, error) {
	if p == "" {
		return "", fmt.Errorf("cannot find source, no path received")
	}
	releaseParts := strings.Split(r.Release, ".")
	var release string
	if len(releaseParts) > 1 {
		for _, v := range releaseParts {
			release += v
		}
	}
	// src_dir/:distro/:release/:build_name/p
	tmpPath := filepath.Join(r.SrcDir, r.Distro, r.Release, r.BuildName, p)
	_, err := os.Stat(tmpPath)
	if err == nil {
		jww.TRACE.Printf("findSource:  %s found", tmpPath)
		return tmpPath, nil
	}
	jww.TRACE.Printf("findSource:  %s not found", tmpPath)
	// src_dir/:distro/release/:build_name/p
	tmpPath = filepath.Join(r.SrcDir, r.Distro, release, r.BuildName, p)
	_, err = os.Stat(tmpPath)
	if err == nil {
		jww.TRACE.Printf("findSource:  %s found", tmpPath)
		return tmpPath, nil
	}
	jww.TRACE.Printf("findSource:  %s not found", tmpPath)
	// src_dir/:distro/releaseBase/:build_name/p
	tmpPath = filepath.Join(r.SrcDir, r.Distro, releaseParts[0], r.BuildName, p)
	_, err = os.Stat(tmpPath)
	if err == nil {
		jww.TRACE.Printf("findSource:  %s found", tmpPath)
		return tmpPath, nil
	}
	jww.TRACE.Printf("findSource:  %s not found", tmpPath)
	// src_dir/:distro/:build_name/p
	tmpPath = filepath.Join(r.SrcDir, r.Distro, r.BuildName, p)
	_, err = os.Stat(tmpPath)
	if err == nil {
		jww.TRACE.Printf("findSource:  %s found", tmpPath)
		return tmpPath, nil
	}
	jww.TRACE.Printf("findSource:  %s not found", tmpPath)
	// src_dir/:build_name/p
	tmpPath = filepath.Join(r.SrcDir, r.BuildName, p)
	_, err = os.Stat(tmpPath)
	if err == nil {
		jww.TRACE.Printf("findSource:  %s found", tmpPath)
		return tmpPath, nil
	}
	jww.TRACE.Printf("findSource:  %s not found", tmpPath)
	// src_dir/:distro/:release/:arch/p
	tmpPath = filepath.Join(r.SrcDir, r.Distro, r.Release, r.Arch, p)
	_, err = os.Stat(tmpPath)
	if err == nil {
		jww.TRACE.Printf("findSource:  %s found", tmpPath)
		return tmpPath, nil
	}
	jww.TRACE.Printf("findSource:  %s not found", tmpPath)
	// src_dir/:distro/release/:arch/p
	tmpPath = filepath.Join(r.SrcDir, r.Distro, release, r.Arch, p)
	_, err = os.Stat(tmpPath)
	if err == nil {
		jww.TRACE.Printf("findSource:  %s found", tmpPath)
		return tmpPath, nil
	}
	jww.TRACE.Printf("findSource:  %s not found", tmpPath)
	// src_dir/:distro/releaseBase/:arch/p
	tmpPath = filepath.Join(r.SrcDir, r.Distro, releaseParts[0], r.Arch, p)
	_, err = os.Stat(tmpPath)
	if err == nil {
		jww.TRACE.Printf("findSource:  %s found", tmpPath)
		return tmpPath, nil
	}
	jww.TRACE.Printf("findSource:  %s not found", tmpPath)
	// src_dir/:distro/:release/p
	tmpPath = filepath.Join(r.SrcDir, r.Distro, r.Release, p)
	_, err = os.Stat(tmpPath)
	if err == nil {
		jww.TRACE.Printf("findSource:  %s found", tmpPath)
		return tmpPath, nil
	}
	jww.TRACE.Printf("findSource:  %s not found", tmpPath)
	// src_dir/:distro/release/p
	tmpPath = filepath.Join(r.SrcDir, r.Distro, release, p)
	_, err = os.Stat(tmpPath)
	if err == nil {
		jww.TRACE.Printf("findSource:  %s found", tmpPath)
		return tmpPath, nil
	}
	jww.TRACE.Printf("findSource:  %s not found", tmpPath)
	// src_dir/:distro/releaseBase/p
	tmpPath = filepath.Join(r.SrcDir, r.Distro, releaseParts[0], p)
	_, err = os.Stat(tmpPath)
	if err == nil {
		jww.TRACE.Printf("findSource:  %s found", tmpPath)
		return tmpPath, nil
	}
	jww.TRACE.Printf("findSource:  %s not found", tmpPath)
	// src_dir/:distro/:arch/p
	tmpPath = filepath.Join(r.SrcDir, r.Distro, r.Arch, p)
	_, err = os.Stat(tmpPath)
	if err == nil {
		jww.TRACE.Printf("findSource:  %s found", tmpPath)
		return tmpPath, nil
	}
	jww.TRACE.Printf("findSource:  %s not found", tmpPath)
	// src_dir/:distro/p
	tmpPath = filepath.Join(r.SrcDir, r.Distro, p)
	_, err = os.Stat(tmpPath)
	if err == nil {
		jww.TRACE.Printf("findSource:  %s found", tmpPath)
		return tmpPath, nil
	}
	jww.TRACE.Printf("findSource:  %s not found", tmpPath)
	// src_dir/p
	tmpPath = filepath.Join(r.SrcDir, p)
	_, err = os.Stat(tmpPath)
	if err == nil {
		jww.TRACE.Printf("findSource:  %s found", tmpPath)
		return tmpPath, nil
	}
	jww.TRACE.Printf("findSource:  %s not found", tmpPath)
	return "", os.ErrNotExist
}

// buildOutPath builds the full output path of the passed path, p, and returns
// that value.  If the template is set to include the component string as the
// parent directory, it is added to the path.
func (r *rawTemplate) buildOutPath(component, p string) string {
	if r.IncludeComponentString && component != "" {
		return filepath.Join(r.OutDir, component, p)
	}
	return filepath.Join(r.OutDir, p)
}

// buildTemplateResourcePath builds the path that will be added to the Packer
// template for the passed path, p, and returns that value.  If the template is
// set to include the component string as the parent directory, it is added to
// the path.
func (r *rawTemplate) buildTemplateResourcePath(component, p string) string {
	if r.IncludeComponentString && component != "" {
		return filepath.Join(component, p)
	}
	return p
}