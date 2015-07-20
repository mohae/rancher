package app

import (
	"strings"
	"testing"
	"time"
)

var testRawTpl = newRawTemplate()

var updatedBuilders = map[string]builder{
	"common": {
		templateSection{
			Settings: []string{
				"ssh_wait_timeout = 300m",
			},
		},
	},
	"virtualbox-iso": {
		templateSection{
			Settings: []string{},
			Arrays: map[string]interface{}{
				"vm_settings": []string{
					"memory=4096",
				},
			},
		},
	},
}

var comparePostProcessors = map[string]postProcessor{
	"vagrant": {
		templateSection{
			Settings: []string{
				"output = :out_dir/packer.box",
			},
			Arrays: map[string]interface{}{
				"except": []string{
					"docker",
				},
				"only": []string{
					"virtualbox-iso",
				},
			},
		},
	},
	"vagrant-cloud": {
		templateSection{
			Settings: []string{
				"access_token = getAValidTokenFrom-VagrantCloud.com",
				"box_tag = foo/bar/baz",
				"no_release = false",
				"version = 1.0.2",
			},
			Arrays: map[string]interface{}{},
		},
	},
}

var compareProvisioners = map[string]provisioner{
	"shell": {
		templateSection{
			Settings: []string{
				"execute_command = execute_test.command",
			},
			Arrays: map[string]interface{}{
				"scripts": []string{
					"setup_test.sh",
					"vagrant_test.sh",
					"cleanup_test.sh",
				},
				"except": []string{
					"docker",
				},
				"only": []string{
					"virtualbox-iso",
				},
			},
		},
	},
}

var testBuildNewTPL = &rawTemplate{
	PackerInf: PackerInf{
		Description: "Test build new template",
	},
	Distro:  "ubuntu",
	Arch:    "amd64",
	Image:   "server",
	Release: "12.04",
	varVals: map[string]string{},
	dirs:    map[string]string{},
	files:   map[string]string{},
	build: build{
		BuilderTypes: []string{
			"virtualbox-iso",
		},
		Builders: map[string]builder{
			"common": {
				templateSection{
					Settings: []string{
						"ssh_wait_timeout = 300m",
					},
				},
			},
			"virtualbox-iso": {
				templateSection{
					Arrays: map[string]interface{}{
						"vm_settings": []string{
							"memory=4096",
						},
					},
				},
			},
		},
		PostProcessorTypes: []string{
			"vagrant",
			"vagrant-cloud",
		},
		PostProcessors: map[string]postProcessor{
			"vagrant": {
				templateSection{
					Settings: []string{
						"output = :out_dir/packer.box",
					},
					Arrays: map[string]interface{}{
						"except": []string{
							"docker",
						},
						"only": []string{
							"virtualbox-iso",
						},
					},
				},
			},
			"vagrant-cloud": {
				templateSection{
					Settings: []string{
						"access_token = getAValidTokenFrom-VagrantCloud.com",
						"box_tag = foo/bar/baz",
						"no_release = false",
						"version = 1.0.2",
					},
				},
			},
		},
		ProvisionerTypes: []string{
			"shell",
		},
		Provisioners: map[string]provisioner{
			"shell": {
				templateSection{
					Settings: []string{
						"execute_command = execute_test.command",
					},
					Arrays: map[string]interface{}{
						"scripts": []string{
							"setup_test.sh",
							"vagrant_test.sh",
							"cleanup_test.sh",
						},
						"except": []string{
							"docker",
						},
						"only": []string{
							"virtualbox-iso",
						},
					},
				},
			},
		},
	},
}

var testRawTemplateBuilderOnly = &rawTemplate{
	PackerInf: PackerInf{MinPackerVersion: "0.4.0", Description: "Test supported distribution template"},
	IODirInf: IODirInf{
		OutputDir: "../test_files/out/:distro/:build_name",
		SourceDir: "../test_files/src/:distro",
	},
	BuildInf: BuildInf{
		Name:      ":build_name",
		BuildName: "",
		BaseURL:   "http://releases.ubuntu.org/",
	},
	date:    today,
	delim:   ":",
	Distro:  "ubuntu",
	Arch:    "amd64",
	Image:   "server",
	Release: "12.04",
	varVals: map[string]string{},
	dirs:    map[string]string{},
	files:   map[string]string{},
	build:   build{},
}

