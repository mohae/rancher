// Copyright 2014 Joel Scoble. All Rights Reserved.
//
// This Source Code Form is subject to the terms of the Mozilla Public
// License, v. 2.0. If a copy of the MPL was not distributed with this
// file, You can obtain one at http://mozilla.org/MPL/2.0/.

// Package ranchr implements the creation of Packer templates from Rancher
// build definitions.
package ranchr

import (
	"errors"
	"fmt"
	"os"
	"strings"
	"sync"

	"github.com/BurntSushi/toml"
	jww "github.com/spf13/jwalterweatherman"
)

// Contains most of the information for Packer templates within a Rancher Build.
type build struct {
	// Targeted builders: the values are consistent with Packer's, e.g.
	// `virtualbox.iso` is used for VirtualBox.
	BuilderType []string `toml:"builder_type"`

	// A map of builder configuration. There should always be a `common`
	// builder, which has settings common to both VMWare and VirtualBox.
	Builders map[string]builder `toml:"builders"`

	// A map of post-processor configurations.
	PostProcessors map[string]postProcessors `toml:"post_processors"`

	// A map of provisioner configurations.
	Provisioners map[string]provisioners `toml:"provisioners"`
}

// Defines a representation of the builder section of a Packer template.
type builder struct {
	// Settings that are common to both builders.
	Settings []string `toml:"Settings"`

	// VM Specific settings. Each VM builder should have its own defined.
	// The 'common' builder does not have this section
	VMSettings []string `toml:"vm_Settings"`
}

// Merge the settings section of a builder. New values supercede existing ones.
func (b *builder) mergeSettings(new []string) {
	b.Settings = mergeSettingsSlices(b.Settings, new)
}

// Merge the VMSettings section of a builder. New values supercede existing ones.
func (b *builder) mergeVMSettings(new []string) {
	b.VMSettings = mergeSettingsSlices(b.VMSettings, new)
}

// Go through all of the Settings and convert them to a map. Each setting
// is parsed into its constituent parts. The value then goes through
// variable replacement to ensure that the settings are properly resolved.
func (b *builder) settingsToMap(r *rawTemplate) map[string]interface{} {
	var k, v string
	m := make(map[string]interface{}, len(b.Settings)+len(b.VMSettings))

	for _, s := range b.Settings {
		k, v = parseVar(s)
		v = r.replaceVariables(v)
		m[k] = v
	}

	return m
}

// Type for handling the post-processor section of the configs.
type postProcessors struct {
	// Rest of the settings in "key=value" format.
	Settings []string
	// Support Packer post-processors' `except` configuration.
	Except []string
	// Support Packer post-processors' `only` configuration.
	Only []string
}

// Merge the settings section of a post-processor. New values supercede
// existing ones.
func (p *postProcessors) mergeSettings(new []string) {
	p.Settings = mergeSettingsSlices(p.Settings, new)
}

// Merge the Except section of a post-processor. If there are new values,
// replace the existing values.
func (p *postProcessors) mergeExcept(new []string) {
	if len(new) > 0 {
		p.Except = new
	}
}

// Merge the Only section of a post-processor. New values supercede
// existing ones.
func (p *postProcessors) mergeOnly(new []string) {
	if len(new) > 0 {
		p.Only = new
	}
}

// Go through all of the Settings and convert them to a map. Each setting
// is parsed into its constituent parts. The value then goes through
// variable replacement to ensure that the settings are properly resolved.
func (p *postProcessors) settingsToMap(Type string, r *rawTemplate) map[string]interface{} {
	var k string
	var v interface{}
	m := make(map[string]interface{}, len(p.Settings))
	m["type"] = Type

	for _, s := range p.Settings {
		k, v = parseVar(s)

		switch k {
		case "keep_input_artifact":
			if strings.ToLower(fmt.Sprint(v)) == "true" {
				v = true
			} else {
				v = false
			}
		default:
			v = r.replaceVariables(fmt.Sprint(v))
		}

		m[k] = v
	}

	return m
}

// Type for handling the provisioners sections of the configs.
type provisioners struct {
	// Rest of the settings in "key=value" format.
	Settings []string `toml:"settings"`
	// Scripts are defined separately because it's simpler that way.
	Scripts []string `toml:"scripts"`
	// Support Packer Provisioner's `except` configuration.
	Except []string `toml:"except"`
	// Support Packer Provisioner's `only` configuration.
	Only []string `toml:"only"`
}

// Merge the Except section of a post-processor. If there are new values,
// replace the existing values.
func (p *provisioners) mergeExcept(new []string) {
	if len(new) > 0 {
		p.Except = new
	}
}

// Merge the Only section of a post-processor. New values supercede
// existing ones.
func (p *provisioners) mergeOnly(new []string) {
	if len(new) > 0 {
		p.Only = new
	}
}

