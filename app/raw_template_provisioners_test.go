// raw_template_provisioners_test.go: tests for provisioners.
package app

import (
	"testing"
)

var testRawTemplateProvisioner = &rawTemplate{
	PackerInf: PackerInf{
		MinPackerVersion: "0.4.0",
		Description:      "Test template config and Rancher options for CentOS",
	},
	IODirInf: IODirInf{
		OutputDir: "../test_files/out/:build_name",
		SourceDir: "../test_files/src",
	},
	BuildInf: BuildInf{
		Name:      ":build_name",
		BuildName: "",
		BaseURL:   "",
	},
	date:    today,
	delim:   ":",
	Distro:  "centos",
	Arch:    "x86_64",
	Image:   "minimal",
	Release: "6",
	varVals: map[string]string{},
	dirs:    map[string]string{},
	files:   map[string]string{},
	build: build{
		BuilderIDs: []string{"virtualbox-iso", "vmware-iso"},
		Builders: map[string]builder{
			"common": {
				templateSection{
					Type: "common",
					Settings: []string{
						"boot_command = boot_test.command",
						"boot_wait = 5s",
						"disk_size = 20000",
						"guest_os_type = ",
						"headless = true",
						"http_directory = http",
						"iso_checksum_type = sha256",
						"shutdown_command = shutdown_test.command",
						"ssh_password = vagrant",
						"ssh_port = 22",
						"ssh_username = vagrant",
						"ssh_wait_timeout = 240m",
					},
				},
			},
			"virtualbox-iso": {
				templateSection{
					Type: "virtualbox-iso",
					Settings: []string{
						"virtualbox_version_file = .vbox_version",
					},
					Arrays: map[string]interface{}{
						"vm_settings": []string{
							"cpus=1",
							"memory=1024",
						},
					},
				},
			},
			"vmware-iso": {
				templateSection{
					Type:     "vmware-iso",
					Settings: []string{},
					Arrays: map[string]interface{}{
						"vm_settings": []string{
							"cpuid.coresPerSocket=1",
							"memsize=1024",
							"numvcpus=1",
						},
					},
				},
			},
		},
		PostProcessorIDs: []string{
			"vagrant",
			"vagrant-cloud",
		},
		PostProcessors: map[string]postProcessor{
			"vagrant": {
				templateSection{
					Type: "vagrant",
					Settings: []string{
						"compression_level = 9",
						"keep_input_artifact = false",
						"output = out/rancher-packer.box",
					},
					Arrays: map[string]interface{}{
						"include": []string{
							"include1",
							"include2",
						},
						"only": []string{
							"virtualbox-iso",
						},
					},
				},
			},
			"vagrant-cloud": {
				templateSection{
					Type: "vagrant-cloud",
					Settings: []string{
						"access_token = getAValidTokenFrom-VagrantCloud.com",
						"box_tag = foo/bar",
						"no_release = true",
						"version = 1.0.1",
					},
				},
			},
		},
		ProvisionerIDs: []string{
			"shell-test",
			"file",
		},
		Provisioners: map[string]provisioner{
			"shell-test": {
				templateSection{
					Type: "shell",
					Settings: []string{
						"execute_command = execute_test.command",
					},
					Arrays: map[string]interface{}{
						"except": []string{
							"docker",
						},
						"only": []string{
							"virtualbox-iso",
						},
						"scripts": []string{
							"setup_test.sh",
							"vagrant_test.sh",
							"sudoers_test.sh",
							"cleanup_test.sh",
						},
					},
				},
			},
			"file": {
				templateSection{
					Type: "file",
					Settings: []string{
						"source = app.tar.gz",
						"destination = /tmp/app.tar.gz",
					},
					Arrays: map[string]interface{}{},
				},
			},
		},
	},
}