var testRawTemplateWOSection = &rawTemplate{
	PackerInf: PackerInf{MinPackerVersion: "0.4.0", Description: "Test supported distribution template"},
	IODirInf: IODirInf{
		OutputDir: "../test_files/out/:distro/:build_name",
		SourceDir: "../test_files/src/:distro",
	},
	BuildInf: BuildInf{
		Name:      ":build_name",
		BuildName: "",
		BaseURL:   "http://releases.ubuntu.org/",
	},
	date:    today,
	delim:   ":",
	Distro:  "ubuntu",
	Arch:    "amd64",
	Image:   "server",
	Release: "12.04",
	varVals: map[string]string{},
	dirs:    map[string]string{},
	files:   map[string]string{},
	build: build{
		BuilderTypes:       []string{"amazon-ebs"},
		Builders:           map[string]builder{},
		PostProcessorTypes: []string{"compress"},
		PostProcessors:     map[string]postProcessor{},
		ProvisionerTypes:   []string{"ansible-local"},
		Provisioners:       map[string]provisioner{},
	},
}

func TestNewRawTemplate(t *testing.T) {
	rawTpl := newRawTemplate()
	if MarshalJSONToString.Get(rawTpl) != MarshalJSONToString.Get(testRawTpl) {
		t.Errorf("Expected %q, got %q", MarshalJSONToString.Get(testRawTpl), MarshalJSONToString.Get(rawTpl))
	}
}

func TestReplaceVariables(t *testing.T) {
	r := newRawTemplate()
	r.varVals = map[string]string{
		":arch":            "amd64",
		":command_src_dir": "commands",
		":image":           "server",
		":name":            ":distro-:release:-:image-:arch",
		":out_dir":         "../test_files/out/:distro",
		":release":         "14.04",
		":src_dir":         "../test_files/src/:distro",
		":distro":          "ubuntu",
	}
	r.delim = ":"
	s := r.replaceVariables("../test_files/src/:distro")
	if s != "../test_files/src/ubuntu" {
		t.Errorf("Expected \"../test_files/src/ubuntu\", got %q", s)
	}
	s = r.replaceVariables("../test_files/src/:distro/command")
	if s != "../test_files/src/ubuntu/command" {
		t.Errorf("Expected \"../test_files/src/ubuntu/command\", got %q", s)
	}
	s = r.replaceVariables("http")
	if s != "http" {
		t.Errorf("Expected \"http\", got %q", s)
	}
	s = r.replaceVariables("../test_files/out/:distro")
	if s != "../test_files/out/ubuntu" {
		t.Errorf("Expected \"../test_files/out/ubuntu\", got %q", s)
	}
}