// Merge the settings section of a post-processor. New values supercede existing ones.
func (p *provisioners) mergeSettings(new []string) {
	p.Settings = mergeSettingsSlices(p.Settings, new)
}

// Go through all of the Settings and convert them to a map. Each setting is
// parsed into its constituent parts. The value then goes through variable
// replacement to ensure that the settings are properly resolved.
func (p *provisioners) settingsToMap(Type string, r *rawTemplate) (map[string]interface{}, error) {
	var k, v string
	var err error

	m := make(map[string]interface{}, len(p.Settings))
	m["type"] = Type

	for _, s := range p.Settings {
		k, v = parseVar(s)
		v = r.replaceVariables(v)

		switch k {
		case "execute_command":
			// Get the command from the specified file
			var c []string

			if c, err = commandsFromFile(v); err != nil {
				jww.ERROR.Print(err.Error())
				return nil, err
			}
			v = c[0]
		}

		m[k] = v
	}

	// Add except array.
	if p.Except != nil {
		m["except"] = p.Except
	}

	// Add only array.
	if p.Only != nil {
		m["only"] = p.Only
	}

	return m, nil
}

func (p *provisioners) setScripts(new []string) {
	// Scripts are only replaced if it has values, otherwise the existing values are used.
	if len(new) > 0 {
		p.Scripts = new
	}
}

// defaults is used to store Rancher application level defaults for Packer templates.
type defaults struct {
	IODirInf
	PackerInf
	BuildInf
	build
	load   sync.Once
	loaded bool
}

// BuildInf is a container for information about a specific build.
type BuildInf struct {
	// Name is the name for the build. This may be an assigned value from
	// a TOML file setting.
	Name      string `toml:"name"`
	// BuildName is the name of the build. This is either the name, as
	// specified in the build.toml file, or a generated name for -distro
	// flag based builds.
	BuildName string `toml:"build_name"`
	BaseURL   string `toml:"base_url"`
}

func (i *BuildInf) update(new BuildInf) {
	if new.Name != "" {
		i.Name = new.Name
	}

	if new.BuildName != "" {
		i.BuildName = new.BuildName
	}

	return
}

// IODirInf is used to store information about where Rancher can find and put
// things. Source files are always in a SrcDir, e.g. HTTPSrcDir is the source
// directory for the HTTP directory. The destination directory is always a Dir,
// e.g. HTTPDir is the destination directory for the HTTP directory.
type IODirInf struct {
	// The directory in which the command files are located
	CommandsSrcDir string `toml:"commands_src_dir"`

	// The directory that will be used for the HTTP setting.
	HTTPDir string `toml:"http_dir"`

	// The directory that is the source for files to be copied to the HTTP
	// directory, HTTPDir
	HTTPSrcDir string `toml:"http_src_dir"`

	// The directory that the output artifacts will be written to.
	OutDir string `toml:"out_dir"`

	// The directory that scripts for the Packer template will be copied to.
	ScriptsDir string `toml:"scripts_dir"`

	// The directory that contains the scripts that will be copied.
	ScriptsSrcDir string `toml:"scripts_src_dir"`

	// The directory that contains the source files for this build.
	SrcDir string `toml:"src_dir"`
}

func (i *IODirInf) update(new IODirInf) {
	if new.CommandsSrcDir != "" {
		i.CommandsSrcDir = appendSlash(new.CommandsSrcDir)
	}

	if new.HTTPDir != "" {
		i.HTTPDir = appendSlash(new.HTTPDir)
	}

	if new.HTTPSrcDir != "" {
		i.HTTPSrcDir = appendSlash(new.HTTPSrcDir)
	}

	if new.OutDir != "" {
		i.OutDir = appendSlash(new.OutDir)
	}

	if new.ScriptsDir != "" {
		i.ScriptsDir = appendSlash(new.ScriptsDir)
	}

	if new.ScriptsSrcDir != "" {
		i.ScriptsSrcDir = appendSlash(new.ScriptsSrcDir)
	}

	if new.SrcDir != "" {
		i.SrcDir = appendSlash(new.SrcDir)
	}

	return
}

// PackerInf is used to store information about a Packer Template. In Packer,
// these fields are optional, put used here because they are always printed out
// in a template as custom creation of template output hasn't been written--it
// may never be written.
type PackerInf struct {
	MinPackerVersion string `toml:"min_packer_version" json:"min_packer_version"`
	Description      string `toml:"description" json:"description"`
}

func (i *PackerInf) update(new PackerInf) {
	if new.MinPackerVersion != "" {
		i.MinPackerVersion = new.MinPackerVersion
	}

	if new.Description != "" {
		i.Description = new.Description
	}

	return
}