var testRawTemplateProvisionersAll = &rawTemplate{
	PackerInf: PackerInf{
		MinPackerVersion: "0.4.0",
		Description:      "Test template config and Rancher options for CentOS",
	},
	IODirInf: IODirInf{
		OutputDir: "../test_files/out/:build_name",
		SourceDir: "../test_files/src",
	},
	BuildInf: BuildInf{
		Name:      ":build_name",
		BuildName: "",
		BaseURL:   "",
	},
	date:    today,
	delim:   ":",
	Distro:  "centos",
	Arch:    "x86_64",
	Image:   "minimal",
	Release: "6",
	varVals: map[string]string{},
	dirs:    map[string]string{},
	files:   map[string]string{},
	build: build{
		BuilderIDs: []string{"virtualbox-iso", "vmware-iso"},
		Builders: map[string]builder{
			"common": {
				templateSection{
					Type: "common",
					Settings: []string{
						"boot_command = boot_test.command",
						"boot_wait = 5s",
						"disk_size = 20000",
						"guest_os_type = ",
						"headless = true",
						"http_directory = http",
						"iso_checksum_type = sha256",
						"shutdown_command = shutdown_test.command",
						"ssh_password = vagrant",
						"ssh_port = 22",
						"ssh_username = vagrant",
						"ssh_wait_timeout = 240m",
					},
				},
			},
			"virtualbox-iso": {
				templateSection{
					Type: "virtualbox-iso",
					Settings: []string{
						"virtualbox_version_file = .vbox_version",
					},
					Arrays: map[string]interface{}{
						"vm_settings": []string{
							"cpus=1",
							"memory=1024",
						},
					},
				},
			},
			"vmware-iso": {
				templateSection{
					Type:     "vmware-iso",
					Settings: []string{},
					Arrays: map[string]interface{}{
						"vm_settings": []string{
							"cpuid.coresPerSocket=1",
							"memsize=1024",
							"numvcpus=1",
						},
					},
				},
			},
		},
		PostProcessorIDs: []string{
			"vagrant",
			"vagrant-cloud",
		},
		PostProcessors: map[string]postProcessor{
			"vagrant": {
				templateSection{
					Type: "vagrant",
					Settings: []string{
						"compression_level = 9",
						"keep_input_artifact = false",
						"output = out/rancher-packer.box",
					},
					Arrays: map[string]interface{}{
						"include": []string{
							"include1",
							"include2",
						},
						"only": []string{
							"virtualbox-iso",
						},
					},
				},
			},
			"vagrant-cloud": {
				templateSection{
					Type: "vagrant-cloud",
					Settings: []string{
						"access_token = getAValidTokenFrom-VagrantCloud.com",
						"box_tag = foo/bar",
						"no_release = true",
						"version = 1.0.1",
					},
				},
			},
		},
		ProvisionerIDs: []string{
			"ansible-local",
			"file",
			"chef-client",
			"chef-solo",
			"puppet-client",
			"salt-masterless",
			"shell",
		},
		Provisioners: map[string]provisioner{
			"ansible-local": {
				templateSection{
					Type: "ansible-local",
					Settings: []string{
						"playbook_file= playbook.yml",
						"command =  ansible_test.command",
						"inventory_file = inventory_file",
						"group_vars = groupvars",
						"host_vars = hostvars",
						"playbook_dir = playbooks",
						"staging_directory = staging/directory",
					},
					Arrays: map[string]interface{}{
						"extra_arguments": []string{
							"arg1",
							"arg2",
						},
						"playbook_paths": []string{
							"playbook1",
							"playbook2",
						},
						"role_paths": []string{
							"roles1",
							"roles2",
						},
					},
				},
			},
			"chef-client": {
				templateSection{
					Type: "chef-client",
					Settings: []string{
						"chef_environment=web",
						"config_template=chef.cfg",
						"execute_command=execute.command",
						"install_command=install.command",
						"node_name=test-chef",
						"prevent_sudo=false",
						"server_url=https://mychefserver.com",
						"skip_clean_client=true",
						"skip_clean_node=false",
						"skip_install=false",
						"staging_directory=/tmp/chef/",
						"validation_client_name=some_value",
						"validation_key_path=/home/user/chef/chef-key",
					},
					Arrays: map[string]interface{}{
						"run_list": []string{
							"recipe[hello::default]",
							"recipe[world::default]",
						},
					},
				},
			},
			"chef-solo": {
				templateSection{
					Type: "chef-solo",
					Settings: []string{
						"config_template=chef.cfg",
						"data_bags_path=data_bag",
						"encrypted_data_bag_secret_path=/home/user/chef/secret_data_bag",
						"environments_path=environments",
						"execute_command=execute.command",
						"install_command=install.command",
						"prevent_sudo=false",
						"roles_path=roles",
						"skip_install=false",
						"staging_directory=/tmp/chef/",
					},
					Arrays: map[string]interface{}{
						"cookbook_paths": []string{
							"cookbook1",
							"cookbook2",
						},
						"remote_cookbook_paths": []string{
							"remote/path1",
							"remote/path2",
						},
						"run_list": []string{
							"recipe[hello::default]",
							"recipe[world::default]",
						},
					},
				},
			},
			"puppet-masterless": {
				templateSection{
					Type: "puppet-masterless",
					Settings: []string{
						"manifest_filet=site.pp",
						"execute_command=execute.command",
						"hiera_config_path=hiera.yaml",
						"manifest_dir=manifests",
						"manifest_file=site.pp",
						"prevent_sudo=false",
						"staging_directory=/tmp/puppet-masterless",
					},
					Arrays: map[string]interface{}{
						"facter": map[string]string{
							"server_role": "webserver",
						},
						"module_paths": []string{
							"/etc/puppetlabs/puppet/modules",
							"/opt/puppet/share/puppet/modules",
						},
					},
				},
			},
			"puppet-server": {
				templateSection{
					Type: "puppet-server",
					Settings: []string{
						"client_cert_path = /etc/puppet/client.pem",
						"client_private_key_path=/home/puppet/.ssh/puppet_id_rsa",
						"options=-v --detailed-exitcodes",
						"prevent_sudo= false",
						"puppet_node=vagrant-puppet-srv01",
						"puppet_server=server",
						"staging_directory=/tmp/puppet-server",
					},
					Arrays: map[string]interface{}{
						"facter": map[string]string{
							"server_role": "webserver",
						},
					},
				},
			},
			"salt-masterless": {
				templateSection{
					Type: "salt-masterless",
					Settings: []string{
						"bootstrap_args = args",
						"local_pillar_roots=pillar",
						"local_state_tree=salt",
						"minion_config=salt",
						"skip_bootstrap=false",
						"temp_config_dir=/tmp",
					},
				},
			},
			"shell": {
				templateSection{
					Type: "shell",
					Settings: []string{
						"binary = false",
						"execute_command = execute_test.command",
						"inline_shebang = /bin/sh",
						"remote_path = /tmp/script.sh",
						"start_retry_timeout = 5m",
					},
					Arrays: map[string]interface{}{
						"except": []string{
							"docker",
						},
						"only": []string{
							"virtualbox-iso",
						},
						"scripts": []string{
							"setup_test.sh",
							"vagrant_test.sh",
							"sudoers_test.sh",
							"cleanup_test.sh",
						},
					},
				},
			},
			"file": {
				templateSection{
					Type: "file",
					Settings: []string{
						"source = app.tar.gz",
						"destination = /tmp/app.tar.gz",
					},
				},
			},
		},
	},
}