func TestRawTemplateUpdateBuildSettings(t *testing.T) {
	r := newRawTemplate()
	r.setDefaults(testSupportedCentOS)
	r.updateBuildSettings(testBuildNewTPL)
	if MarshalJSONToString.Get(r.IODirInf) != MarshalJSONToString.Get(testSupportedCentOS.IODirInf) {
		t.Errorf("Expected %q, got %q", MarshalJSONToString.Get(testSupportedCentOS.IODirInf), MarshalJSONToString.Get(r.IODirInf))
	}
	if MarshalJSONToString.Get(r.PackerInf) != MarshalJSONToString.Get(testBuildNewTPL.PackerInf) {
		t.Errorf("Expected %q, got %q", MarshalJSONToString.Get(testBuildNewTPL.PackerInf), MarshalJSONToString.Get(r.PackerInf))
	}
	if MarshalJSONToString.Get(r.BuildInf) != MarshalJSONToString.Get(testSupportedCentOS.BuildInf) {
		t.Errorf("Expected %q, got %q", MarshalJSONToString.Get(testSupportedCentOS.BuildInf), MarshalJSONToString.Get(r.BuildInf))
	}
	if MarshalJSONToString.Get(r.BuilderTypes) != MarshalJSONToString.Get(testBuildNewTPL.BuilderTypes) {
		t.Errorf("Expected %q, got %q", MarshalJSONToString.Get(testBuildNewTPL.BuilderTypes), MarshalJSONToString.Get(r.BuilderTypes))
	}
	if MarshalJSONToString.Get(r.PostProcessorTypes) != MarshalJSONToString.Get(testBuildNewTPL.PostProcessorTypes) {
		t.Errorf("Expected %q, got %q", MarshalJSONToString.Get(testBuildNewTPL.PostProcessorTypes), MarshalJSONToString.Get(r.PostProcessorTypes))
	}
	if MarshalJSONToString.Get(r.ProvisionerTypes) != MarshalJSONToString.Get(testBuildNewTPL.ProvisionerTypes) {
		t.Errorf("Expected %q, got %q", MarshalJSONToString.Get(testBuildNewTPL.ProvisionerTypes), MarshalJSONToString.Get(r.ProvisionerTypes))
	}
	if MarshalJSONToString.Get(r.Builders) != MarshalJSONToString.Get(updatedBuilders) {
		t.Errorf("Expected %q, got %q", MarshalJSONToString.Get(updatedBuilders), MarshalJSONToString.Get(r.Builders))
	}
	if MarshalJSONToString.Get(r.PostProcessors) != MarshalJSONToString.Get(comparePostProcessors) {
		t.Errorf("Expected %q, got %q", MarshalJSONToString.Get(comparePostProcessors), MarshalJSONToString.Get(r.PostProcessors))
	}
	if MarshalJSONToString.Get(r.Provisioners) != MarshalJSONToString.Get(compareProvisioners) {
		t.Errorf("Expected %q, got %q", MarshalJSONToString.Get(compareProvisioners), MarshalJSONToString.Get(r.Provisioners))
	}
}

func TestMergeVariables(t *testing.T) {
	r := testDistroDefaults.Templates[Ubuntu]
	r.mergeVariables()
	if r.OutputDir != "../test_files/out/ubuntu" {
		t.Errorf("Expected \"../test_files/out/ubuntu\", got %q", r.OutputDir)
	}
	if r.SourceDir != "../test_files/src/ubuntu" {
		t.Errorf("Expected \"../test_files/src/ubuntu\", got %q", r.SourceDir)
	}
}

func TestPackerInf(t *testing.T) {
	oldPackerInf := PackerInf{MinPackerVersion: "0.40", Description: "test info"}
	newPackerInf := PackerInf{}
	oldPackerInf.update(newPackerInf)
	if oldPackerInf.MinPackerVersion != "0.40" {
		t.Errorf("Expected \"0.40\", got %q", oldPackerInf.MinPackerVersion)
	}
	if oldPackerInf.Description != "test info" {
		t.Errorf("Expected \"test info\", got %q", oldPackerInf.Description)
	}

	oldPackerInf = PackerInf{MinPackerVersion: "0.40", Description: "test info"}
	newPackerInf = PackerInf{MinPackerVersion: "0.50"}
	oldPackerInf.update(newPackerInf)
	if oldPackerInf.MinPackerVersion != "0.50" {
		t.Errorf("Expected \"0.50\", got %q", oldPackerInf.MinPackerVersion)
	}
	if oldPackerInf.Description != "test info" {
		t.Errorf("Expected \"test info\", got %q", oldPackerInf.Description)
	}

	oldPackerInf = PackerInf{MinPackerVersion: "0.40", Description: "test info"}
	newPackerInf = PackerInf{Description: "new test info"}
	oldPackerInf.update(newPackerInf)
	if oldPackerInf.MinPackerVersion != "0.40" {
		t.Errorf("Expected \"0.40\", got %q", oldPackerInf.MinPackerVersion)
	}
	if oldPackerInf.Description != "new test info" {
		t.Errorf("Expected \"new test info\", got %q", oldPackerInf.Description)
	}

	oldPackerInf = PackerInf{MinPackerVersion: "0.40", Description: "test info"}
	newPackerInf = PackerInf{MinPackerVersion: "0.5.1", Description: "updated"}
	oldPackerInf.update(newPackerInf)
	if oldPackerInf.MinPackerVersion != "0.5.1" {
		t.Errorf("Expected \"0.5.1\", got %q", oldPackerInf.MinPackerVersion)
	}
	if oldPackerInf.Description != "updated" {
		t.Errorf("Expected \"updated\", got %q", oldPackerInf.Description)
	}
}