// Ensures that the default configs get loaded once. Uses a mutex to prevent
// race conditions as there can be concurrent processing of Packer templates.
// When loaded, it sets the loaded boolean so that it only needs to be called
// when it hasn't been loaded.
func (d *defaults) LoadOnce() error {
	var err error

	loadFunc := func() {
		name := os.Getenv(EnvDefaultsFile)

		if name == "" {
			err = errors.New("could not retrieve the default Settings because the " + EnvDefaultsFile + " environment variable was not set. Either set it or check your rancher.cfg setting")
			jww.CRITICAL.Print(err.Error())
			return
		}

		if _, err = toml.DecodeFile(name, &d); err != nil {
			jww.CRITICAL.Print(err.Error())
			return
		}

		d.loaded = true
		return
	}

	d.load.Do(loadFunc)

	// Don't need to log this as the loadFunc logged already logged it
	if err != nil {
		return err
	}

	jww.DEBUG.Printf("defaults loaded: %v", d)

	return nil
}

// To add support for a distribution, the information about it must be added to
// the supported. file, in addition to adding the code to support it to the
// application.
type supported struct {
	Distro map[string]*distro
	load   sync.Once
	loaded bool
}

// Struct to hold the details of supported distros. From this information a
// user should be able to build a Packer template by only executing the
// following, at minimum:
//	$ rancher build -distro=ubuntu
// All settings can be overridden. The information here represents the standard
// box configuration for its respective distribution.
type distro struct {
	IODirInf
	PackerInf
	BuildInf

	// The supported Architectures, which can differ per distro. The labels can also
	// differ, e.g. amd64 and x86_64.
	Arch []string `toml:"Arch"`

	// Supported iso Images, e.g. server, minimal, etc.
	Image []string `toml:"Image"`

	// Supported Releases: the supported Releases are the Releases available for
	// download from that distribution's download page. Archived and unsupported
	// Releases are not used.
	Release []string `toml:"Release"`

	// The default Image configuration for this distribution. This usually consists of
	// things like Release, Architecture, Image type, etc.
	DefImage []string `toml:"default_Image"`

	// The configurations needed to generate the default settings for a build for this
	// distribution.
	build
}

// Ensures that the supported distro information only get loaded once. Uses a
// mutex to prevent race conditions as there can be concurrent processing of
// Packer templates. When loaded, it sets the loaded boolean so that it only
// needs to be called when it hasn't been loaded.
func (s *supported) LoadOnce() error {
	var err error

	loadFunc := func() {
		name := os.Getenv(EnvSupportedFile)

		if name == "" {
			err = errors.New("could not retrieve the Supported information because the " + EnvSupportedFile + " environment variable was not set. Either set it or check your rancher.cfg setting")
			jww.CRITICAL.Print(err.Error())
			return
		}

		if _, err = toml.DecodeFile(name, &s); err != nil {
			jww.CRITICAL.Print(err.Error())
			return
		}

		s.loaded = true
		return
	}

	s.load.Do(loadFunc)

	// Don't need to log the error as loadFunc already did it.
	if err != nil {
		return err
	}

	jww.DEBUG.Printf("supported loaded: %v", s)
	return nil
}

// Struct to hold the builds.
type builds struct {
	Build  map[string]rawTemplate
	load   sync.Once
	loaded bool
}

// Ensures that the build information only get loaded once. Uses a mutex to
// prevent race conditions as there can be concurrent processing of Packer
// templates. When loaded, it sets the loaded boolean so that it only needs to
// be called when it hasn't been loaded.
func (b *builds) LoadOnce() error {
	var err error

	loadFunc := func() {
		name := os.Getenv(EnvBuildsFile)

		if name == "" {
			err = errors.New("could not retrieve the Builds configurations because the " + EnvBuildsFile + " environment variable was not set. Either set it or check your rancher.cfg setting")
			jww.CRITICAL.Print(err.Error())
			return
		}

		if _, err = toml.DecodeFile(name, &b); err != nil {
			jww.CRITICAL.Print(err.Error())
			return
		}

		b.loaded = true
		return
	}

	b.load.Do(loadFunc)

	// Don't need to log the error as loadFunc already handled that.
	if err != nil {
		return err
	}

	jww.DEBUG.Printf("builds loaded: %v", b)
	return nil
}

// Contains lists of builds.
type buildLists struct {
	List map[string]list
}

// A list contains 1 or more builds.
type list struct {
	Builds []string
}

// This is a normal load, no mutex, as this is only called once.
func (b *buildLists) Load() error {
	// Load the build lists.
	name := os.Getenv(EnvBuildListsFile)

	if name == "" {
		err := errors.New("could not retrieve the BuildLists file because the " + EnvBuildListsFile + " environment variable was not set. Either set it or check your rancher.cfg setting")
		jww.ERROR.Print(err.Error())
		return err
	}

	if _, err := toml.DecodeFile(name, &b); err != nil {
		jww.ERROR.Print(err.Error())
		return err
	}

	jww.DEBUG.Printf("buildLists loaded: %v", b)
	return nil
}