var pr = &provisioner{
	templateSection{
		Settings: []string{
			"execute_command= echo 'vagrant' | sudo -S sh '{{.Path}}'",
			"type = shell",
		},
		Arrays: map[string]interface{}{
			"override": map[string]interface{}{
				"virtualbox-iso": map[string]interface{}{
					"scripts": []string{
						"base.sh",
						"vagrant.sh",
						"virtualbox.sh",
						"cleanup.sh",
					},
				},
				"scripts": []string{
					"base.sh",
					"vagrant.sh",
					"cleanup.sh",
				},
			},
		},
	},
}

var prOrig = map[string]provisioner{
	"shell-test": provisioner{
		templateSection{
			Type: "shell",
			Settings: []string{
				"execute_command = execute_test.command",
			},
			Arrays: map[string]interface{}{
				"except": []string{
					"docker",
				},
				"only": []string{
					"virtualbox-iso",
				},
				"scripts": []string{
					"setup_test.sh",
					"vagrant_test.sh",
					"sudoers_test.sh",
					"cleanup_test.sh",
				},
			},
		},
	},
	"file": {
		templateSection{
			Type: "file",
			Settings: []string{
				"source = app.tar.gz",
				"destination = /tmp/app.tar.gz",
			},
			Arrays: map[string]interface{}{},
		},
	},
}