func TestBuildInf(t *testing.T) {
	oldBuildInf := BuildInf{Name: "old Name", BuildName: "old BuildName"}
	newBuildInf := BuildInf{}
	oldBuildInf.update(newBuildInf)
	if oldBuildInf.Name != "old Name" {
		t.Errorf("Expected \"old Name\", got %q", oldBuildInf.Name)
	}
	if oldBuildInf.BuildName != "old BuildName" {
		t.Errorf("Expected \"old BuildName\", got %q", oldBuildInf.BuildName)
		t.Errorf("Expected \"old BuildName\", got %q", oldBuildInf.BuildName)
	}

	newBuildInf.Name = "new Name"
	oldBuildInf.update(newBuildInf)
	if oldBuildInf.Name != "new Name" {
		t.Errorf("Expected \"new Name\", got %q", oldBuildInf.Name)
	}
	if oldBuildInf.BuildName != "old BuildName" {
		t.Errorf("Expected \"old BuildName\", got %q", oldBuildInf.BuildName)
	}

	newBuildInf.BuildName = "new BuildName"
	oldBuildInf.update(newBuildInf)
	if oldBuildInf.Name != "new Name" {
		t.Errorf("Expected \"new Name\", got %q", oldBuildInf.Name)
	}
	if oldBuildInf.BuildName != "new BuildName" {
		t.Errorf("Expected \"new BuildName\", got %q", oldBuildInf.BuildName)
	}
}

func TestRawTemplateMergeSrcDir(t *testing.T) {
	tests := []struct {
		SrcDir         string
		ExpectedSrcDir string
	}{
		{"src/", "src"},
		{"src/custom/", "src/custom"},
		{"src/:distro/", "src/ubuntu"},
		{"src/:distro/", "src/ubuntu"},
		{"src/files/", "src/files"},
	}
	rawTpl := newRawTemplate()
	rawTpl.delim = ":"
	rawTpl.Distro = "ubuntu"
	rawTpl.setBaseVarVals()
	for i, test := range tests {
		rawTpl.SourceDir = test.SrcDir
		rawTpl.mergeSourceDir()
		if rawTpl.SourceDir != test.ExpectedSrcDir {
			t.Errorf("MergeSrcDir test %d: expected SrcDir to be %s; got %s", i, test.ExpectedSrcDir, rawTpl.SourceDir)
		}
	}
}

func TestRawTemplateMergeOutDir(t *testing.T) {
	tests := []struct {
		OutDir         string
		ExpectedOutDir string
	}{
		{"out", "out"},
		{"out/custom/", "out/custom"},
		{"out/:distro/", "out/ubuntu"},
		{"out/:distro/", "out/ubuntu"},
		{"out/files/", "out/files"},
	}
	rawTpl := newRawTemplate()
	rawTpl.delim = ":"
	rawTpl.Distro = "ubuntu"
	rawTpl.setBaseVarVals()
	for i, test := range tests {
		rawTpl.OutputDir = test.OutDir
		rawTpl.mergeOutDir()
		if rawTpl.OutputDir != test.ExpectedOutDir {
			t.Errorf("MergeOutDirtest %d: expected OutDir to be %s; got %s", i, test.ExpectedOutDir, rawTpl.OutputDir)
		}
	}
}