var prNew = map[string]provisioner{
	"shell-test": provisioner{
		templateSection{
			Type:     "shell",
			Settings: []string{},
			Arrays: map[string]interface{}{
				"only": []string{
					"vmware-iso",
				},
				"except": []string{
					"digitalocean",
				},
				"override": map[string]interface{}{
					"vmware-iso": map[string]interface{}{
						"scripts": []string{
							"setup_test.sh",
							"vagrant_test.sh",
							"vmware_test.sh",
							"cleanup_test.sh",
						},
					},
				},
				"scripts": []string{
					"setup_test.sh",
					"vagrant_test.sh",
					"sudoers_test.sh",
					"cleanup_test.sh",
				},
			},
		},
	},
}

var prMerged = map[string]provisioner{
	"shell-test": provisioner{
		templateSection{
			Type: "shell",
			Settings: []string{
				"execute_command = execute_test.command",
			},
			Arrays: map[string]interface{}{
				"except": []string{
					"digitalocean",
				},
				"only": []string{
					"vmware-iso",
				},
				"override": map[string]interface{}{
					"vmware-iso": map[string]interface{}{
						"scripts": []string{
							"setup_test.sh",
							"vagrant_test.sh",
							"vmware_test.sh",
							"cleanup_test.sh",
						},
					},
				},
				"scripts": []string{
					"setup_test.sh",
					"vagrant_test.sh",
					"sudoers_test.sh",
					"cleanup_test.sh",
				},
			},
		},
	},
	"file": {
		templateSection{
			Type: "file",
			Settings: []string{
				"source = app.tar.gz",
				"destination = /tmp/app.tar.gz",
			},
			Arrays: map[string]interface{}{},
		},
	},
}

func init() {
	b := true
	testRawTemplateProvisionersAll.IncludeComponentString = &b
}

func TestRawTemplateUpdateProvisioners(t *testing.T) {
	err := testRawTemplateProvisioner.updateProvisioners(nil)
	if err != nil {
		t.Errorf("expected error to be nil, got %q", err)
	}
	if MarshalJSONToString.Get(testRawTemplateProvisioner.Provisioners) != MarshalJSONToString.Get(prOrig) {
		t.Errorf("Got %q, want %q", MarshalJSONToString.Get(testRawTemplateProvisioner.Provisioners), MarshalJSONToString.Get(prOrig))
	}

	err = testRawTemplateProvisioner.updateProvisioners(prNew)
	if err != nil {
		t.Errorf("expected error to be nil, got %q", err)
	}
	if MarshalJSONToString.GetIndented(testRawTemplateProvisioner.Provisioners) != MarshalJSONToString.GetIndented(prMerged) {
		t.Errorf("Got %q, want %q", MarshalJSONToString.Get(testRawTemplateProvisioner.Provisioners), MarshalJSONToString.Get(prMerged))
	}
}

func TestProvisionersSettingsToMap(t *testing.T) {
	res := pr.settingsToMap("shell", testRawTpl)
	compare := map[string]interface{}{"type": "shell", "execute_command": "echo 'vagrant' | sudo -S sh '{{.Path}}'"}
	for k, v := range res {
		val, ok := compare[k]
		if !ok {
			t.Errorf("Expected to find entry for Key %s, none found", k)
			continue
		}
		if val != v {
			t.Errorf("Got %q, want %q", v, val)
		}
	}
}

func TestAnsibleProvisioner(t *testing.T) {
	expected := map[string]interface{}{
		"command": "ansible_test.command",
		"extra_arguments": []string{
			"arg1",
			"arg2",
		},
		"group_vars":     "ansible-local/groupvars",
		"host_vars":      "ansible-local/hostvars",
		"inventory_file": "ansible-local/inventory_file",
		"playbook_dir":   "ansible-local/playbooks",
		"playbook_file":  "ansible-local/playbook.yml",
		"playbook_paths": []string{
			"ansible-local/playbook1",
			"ansible-local/playbook2",
		},
		"role_paths": []string{
			"ansible-local/roles1",
			"ansible-local/roles2",
		},
		"staging_directory": "staging/directory",
		"type":              "ansible-local",
	}
	settings, err := testRawTemplateProvisionersAll.createAnsible("ansible-local")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if MarshalJSONToString.Get(settings) != MarshalJSONToString.Get(expected) {
			t.Errorf("Expected %q, got %q", MarshalJSONToString.Get(expected), MarshalJSONToString.Get(settings))
		}
	}
}

func TestChefClientProvisioner(t *testing.T) {
	expected := map[string]interface{}{
		"chef_environment": "web",
		"config_template":  "chef-client/chef.cfg",
		"execute_command":  "{{if .Sudo}}sudo {{end}}chef-client --no-color -c {{.ConfigPath}} -j {{.JsonPath}}",
		"install_command":  "curl -L https://www.opscode.com/chef/install.sh | {{if .Sudo}}sudo{{end}} bash",
		"node_name":        "test-chef",
		"prevent_sudo":     false,
		"run_list": []string{
			"recipe[hello::default]",
			"recipe[world::default]",
		},
		"server_url":        "https://mychefserver.com",
		"skip_clean_client": true,
		"skip_clean_node":   false,
		"skip_install":      false,
		"staging_directory": "/tmp/chef/",
		"type":              "chef-client",
		"validation_client_name": "some_value",
		"validation_key_path":    "/home/user/chef/chef-key",
	}
	settings, err := testRawTemplateProvisionersAll.createChefClient("chef-client")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if MarshalJSONToString.Get(settings) != MarshalJSONToString.Get(expected) {
			t.Errorf("Expected %q, got %q", MarshalJSONToString.Get(expected), MarshalJSONToString.Get(settings))
		}
	}
}

func TestChefSoloProvisioner(t *testing.T) {
	expected := map[string]interface{}{
		"config_template": "chef-solo/chef.cfg",
		"cookbook_paths": []string{
			"chef-solo/cookbook1",
			"chef-solo/cookbook2",
		},
		"data_bags_path":                 "chef-solo/data_bag",
		"encrypted_data_bag_secret_path": "/home/user/chef/secret_data_bag",
		"environments_path":              "chef-solo/environments",
		"execute_command":                "{{if .Sudo}}sudo {{end}}chef-client --no-color -c {{.ConfigPath}} -j {{.JsonPath}}",
		"install_command":                "curl -L https://www.opscode.com/chef/install.sh | {{if .Sudo}}sudo{{end}} bash",
		"prevent_sudo":                   false,
		"roles_path":                     "chef-solo/roles",
		"remote_cookbook_paths": []string{
			"remote/path1",
			"remote/path2",
		},
		"run_list": []string{
			"recipe[hello::default]",
			"recipe[world::default]",
		},
		"skip_install":      false,
		"staging_directory": "/tmp/chef/",
		"type":              "chef-solo",
	}
	settings, err := testRawTemplateProvisionersAll.createChefSolo("chef-solo")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if MarshalJSONToString.Get(settings) != MarshalJSONToString.Get(expected) {
			t.Errorf("Expected %q, got %q", MarshalJSONToString.Get(expected), MarshalJSONToString.Get(settings))
		}
	}
}