func TestRawTemplateSetBaseVarVals(t *testing.T) {
	now := time.Now()
	splitDate := strings.Split(now.String(), " ")
	tests := []struct {
		Distro    string
		Release   string
		Arch      string
		Image     string
		BuildName string
	}{
		{"ubuntu", "14.04", "amd64", "server", "14.04-test"},
		{"centos", "7", "x86_64", "minimal", "7-test"},
	}

	r := newRawTemplate()
	r.delim = ":"
	for i, test := range tests {
		r.Distro = test.Distro
		r.Release = test.Release
		r.Arch = test.Arch
		r.Image = test.Image
		r.BuildName = test.BuildName
		// make the map empty
		r.varVals = map[string]string{}
		r.setBaseVarVals()
		tmp, ok := r.varVals[":distro"]
		if !ok {
			t.Errorf("%d: expected :distro to be in map, it wasn't", i)
		} else {
			if tmp != test.Distro {
				t.Errorf("%d: expected :distro to be %q, got %q", i, test.Distro, tmp)
			}
		}
		tmp, ok = r.varVals[":release"]
		if !ok {
			t.Errorf("%d: expected :release to be in map, it wasn't", i)
		} else {
			if tmp != test.Release {
				t.Errorf("%d: expected :release to be %q, got %q", i, test.Release, tmp)
			}
		}
		tmp, ok = r.varVals[":arch"]
		if !ok {
			t.Errorf("%d: expected :arch to be in map, it wasn't", i)
		} else {
			if tmp != test.Arch {
				t.Errorf("%d: expected :arch to be %q, got %q", i, test.Arch, tmp)
			}
		}
		tmp, ok = r.varVals[":image"]
		if !ok {
			t.Errorf("%d: expected :image to be in map, it wasn't", i)
		} else {
			if tmp != test.Image {
				t.Errorf("%d: expected :image to be %q, got %q", i, test.Image, tmp)
			}
		}
		tmp, ok = r.varVals[":date"]
		if !ok {
			t.Errorf("%d: expected :date to be in map, it wasn't", i)
		} else {
			if tmp != splitDate[0] {
				t.Errorf("%d: expected :date to be %q, got %q", i, splitDate[0], tmp)
			}
		}
		tmp, ok = r.varVals[":build_name"]
		if !ok {
			t.Errorf("%d: expected :build_name to be in map, it wasn't", i)
		} else {
			if tmp != test.BuildName {
				t.Errorf("%d: expected :build_name to be %q, got %q", i, test.BuildName, tmp)
			}
		}
	}
}

func TestRawTemplateMergeString(t *testing.T) {
	tests := []struct {
		value    string
		dflt     string
		expected string
	}{
		{"", "", ""},
		{"", "src", "src"},
		{"dir", "src", "dir"},
		{"dir/", "src", "dir"},
		{"dir", "", "dir"},
		{"dir/", "", "dir"},
	}
	r := newRawTemplate()
	for i, test := range tests {
		v := r.mergeString(test.value, test.dflt)
		if v != test.expected {
			t.Errorf("mergeString %d: expected %q, got %q", i, test.expected, v)
		}
	}
}

func TestFindSource(t *testing.T) {
	tests := []struct {
		p           string
		isDir       bool
		src         string
		expectedErr string
	}{
		{"", false, "", "cannot find source, no path received"},
		{"something", false, "", "file does not exist"},
		{"http/preseed.cfg", false, "../test_files/src/ubuntu/http/preseed.cfg", ""},
		{"chef-solo/cookbook1", true, "../test_files/src/chef-solo/cookbook1", ""},
		{"14.04_ubuntu_build.txt", false, "../test_files/src/ubuntu/14.04/ubuntu_build/14.04_ubuntu_build.txt", ""},
		{"1404_ubuntu_build.txt", false, "../test_files/src/ubuntu/1404/ubuntu_build/1404_ubuntu_build.txt", ""},
		{"14_ubuntu_build.txt", false, "../test_files/src/ubuntu/14/ubuntu_build/14_ubuntu_build.txt", ""},
		{"ubuntu_build_text.txt", false, "../test_files/src/ubuntu/ubuntu_build/ubuntu_build_text.txt", ""},
		{"ubuntu_build.txt", false, "../test_files/src/ubuntu_build/ubuntu_build.txt", ""},
		{"14.04_amd64_build_text.txt", false, "../test_files/src/ubuntu/14.04/amd64/14.04_amd64_build_text.txt", ""},
		{"1404_amd64_build_text.txt", false, "../test_files/src/ubuntu/1404/amd64/1404_amd64_build_text.txt", ""},
		{"14_amd64_build_text.txt", false, "../test_files/src/ubuntu/14/amd64/14_amd64_build_text.txt", ""},
		{"14.04_text.txt", false, "../test_files/src/ubuntu/14.04/14.04_text.txt", ""},
		{"1404_text.txt", false, "../test_files/src/ubuntu/1404/1404_text.txt", ""},
		{"14_text.txt", false, "../test_files/src/ubuntu/14/14_text.txt", ""},
		{"amd64_text.txt", false, "../test_files/src/ubuntu/amd64/amd64_text.txt", ""},
		{"ubuntu_text.txt", false, "../test_files/src/ubuntu/ubuntu_text.txt", ""},
	}
	r := newRawTemplate()
	r.Distro = "ubuntu"
	r.Arch = "amd64"
	r.Release = "14.04"
	r.Image = "server"
	r.SourceDir = "../test_files/src"
	r.BuildName = "ubuntu_build"
	for i, test := range tests {
		src, err := r.findSource(test.p, test.isDir)
		if err != nil {
			if err.Error() != test.expectedErr {
				t.Errorf("TestFindSource %d: expected %q got %q", i, test.expectedErr, err.Error())
			}
			continue
		}
		if test.expectedErr != "" {
			t.Errorf("TestFindSource %d: expected %q, got no error", i, test.expectedErr)
			continue
		}
		if test.src != src {
			t.Errorf("TestFindSource %d: expected %q, got %q", i, test.src, src)
		}
	}
}

func TestFindComponentSource(t *testing.T) {
	tests := []struct {
		component   string
		p           string
		isDir       bool
		src         string
		expectedErr string
	}{
		{"", "", false, "", "cannot find source, no path received"},
		{"", "chef.cfg", false, "", " file \"chef.cfg\": file does not exist"},
		{"salt", "minion", false, "", "salt file \"minion\": file does not exist"},
		{"salt-masterless", "master", false, "", "salt-masterless file \"master\": file does not exist"},
		{"chef-solo", "chef.cfg", false, "../test_files/src/chef-solo/chef.cfg", ""},
		{"chef-client", "chef.cfg", false, "../test_files/src/chef-client/chef.cfg", ""},
		{"chef", "chef.cfg", false, "../test_files/src/chef/chef.cfg", ""},
		{"shell", "commands", true, "../test_files/src/shell/commands", ""},
		{"", "ubuntu_build.txt", false, "../test_files/src/ubuntu_build/ubuntu_build.txt", ""},
	}
	r := newRawTemplate()
	r.Distro = "ubuntu"
	r.Arch = "amd64"
	r.Release = "14.04"
	r.Image = "server"
	r.SourceDir = "../test_files/src"
	r.BuildName = "ubuntu_build"
	for i, test := range tests {
		src, err := r.findComponentSource(test.component, test.p, test.isDir)
		if err != nil {
			if err.Error() != test.expectedErr {
				t.Errorf("TestFindSource %d: expected %q got %q", i, test.expectedErr, err.Error())
			}
			continue
		}
		if test.expectedErr != "" {
			t.Errorf("TestFindSource %d: expected %q, got no error", i, test.expectedErr)
			continue
		}
		if test.src != src {
			t.Errorf("TestFindSource %d: expected %q, got %q", i, test.src, src)
		}
	}
}

func TestFindCommandFile(t *testing.T) {
	tests := []struct {
		component   string
		p           string
		src         string
		expectedErr string
	}{
		{"", "", "", "the passed command filename was empty"},
		{"", "test.command", "", " file \"commands/test.command\": file does not exist"},
		{"", "execute.command", "../test_files/src/commands/execute.command", ""},
		{"shell", "execute_test.command", "../test_files/src/shell/commands/execute_test.command", ""},
		{"chef-solo", "execute.command", "../test_files/src/chef-solo/commands/execute.command", ""},
		{"chef-solo", "chef.command", "../test_files/src/chef/commands/chef.command", ""},
		{"shell", "ubuntu.command", "../test_files/src/ubuntu/commands/ubuntu.command", ""},
		{"shell", "ubuntu-14.command", "../test_files/src/ubuntu/14/commands/ubuntu-14.command", ""},
	}
	r := newRawTemplate()
	r.Distro = "ubuntu"
	r.Arch = "amd64"
	r.Release = "14.04"
	r.Image = "server"
	r.SourceDir = "../test_files/src"
	r.BuildName = "ubuntu_build"
	for i, test := range tests {
		src, err := r.findCommandFile(test.component, test.p)
		if err != nil {
			if err.Error() != test.expectedErr {
				t.Errorf("TestFindCommandFile %d: expected %q got %q", i, test.expectedErr, err.Error())
			}
			continue
		}
		if test.expectedErr != "" {
			t.Errorf("TestFindCommandFile %d: expected %q, got no error", i, test.expectedErr)
			continue
		}
		if test.src != src {
			t.Errorf("TestFindCommandFile %d: expected %q, got %q", i, test.src, src)
		}
	}
}