func TestPuppetMasterlessProvisioner(t *testing.T) {
	expected := map[string]interface{}{
		"execute_command":   "echo 'vagrant'|sudo -S sh '{{.Path}}'",
		"hiera_config_path": "puppet-masterless/hiera.yaml",
		"facter": map[string]string{
			"server_role": "webserver",
		},
		"manifest_dir":  "puppet-masterless/manifests",
		"manifest_file": "puppet-masterless/site.pp",
		"module_paths": []string{
			"/etc/puppetlabs/puppet/modules",
			"/opt/puppet/share/puppet/modules",
		},
		"prevent_sudo":      false,
		"staging_directory": "/tmp/puppet-masterless",
		"type":              "puppet-masterless",
	}
	settings, err := testRawTemplateProvisionersAll.createPuppetMasterless("puppet-masterless")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if MarshalJSONToString.Get(settings) != MarshalJSONToString.Get(expected) {
			t.Errorf("Expected %q, got %q", MarshalJSONToString.Get(expected), MarshalJSONToString.Get(settings))
		}
	}
}

func TestPuppetServerProvisioner(t *testing.T) {
	expected := map[string]interface{}{
		"client_cert_path":        "/etc/puppet/client.pem",
		"client_private_key_path": "/home/puppet/.ssh/puppet_id_rsa",
		"facter": map[string]string{
			"server_role": "webserver",
		},
		"options":           "-v --detailed-exitcodes",
		"prevent_sudo":      false,
		"puppet_node":       "vagrant-puppet-srv01",
		"puppet_server":     "server",
		"staging_directory": "/tmp/puppet-server",
		"type":              "puppet-server",
	}
	settings, err := testRawTemplateProvisionersAll.createPuppetServer("puppet-server")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if MarshalJSONToString.Get(settings) != MarshalJSONToString.Get(expected) {
			t.Errorf("Expected %q, got %q", MarshalJSONToString.Get(expected), MarshalJSONToString.Get(settings))
		}
	}
}

func TestSaltProvisioner(t *testing.T) {
	expected := map[string]interface{}{
		"bootstrap_args":     "args",
		"local_pillar_roots": "salt-masterless/pillar",
		"local_state_tree":   "salt-masterless/salt",
		"minion_config":      "salt-masterless/salt",
		"skip_bootstrap":     false,
		"temp_config_dir":    "/tmp",
		"type":               "salt-masterless",
	}
	settings, err := testRawTemplateProvisionersAll.createSalt("salt-masterless")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if MarshalJSONToString.Get(settings) != MarshalJSONToString.Get(expected) {
			t.Errorf("Expected %q, got %q", MarshalJSONToString.Get(expected), MarshalJSONToString.Get(settings))
		}
	}
}

func TestShellProvisioner(t *testing.T) {
	expected := map[string]interface{}{
		"binary": false,
		"except": []string{
			"docker",
		},
		"execute_command": "echo 'vagrant'|sudo -S sh '{{.Path}}'",
		"inline_shebang":  "/bin/sh",
		"only": []string{
			"virtualbox-iso",
		},
		"remote_path": "/tmp/script.sh",
		"scripts": []string{
			"shell/setup_test.sh",
			"shell/vagrant_test.sh",
			"shell/sudoers_test.sh",
			"shell/cleanup_test.sh",
		},
		"start_retry_timeout": "5m",
		"type":                "shell",
	}
	settings, err := testRawTemplateProvisionersAll.createShellScript("shell")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if MarshalJSONToString.Get(settings) != MarshalJSONToString.Get(expected) {
			t.Errorf("Expected %q, got %q", MarshalJSONToString.Get(expected), MarshalJSONToString.Get(settings))
		}
	}
}

func TestFileUploadsProvisioner(t *testing.T) {
	expected := map[string]interface{}{
		"destination": "/tmp/app.tar.gz",
		"source":      "file/app.tar.gz",
		"type":        "file",
	}
	settings, err := testRawTemplateProvisionersAll.createFileUploads("file")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if MarshalJSONToString.Get(settings) != MarshalJSONToString.Get(expected) {
			t.Errorf("Expected %q, got %q", MarshalJSONToString.Get(expected), MarshalJSONToString.Get(settings))
		}
	}
}