func TestCommandsFromFile(t *testing.T) {
	tests := []struct {
		component   string
		p           string
		expected    []string
		expectedErr string
	}{
		{"", "", []string{}, "the passed command filename was empty"},
		{"", "test.command", []string{}, " file \"commands/test.command\": file does not exist"},
		{"shell", "execute.command", []string{"echo 'vagrant'|sudo -S sh '{{.Path}}'"}, ""},
		{"shell", "boot.command", []string{"<esc><wait>", "<esc><wait>", "<enter><wait>"}, ""},
	}
	r := newRawTemplate()
	r.Distro = "ubuntu"
	r.Arch = "amd64"
	r.Release = "14.04"
	r.Image = "server"
	r.SourceDir = "../test_files/src"
	r.BuildName = "ubuntu_build"
	for i, test := range tests {
		commands, err := r.commandsFromFile(test.component, test.p)
		if err != nil {
			if err.Error() != test.expectedErr {
				t.Errorf("TestCommandsFromFile %d: expected %q got %q", i, test.expectedErr, err.Error())
			}
			continue
		}
		if test.expectedErr != "" {
			t.Errorf("TestCommandsFromFile %d: expected %q, got no error", i, test.expectedErr)
			continue
		}
		if len(commands) != len(test.expected) {
			t.Errorf("TestCommandsFromFile %d: expected commands slice to have a len of %d got %d", i, len(test.expected), len(commands))
			continue
		}
		for i, v := range commands {
			if v != test.expected[i] {
				t.Errorf("TestCommandsFromFile %d: expected commands slice to be %v, got %v", i, test.expected, commands)
				break
			}
		}
	}

}

func TestBuildOutPath(t *testing.T) {
	tests := []struct {
		includeComponent string
		component        string
		path             string
		expected         string
	}{
		{"false", "", "", "out"},
		{"true", "", "", "out"},
		{"false", "vagrant", "", "out"},
		{"true", "vagrant", "", "out/vagrant"},
		{"false", "", "file.txt", "out/file.txt"},
		{"false", "", "path/to/file.txt", "out/path/to/file.txt"},
		{"false", "shell", "file.txt", "out/file.txt"},
		{"false", "shell", "path/to/file.txt", "out/path/to/file.txt"},
		{"true", "", "file.txt", "out/file.txt"},
		{"true", "", "path/to/file.txt", "out/path/to/file.txt"},
		{"true", "shell", "file.txt", "out/shell/file.txt"},
		{"true", "shell", "path/to/file.txt", "out/shell/path/to/file.txt"},
	}
	r := newRawTemplate()
	r.OutputDir = "out"
	for i, test := range tests {
		r.IncludeComponentString = test.includeComponent
		p := r.buildOutPath(test.component, test.path)
		if p != test.expected {
			t.Errorf("TestBuildOutPath %d: expected %q, got %q", i, test.expected, p)
		}
	}
}

func TestBuildTemplateResourcePath(t *testing.T) {
	tests := []struct {
		includeComponent string
		component        string
		path             string
		expected         string
	}{
		{"false", "", "", ""},
		{"true", "", "", ""},
		{"false", "vagrant", "", ""},
		{"true", "vagrant", "", "vagrant"},
		{"false", "", "file.txt", "file.txt"},
		{"false", "", "path/to/file.txt", "path/to/file.txt"},
		{"false", "shell", "file.txt", "file.txt"},
		{"false", "shell", "path/to/file.txt", "path/to/file.txt"},
		{"true", "", "file.txt", "file.txt"},
		{"true", "", "path/to/file.txt", "path/to/file.txt"},
		{"true", "shell", "file.txt", "shell/file.txt"},
		{"true", "shell", "path/to/file.txt", "shell/path/to/file.txt"},
	}
	r := newRawTemplate()
	r.OutputDir = "out"
	for i, test := range tests {
		r.IncludeComponentString = test.includeComponent
		p := r.buildTemplateResourcePath(test.component, test.path)
		if p != test.expected {
			t.Errorf("TestBuildTemplateResourcePath %d: expected %q, got %q", i, test.expected, p)
		}
	}

}
