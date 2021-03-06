package app

import (
	"reflect"
	"testing"

	"github.com/mohae/contour"
)

var testUbuntu = RawTemplate{
	IODirInf: IODirInf{
		TemplateOutputDir: "../test_files/ubuntu/out/ubuntu",
		PackerOutputDir:   "boxes/:distro/:build_name",
		SourceDir:         "../test_files/src/ubuntu",
	},
	PackerInf: PackerInf{
		MinPackerVersion: "",
		Description:      "Test build template",
	},
	BuildInf: BuildInf{
		Name:      ":type-:release-:image-:arch",
		BuildName: "",
		BaseURL:   "http://releases.ubuntu.com/",
	},
	Distro:  "ubuntu",
	Arch:    "amd64",
	Image:   "server",
	Release: "14.04",
	VarVals: map[string]string{},
	Dirs:    map[string]string{},
	Files:   map[string]string{},
	Build: Build{
		BuilderIDs: []string{
			"virtualbox-iso",
			"vmware-iso",
		},
		Builders: map[string]BuilderC{
			"common": {
				TemplateSection{
					Settings: []string{
						"boot_command = boot_test.command",
						"boot_wait = 5s",
						"disk_size = 20000",
						"http_directory = http",
						"iso_checksum_type = sha256",
						"shutdown_command = shutdown_test.command",
						"ssh_password = vagrant",
						"ssh_port = 22",
						"ssh_username = vagrant",
						"ssh_timeout = 30m",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"virtualbox-iso": {
				TemplateSection{
					Arrays: map[string]interface{}{
						"vm_settings": []string{
							"cpus=1",
							"memory=4096",
						},
					},
				},
			},
			"vmware-iso": {
				TemplateSection{
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
		},
		PostProcessors: map[string]PostProcessorC{
			"vagrant": {
				TemplateSection{
					Settings: []string{
						"keep_input_artifact = false",
						"output = out/someComposedBoxName.box",
					},
				},
			},
		},
		ProvisionerIDs: []string{
			"shell",
		},
		Provisioners: map[string]ProvisionerC{
			"shell": {
				TemplateSection{
					Settings: []string{
						"execute_command = execute_test.command",
					},
					Arrays: map[string]interface{}{
						"scripts": []string{
							"setup_test.sh",
							"base_test.sh",
							"vagrant_test.sh",
							"cleanup_test.sh",
							"zerodisk_test.sh",
						},
					},
				},
			},
		},
	},
}

var testCentOS = RawTemplate{
	IODirInf: IODirInf{
		TemplateOutputDir: "../test_files/out/centos",
		PackerOutputDir:   "boxes/:distro/:build_name",
		SourceDir:         "../test_files/src/centos",
	},
	PackerInf: PackerInf{
		MinPackerVersion: "",
		Description:      "Test build template for salt provisioner using CentOS6",
	},
	BuildInf: BuildInf{
		Name:      ":type-:release-:image-:arch",
		BuildName: "",
		BaseURL:   "",
	},
	Distro:  "centos",
	Arch:    "x86_64",
	Image:   "minimal",
	Release: "6",
	VarVals: map[string]string{},
	Dirs:    map[string]string{},
	Files:   map[string]string{},
	Build: Build{
		BuilderIDs: []string{
			"virtualbox-iso",
			"virtualbox-ovf",
			"vmware-iso",
			"vmware-vmx",
		},
		Builders: map[string]BuilderC{
			"common": {
				TemplateSection{
					Settings: []string{
						"boot_command = boot_test.command",
						"boot_wait = 5s",
						"disk_size = 20000",
						"http_directory = http",
						"iso_checksum_type = sha256",
						"shutdown_command = shutdown_test.command",
						"ssh_password = vagrant",
						"ssh_port = 22",
						"ssh_username = vagrant",
						"ssh_timeout = 30m",
					},
				},
			},
			"virtualbox-iso": {
				TemplateSection{
					Arrays: map[string]interface{}{
						"vm_settings": []string{
							"--cpus=1",
							"memory=4096",
						},
					},
				},
			},
			"virtualbox-ovf": {
				TemplateSection{
					Arrays: map[string]interface{}{
						"vm_settings": []string{
							"cpus=1",
							"--memory=4096",
						},
					},
				},
			},
			"vmware-iso": {
				TemplateSection{
					Arrays: map[string]interface{}{
						"vm_settings": []string{
							"cpuid.coresPerSocket=1",
							"memsize=1024",
							"numvcpus=1",
						},
					},
				},
			},
			"vmware-vmx": {
				TemplateSection{
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
		},
		PostProcessors: map[string]PostProcessorC{
			"vagrant": {
				TemplateSection{
					Settings: []string{
						"keep_input_artifact = false",
						"output = out/someComposedBoxName.box",
					},
				},
			},
		},
		ProvisionerIDs: []string{
			"shell",
			"salt",
		},
		Provisioners: map[string]ProvisionerC{
			"salt": {
				TemplateSection{
					Settings: []string{
						"local_state_tree = ~/saltstates/centos6/salt",
						"skip_bootstrap = true",
					},
				},
			},
			"shell": {
				TemplateSection{
					Settings: []string{
						"execute_command = execute_test.command",
					},
					Arrays: map[string]interface{}{
						"scripts": []string{
							"setup_test.sh",
							"base_test.sh",
							"vagrant_test.sh",
							"cleanup_test.sh",
							"zerodisk_test.sh",
						},
					},
				},
			},
		},
	},
}

// Not all the settings in are valid for winrm, the invalid ones should not be included in the
// output.
var testAllBuilders = RawTemplate{
	IODirInf: IODirInf{
		TemplateOutputDir: "../test_files/out",
		PackerOutputDir:   "boxes/:distro/:build_name",
		SourceDir:         "../test_files/src",
	},
	PackerInf: PackerInf{
		MinPackerVersion: "",
		Description:      "Test build template for all builders",
	},
	BuildInf: BuildInf{
		Name:      "docker-alt",
		BuildName: "",
		BaseURL:   "",
	},
	Distro:  "ubuntu",
	Arch:    "amd64",
	Image:   "server",
	Release: "14.04",
	VarVals: map[string]string{},
	Dirs:    map[string]string{},
	Files:   map[string]string{},
	Build: Build{
		BuilderIDs: []string{
			"amazon-ebs",
			"amazon-instance",
			"digitalocean",
			"docker",
			"googlecompute",
			"null",
			"openstack1",
			"openstack2",
			"parallels-iso",
			"parallels-pvm",
			"virtualbox-iso",
			"virtualbox-ovf",
			"vmware-iso",
			"vmware-vmx",
		},
		Builders: map[string]BuilderC{
			"common": {
				TemplateSection{
					Type: "common",
					Settings: []string{
						"boot_wait = 5s",
						"disk_size = 20000",
						"http_directory = http",
						"iso_checksum_type = sha256",
						"shutdown_command = echo 'shutdown -P now' > /tmp/shutdown.sh; echo 'vagrant'|sudo -S sh '/tmp/shutdown.sh'",
						"ssh_password = vagrant",
						"ssh_port = 22",
						"ssh_username = vagrant",
						"ssh_timeout = 30m",
					},
				},
			},
			"amazon-chroot": {
				TemplateSection{
					Type: "amazon-chroot",
					Settings: []string{
						"access_key=AWS_ACCESS_KEY",
						"ami_description=AMI_DESCRIPTION",
						"ami_name=AMI_NAME",
						"ami_virtualization_type=paravirtual",
						"command_wrapper={{.Command}}",
						"device_path=/dev/xvdf",
						"enhanced_networking=false",
						"mount_path=packer-amazon-chroot-volumes/{{.Device}}",
						"secret_key=AWS_SECRET_ACCESS_KEY",
						"source_ami=SOURCE_AMI",
					},
					Arrays: map[string]interface{}{
						"ami_groups": []string{
							"AGroup",
						},
						"ami_product_codes": []string{
							"ami-d4e356aa",
						},
						"ami_regions": []string{
							"us-east-1",
						},
						"ami_users": []string{
							"aws-account-1",
						},
						"chroot_mounts": []interface{}{
							[]string{
								"proc",
								"proc",
								"/proc",
							},
							[]string{
								"bind",
								"/dev",
								"/dev",
							},
						},
						"copy_files": []string{
							"/etc/resolv.conf",
						},
						"tags": map[string]string{
							"OS_Version": "Ubuntu",
							"Release":    "Latest",
						},
					},
				},
			},
			"amazon-ebs": {
				TemplateSection{
					Type: "amazon-ebs",
					Settings: []string{
						"access_key=AWS_ACCESS_KEY",
						"ami_description=AMI_DESCRIPTION",
						"ami_name=AMI_NAME",
						"associate_public_ip_address=false",
						"availability_zone=us-east-1b",
						"enhanced_networking=false",
						"iam_instance_profile=INSTANCE_PROFILE",
						"instance_type=m3.medium",
						"region=us-east-1",
						"secret_key=AWS_SECRET_ACCESS_KEY",
						"security_group_id=GROUP_ID",
						"source_ami=SOURCE_AMI",
						"spot_price=auto",
						"spot_price_auto_product=Linux/Unix",
						"ssh_private_key_file=myKey",
						"ssh_username=vagrant",
						"subnet_id=subnet-12345def",
						"temporary_key_pair_name=TMP_KEYPAIR",
						"token=AWS_SECURITY_TOKEN",
						"user_data=SOME_USER_DATA",
						"user_data_file=amazon.userdata",
						"vpc_id=VPC_ID",
						"windows_password_timeout=10m",
					},
					Arrays: map[string]interface{}{
						"ami_block_device_mappings": []map[string]interface{}{
							{
								"device_name":  "/dev/sdb",
								"virtual_name": "/ephemeral0",
							},
							{
								"device_name":  "/dev/sdc",
								"virtual_name": "/ephemeral1",
							},
						},
						"ami_groups": []string{
							"AGroup",
						},
						"ami_product_codes": []string{
							"ami-d4e356aa",
						},
						"ami_regions": []string{
							"us-east-1",
						},
						"ami_users": []string{
							"ami-account",
						},
						"launch_block_device_mappings": []map[string]string{
							{
								"device_name":  "/dev/sdd",
								"virtual_name": "/ephemeral2",
							},
							{
								"device_name":  "/dev/sde",
								"virtual_name": "/ephemeral3",
							},
						},
						"security_group_ids": []string{
							"SECURITY_GROUP",
						},
						"run_tags": map[string]string{
							"foo": "bar",
							"fiz": "baz",
						},
						"tags": map[string]string{
							"OS_Version": "Ubuntu",
							"Release":    "Latest",
						},
					},
				},
			},
			"amazon-instance": {
				TemplateSection{
					Type: "amazon-instance",
					Settings: []string{
						"access_key=AWS_ACCESS_KEY",
						"account_id=YOUR_ACCOUNT_ID",
						"ami_description=AMI_DESCRIPTION",
						"ami_name=AMI_NAME",
						"ami_virtualization_type=paravirtual",
						"associate_public_ip_address=false",
						"availability_zone=us-east-1b",
						"bundle_destination=/tmp",
						"bundle_prefix=image--{{timestamp}}",
						"bundle_upload_command=bundle_upload.command",
						"bundle_vol_command=bundle_vol.command",
						"ebs_optimized=true",
						"enhanced_networking=false",
						"force_deregister=false",
						"iam_instance_profile=INSTANCE_PROFILE",
						"instance_type=m3.medium",
						"region=us-east-1",
						"s3_bucket=packer_bucket",
						"secret_key=AWS_SECRET_ACCESS_KEY",
						"security_group_id=GROUP_ID",
						"source_ami=SOURCE_AMI",
						"spot_price=auto",
						"spot_price_auto_product=Linux/Unix",
						"ssh_keypair_name=myKeyPair",
						"ssh_private_ip=true",
						"ssh_private_key_file=myKey",
						"ssh_username=vagrant",
						"subnet_id=subnet-12345def",
						"temporary_key_pair_name=TMP_KEYPAIR",
						"user_data=SOME_USER_DATA",
						"user_data_file=amazon.userdata",
						"vpc_id=VPC_ID",
						"windows_password_timeout=10m",
						"x509_cert_path=/path/to/x509/cert",
						"x509_key_path=/path/to/x509/key",
						"x509_upload_path=/etc/x509",
					},
					Arrays: map[string]interface{}{
						"ami_block_device_mappings": [][]string{
							[]string{
								"delete_on_termination=true",
								"device_name=/dev/sdb",
								"encrypted=true",
								"iops=1000",
								"no_device=false",
								"snapshot_id=SNAPSHOT",
								"virtual_name=ephemeral0",
								"volume_type=io1",
								"volume_size=10",
							},
							[]string{
								"device_name=/dev/sdc",
								"volume_type=io1",
								"volume_size=10",
							},
						},
						"ami_groups": []string{
							"AGroup",
						},
						"ami_product_codes": []string{
							"ami-d4e356aa",
						},
						"ami_regions": []string{
							"us-east-1",
						},
						"ami_users": []string{
							"ami-account",
						},
						"launch_block_device_mappings": []map[string]string{
							{
								"device_name":  "/dev/sdd",
								"virtual_name": "/ephemeral2",
							},
							{
								"device_name":  "/dev/sde",
								"virtual_name": "/ephemeral3",
							},
						},
						"run_tags": map[string]string{
							"foo": "bar",
							"fiz": "baz",
						},
						"security_group_ids": []string{
							"SECURITY_GROUP",
						},
						"tags": map[string]string{
							"OS_Version": "Ubuntu",
							"Release":    "Latest",
						},
					},
				},
			},
			"digitalocean": {
				TemplateSection{
					Type: "digitalocean",
					Settings: []string{
						"api_token=DIGITALOCEAN_API_TOKEN",
						"droplet_name=ocean-drop",
						"image=ubuntu-12-04-x64",
						"private_networking=false",
						"region=nyc3",
						"size=512mb",
						"snapshot_name=my-snapshot",
						"state_timeout=6m",
						"user_data=userdata",
					},
				},
			},
			"docker": {
				TemplateSection{
					Type: "docker",
					Settings: []string{
						"commit=true",
						"discard=false",
						"export_path=export/path",
						"image=baseImage",
						"login=true",
						"login_email=test@test.com",
						"login_username=username",
						"login_password=password",
						"login_server=127.0.0.1",
						"pull=true",
					},
					Arrays: map[string]interface{}{
						"run_command": []string{
							"-d",
							"-i",
							"-t",
							"{{.Image}}",
							"/bin/bash",
						},
						"volumes": map[string]string{
							"/var/data1": "/var/data",
							"/var/www":   "/var/www",
						},
					},
				},
			},
			"googlecompute": {
				TemplateSection{
					Type: "googlecompute",
					Settings: []string{
						"account_file=account.json",
						"address=ext-static",
						"disk_size=20",
						"image_name=packer-{{timestamp}}",
						"image_description=test image",
						"instance_name=packer-{{uuid}}",
						"machine_type=nl-standard-1",
						"network=default",
						"preemtible=true",
						"project_id=projectID",
						"source_image=centos-6",
						"state_timeout=5m",
						"use_internal_ip=true",
						"zone=us-central1-a",
					},
					Arrays: map[string]interface{}{
						"metadata": map[string]string{
							"key-1": "value-1",
							"key-2": "value-2",
						},
						"tags": []string{
							"tag1",
						},
					},
				},
			},
			"null": {
				TemplateSection{
					Type:     "null",
					Settings: []string{},
					Arrays:   map[string]interface{}{},
				},
			},
			"openstack1": {
				TemplateSection{
					Type: "openstack",
					Settings: []string{
						"api_key=APIKEY",
						"availability_zone=zone1",
						"config_drive=true",
						"flavor=2",
						"floating_ip=192.168.1.1",
						"floating_ip_pool=192.168.100.1/24",
						"image_name=test image",
						"insecure=true",
						"password=packer",
						"rackconnect_wait",
						"region=DFW",
						"source_image=23b564c9-c3e6-49f9-bc68-86c7a9ab5018",
						"ssh_interface=private",
						"tenant_id=123",
						"use_floating_ip=true",
						"username=packer",
					},
					Arrays: map[string]interface{}{
						"networks": []string{
							"de305d54-75b4-431b-adb2-eb6b9e546014",
						},
						"security_groups": []string{
							"admins",
						},
						"metadata": map[string]interface{}{
							"quota_metadata_items": 128,
							"metadata_listen":      "0.0.0.0",
						},
					},
				},
			},
			"openstack2": {
				TemplateSection{
					Type: "openstack",
					Settings: []string{
						"api_key=APIKEY",
						"availability_zone=zone1",
						"config_drive=true",
						"flavor=2",
						"floating_ip=192.168.1.1",
						"floating_ip_pool=192.168.100.1/24",
						"image_name=test image",
						"insecure=true",
						"password=packer",
						"rackconnect_wait",
						"region=DFW",
						"source_image=23b564c9-c3e6-49f9-bc68-86c7a9ab5018",
						"ssh_interface=private",
						"tenant_name=acme",
						"use_floating_ip=true",
						"username=packer",
					},
					Arrays: map[string]interface{}{
						"networks": []string{
							"de305d54-75b4-431b-adb2-eb6b9e546014",
						},
						"security_groups": []string{
							"admins",
						},
						"metadata": map[string]interface{}{
							"quota_metadata_items": 128,
							"metadata_listen":      "0.0.0.0",
						},
					},
				},
			},
			"parallels-iso": {
				TemplateSection{
					Type: "parallels-iso",
					Settings: []string{
						"boot_wait=30s",
						"disk_size=20000",
						"guest_os_type=ubuntu",
						"hard_drive_interface=ide",
						"http_directory=http",
						"http_port_min=8000",
						"http_port_max=9000",
						"iso_checksum=ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
						"iso_target_path=packer_cache",
						"output_directory=out/dir",
						"parallels_tools_flavor=lin",
						"parallels_tools_guest_path=ptools",
						"prlctl_version_file=.prlctl_version",
						"shutdown_command=shutdown.command",
						"shutdown_timeout=5m",
						"skip_compaction=true",
						"vm_name=test-iso",
					},
					Arrays: map[string]interface{}{
						"boot_command": []string{
							"<bs>",
							"<del>",
							"<enter><return>",
							"<esc>",
						},
						"floppy_files": []string{
							"disk1",
						},
						"iso_urls": []string{
							"http://releases.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
							"http://2.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
						},
						"prlctl": [][]string{
							[]string{"set", "{{.Name}}", "--shf-host-add", "log", "--path", "{{pwd}}/log", "--mode", "rw", "--enable"},
							[]string{"set", "{{.Name}}", "--cpus", "1"},
						},

						"prlctl_post": [][]string{
							[]string{"set", "{{.Name}}", "--shf-host-del", "log"},
						},
					},
				},
			},
			"parallels-pvm": {
				TemplateSection{
					Type: "parallels-pvm",
					Settings: []string{
						"boot_wait=30s",
						"disk_size=20000",
						"output_directory=out/dir",
						"parallels_tools_flavor=lin",
						"parallels_tools_guest_path=ptools",
						"parallels_tools_mode=upload",
						"parallels_tools_path=prl-tools.iso",
						"prlctl_version_file=.prlctl_version",
						"reassign_mac=true",
						"shutdown_command=shutdown.command",
						"shutdown_timeout=5m",
						"skip_compaction=true",
						"source_path=source.pvm",
						"vm_name=test-iso",
					},
					Arrays: map[string]interface{}{
						"boot_command": []string{
							"<bs>",
							"<del>",
							"<enter><return>",
							"<esc>",
						},
						"floppy_files": []string{
							"disk1",
						},
						"iso_urls": []string{
							"http://releases.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
							"http://2.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
						},
						"prlctl": [][]string{
							[]string{"set", "{{.Name}}", "--shf-host-add", "log", "--path", "{{pwd}}/log", "--mode", "rw", "--enable"},
							[]string{"set", "{{.Name}}", "--cpus", "1"},
						},

						"prlctl_post": [][]string{
							[]string{"set", "{{.Name}}", "--shf-host-del", "log"},
						},
					},
				},
			},
			"qemu": {
				TemplateSection{
					Type: "qemu",
					Settings: []string{
						"accelerator=kvm",
						"boot_wait=10s",
						"disk_cache=writeback",
						"disk_compression=true",
						"disk_discard=ignore",
						"disk_image=true",
						"disk_interface=ide",
						"disk_size=40000",
						"format = ovf",
						"headless=true",
						"http_directory=http",
						"http_port_min=8000",
						"http_port_max=9000",
						"iso_checksum=ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
						"iso_target_path=isocache",
						"net_device=i82551",
						"output_directory=out/dir",
						"qemu_binary=qemu-system-x86_64",
						"skip_compaction=true",
						"ssh_username=vagrant",
					},
					Arrays: map[string]interface{}{
						"boot_command": []string{
							"<bs>",
							"<del>",
							"<enter><return>",
							"<esc>",
						},
						"floppy_files": []string{
							"disk1",
						},
						"iso_urls": []string{
							"http://releases.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
							"http://2.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
						},
						"qemuargs": [][]string{
							[]string{
								"-m",
								"1024m",
							},
							[]string{
								"--no-acpi",
								"",
							},
						},
					},
				},
			},
			"virtualbox-iso": {
				TemplateSection{
					Type: "virtualbox-iso",
					Settings: []string{
						"format = ovf",
						"guest_additions_mode=upload",
						"guest_additions_path=path/to/additions",
						"guest_additions_sha256=89dac78769b26f8facf98ce85020a605b7601fec1946b0597e22ced5498b3597",
						"guest_additions_url=file://guest-additions",
						"guest_os_type=Ubuntu_64",
						"hard_drive_interface=ide",
						"headless=true",
						"http_port_min=8000",
						"http_port_max=9000",
						"iso_checksum=ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
						"iso_interface=ide",
						"output_directory=out/dir",
						"shutdown_timeout=5m",
						"ssh_host_port_min=22",
						"ssh_host_port_max=40",
						"ssh_private_key_file=key/path",
						"virtualbox_version_file=.vbox_version",
						"vm_name=test-vb-iso",
					},
					Arrays: map[string]interface{}{
						"boot_command": []string{
							"<bs>",
							"<del>",
							"<enter><return>",
							"<esc>",
						},
						"export_opts": []string{
							"opt1",
						},
						"floppy_files": []string{
							"disk1",
						},
						"iso_urls": []string{
							"http://releases.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
							"http://2.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
						},
						"vboxmanage": []string{
							"--cpus=1",
							"memory=4096",
						},
						"vboxmanage_post": []string{
							"something=value",
						},
					},
				},
			},
			"virtualbox-ovf": {
				TemplateSection{
					Type: "virtualbox-ovf",
					Settings: []string{
						"format = ovf",
						"guest_additions_mode=upload",
						"guest_additions_path=path/to/additions",
						"guest_additions_sha256=89dac78769b26f8facf98ce85020a605b7601fec1946b0597e22ced5498b3597",
						"guest_additions_url=file://guest-additions",
						"headless=true",
						"http_port_min=8000",
						"http_port_max=9000",
						"import_opts=keepallmacs",
						"iso_checksum=ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
						"output_directory=out/dir",
						"shutdown_timeout=5m",
						"ssh_skip_nat_mapping=false",
						"source_path=source.ova",
						"ssh_host_port_min=22",
						"ssh_host_port_max=40",
						"ssh_private_key_file=key/path",
						"ssh_skip_nat_mapping=true",
						"virtualbox_version_file=.vbox_version",
						"vm_name=test-vb-ovf",
					},
					Arrays: map[string]interface{}{
						"boot_command": []string{
							"<bs>",
							"<del>",
							"<enter><return>",
							"<esc>",
						},
						"export_opts": []string{
							"opt1",
						},
						"import_flags": []string{
							"--eula-accept",
						},
						"floppy_files": []string{
							"disk1",
						},
						"vboxmanage": []string{
							"cpus=1",
							"--memory=4096",
						},
						"vboxmanage_post": []string{
							"something=value",
						},
					},
				},
			},
			"vmware-iso": {
				TemplateSection{
					Type: "vmware-iso",
					Settings: []string{
						"communicator=none",
						"disk_type_id=1",
						"fusion_app_path=/Applications/VMware Fusion.app",
						"headless=true",
						"http_port_min=8000",
						"http_port_max=9000",
						"iso_checksum=ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
						"iso_target_path=../isocache/",
						"output_directory=out/dir",
						"remote_cache_datastore=datastore1",
						"remote_cache_directory=packer_cache",
						"remote_datastore=datastore1",
						"remote_host=remoteHost",
						"remote_password=rpassword",
						"remote_private_key_file=secret",
						"remote_type=esx5",
						"shutdown_timeout=5m",
						"skip_compaction=true",
						"ssh_host=127.0.0.1",
						"tools_upload_flavor=linux",
						"tools_upload_path={{.Flavor}}.iso",
						"version=9",
						"vm_name=packer-BUILDNAME",
						"vmdk_name=packer",
						"vmx_template_path=template/path",
						"vnc_port_min=5900",
						"vnc_port_max=6000",
					},
					Arrays: map[string]interface{}{
						"boot_command": []string{
							"<bs>",
							"<del>",
							"<enter><return>",
							"<esc>",
						},
						"disk_additional_size": []string{
							"10000",
						},
						"floppy_files": []string{
							"disk1",
						},
						"iso_urls": []string{
							"http://releases.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
							"http://2.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
						},
						"vmx_data": []string{
							"cpuid.coresPerSocket=1",
							"memsize=1024",
							"numvcpus=1",
						},
						"vmx_data_post": []string{
							"something=value",
						},
					},
				},
			},
			"vmware-vmx": {
				TemplateSection{
					Type: "vmware-vmx",
					Settings: []string{
						"fusion_app_path=/Applications/VMware Fusion.app",
						"headless=true",
						"http_port_min=8000",
						"http_port_max=9000",
						"output_directory=out/dir",
						"shutdown_timeout=5m",
						"skip_compaction=false",
						"source_path=source.vmx",
						"vm_name=packer-BUILDNAME",
						"vnc_port_min=5900",
						"vnc_port_max=6000",
					},
					Arrays: map[string]interface{}{
						"boot_command": []string{
							"<bs>",
							"<del>",
							"<enter><return>",
							"<esc>",
						},
						"floppy_files": []string{
							"disk1",
						},
						"vmx_data": []string{
							"cpuid.coresPerSocket=1",
							"memsize=1024",
							"numvcpus=1",
						},
						"vmx_data_post": []string{
							"something=value",
						},
					},
				},
			},
		},
		PostProcessorIDs: []string{
			"vagrant",
		},
		PostProcessors: map[string]PostProcessorC{
			"vagrant": {
				TemplateSection{
					Type: "vagrant",
					Settings: []string{
						"keep_input_artifact = false",
						"output = out/someComposedBoxName.box",
					},
				},
			},
		},
		ProvisionerIDs: []string{
			"salt",
		},
		Provisioners: map[string]ProvisionerC{
			"salt": {
				TemplateSection{
					Type: "salt",
					Settings: []string{
						"local_state_tree = ~/saltstates/centos6/salt",
						"skip_bootstrap = true",
					},
				},
			},
		},
	},
}

// Not all the settings in are valid for winrm, the invalid ones should not be included in the
// output.
var testAllBuildersSSH = RawTemplate{
	IODirInf: IODirInf{
		TemplateOutputDir: "../test_files/out",
		PackerOutputDir:   "boxes/:distro/:build_name",
		SourceDir:         "../test_files/src",
	},
	PackerInf: PackerInf{
		MinPackerVersion: "",
		Description:      "Test build template for all builders",
	},
	BuildInf: BuildInf{
		Name:      "docker-alt",
		BuildName: "",
		BaseURL:   "",
	},
	Distro:  "ubuntu",
	Arch:    "amd64",
	Image:   "server",
	Release: "14.04",
	VarVals: map[string]string{},
	Dirs:    map[string]string{},
	Files:   map[string]string{},
	Build: Build{
		BuilderIDs: []string{
			"amazon-ebs",
			"amazon-instance",
			"digitalocean",
			"docker",
			"googlecompute",
			"null",
			"openstack",
			"parallels-iso",
			"parallels-pvm",
			"virtualbox-iso",
			"virtualbox-ovf",
			"vmware-iso",
			"vmware-vmx",
		},
		Builders: map[string]BuilderC{
			"common": {
				TemplateSection{
					Type: "common",
					Settings: []string{
						"boot_wait = 5s",
						"disk_size = 20000",
						"http_directory = http",
						"iso_checksum_type = sha256",
						"shutdown_command = echo 'shutdown -P now' > /tmp/shutdown.sh; echo 'vagrant'|sudo -S sh '/tmp/shutdown.sh'",
						"ssh_bastion_host=bastion.host",
						"ssh_bastion_port=2222",
						"ssh_bastion_username=packer",
						"ssh_bastion_password=packer",
						"ssh_bastion_private_key_file=secret",
						"ssh_disable_agent=true",
						"ssh_handshake_attempts=10",
						"ssh_host=127.0.0.1",
						"ssh_password=vagrant",
						"ssh_port=22",
						"ssh_private_key_file=key/path",
						"ssh_pty=true",
						"ssh_timeout=10m",
						"ssh_username=vagrant",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"amazon-chroot": {
				TemplateSection{
					Type: "amazon-chroot",
					Settings: []string{
						"access_key=AWS_ACCESS_KEY",
						"ami_description=AMI_DESCRIPTION",
						"ami_name=AMI_NAME",
						"ami_virtualization_type=paravirtual",
						"communicator=ssh",
						"command_wrapper={{.Command}}",
						"device_path=/dev/xvdf",
						"enhanced_networking=false",
						"mount_path=packer-amazon-chroot-volumes/{{.Device}}",
						"secret_key=AWS_SECRET_ACCESS_KEY",
						"source_ami=SOURCE_AMI",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"amazon-ebs": {
				TemplateSection{
					Type: "amazon-ebs",
					Settings: []string{
						"access_key=AWS_ACCESS_KEY",
						"ami_description=AMI_DESCRIPTION",
						"ami_name=AMI_NAME",
						"associate_public_ip_address=false",
						"availability_zone=us-east-1b",
						"communicator=ssh",
						"enhanced_networking=false",
						"iam_instance_profile=INSTANCE_PROFILE",
						"instance_type=m3.medium",
						"region=us-east-1",
						"secret_key=AWS_SECRET_ACCESS_KEY",
						"security_group_id=GROUP_ID",
						"source_ami=SOURCE_AMI",
						"spot_price=auto",
						"spot_price_auto_product=Linux/Unix",
						"ssh_private_key_file=myKey",
						"ssh_username=vagrant",
						"subnet_id=subnet-12345def",
						"temporary_key_pair_name=TMP_KEYPAIR",
						"token=AWS_SECURITY_TOKEN",
						"user_data=SOME_USER_DATA",
						"user_data_file=amazon.userdata",
						"vpc_id=VPC_ID",
						"windows_password_timeout=10m",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"amazon-instance": {
				TemplateSection{
					Type: "amazon-instance",
					Settings: []string{
						"access_key=AWS_ACCESS_KEY",
						"account_id=YOUR_ACCOUNT_ID",
						"ami_description=AMI_DESCRIPTION",
						"ami_name=AMI_NAME",
						"ami_virtualization_type=paravirtual",
						"associate_public_ip_address=false",
						"availability_zone=us-east-1b",
						"bundle_destination=/tmp",
						"bundle_prefix=image--{{timestamp}}",
						"bundle_upload_command=bundle_upload.command",
						"bundle_vol_command=bundle_vol.command",
						"communicator=ssh",
						"ebs_optimized=true",
						"enhanced_networking=false",
						"force_deregister=false",
						"iam_instance_profile=INSTANCE_PROFILE",
						"instance_type=m3.medium",
						"region=us-east-1",
						"s3_bucket=packer_bucket",
						"secret_key=AWS_SECRET_ACCESS_KEY",
						"security_group_id=GROUP_ID",
						"source_ami=SOURCE_AMI",
						"spot_price=auto",
						"spot_price_auto_product=Linux/Unix",
						"ssh_keypair_name=myKeyPair",
						"ssh_private_ip=true",
						"ssh_private_key_file=myKey",
						"ssh_username=vagrant",
						"subnet_id=subnet-12345def",
						"temporary_key_pair_name=TMP_KEYPAIR",
						"user_data=SOME_USER_DATA",
						"user_data_file=amazon.userdata",
						"vpc_id=VPC_ID",
						"windows_password_timeout=10m",
						"x509_cert_path=/path/to/x509/cert",
						"x509_key_path=/path/to/x509/key",
						"x509_upload_path=/etc/x509",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"digitalocean": {
				TemplateSection{
					Type: "digitalocean",
					Settings: []string{
						"api_token=DIGITALOCEAN_API_TOKEN",
						"communicator=ssh",
						"droplet_name=ocean-drop",
						"image=ubuntu-12-04-x64",
						"private_networking=false",
						"region=nyc3",
						"size=512mb",
						"snapshot_name=my-snapshot",
						"state_timeout=6m",
						"user_data=userdata",
					},
				},
			},
			"docker": {
				TemplateSection{
					Type: "docker",
					Settings: []string{
						"commit=true",
						"communicator=ssh",
						"discard=false",
						"export_path=export/path",
						"image=baseImage",
						"login=true",
						"login_email=test@test.com",
						"login_username=username",
						"login_password=password",
						"login_server=127.0.0.1",
						"pull=true",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"googlecompute": {
				TemplateSection{
					Type: "googlecompute",
					Settings: []string{
						"account_file=account.json",
						"address=ext-static",
						"communicator=ssh",
						"disk_size=20",
						"image_name=packer-{{timestamp}}",
						"image_description=test image",
						"instance_name=packer-{{uuid}}",
						"machine_type=nl-standard-1",
						"network=default",
						"preemtible=true",
						"project_id=projectID",
						"source_image=centos-6",
						"state_timeout=5m",
						"use_internal_ip=true",
						"zone=us-central1-a",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"null": {
				TemplateSection{
					Type: "null",
					Settings: []string{
						"communicator=ssh",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"openstack": {
				TemplateSection{
					Type: "openstack",
					Settings: []string{
						"api_key=APIKEY",
						"availability_zone=zone1",
						"communicator=ssh",
						"config_drive=true",
						"flavor=2",
						"floating_ip=192.168.1.1",
						"floating_ip_pool=192.168.100.1/24",
						"image_name=test image",
						"insecure=true",
						"password=packer",
						"rackconnect_wait",
						"region=DFW",
						"source_image=23b564c9-c3e6-49f9-bc68-86c7a9ab5018",
						"ssh_interface=private",
						"tenant_name=acme",
						"use_floating_ip=true",
						"username=packer",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"parallels-iso": {
				TemplateSection{
					Type: "parallels-iso",
					Settings: []string{
						"boot_wait=30s",
						"communicator=ssh",
						"disk_size=20000",
						"guest_os_type=ubuntu",
						"hard_drive_interface=ide",
						"http_directory=http",
						"http_port_min=8000",
						"http_port_max=9000",
						"iso_checksum=ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
						"iso_target_path=packer_cache",
						"iso_url=http://releases.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
						"output_directory=out/dir",
						"parallels_tools_flavor=lin",
						"parallels_tools_guest_path=ptools",
						"prlctl_version_file=.prlctl_version",
						"shutdown_command=shutdown.command",
						"shutdown_timeout=5m",
						"skip_compaction=true",
						"vm_name=test-iso",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"parallels-pvm": {
				TemplateSection{
					Type: "parallels-pvm",
					Settings: []string{
						"boot_wait=30s",
						"communicator=ssh",
						"disk_size=20000",
						"output_directory=out/dir",
						"parallels_tools_flavor=lin",
						"parallels_tools_guest_path=ptools",
						"parallels_tools_mode=upload",
						"parallels_tools_path=prl-tools.iso",
						"prlctl_version_file=.prlctl_version",
						"reassign_mac=true",
						"shutdown_command=shutdown.command",
						"shutdown_timeout=5m",
						"source_path=source.pvm",
						"skip_compaction=true",
						"vm_name=test-iso",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"qemu": {
				TemplateSection{
					Type: "qemu",
					Settings: []string{
						"accelerator=kvm",
						"boot_wait=10s",
						"communicator=ssh",
						"disk_cache=writeback",
						"disk_compression=true",
						"disk_discard=ignore",
						"disk_image=true",
						"disk_interface=ide",
						"disk_size=40000",
						"format = ovf",
						"headless=true",
						"http_directory=http",
						"http_port_min=8000",
						"http_port_max=9000",
						"iso_checksum=ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
						"iso_target_path=isocache",
						"iso_url=http://releases.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
						"net_device=i82551",
						"output_directory=out/dir",
						"qemu_binary=qemu-system-x86_64",
						"skip_compaction=true",
						"ssh_username=vagrant",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"virtualbox-iso": {
				TemplateSection{
					Type: "virtualbox-iso",
					Settings: []string{
						"communicator=ssh",
						"format = ovf",
						"guest_additions_mode=upload",
						"guest_additions_path=path/to/additions",
						"guest_additions_sha256=89dac78769b26f8facf98ce85020a605b7601fec1946b0597e22ced5498b3597",
						"guest_additions_url=file://guest-additions",
						"guest_os_type=Ubuntu_64",
						"hard_drive_interface=ide",
						"headless=true",
						"http_port_min=8000",
						"http_port_max=9000",
						"iso_checksum=ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
						"iso_interface=ide",
						"iso_url=http://releases.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
						"output_directory=out/dir",
						"shutdown_timeout=5m",
						"ssh_host_port_min=22",
						"ssh_host_port_max=40",
						"ssh_private_key_file=key/path",
						"virtualbox_version_file=.vbox_version",
						"vm_name=test-vb-iso",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"virtualbox-ovf": {
				TemplateSection{
					Type: "virtualbox-ovf",
					Settings: []string{
						"communicator=ssh",
						"format = ovf",
						"guest_additions_mode=upload",
						"guest_additions_path=path/to/additions",
						"guest_additions_sha256=89dac78769b26f8facf98ce85020a605b7601fec1946b0597e22ced5498b3597",
						"guest_additions_url=file://guest-additions",
						"headless=true",
						"http_port_min=8000",
						"http_port_max=9000",
						"import_opts=keepallmacs",
						"output_directory=out/dir",
						"shutdown_timeout=5m",
						"ssh_private_key_file=key/path",
						"ssh_skip_nat_mapping=false",
						"source_path=source.ova",
						"ssh_host_port_min=22",
						"ssh_host_port_max=40",
						"ssh_private_key_file=key/path",
						"ssh_skip_nat_mapping=true",
						"virtualbox_version_file=.vbox_version",
						"vm_name=test-vb-ovf",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"vmware-iso": {
				TemplateSection{
					Type: "vmware-iso",
					Settings: []string{
						"communicator=ssh",
						"disk_type_id=1",
						"fusion_app_path=/Applications/VMware Fusion.app",
						"headless=true",
						"http_port_min=8000",
						"http_port_max=9000",
						"iso_checksum=ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
						"iso_target_path=../isocache/",
						"iso_url=http://releases.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
						"output_directory=out/dir",
						"remote_cache_datastore=datastore1",
						"remote_cache_directory=packer_cache",
						"remote_datastore=datastore1",
						"remote_host=remoteHost",
						"remote_password=rpassword",
						"remote_private_key_file=secret",
						"remote_type=esx5",
						"shutdown_timeout=5m",
						"skip_compaction=true",
						"ssh_skip_nat_mapping=false",
						"ssh_host_port_min=22",
						"ssh_host_port_max=40",
						"tools_upload_flavor=linux",
						"tools_upload_path={{.Flavor}}.iso",
						"version=9",
						"vm_name=packer-BUILDNAME",
						"vmdk_name=packer",
						"vmx_template_path=template/path",
						"vnc_port_min=5900",
						"vnc_port_max=6000",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"vmware-vmx": {
				TemplateSection{
					Type: "vmware-vmx",
					Settings: []string{
						"communicator=ssh",
						"fusion_app_path=/Applications/VMware Fusion.app",
						"headless=true",
						"http_port_min=8000",
						"http_port_max=9000",
						"output_directory=out/dir",
						"shutdown_timeout=5m",
						"skip_compaction=false",
						"source_path=source.vmx",
						"vm_name=packer-BUILDNAME",
						"vnc_port_min=5900",
						"vnc_port_max=6000",
					},
					Arrays: map[string]interface{}{},
				},
			},
		},
	},
}

// Not all the settings in are valid for winrm, the invalid ones should not be included in the
// output.
var testAllBuildersWinRM = RawTemplate{
	IODirInf: IODirInf{
		TemplateOutputDir: "../test_files/out",
		PackerOutputDir:   "boxes/:distro/:build_name",
		SourceDir:         "../test_files/src",
	},
	PackerInf: PackerInf{
		MinPackerVersion: "",
		Description:      "Test build template for all builders",
	},
	BuildInf: BuildInf{
		Name:      "docker-alt",
		BuildName: "",
		BaseURL:   "",
	},
	Distro:  "ubuntu",
	Arch:    "amd64",
	Image:   "server",
	Release: "14.04",
	VarVals: map[string]string{},
	Dirs:    map[string]string{},
	Files:   map[string]string{},
	Build: Build{
		BuilderIDs: []string{
			"amazon-ebs",
			"amazon-instance",
			"digitalocean",
			"docker",
			"googlecompute",
			"null",
			"openstack",
			"parallels-iso",
			"parallels-pvm",
			"virtualbox-iso",
			"virtualbox-ovf",
			"vmware-iso",
			"vmware-vmx",
		},
		Builders: map[string]BuilderC{
			"common": {
				TemplateSection{
					Type: "common",
					Settings: []string{
						"boot_wait = 5s",
						"disk_size = 20000",
						"http_directory = http",
						"iso_checksum_type = sha256",
						"shutdown_command = echo 'shutdown -P now' > /tmp/shutdown.sh; echo 'vagrant'|sudo -S sh '/tmp/shutdown.sh'",
						"winrm_host=host",
						"winrm_password = vagrant",
						"winrm_port = 22",
						"winrm_username = vagrant",
						"winrm_timeout=10m",
						"winrm_use_ssl=true",
						"winrm_insecure=true",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"amazon-chroot": {
				TemplateSection{
					Type: "amazon-chroot",
					Settings: []string{
						"access_key=AWS_ACCESS_KEY",
						"ami_description=AMI_DESCRIPTION",
						"ami_name=AMI_NAME",
						"ami_virtualization_type=paravirtual",
						"communicator=winrm",
						"command_wrapper={{.Command}}",
						"device_path=/dev/xvdf",
						"enhanced_networking=false",
						"mount_path=packer-amazon-chroot-volumes/{{.Device}}",
						"secret_key=AWS_SECRET_ACCESS_KEY",
						"source_ami=SOURCE_AMI",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"amazon-ebs": {
				TemplateSection{
					Type: "amazon-ebs",
					Settings: []string{
						"access_key=AWS_ACCESS_KEY",
						"ami_description=AMI_DESCRIPTION",
						"ami_name=AMI_NAME",
						"associate_public_ip_address=false",
						"availability_zone=us-east-1b",
						"communicator=winrm",
						"enhanced_networking=false",
						"iam_instance_profile=INSTANCE_PROFILE",
						"instance_type=m3.medium",
						"region=us-east-1",
						"secret_key=AWS_SECRET_ACCESS_KEY",
						"security_group_id=GROUP_ID",
						"source_ami=SOURCE_AMI",
						"spot_price=auto",
						"spot_price_auto_product=Linux/Unix",
						"ssh_private_key_file=myKey",
						"ssh_username=vagrant",
						"subnet_id=subnet-12345def",
						"temporary_key_pair_name=TMP_KEYPAIR",
						"token=AWS_SECURITY_TOKEN",
						"user_data=SOME_USER_DATA",
						"user_data_file=amazon.userdata",
						"vpc_id=VPC_ID",
						"windows_password_timeout=10m",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"amazon-instance": {
				TemplateSection{
					Type: "amazon-instance",
					Settings: []string{
						"access_key=AWS_ACCESS_KEY",
						"account_id=YOUR_ACCOUNT_ID",
						"ami_description=AMI_DESCRIPTION",
						"ami_name=AMI_NAME",
						"ami_virtualization_type=paravirtual",
						"associate_public_ip_address=false",
						"availability_zone=us-east-1b",
						"bundle_destination=/tmp",
						"bundle_prefix=image--{{timestamp}}",
						"bundle_upload_command=bundle_upload.command",
						"bundle_vol_command=bundle_vol.command",
						"communicator=winrm",
						"ebs_optimized=true",
						"enhanced_networking=false",
						"force_deregister=false",
						"iam_instance_profile=INSTANCE_PROFILE",
						"instance_type=m3.medium",
						"region=us-east-1",
						"s3_bucket=packer_bucket",
						"secret_key=AWS_SECRET_ACCESS_KEY",
						"security_group_id=GROUP_ID",
						"source_ami=SOURCE_AMI",
						"spot_price=auto",
						"spot_price_auto_product=Linux/Unix",
						"ssh_keypair_name=myKeyPair",
						"ssh_private_ip=true",
						"ssh_private_key_file=myKey",
						"ssh_username=vagrant",
						"subnet_id=subnet-12345def",
						"temporary_key_pair_name=TMP_KEYPAIR",
						"user_data=SOME_USER_DATA",
						"user_data_file=amazon.userdata",
						"vpc_id=VPC_ID",
						"windows_password_timeout=10m",
						"x509_cert_path=/path/to/x509/cert",
						"x509_key_path=/path/to/x509/key",
						"x509_upload_path=/etc/x509",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"digitalocean": {
				TemplateSection{
					Type: "digitalocean",
					Settings: []string{
						"api_token=DIGITALOCEAN_API_TOKEN",
						"communicator=winrm",
						"droplet_name=ocean-drop",
						"image=ubuntu-12-04-x64",
						"private_networking=false",
						"region=nyc3",
						"size=512mb",
						"snapshot_name=my-snapshot",
						"state_timeout=6m",
						"user_data=userdata",
					},
				},
			},
			"docker": {
				TemplateSection{
					Type: "docker",
					Settings: []string{
						"commit=true",
						"communicator=winrm",
						"discard=false",
						"export_path=export/path",
						"image=baseImage",
						"login=true",
						"login_email=test@test.com",
						"login_username=username",
						"login_password=password",
						"login_server=127.0.0.1",
						"pull=true",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"googlecompute": {
				TemplateSection{
					Type: "googlecompute",
					Settings: []string{
						"account_file=account.json",
						"address=ext-static",
						"communicator=winrm",
						"disk_size=20",
						"image_name=packer-{{timestamp}}",
						"image_description=test image",
						"instance_name=packer-{{uuid}}",
						"machine_type=nl-standard-1",
						"network=default",
						"preemtible=true",
						"project_id=projectID",
						"source_image=centos-6",
						"state_timeout=5m",
						"use_internal_ip=true",
						"zone=us-central1-a",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"null": {
				TemplateSection{
					Type: "null",
					Settings: []string{
						"communicator=winrm",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"openstack": {
				TemplateSection{
					Type: "openstack",
					Settings: []string{
						"api_key=APIKEY",
						"availability_zone=zone1",
						"communicator=winrm",
						"config_drive=true",
						"flavor=2",
						"floating_ip=192.168.1.1",
						"floating_ip_pool=192.168.100.1/24",
						"image_name=test image",
						"insecure=true",
						"password=packer",
						"rackconnect_wait",
						"region=DFW",
						"source_image=23b564c9-c3e6-49f9-bc68-86c7a9ab5018",
						"ssh_interface=private",
						"tenant_name=acme",
						"use_floating_ip=true",
						"username=packer",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"parallels-iso": {
				TemplateSection{
					Type: "parallels-iso",
					Settings: []string{
						"boot_wait=30s",
						"communicator=winrm",
						"guest_os_type=ubuntu",
						"hard_drive_interface=ide",
						"http_directory=http",
						"http_port_min=8000",
						"http_port_max=9000",
						"iso_checksum=ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
						"iso_target_path=packer_cache",
						"iso_url=http://releases.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
						"output_directory=out/dir",
						"parallels_tools_flavor=lin",
						"parallels_tools_guest_path=ptools",
						"prlctl_version_file=.prlctl_version",
						"shutdown_command=shutdown.command",
						"shutdown_timeout=5m",
						"skip_compaction=true",
						"vm_name=test-iso",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"parallels-pvm": {
				TemplateSection{
					Type: "parallels-pvm",
					Settings: []string{
						"boot_wait=30s",
						"communicator=winrm",
						"output_directory=out/dir",
						"parallels_tools_flavor=lin",
						"parallels_tools_guest_path=ptools",
						"parallels_tools_mode=upload",
						"parallels_tools_path=prl-tools.iso",
						"prlctl_version_file=.prlctl_version",
						"reassign_mac=true",
						"shutdown_command=shutdown.command",
						"shutdown_timeout=5m",
						"skip_compaction=true",
						"source_path=source.pvm",
						"vm_name=test-iso",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"qemu": {
				TemplateSection{
					Type: "qemu",
					Settings: []string{
						"accelerator=kvm",
						"boot_wait=10s",
						"communicator=winrm",
						"disk_cache=writeback",
						"disk_compression=true",
						"disk_discard=ignore",
						"disk_image=true",
						"disk_interface=ide",
						"disk_size=40000",
						"format = ovf",
						"headless=true",
						"http_directory=http",
						"http_port_min=8000",
						"http_port_max=9000",
						"iso_checksum=ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
						"iso_target_path=isocache",
						"iso_url=http://releases.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
						"net_device=i82551",
						"output_directory=out/dir",
						"qemu_binary=qemu-system-x86_64",
						"skip_compaction=true",
						"ssh_username=vagrant",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"virtualbox-iso": {
				TemplateSection{
					Type: "virtualbox-iso",
					Settings: []string{
						"communicator=winrm",
						"format = ovf",
						"guest_additions_mode=upload",
						"guest_additions_path=path/to/additions",
						"guest_additions_sha256=89dac78769b26f8facf98ce85020a605b7601fec1946b0597e22ced5498b3597",
						"guest_additions_url=file://guest-additions",
						"guest_os_type=Ubuntu_64",
						"hard_drive_interface=ide",
						"headless=true",
						"http_port_min=8000",
						"http_port_max=9000",
						"iso_checksum=ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
						"iso_interface=ide",
						"iso_url=http://releases.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
						"output_directory=out/dir",
						"shutdown_timeout=5m",
						"ssh_skip_nat_mapping=false",
						"ssh_host_port_min=22",
						"ssh_host_port_max=40",
						"ssh_private_key_file=key/path",
						"virtualbox_version_file=.vbox_version",
						"vm_name=test-vb-iso",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"virtualbox-ovf": {
				TemplateSection{
					Type: "virtualbox-ovf",
					Settings: []string{
						"communicator=winrm",
						"format = ovf",
						"guest_additions_mode=upload",
						"guest_additions_path=path/to/additions",
						"guest_additions_sha256=89dac78769b26f8facf98ce85020a605b7601fec1946b0597e22ced5498b3597",
						"guest_additions_url=file://guest-additions",
						"headless=true",
						"http_port_min=8000",
						"http_port_max=9000",
						"import_opts=keepallmacs",
						"iso_checksum=ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
						"output_directory=out/dir",
						"shutdown_timeout=5m",
						"source_path=source.ova",
						"ssh_host_port_min=22",
						"ssh_host_port_max=40",
						"ssh_private_key_file=key/path",
						"ssh_skip_nat_mapping=true",
						"virtualbox_version_file=.vbox_version",
						"vm_name=test-vb-ovf",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"vmware-iso": {
				TemplateSection{
					Type: "vmware-iso",
					Settings: []string{
						"communicator=winrm",
						"disk_type_id=1",
						"fusion_app_path=/Applications/VMware Fusion.app",
						"headless=true",
						"http_port_min=8000",
						"http_port_max=9000",
						"iso_checksum=ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
						"iso_target_path=../isocache/",
						"iso_url=http://releases.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
						"output_directory=out/dir",
						"remote_cache_datastore=datastore1",
						"remote_cache_directory=packer_cache",
						"remote_datastore=datastore1",
						"remote_host=remoteHost",
						"remote_password=rpassword",
						"remote_private_key_file=secret",
						"remote_type=esx5",
						"shutdown_timeout=5m",
						"skip_compaction=true",
						"ssh_skip_nat_mapping=false",
						"ssh_host_port_min=22",
						"ssh_host_port_max=40",
						"tools_upload_flavor=linux",
						"tools_upload_path={{.Flavor}}.iso",
						"version=9",
						"vm_name=packer-BUILDNAME",
						"vmdk_name=packer",
						"vmx_template_path=template/path",
						"vnc_port_min=5900",
						"vnc_port_max=6000",
					},
					Arrays: map[string]interface{}{},
				},
			},
			"vmware-vmx": {
				TemplateSection{
					Type: "vmware-vmx",
					Settings: []string{
						"communicator=winrm",
						"fusion_app_path=/Applications/VMware Fusion.app",
						"headless=true",
						"http_port_min=8000",
						"http_port_max=9000",
						"output_directory=out/dir",
						"shutdown_timeout=5m",
						"skip_compaction=false",
						"source_path=source.vmx",
						"vm_name=packer-BUILDNAME",
						"vnc_port_min=5900",
						"vnc_port_max=6000",
					},
					Arrays: map[string]interface{}{},
				},
			},
		},
	},
}

var testDockerRunComandFile = RawTemplate{
	IODirInf: IODirInf{
		TemplateOutputDir: "../test_files/out",
		PackerOutputDir:   "boxes/:distro/:build_name",
		SourceDir:         "../test_files/src",
	},
	PackerInf: PackerInf{
		MinPackerVersion: "",
		Description:      "Test build template for all builders",
	},
	BuildInf: BuildInf{
		Name:      ":type-:release-:image-:arch",
		BuildName: "",
		BaseURL:   "",
	},
	Distro:  "ubuntu",
	Arch:    "amd64",
	Image:   "minimal",
	Release: "14.04",
	VarVals: map[string]string{},
	Dirs:    map[string]string{},
	Files:   map[string]string{},
	Build: Build{
		BuilderIDs: []string{
			"docker",
		},
		Builders: map[string]BuilderC{
			"docker": {
				TemplateSection{
					Settings: []string{
						"commit=true",
						"discard=false",
						"export_path=export/path",
						"image=baseImage",
						"login=true",
						"login_email=test@test.com",
						"login_username=username",
						"login_password=password",
						"login_server=127.0.0.1",
						"pull=true",
						"run_command=docker.command",
					},
					Arrays: map[string]interface{}{},
				},
			},
		},
	},
}

// This should still result in only 1 command array, using the array value and not the
// file
var testDockerRunComand = RawTemplate{
	IODirInf: IODirInf{
		TemplateOutputDir: "../test_files/out",
		PackerOutputDir:   "boxes/:distro/:build_name",
		SourceDir:         "../test_files/src",
	},
	PackerInf: PackerInf{
		MinPackerVersion: "",
		Description:      "Test build template for all builders",
	},
	BuildInf: BuildInf{
		Name:      ":type-:release-:image-:arch",
		BuildName: "",
		BaseURL:   "",
	},
	Distro:  "ubuntu",
	Arch:    "amd64",
	Image:   "minimal",
	Release: "14.04",
	VarVals: map[string]string{},
	Dirs:    map[string]string{},
	Files:   map[string]string{},
	Build: Build{
		BuilderIDs: []string{
			"docker",
		},
		Builders: map[string]BuilderC{
			"docker": {
				TemplateSection{
					Settings: []string{
						"commit=true",
						"discard=false",
						"export_path=export/path",
						"image=baseImage",
						"login=true",
						"login_email=test@test.com",
						"login_username=username",
						"login_password=password",
						"login_server=127.0.0.1",
						"pull=true",
						"run_command=docker.command",
					},
					Arrays: map[string]interface{}{
						"run_command": []string{
							"-d",
							"-i",
							"-t",
							"{{.Image}}",
							"/bin/bash",
						},
					},
				},
			},
		},
	},
}
var builderOrig = map[string]BuilderC{
	"common": {
		TemplateSection{
			Settings: []string{
				"boot_command = boot_test.command",
				"boot_wait = 5s",
				"disk_size = 20000",
				"http_directory = http",
				"iso_checksum_type = sha256",
				"shutdown_command = shutdown_test.command",
				"ssh_password = vagrant",
				"ssh_port = 22",
				"ssh_username = vagrant",
				"ssh_timeout = 30m",
			},
			Arrays: map[string]interface{}{},
		},
	},
	"virtualbox-iso": {
		TemplateSection{
			Arrays: map[string]interface{}{
				"vm_settings": []string{
					"cpus=1",
					"memory=4096",
				},
			},
		},
	},
	"vmware-iso": {
		TemplateSection{
			Arrays: map[string]interface{}{
				"vm_settings": []string{
					"cpuid.coresPerSocket=1",
					"memsize=1024",
					"numvcpus=1",
				},
			},
		},
	},
}

var builderNew = map[string]BuilderC{
	"common": {
		TemplateSection{
			Settings: []string{
				"boot_command = boot_test.command",
				"boot_wait = 15s",
				"disk_size = 20000",
				"http_directory = http",
				"iso_checksum_type = sha256",
				"shutdown_command = shutdown_test.command",
				"ssh_password = vagrant",
				"ssh_port = 22",
				"ssh_username = vagrant",
				"ssh_timeout = 240m",
			},
		},
	},
	"virtualbox-iso": {
		TemplateSection{
			Arrays: map[string]interface{}{
				"vm_settings": []string{
					"cpus=1",
					"memory=2048",
				},
			},
		},
	},
}

var builderMerged = map[string]BuilderC{
	"common": {
		TemplateSection{
			Settings: []string{
				"boot_command = boot_test.command",
				"boot_wait = 15s",
				"disk_size = 20000",
				"http_directory = http",
				"iso_checksum_type = sha256",
				"shutdown_command = shutdown_test.command",
				"ssh_password = vagrant",
				"ssh_port = 22",
				"ssh_username = vagrant",
				"ssh_timeout = 240m",
			},
			Arrays: map[string]interface{}{},
		},
	},
	"virtualbox-iso": {
		TemplateSection{
			Arrays: map[string]interface{}{
				"vm_settings": []string{
					"cpus=1",
					"memory=2048",
				},
			},
		},
	},
	"vmware-iso": {
		TemplateSection{
			Arrays: map[string]interface{}{
				"vm_settings": []string{
					"cpuid.coresPerSocket=1",
					"memsize=1024",
					"numvcpus=1",
				},
			},
		},
	},
}

var vbB = BuilderC{
	TemplateSection{
		Settings: []string{
			"boot_wait=5s",
			"disk_size = 20000",
			"ssh_port= 22",
			"ssh_username =vagrant",
		},
		Arrays: map[string]interface{}{
			"vm_settings": []string{
				"cpuid.coresPerSocket=1",
				"memsize=2048",
			},
		},
	},
}

func init() {
	b := true
	testAllBuilders.IncludeComponentString = &b
	testAllBuildersSSH.IncludeComponentString = &b
	testAllBuildersWinRM.IncludeComponentString = &b
}

func TestCreateBuilders(t *testing.T) {
	_, err := testRawTemplateBuilderOnly.createBuilders()
	if err == nil {
		t.Error("Expected error \"no builders specified\", got nil")
	} else {
		if err.Error() != "no builders specified" {
			t.Errorf("Expected \"no builders specified\", got %q", err)
		}
	}

	xerr := BuilderErr{Builder: AmazonEBS, Err: ErrBuilderNotFound}
	_, err = testRawTemplateWOSection.createBuilders()
	if err == nil {
		t.Errorf("Expected %s, got nil", xerr)
	} else {
		if err.Error() != xerr.Error() {
			t.Errorf("got %q, want %q", err, xerr)
		}
	}

	xerr.Builder = DigitalOcean
	testRawTemplateWOSection.Build.BuilderIDs[0] = "digitalocean"
	_, err = testRawTemplateWOSection.createBuilders()
	if err == nil {
		t.Errorf("Expected %s, got nil", xerr)
	} else {
		if err.Error() != xerr.Error() {
			t.Errorf("got %q, want %q", err, xerr)
		}
	}

	xerr.Builder = Docker
	testRawTemplateWOSection.Build.BuilderIDs[0] = "docker"
	_, err = testRawTemplateWOSection.createBuilders()
	if err == nil {
		t.Errorf("Expected %s, got nil", xerr)
	} else {
		if err.Error() != xerr.Error() {
			t.Errorf("got %q, want %q", err, xerr)
		}
	}

	xerr.Builder = GoogleCompute
	testRawTemplateWOSection.Build.BuilderIDs[0] = "googlecompute"
	_, err = testRawTemplateWOSection.createBuilders()
	if err == nil {
		t.Errorf("Expected %s, got nil", xerr)
	} else {
		if err.Error() != xerr.Error() {
			t.Errorf("got %q, want %q", err, xerr)
		}
	}

	xerr.Builder = VirtualBoxISO
	testRawTemplateWOSection.Build.BuilderIDs[0] = "virtualbox-iso"
	_, err = testRawTemplateWOSection.createBuilders()
	if err == nil {
		t.Errorf("Expected %s, got nil", xerr)
	} else {
		if err.Error() != xerr.Error() {
			t.Errorf("got %q, want %q", err, xerr)
		}
	}

	xerr.Builder = VirtualBoxOVF
	testRawTemplateWOSection.Build.BuilderIDs[0] = "virtualbox-ovf"
	_, err = testRawTemplateWOSection.createBuilders()
	if err == nil {
		t.Errorf("Expected %s, got nil", xerr)
	} else {
		if err.Error() != xerr.Error() {
			t.Errorf("got %q, want %q", err, xerr)
		}
	}

	xerr.Builder = VMWareISO
	testRawTemplateWOSection.Build.BuilderIDs[0] = "vmware-iso"
	_, err = testRawTemplateWOSection.createBuilders()
	if err == nil {
		t.Errorf("Expected %s, got nil", xerr)
	} else {
		if err.Error() != xerr.Error() {
			t.Errorf("got %q, want %q", err, xerr)
		}
	}

	xerr.Builder = VMWareVMX
	testRawTemplateWOSection.Build.BuilderIDs[0] = "vmware-vmx"
	_, err = testRawTemplateWOSection.createBuilders()
	if err == nil {
		t.Errorf("Expected %s, got nil", xerr)
	} else {
		if err.Error() != xerr.Error() {
			t.Errorf("got %q, want %q", err, xerr)
		}
	}

	r := testDistroDefaultUbuntu
	r.BuilderIDs = nil
	_, err = r.createBuilders()
	if err == nil {
		t.Error("Expected an error, got nil")
	} else {
		if err.Error() != "no builders specified" {
			t.Errorf("Expected \"no builders specified\"), got %q", err)
		}
	}
}

func TestRawTemplateUpdatebuilders(t *testing.T) {
	err := testUbuntu.updateBuilders(nil)
	if err != nil {
		t.Errorf("expected error to be nil, got %q", err)
	}
	msg, ok := EvalBuilders(testUbuntu.Builders, builderOrig)
	if !ok {
		t.Error(msg)
	}

	err = testUbuntu.updateBuilders(builderNew)
	if err != nil {
		t.Errorf("expected error to be nil, got %q", err)
	}

	msg, ok = EvalBuilders(testUbuntu.Builders, builderMerged)
	if !ok {
		t.Error(msg)
	}
}

func TestRawTemplateUpdateBuilderCommon(t *testing.T) {
	testUbuntu.updateCommon(builderNew["common"])
	new := testUbuntu.Builders["common"]
	old := builderMerged["common"]
	msg, ok := EvalTemplateSection(&new.TemplateSection, &old.TemplateSection)
	if !ok {
		t.Error(msg)
	}
}

func TestRawTemplateBuildersSettingsToMap(t *testing.T) {
	settings := vbB.settingsToMap(testRawTpl)
	if settings["boot_wait"] != "5s" {
		t.Errorf("Expected \"5s\", got %q", settings["boot_wait"])
	}
	if settings["disk_size"] != "20000" {
		t.Errorf("Expected \"20000\", got %q", settings["disk_size"])
	}
	if settings["ssh_port"] != "22" {
		t.Errorf("Expected \"22\", got %q", settings["ssh_port"])
	}
	if settings["ssh_username"] != "vagrant" {
		t.Errorf("Expected \"vagrant\", got %q", settings["ssh_username"])
	}
}

func TestCreateAmazonChroot(t *testing.T) {
	expected := map[string]interface{}{
		"access_key":      "AWS_ACCESS_KEY",
		"ami_description": "AMI_DESCRIPTION",
		"ami_groups": []string{
			"AGroup",
		},
		"ami_name": "AMI_NAME",
		"ami_product_codes": []string{
			"ami-d4e356aa",
		},
		"ami_regions": []string{
			"us-east-1",
		},
		"ami_users": []string{
			"aws-account-1",
		},
		"ami_virtualization_type": "paravirtual",
		"chroot_mounts": []interface{}{
			[]string{
				"proc",
				"proc",
				"/proc",
			},
			[]string{
				"bind",
				"/dev",
				"/dev",
			},
		},
		"command_wrapper": "{{.Command}}",
		"copy_files": []string{
			"/etc/resolv.conf",
		},
		"device_path":         "/dev/xvdf",
		"enhanced_networking": false,
		"mount_path":          "packer-amazon-chroot-volumes/{{.Device}}",
		"secret_key":          "AWS_SECRET_ACCESS_KEY",
		"source_ami":          "SOURCE_AMI",
		"tags": map[string]string{
			"OS_Version": "Ubuntu",
			"Release":    "Latest",
		},
		"type": "amazon-chroot",
	}
	bldr, err := testAllBuilders.createAmazonChroot("amazon-chroot")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(bldr, expected) {
			t.Errorf("Expected %#v, got %#v", expected, bldr)
		}
	}
	// SSH
	expectedSSH := map[string]interface{}{
		"access_key":                   "AWS_ACCESS_KEY",
		"ami_description":              "AMI_DESCRIPTION",
		"ami_name":                     "AMI_NAME",
		"ami_virtualization_type":      "paravirtual",
		"command_wrapper":              "{{.Command}}",
		"communicator":                 "ssh",
		"device_path":                  "/dev/xvdf",
		"enhanced_networking":          false,
		"mount_path":                   "packer-amazon-chroot-volumes/{{.Device}}",
		"secret_key":                   "AWS_SECRET_ACCESS_KEY",
		"source_ami":                   "SOURCE_AMI",
		"ssh_bastion_host":             "bastion.host",
		"ssh_bastion_port":             2222,
		"ssh_bastion_username":         "packer",
		"ssh_bastion_password":         "packer",
		"ssh_bastion_private_key_file": "secret",
		"ssh_disable_agent":            true,
		"ssh_handshake_attempts":       10,
		"ssh_host":                     "127.0.0.1",
		"ssh_password":                 "vagrant",
		"ssh_port":                     22,
		"ssh_private_key_file":         "key/path",
		"ssh_pty":                      true,
		"ssh_username":                 "vagrant",
		"ssh_timeout":                  "10m",
		"type":                         "amazon-chroot",
	}
	bldr, err = testAllBuildersSSH.createAmazonChroot("amazon-chroot")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(bldr, expectedSSH) {
			t.Errorf("Expected %#v, got %#v", expectedSSH, bldr)
		}
	}
	// WinRM
	expectedWinRM := map[string]interface{}{
		"access_key":              "AWS_ACCESS_KEY",
		"ami_description":         "AMI_DESCRIPTION",
		"ami_name":                "AMI_NAME",
		"ami_virtualization_type": "paravirtual",
		"command_wrapper":         "{{.Command}}",
		"communicator":            "winrm",
		"device_path":             "/dev/xvdf",
		"enhanced_networking":     false,
		"mount_path":              "packer-amazon-chroot-volumes/{{.Device}}",
		"secret_key":              "AWS_SECRET_ACCESS_KEY",
		"source_ami":              "SOURCE_AMI",
		"type":                    "amazon-chroot",
		"winrm_host":              "host",
		"winrm_password":          "vagrant",
		"winrm_port":              22,
		"winrm_timeout":           "10m",
		"winrm_username":          "vagrant",
		"winrm_use_ssl":           true,
		"winrm_insecure":          true,
	}
	bldr, err = testAllBuildersWinRM.createAmazonChroot("amazon-chroot")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(bldr, expectedWinRM) {
			t.Errorf("Expected %#v, got %#v", expectedWinRM, bldr)
		}
	}
}

func TestCreateAmazonEBS(t *testing.T) {
	expected := map[string]interface{}{
		"access_key": "AWS_ACCESS_KEY",
		"ami_block_device_mappings": []map[string]interface{}{
			{
				"device_name":  "/dev/sdb",
				"virtual_name": "/ephemeral0",
			},
			{
				"device_name":  "/dev/sdc",
				"virtual_name": "/ephemeral1",
			},
		},
		"ami_description": "AMI_DESCRIPTION",
		"ami_groups": []string{
			"AGroup",
		},
		"ami_name": "AMI_NAME",
		"ami_product_codes": []string{
			"ami-d4e356aa",
		},
		"ami_regions": []string{
			"us-east-1",
		},
		"ami_users": []string{
			"ami-account",
		},
		"associate_public_ip_address": false,
		"availability_zone":           "us-east-1b",
		"enhanced_networking":         false,
		"iam_instance_profile":        "INSTANCE_PROFILE",
		"instance_type":               "m3.medium",
		"launch_block_device_mappings": []map[string]string{
			{
				"device_name":  "/dev/sdd",
				"virtual_name": "/ephemeral2",
			},
			{
				"device_name":  "/dev/sde",
				"virtual_name": "/ephemeral3",
			},
		},
		"region": "us-east-1",
		"run_tags": map[string]string{
			"foo": "bar",
			"fiz": "baz",
		},
		"secret_key":        "AWS_SECRET_ACCESS_KEY",
		"security_group_id": "GROUP_ID",
		"security_group_ids": []string{
			"SECURITY_GROUP",
		},
		"source_ami":              "SOURCE_AMI",
		"spot_price":              "auto",
		"spot_price_auto_product": "Linux/Unix",
		"ssh_private_key_file":    "myKey",
		"ssh_username":            "vagrant",
		"subnet_id":               "subnet-12345def",
		"tags": map[string]string{
			"OS_Version": "Ubuntu",
			"Release":    "Latest",
		},
		"temporary_key_pair_name":  "TMP_KEYPAIR",
		"token":                    "AWS_SECURITY_TOKEN",
		"type":                     "amazon-ebs",
		"user_data":                "SOME_USER_DATA",
		"user_data_file":           "amazon-ebs/amazon.userdata",
		"vpc_id":                   "VPC_ID",
		"windows_password_timeout": "10m",
	}
	bldr, err := testAllBuilders.createAmazonEBS("amazon-ebs")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(bldr, expected) {
			t.Errorf("Expected %#v, got %#v", expected, bldr)
		}
	}
	// SSH
	expectedSSH := map[string]interface{}{
		"access_key":                   "AWS_ACCESS_KEY",
		"ami_description":              "AMI_DESCRIPTION",
		"ami_name":                     "AMI_NAME",
		"associate_public_ip_address":  false,
		"availability_zone":            "us-east-1b",
		"communicator":                 "ssh",
		"enhanced_networking":          false,
		"iam_instance_profile":         "INSTANCE_PROFILE",
		"instance_type":                "m3.medium",
		"region":                       "us-east-1",
		"secret_key":                   "AWS_SECRET_ACCESS_KEY",
		"security_group_id":            "GROUP_ID",
		"source_ami":                   "SOURCE_AMI",
		"spot_price":                   "auto",
		"spot_price_auto_product":      "Linux/Unix",
		"ssh_bastion_host":             "bastion.host",
		"ssh_bastion_port":             2222,
		"ssh_bastion_username":         "packer",
		"ssh_bastion_password":         "packer",
		"ssh_bastion_private_key_file": "secret",
		"ssh_disable_agent":            true,
		"ssh_handshake_attempts":       10,
		"ssh_host":                     "127.0.0.1",
		"ssh_password":                 "vagrant",
		"ssh_port":                     22,
		"ssh_private_key_file":         "myKey",
		"ssh_pty":                      true,
		"ssh_timeout":                  "10m",
		"ssh_username":                 "vagrant",
		"subnet_id":                    "subnet-12345def",
		"temporary_key_pair_name":      "TMP_KEYPAIR",
		"token":                        "AWS_SECURITY_TOKEN",
		"type":                         "amazon-ebs",
		"user_data":                    "SOME_USER_DATA",
		"user_data_file":               "amazon-ebs/amazon.userdata",
		"vpc_id":                       "VPC_ID",
	}
	bldr, err = testAllBuildersSSH.createAmazonEBS("amazon-ebs")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(bldr, expectedSSH) {
			t.Errorf("Expected %#v, got %#v", expectedSSH, bldr)
		}
	}
	// WinRM
	expectedWinRM := map[string]interface{}{
		"access_key":                  "AWS_ACCESS_KEY",
		"ami_description":             "AMI_DESCRIPTION",
		"ami_name":                    "AMI_NAME",
		"associate_public_ip_address": false,
		"availability_zone":           "us-east-1b",
		"communicator":                "winrm",
		"enhanced_networking":         false,
		"iam_instance_profile":        "INSTANCE_PROFILE",
		"instance_type":               "m3.medium",
		"region":                      "us-east-1",
		"secret_key":                  "AWS_SECRET_ACCESS_KEY",
		"security_group_id":           "GROUP_ID",
		"source_ami":                  "SOURCE_AMI",
		"spot_price":                  "auto",
		"spot_price_auto_product":     "Linux/Unix",
		"subnet_id":                   "subnet-12345def",
		"temporary_key_pair_name":     "TMP_KEYPAIR",
		"token":                       "AWS_SECURITY_TOKEN",
		"type":                        "amazon-ebs",
		"user_data":                   "SOME_USER_DATA",
		"user_data_file":              "amazon-ebs/amazon.userdata",
		"vpc_id":                      "VPC_ID",
		"windows_password_timeout":    "10m",
		"winrm_host":                  "host",
		"winrm_password":              "vagrant",
		"winrm_port":                  22,
		"winrm_timeout":               "10m",
		"winrm_username":              "vagrant",
		"winrm_use_ssl":               true,
		"winrm_insecure":              true,
	}
	bldr, err = testAllBuildersWinRM.createAmazonEBS("amazon-ebs")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(bldr, expectedWinRM) {
			t.Errorf("Expected %#v, got %#v", expectedWinRM, bldr)
		}
	}
}

func TestCreateAmazonInstance(t *testing.T) {
	expected := map[string]interface{}{
		"access_key":      "AWS_ACCESS_KEY",
		"account_id":      "YOUR_ACCOUNT_ID",
		"ami_description": "AMI_DESCRIPTION",
		"ami_block_device_mappings": []map[string]interface{}{
			{
				"delete_on_termination": true,
				"iops":                  1000,
				"device_name":           "/dev/sdb",
				"encrypted":             true,
				"no_device":             false,
				"snapshot_id":           "SNAPSHOT",
				"virtual_name":          "ephemeral0",
				"volume_type":           "io1",
				"volume_size":           10,
			},
			{
				"device_name": "/dev/sdc",
				"volume_type": "io1",
				"volume_size": 10,
			},
		},
		"ami_groups": []string{
			"AGroup",
		},
		"ami_name": "AMI_NAME",
		"ami_product_codes": []string{
			"ami-d4e356aa",
		},
		"ami_regions": []string{
			"us-east-1",
		},
		"ami_users": []string{
			"ami-account",
		},
		"ami_virtualization_type":     "paravirtual",
		"associate_public_ip_address": false,
		"availability_zone":           "us-east-1b",
		"bundle_destination":          "/tmp",
		"bundle_prefix":               "image--{{timestamp}}",
		"bundle_upload_command":       "sudo -n ec2-bundle-vol -k {{.KeyPath}} -u {{.AccountId}} -c {{.CertPath}} -r {{.Architecture}} -e {{.PrivatePath}} -d {{.Destination}} -p {{.Prefix}} --batch --no-filter",
		"bundle_vol_command":          "sudo -n ec2-upload-bundle -b {{.BucketName}} -m {{.ManifestPath}} -a {{.AccessKey}} -s {{.SecretKey}} -d {{.BundleDirectory}} --batch --region {{.Region}} --retry",
		"ebs_optimized":               true,
		"enhanced_networking":         false,
		"force_deregister":            false,
		"iam_instance_profile":        "INSTANCE_PROFILE",
		"instance_type":               "m3.medium",
		"launch_block_device_mappings": []map[string]string{
			{
				"device_name":  "/dev/sdd",
				"virtual_name": "/ephemeral2",
			},
			{
				"device_name":  "/dev/sde",
				"virtual_name": "/ephemeral3",
			},
		},
		"region": "us-east-1",
		"run_tags": map[string]string{
			"foo": "bar",
			"fiz": "baz",
		},
		"s3_bucket":         "packer_bucket",
		"secret_key":        "AWS_SECRET_ACCESS_KEY",
		"security_group_id": "GROUP_ID",
		"security_group_ids": []string{
			"SECURITY_GROUP",
		},
		"source_ami":              "SOURCE_AMI",
		"spot_price":              "auto",
		"spot_price_auto_product": "Linux/Unix",
		"ssh_keypair_name":        "myKeyPair",
		"ssh_private_ip":          true,
		"ssh_private_key_file":    "myKey",
		"ssh_username":            "vagrant",
		"subnet_id":               "subnet-12345def",
		"temporary_key_pair_name": "TMP_KEYPAIR",
		"tags": map[string]string{
			"OS_Version": "Ubuntu",
			"Release":    "Latest",
		},
		"type":                     "amazon-instance",
		"user_data":                "SOME_USER_DATA",
		"user_data_file":           "amazon-instance/amazon.userdata",
		"vpc_id":                   "VPC_ID",
		"windows_password_timeout": "10m",
		"x509_cert_path":           "/path/to/x509/cert",
		"x509_key_path":            "/path/to/x509/key",
		"x509_upload_path":         "/etc/x509",
	}
	contour.UpdateString("source_dir", "../test_files/src")
	bldr, err := testAllBuilders.createAmazonInstance("amazon-instance")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(bldr, expected) {
			t.Errorf("Expected %#v, got %#v", expected, bldr)
		}
	}
	// SSH
	expectedSSH := map[string]interface{}{
		"access_key":                   "AWS_ACCESS_KEY",
		"account_id":                   "YOUR_ACCOUNT_ID",
		"ami_description":              "AMI_DESCRIPTION",
		"ami_name":                     "AMI_NAME",
		"ami_virtualization_type":      "paravirtual",
		"associate_public_ip_address":  false,
		"availability_zone":            "us-east-1b",
		"bundle_destination":           "/tmp",
		"bundle_prefix":                "image--{{timestamp}}",
		"bundle_upload_command":        "sudo -n ec2-bundle-vol -k {{.KeyPath}} -u {{.AccountId}} -c {{.CertPath}} -r {{.Architecture}} -e {{.PrivatePath}} -d {{.Destination}} -p {{.Prefix}} --batch --no-filter",
		"bundle_vol_command":           "sudo -n ec2-upload-bundle -b {{.BucketName}} -m {{.ManifestPath}} -a {{.AccessKey}} -s {{.SecretKey}} -d {{.BundleDirectory}} --batch --region {{.Region}} --retry",
		"communicator":                 "ssh",
		"ebs_optimized":                true,
		"enhanced_networking":          false,
		"force_deregister":             false,
		"iam_instance_profile":         "INSTANCE_PROFILE",
		"instance_type":                "m3.medium",
		"region":                       "us-east-1",
		"s3_bucket":                    "packer_bucket",
		"secret_key":                   "AWS_SECRET_ACCESS_KEY",
		"security_group_id":            "GROUP_ID",
		"source_ami":                   "SOURCE_AMI",
		"spot_price":                   "auto",
		"spot_price_auto_product":      "Linux/Unix",
		"ssh_keypair_name":             "myKeyPair",
		"ssh_private_ip":               true,
		"ssh_bastion_host":             "bastion.host",
		"ssh_bastion_port":             2222,
		"ssh_bastion_username":         "packer",
		"ssh_bastion_password":         "packer",
		"ssh_bastion_private_key_file": "secret",
		"ssh_disable_agent":            true,
		"ssh_handshake_attempts":       10,
		"ssh_host":                     "127.0.0.1",
		"ssh_password":                 "vagrant",
		"ssh_port":                     22,
		"ssh_private_key_file":         "myKey",
		"ssh_pty":                      true,
		"ssh_username":                 "vagrant",
		"ssh_timeout":                  "10m",
		"subnet_id":                    "subnet-12345def",
		"temporary_key_pair_name":      "TMP_KEYPAIR",
		"type":             "amazon-instance",
		"user_data":        "SOME_USER_DATA",
		"user_data_file":   "amazon-instance/amazon.userdata",
		"vpc_id":           "VPC_ID",
		"x509_cert_path":   "/path/to/x509/cert",
		"x509_key_path":    "/path/to/x509/key",
		"x509_upload_path": "/etc/x509",
	}
	bldr, err = testAllBuildersSSH.createAmazonInstance("amazon-instance")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(bldr, expectedSSH) {
			t.Errorf("Expected %#v, got %#v", expectedSSH, bldr)
		}
	}
	// WinRM
	expectedWinRM := map[string]interface{}{
		"access_key":                  "AWS_ACCESS_KEY",
		"account_id":                  "YOUR_ACCOUNT_ID",
		"ami_description":             "AMI_DESCRIPTION",
		"ami_name":                    "AMI_NAME",
		"ami_virtualization_type":     "paravirtual",
		"associate_public_ip_address": false,
		"availability_zone":           "us-east-1b",
		"bundle_destination":          "/tmp",
		"bundle_prefix":               "image--{{timestamp}}",
		"bundle_upload_command":       "sudo -n ec2-bundle-vol -k {{.KeyPath}} -u {{.AccountId}} -c {{.CertPath}} -r {{.Architecture}} -e {{.PrivatePath}} -d {{.Destination}} -p {{.Prefix}} --batch --no-filter",
		"bundle_vol_command":          "sudo -n ec2-upload-bundle -b {{.BucketName}} -m {{.ManifestPath}} -a {{.AccessKey}} -s {{.SecretKey}} -d {{.BundleDirectory}} --batch --region {{.Region}} --retry",
		"communicator":                "winrm",
		"ebs_optimized":               true,
		"enhanced_networking":         false,
		"force_deregister":            false,
		"iam_instance_profile":        "INSTANCE_PROFILE",
		"instance_type":               "m3.medium",
		"region":                      "us-east-1",
		"s3_bucket":                   "packer_bucket",
		"secret_key":                  "AWS_SECRET_ACCESS_KEY",
		"security_group_id":           "GROUP_ID",
		"source_ami":                  "SOURCE_AMI",
		"spot_price":                  "auto",
		"spot_price_auto_product":     "Linux/Unix",
		"subnet_id":                   "subnet-12345def",
		"temporary_key_pair_name":     "TMP_KEYPAIR",
		"type":                     "amazon-instance",
		"user_data":                "SOME_USER_DATA",
		"user_data_file":           "amazon-instance/amazon.userdata",
		"vpc_id":                   "VPC_ID",
		"windows_password_timeout": "10m",
		"winrm_host":               "host",
		"winrm_password":           "vagrant",
		"winrm_port":               22,
		"winrm_timeout":            "10m",
		"winrm_username":           "vagrant",
		"winrm_use_ssl":            true,
		"winrm_insecure":           true,
		"x509_cert_path":           "/path/to/x509/cert",
		"x509_key_path":            "/path/to/x509/key",
		"x509_upload_path":         "/etc/x509",
	}
	bldr, err = testAllBuildersWinRM.createAmazonInstance("amazon-instance")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(bldr, expectedWinRM) {
			t.Errorf("Expected %#v, got %#v", expectedWinRM, bldr)
		}
	}
}

func TestCreateDigitalOcean(t *testing.T) {
	expected := map[string]interface{}{
		"api_token":          "DIGITALOCEAN_API_TOKEN",
		"droplet_name":       "ocean-drop",
		"image":              "ubuntu-12-04-x64",
		"private_networking": false,
		"region":             "nyc3",
		"size":               "512mb",
		"snapshot_name":      "my-snapshot",
		"state_timeout":      "6m",
		"type":               "digitalocean",
		"user_data":          "userdata",
	}
	bldr, err := testAllBuilders.createDigitalOcean("digitalocean")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(bldr, expected) {
			t.Errorf("Expected %#v, got %#v", expected, bldr)
		}
	}
	// SSH
	expectedSSH := map[string]interface{}{
		"api_token":                    "DIGITALOCEAN_API_TOKEN",
		"communicator":                 "ssh",
		"droplet_name":                 "ocean-drop",
		"image":                        "ubuntu-12-04-x64",
		"private_networking":           false,
		"region":                       "nyc3",
		"size":                         "512mb",
		"snapshot_name":                "my-snapshot",
		"ssh_bastion_host":             "bastion.host",
		"ssh_bastion_port":             2222,
		"ssh_bastion_username":         "packer",
		"ssh_bastion_password":         "packer",
		"ssh_bastion_private_key_file": "secret",
		"ssh_disable_agent":            true,
		"ssh_handshake_attempts":       10,
		"ssh_host":                     "127.0.0.1",
		"ssh_password":                 "vagrant",
		"ssh_port":                     22,
		"ssh_private_key_file":         "key/path",
		"ssh_pty":                      true,
		"ssh_username":                 "vagrant",
		"ssh_timeout":                  "10m",
		"state_timeout":                "6m",
		"type":                         "digitalocean",
		"user_data":                    "userdata",
	}
	bldr, err = testAllBuildersSSH.createDigitalOcean("digitalocean")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(bldr, expectedSSH) {
			t.Errorf("Expected %#v, got %#v", expectedSSH, bldr)
		}
	}

	// WinRM
	expectedWinRM := map[string]interface{}{
		"api_token":          "DIGITALOCEAN_API_TOKEN",
		"communicator":       "winrm",
		"droplet_name":       "ocean-drop",
		"image":              "ubuntu-12-04-x64",
		"private_networking": false,
		"region":             "nyc3",
		"size":               "512mb",
		"snapshot_name":      "my-snapshot",
		"state_timeout":      "6m",
		"type":               "digitalocean",
		"user_data":          "userdata",
		"winrm_host":         "host",
		"winrm_password":     "vagrant",
		"winrm_port":         22,
		"winrm_timeout":      "10m",
		"winrm_username":     "vagrant",
		"winrm_use_ssl":      true,
		"winrm_insecure":     true,
	}
	bldr, err = testAllBuildersWinRM.createDigitalOcean("digitalocean")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(bldr, expectedWinRM) {
			t.Errorf("Expected %#v, got %#v", expectedWinRM, bldr)
		}
	}
}

func TestCreateDocker(t *testing.T) {
	expected := map[string]interface{}{
		"commit":         true,
		"discard":        false,
		"export_path":    "export/path",
		"image":          "baseImage",
		"login":          true,
		"login_email":    "test@test.com",
		"login_username": "username",
		"login_password": "password",
		"login_server":   "127.0.0.1",
		"pull":           true,
		"run_command": []string{
			"-d",
			"-i",
			"-t",
			"{{.Image}}",
			"/bin/bash",
		},
		"type": "docker",
		"volumes": map[string]string{
			"/var/data1": "/var/data",
			"/var/www":   "/var/www",
		},
	}
	expectedCommand := map[string]interface{}{
		"commit":         true,
		"discard":        false,
		"export_path":    "export/path",
		"image":          "baseImage",
		"login":          true,
		"login_email":    "test@test.com",
		"login_username": "username",
		"login_password": "password",
		"login_server":   "127.0.0.1",
		"pull":           true,
		"run_command": []string{
			"-d",
			"-i",
			"-t",
			"{{.Image}}",
			"/bin/bash",
		},
		"type": "docker",
	}
	expectedCommandFile := map[string]interface{}{
		"commit":         true,
		"discard":        false,
		"export_path":    "export/path",
		"image":          "baseImage",
		"login":          true,
		"login_email":    "test@test.com",
		"login_username": "username",
		"login_password": "password",
		"login_server":   "127.0.0.1",
		"pull":           true,
		"run_command": []string{
			"-d",
			"-i",
			"-t",
			"{{.Image}}",
			"/bin/bash",
			"/invalid",
		},
		"type": "docker",
	}
	bldr, err := testAllBuilders.createDocker("docker")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(bldr, expected) {
			t.Errorf("Expected %#v, got %#v", expected, bldr)
		}
	}
	bldr, err = testDockerRunComandFile.createDocker("docker")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(bldr, expectedCommandFile) {
			t.Errorf("Expected %#v, got %#v", expectedCommandFile, bldr)
		}
	}
	bldr, err = testDockerRunComand.createDocker("docker")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(bldr, expectedCommand) {
			t.Errorf("Expected %#v, got %#v", expectedCommand, bldr)
		}
	}
	expectedSSH := map[string]interface{}{
		"commit":                       true,
		"communicator":                 "ssh",
		"discard":                      false,
		"export_path":                  "export/path",
		"image":                        "baseImage",
		"login":                        true,
		"login_email":                  "test@test.com",
		"login_username":               "username",
		"login_password":               "password",
		"login_server":                 "127.0.0.1",
		"pull":                         true,
		"ssh_bastion_host":             "bastion.host",
		"ssh_bastion_port":             2222,
		"ssh_bastion_username":         "packer",
		"ssh_bastion_password":         "packer",
		"ssh_bastion_private_key_file": "secret",
		"ssh_disable_agent":            true,
		"ssh_handshake_attempts":       10,
		"ssh_host":                     "127.0.0.1",
		"ssh_password":                 "vagrant",
		"ssh_port":                     22,
		"ssh_private_key_file":         "key/path",
		"ssh_pty":                      true,
		"ssh_username":                 "vagrant",
		"ssh_timeout":                  "10m",
		"type":                         "docker",
	}
	bldr, err = testAllBuildersSSH.createDocker("docker")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(bldr, expectedSSH) {
			t.Errorf("Expected %#v, got %#v", expectedSSH, bldr)
		}
	}
	expectedWinRM := map[string]interface{}{
		"commit":         true,
		"communicator":   "winrm",
		"discard":        false,
		"export_path":    "export/path",
		"image":          "baseImage",
		"login":          true,
		"login_email":    "test@test.com",
		"login_username": "username",
		"login_password": "password",
		"login_server":   "127.0.0.1",
		"pull":           true,
		"winrm_host":     "host",
		"winrm_password": "vagrant",
		"winrm_port":     22,
		"winrm_timeout":  "10m",
		"winrm_username": "vagrant",
		"winrm_use_ssl":  true,
		"winrm_insecure": true,
		"type":           "docker",
	}
	bldr, err = testAllBuildersWinRM.createDocker("docker")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(bldr, expectedWinRM) {
			t.Errorf("Expected %#v, got %#v", expectedWinRM, bldr)
		}
	}
}

func TestCreateGoogleCompute(t *testing.T) {
	expected := map[string]interface{}{
		"account_file":      "account.json",
		"address":           "ext-static",
		"disk_size":         20,
		"image_name":        "packer-{{timestamp}}",
		"image_description": "test image",
		"instance_name":     "packer-{{uuid}}",
		"machine_type":      "nl-standard-1",
		"metadata": map[string]string{
			"key-1": "value-1",
			"key-2": "value-2",
		},
		"network":       "default",
		"preemtible":    true,
		"project_id":    "projectID",
		"source_image":  "centos-6",
		"state_timeout": "5m",
		"tags": []string{
			"tag1",
		},
		"type":            "googlecompute",
		"use_internal_ip": true,
		"zone":            "us-central1-a",
	}

	bldr, err := testAllBuilders.createGoogleCompute("googlecompute")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(bldr, expected) {
			t.Errorf("Expected %#v, got %#v", expected, bldr)
		}
	}
	// ssh
	expectedSSH := map[string]interface{}{
		"account_file":                 "account.json",
		"address":                      "ext-static",
		"communicator":                 "ssh",
		"disk_size":                    20,
		"image_name":                   "packer-{{timestamp}}",
		"image_description":            "test image",
		"instance_name":                "packer-{{uuid}}",
		"machine_type":                 "nl-standard-1",
		"network":                      "default",
		"preemtible":                   true,
		"project_id":                   "projectID",
		"source_image":                 "centos-6",
		"ssh_bastion_host":             "bastion.host",
		"ssh_bastion_port":             2222,
		"ssh_bastion_username":         "packer",
		"ssh_bastion_password":         "packer",
		"ssh_bastion_private_key_file": "secret",
		"ssh_disable_agent":            true,
		"ssh_handshake_attempts":       10,
		"ssh_host":                     "127.0.0.1",
		"ssh_password":                 "vagrant",
		"ssh_port":                     22,
		"ssh_private_key_file":         "key/path",
		"ssh_pty":                      true,
		"ssh_username":                 "vagrant",
		"ssh_timeout":                  "10m",
		"state_timeout":                "5m",
		"type":                         "googlecompute",
		"use_internal_ip":              true,
		"zone":                         "us-central1-a",
	}

	bldr, err = testAllBuildersSSH.createGoogleCompute("googlecompute")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(bldr, expectedSSH) {
			t.Errorf("Expected %#v, got %#v", expectedSSH, bldr)
		}
	}

	expectedWinRM := map[string]interface{}{
		"account_file":      "account.json",
		"address":           "ext-static",
		"communicator":      "winrm",
		"disk_size":         20,
		"image_name":        "packer-{{timestamp}}",
		"image_description": "test image",
		"instance_name":     "packer-{{uuid}}",
		"machine_type":      "nl-standard-1",
		"network":           "default",
		"preemtible":        true,
		"project_id":        "projectID",
		"source_image":      "centos-6",
		"state_timeout":     "5m",
		"type":              "googlecompute",
		"use_internal_ip":   true,
		"winrm_host":        "host",
		"winrm_password":    "vagrant",
		"winrm_port":        22,
		"winrm_timeout":     "10m",
		"winrm_username":    "vagrant",
		"winrm_use_ssl":     true,
		"winrm_insecure":    true,
		"zone":              "us-central1-a",
	}

	bldr, err = testAllBuildersWinRM.createGoogleCompute("googlecompute")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(bldr, expectedWinRM) {
			t.Errorf("Expected %#v, got %#v", expectedWinRM, bldr)
		}
	}
}

func TestBuilderNull(t *testing.T) {
	// a communicator of none or no communicator setting should result in an error
	expected := "null: null: communicator: required setting not found"
	_, err := testAllBuilders.createNull("null")
	if err == nil {
		t.Errorf("expected an error, got none")
	} else {
		if err.Error() != expected {
			t.Errorf("got %q, want %q", err, expected)
		}
	}
	// ssh
	expectedSSH := map[string]interface{}{
		"communicator":                 "ssh",
		"ssh_bastion_host":             "bastion.host",
		"ssh_bastion_port":             2222,
		"ssh_bastion_username":         "packer",
		"ssh_bastion_password":         "packer",
		"ssh_bastion_private_key_file": "secret",
		"ssh_disable_agent":            true,
		"ssh_handshake_attempts":       10,
		"ssh_host":                     "127.0.0.1",
		"ssh_password":                 "vagrant",
		"ssh_port":                     22,
		"ssh_private_key_file":         "key/path",
		"ssh_pty":                      true,
		"ssh_username":                 "vagrant",
		"ssh_timeout":                  "10m",
		"type":                         "null",
	}
	bldr, err := testAllBuildersSSH.createNull("null")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(bldr, expectedSSH) {
			t.Errorf("Expected %#v, got %#v", expectedSSH, bldr)
		}
	}
	// winrm
	expectedWinRM := map[string]interface{}{
		"communicator":   "winrm",
		"winrm_host":     "host",
		"winrm_password": "vagrant",
		"winrm_port":     22,
		"winrm_timeout":  "10m",
		"winrm_username": "vagrant",
		"winrm_use_ssl":  true,
		"winrm_insecure": true,
		"type":           "null",
	}
	bldr, err = testAllBuildersWinRM.createNull("null")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(bldr, expectedWinRM) {
			t.Errorf("Expected %#v, got %#v", expectedWinRM, bldr)
		}
	}

}

func TestCreateOpenstack(t *testing.T) {
	// openstack1 uses tenant_id
	expected1 := map[string]interface{}{
		"api_key":           "APIKEY",
		"availability_zone": "zone1",
		"config_drive":      true,
		"flavor":            "2",
		"floating_ip":       "192.168.1.1",
		"floating_ip_pool":  "192.168.100.1/24",
		"image_name":        "test image",
		"insecure":          true,
		"networks": []string{
			"de305d54-75b4-431b-adb2-eb6b9e546014",
		},
		"metadata": map[string]interface{}{
			"metadata_listen":      "0.0.0.0",
			"quota_metadata_items": 128,
		},
		"password":         "packer",
		"rackconnect_wait": false,
		"region":           "DFW",
		"security_groups": []string{
			"admins",
		},
		"source_image":    "23b564c9-c3e6-49f9-bc68-86c7a9ab5018",
		"ssh_interface":   "private",
		"tenant_id":       "123",
		"type":            "openstack",
		"use_floating_ip": true,
		"username":        "packer",
	}
	ret, err := testAllBuilders.createOpenStack("openstack1")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(ret, expected1) {
			t.Errorf("Expected %#v, got %#v", expected1, ret)
		}
	}

	// openstack2 uses tenant_name
	expected2 := map[string]interface{}{
		"api_key":           "APIKEY",
		"availability_zone": "zone1",
		"config_drive":      true,
		"flavor":            "2",
		"floating_ip":       "192.168.1.1",
		"floating_ip_pool":  "192.168.100.1/24",
		"image_name":        "test image",
		"insecure":          true,
		"networks": []string{
			"de305d54-75b4-431b-adb2-eb6b9e546014",
		},
		"metadata": map[string]interface{}{
			"metadata_listen":      "0.0.0.0",
			"quota_metadata_items": 128,
		},
		"password":         "packer",
		"rackconnect_wait": false,
		"region":           "DFW",
		"security_groups": []string{
			"admins",
		},
		"source_image":    "23b564c9-c3e6-49f9-bc68-86c7a9ab5018",
		"ssh_interface":   "private",
		"tenant_name":     "acme",
		"type":            "openstack",
		"use_floating_ip": true,
		"username":        "packer",
	}
	ret, err = testAllBuilders.createOpenStack("openstack2")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(ret, expected2) {
			t.Errorf("Expected %#v, got %#v", expected2, ret)
		}
	}

	// ssh
	expectedSSH := map[string]interface{}{
		"api_key":                      "APIKEY",
		"availability_zone":            "zone1",
		"communicator":                 "ssh",
		"config_drive":                 true,
		"flavor":                       "2",
		"floating_ip":                  "192.168.1.1",
		"floating_ip_pool":             "192.168.100.1/24",
		"image_name":                   "test image",
		"insecure":                     true,
		"rackconnect_wait":             false,
		"region":                       "DFW",
		"source_image":                 "23b564c9-c3e6-49f9-bc68-86c7a9ab5018",
		"ssh_bastion_host":             "bastion.host",
		"ssh_bastion_port":             2222,
		"ssh_bastion_username":         "packer",
		"ssh_bastion_password":         "packer",
		"ssh_bastion_private_key_file": "secret",
		"ssh_disable_agent":            true,
		"ssh_handshake_attempts":       10,
		"ssh_host":                     "127.0.0.1",
		"ssh_interface":                "private",
		"ssh_password":                 "vagrant",
		"ssh_port":                     22,
		"ssh_private_key_file":         "key/path",
		"ssh_pty":                      true,
		"ssh_username":                 "vagrant",
		"ssh_timeout":                  "10m",
		"tenant_name":                  "acme",
		"type":                         "openstack",
		"use_floating_ip":              true,
	}
	ret, err = testAllBuildersSSH.createOpenStack("openstack")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(ret, expectedSSH) {
			t.Errorf("Expected %#v, got %#v", expectedSSH, ret)
		}
	}
	// winrm
	expectedWinRM := map[string]interface{}{
		"api_key":           "APIKEY",
		"availability_zone": "zone1",
		"communicator":      "winrm",
		"config_drive":      true,
		"flavor":            "2",
		"floating_ip":       "192.168.1.1",
		"floating_ip_pool":  "192.168.100.1/24",
		"image_name":        "test image",
		"insecure":          true,
		"rackconnect_wait":  false,
		"region":            "DFW",
		"source_image":      "23b564c9-c3e6-49f9-bc68-86c7a9ab5018",
		"tenant_name":       "acme",
		"type":              "openstack",
		"use_floating_ip":   true,
		"winrm_host":        "host",
		"winrm_password":    "vagrant",
		"winrm_port":        22,
		"winrm_timeout":     "10m",
		"winrm_username":    "vagrant",
		"winrm_use_ssl":     true,
		"winrm_insecure":    true,
	}
	ret, err = testAllBuildersWinRM.createOpenStack("openstack")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(ret, expectedWinRM) {
			t.Errorf("Expected %#v, got %#v", expectedWinRM, ret)
		}
	}
}

func TestCreateParallelsISO(t *testing.T) {
	expected := map[string]interface{}{
		"boot_command": []string{
			"<bs>",
			"<del>",
			"<enter><return>",
			"<esc>",
		},
		"boot_wait": "30s",
		"disk_size": 20000,
		"floppy_files": []string{
			"disk1",
		},
		"guest_os_type":        "ubuntu",
		"hard_drive_interface": "ide",
		"http_directory":       "http",
		"http_port_max":        9000,
		"http_port_min":        8000,
		"iso_checksum":         "ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
		"iso_checksum_type":    "sha256",
		"iso_target_path":      "packer_cache",
		"iso_urls": []string{
			"http://releases.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
			"http://2.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
		},
		"output_directory":           "out/dir",
		"parallels_tools_flavor":     "lin",
		"parallels_tools_guest_path": "ptools",
		"prlctl": [][]string{
			[]string{
				"set",
				"{{.Name}}",
				"--shf-host-add",
				"log",
				"--path",
				"{{pwd}}/log",
				"--mode",
				"rw",
				"--enable",
			},
			[]string{
				"set",
				"{{.Name}}",
				"--cpus",
				"1",
			},
		},
		"prlctl_post": [][]string{
			[]string{
				"set",
				"{{.Name}}",
				"--shf-host-del",
				"log",
			},
		},
		"prlctl_version_file": ".prlctl_version",
		"shutdown_command":    `shutdown /s /t 10 /f /d p:4:1 /c \"Packer Shutdown\"`,
		"shutdown_timeout":    "5m",
		"skip_compaction":     true,
		"ssh_username":        "vagrant",
		"type":                "parallels-iso",
		"vm_name":             "test-iso",
	}
	testAllBuilders.BaseURL = "http://releases.ubuntu.com/"
	settings, err := testAllBuilders.createParallelsISO("parallels-iso")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if !reflect.DeepEqual(settings, expected) {
			t.Errorf("Expected %#v, got %#v", expected, settings)
		}
	}
	// SSH
	expectedSSH := map[string]interface{}{
		"boot_wait":                    "30s",
		"communicator":                 "ssh",
		"disk_size":                    20000,
		"guest_os_type":                "ubuntu",
		"hard_drive_interface":         "ide",
		"http_directory":               "http",
		"http_port_max":                9000,
		"http_port_min":                8000,
		"iso_checksum":                 "ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
		"iso_checksum_type":            "sha256",
		"iso_target_path":              "packer_cache",
		"iso_url":                      "http://releases.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
		"output_directory":             "out/dir",
		"parallels_tools_flavor":       "lin",
		"parallels_tools_guest_path":   "ptools",
		"prlctl_version_file":          ".prlctl_version",
		"shutdown_command":             `shutdown /s /t 10 /f /d p:4:1 /c \"Packer Shutdown\"`,
		"shutdown_timeout":             "5m",
		"skip_compaction":              true,
		"ssh_bastion_host":             "bastion.host",
		"ssh_bastion_port":             2222,
		"ssh_bastion_username":         "packer",
		"ssh_bastion_password":         "packer",
		"ssh_bastion_private_key_file": "secret",
		"ssh_disable_agent":            true,
		"ssh_handshake_attempts":       10,
		"ssh_host":                     "127.0.0.1",
		"ssh_password":                 "vagrant",
		"ssh_port":                     22,
		"ssh_private_key_file":         "key/path",
		"ssh_pty":                      true,
		"ssh_username":                 "vagrant",
		"ssh_timeout":                  "10m", "type": "parallels-iso",
		"vm_name": "test-iso",
	}
	testAllBuildersSSH.BaseURL = "http://releases.ubuntu.com/"
	settings, err = testAllBuildersSSH.createParallelsISO("parallels-iso")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if !reflect.DeepEqual(settings, expectedSSH) {
			t.Errorf("Expected %#v, got %#v", expectedSSH, settings)
		}
	}
	// WinRM
	expectedWinRM := map[string]interface{}{
		"boot_wait":                  "30s",
		"communicator":               "winrm",
		"disk_size":                  20000,
		"guest_os_type":              "ubuntu",
		"hard_drive_interface":       "ide",
		"http_directory":             "http",
		"http_port_max":              9000,
		"http_port_min":              8000,
		"iso_checksum":               "ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
		"iso_checksum_type":          "sha256",
		"iso_target_path":            "packer_cache",
		"iso_url":                    "http://releases.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
		"output_directory":           "out/dir",
		"parallels_tools_flavor":     "lin",
		"parallels_tools_guest_path": "ptools",
		"prlctl_version_file":        ".prlctl_version",
		"shutdown_command":           `shutdown /s /t 10 /f /d p:4:1 /c \"Packer Shutdown\"`,
		"shutdown_timeout":           "5m",
		"skip_compaction":            true,
		"type":                       "parallels-iso",
		"vm_name":                    "test-iso",
		"winrm_host":                 "host",
		"winrm_password":             "vagrant",
		"winrm_port":                 22,
		"winrm_timeout":              "10m",
		"winrm_username":             "vagrant",
		"winrm_use_ssl":              true,
		"winrm_insecure":             true,
	}
	testAllBuildersWinRM.BaseURL = "http://releases.ubuntu.com/"
	settings, err = testAllBuildersWinRM.createParallelsISO("parallels-iso")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if !reflect.DeepEqual(settings, expectedWinRM) {
			t.Errorf("Expected %#v, got %#v", expectedWinRM, settings)
		}
	}
}

func TestCreateParallelsPVM(t *testing.T) {
	expected := map[string]interface{}{
		"boot_command": []string{
			"<bs>",
			"<del>",
			"<enter><return>",
			"<esc>",
		},
		"boot_wait": "30s",
		"floppy_files": []string{
			"disk1",
		},
		"output_directory":           "out/dir",
		"parallels_tools_flavor":     "lin",
		"parallels_tools_guest_path": "ptools",
		"parallels_tools_mode":       "upload",
		"parallels_tools_path":       "prl-tools.iso",
		"prlctl": [][]string{
			[]string{
				"set",
				"{{.Name}}",
				"--shf-host-add",
				"log",
				"--path",
				"{{pwd}}/log",
				"--mode",
				"rw",
				"--enable",
			},
			[]string{
				"set",
				"{{.Name}}",
				"--cpus",
				"1",
			},
		},
		"prlctl_post": [][]string{
			[]string{
				"set",
				"{{.Name}}",
				"--shf-host-del",
				"log",
			},
		},
		"prlctl_version_file": ".prlctl_version",
		"reassign_mac":        true,
		"shutdown_command":    `shutdown /s /t 10 /f /d p:4:1 /c \"Packer Shutdown\"`,
		"shutdown_timeout":    "5m",
		"skip_compaction":     true,
		"source_path":         "parallels-pvm/source.pvm",
		"ssh_username":        "vagrant",
		"type":                "parallels-pvm",
		"vm_name":             "test-iso",
	}
	testAllBuilders.BaseURL = "http://releases.ubuntu.com/"
	settings, err := testAllBuilders.createParallelsPVM("parallels-pvm")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if !reflect.DeepEqual(settings, expected) {
			t.Errorf("Expected %#v, got %#v", expected, settings)
		}
	}
	// SSH
	expectedSSH := map[string]interface{}{
		"boot_wait":                    "30s",
		"communicator":                 "ssh",
		"output_directory":             "out/dir",
		"parallels_tools_flavor":       "lin",
		"parallels_tools_guest_path":   "ptools",
		"parallels_tools_mode":         "upload",
		"parallels_tools_path":         "prl-tools.iso",
		"prlctl_version_file":          ".prlctl_version",
		"reassign_mac":                 true,
		"shutdown_command":             `shutdown /s /t 10 /f /d p:4:1 /c \"Packer Shutdown\"`,
		"shutdown_timeout":             "5m",
		"skip_compaction":              true,
		"source_path":                  "parallels-pvm/source.pvm",
		"ssh_bastion_host":             "bastion.host",
		"ssh_bastion_port":             2222,
		"ssh_bastion_username":         "packer",
		"ssh_bastion_password":         "packer",
		"ssh_bastion_private_key_file": "secret",
		"ssh_disable_agent":            true,
		"ssh_handshake_attempts":       10,
		"ssh_host":                     "127.0.0.1",
		"ssh_password":                 "vagrant",
		"ssh_port":                     22,
		"ssh_private_key_file":         "key/path",
		"ssh_pty":                      true,
		"ssh_timeout":                  "10m",
		"ssh_username":                 "vagrant",
		"type":                         "parallels-pvm",
		"vm_name":                      "test-iso",
	}
	testAllBuildersSSH.BaseURL = "http://releases.ubuntu.com/"
	settings, err = testAllBuildersSSH.createParallelsPVM("parallels-pvm")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if !reflect.DeepEqual(settings, expectedSSH) {
			t.Errorf("Expected %#v, got %#v", expectedSSH, settings)
		}
	}
	// WinRM
	expectedWinRM := map[string]interface{}{
		"boot_wait":                  "30s",
		"communicator":               "winrm",
		"output_directory":           "out/dir",
		"parallels_tools_flavor":     "lin",
		"parallels_tools_guest_path": "ptools",
		"parallels_tools_mode":       "upload",
		"parallels_tools_path":       "prl-tools.iso",
		"prlctl_version_file":        ".prlctl_version",
		"reassign_mac":               true,
		"shutdown_command":           `shutdown /s /t 10 /f /d p:4:1 /c \"Packer Shutdown\"`,
		"shutdown_timeout":           "5m",
		"skip_compaction":            true,
		"source_path":                "parallels-pvm/source.pvm",
		"type":                       "parallels-pvm",
		"vm_name":                    "test-iso",
		"winrm_host":                 "host",
		"winrm_password":             "vagrant",
		"winrm_port":                 22,
		"winrm_timeout":              "10m",
		"winrm_username":             "vagrant",
		"winrm_use_ssl":              true,
		"winrm_insecure":             true,
	}
	testAllBuildersWinRM.BaseURL = "http://releases.ubuntu.com/"
	settings, err = testAllBuildersWinRM.createParallelsPVM("parallels-pvm")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if !reflect.DeepEqual(settings, expectedWinRM) {
			t.Errorf("Expected %#v, got %#v", expectedWinRM, settings)
		}
	}
}

func TestCreateQEMU(t *testing.T) {
	expected := map[string]interface{}{
		"accelerator": "kvm",
		"boot_command": []string{
			"<bs>",
			"<del>",
			"<enter><return>",
			"<esc>",
		},
		"boot_wait":        "10s",
		"disk_cache":       "writeback",
		"disk_compression": true,
		"disk_discard":     "ignore",
		"disk_image":       true,
		"disk_interface":   "ide",
		"disk_size":        40000,
		"floppy_files": []string{
			"disk1",
		},
		"format":            "ovf",
		"headless":          true,
		"http_directory":    "http",
		"http_port_max":     9000,
		"http_port_min":     8000,
		"iso_checksum":      "ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
		"iso_checksum_type": "sha256",
		"iso_target_path":   "isocache",
		"iso_urls": []string{
			"http://releases.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
			"http://2.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
		},
		"net_device":       "i82551",
		"output_directory": "out/dir",
		"qemu_binary":      "qemu-system-x86_64",
		"qemuargs": [][]string{
			[]string{
				"-m",
				"1024m",
			},
			[]string{
				"--no-acpi",
				"",
			},
		},
		"skip_compaction": true,
		"ssh_username":    "vagrant",
		"type":            "qemu",
	}
	testAllBuilders.BaseURL = "http://releases.ubuntu.com/"
	settings, err := testAllBuilders.createQEMU("qemu")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if !reflect.DeepEqual(settings, expected) {
			t.Errorf("Expected %#v, got %#v", expected, settings)
		}
	}

	// SSH
	expectedSSH := map[string]interface{}{
		"accelerator":                  "kvm",
		"boot_wait":                    "10s",
		"communicator":                 "ssh",
		"disk_cache":                   "writeback",
		"disk_compression":             true,
		"disk_discard":                 "ignore",
		"disk_image":                   true,
		"disk_interface":               "ide",
		"disk_size":                    40000,
		"format":                       "ovf",
		"headless":                     true,
		"http_directory":               "http",
		"http_port_max":                9000,
		"http_port_min":                8000,
		"iso_checksum":                 "ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
		"iso_checksum_type":            "sha256",
		"iso_target_path":              "isocache",
		"iso_url":                      "http://releases.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
		"net_device":                   "i82551",
		"output_directory":             "out/dir",
		"qemu_binary":                  "qemu-system-x86_64",
		"skip_compaction":              true,
		"ssh_bastion_host":             "bastion.host",
		"ssh_bastion_port":             2222,
		"ssh_bastion_username":         "packer",
		"ssh_bastion_password":         "packer",
		"ssh_bastion_private_key_file": "secret",
		"ssh_disable_agent":            true,
		"ssh_handshake_attempts":       10,
		"ssh_host":                     "127.0.0.1",
		"ssh_password":                 "vagrant",
		"ssh_port":                     22,
		"ssh_private_key_file":         "key/path",
		"ssh_pty":                      true,
		"ssh_timeout":                  "10m",
		"ssh_username":                 "vagrant",
		"type":                         "qemu",
	}
	testAllBuilders.BaseURL = "http://releases.ubuntu.com/"
	settings, err = testAllBuildersSSH.createQEMU("qemu")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if !reflect.DeepEqual(settings, expectedSSH) {
			t.Errorf("Expected %#v, got %#v", expectedSSH, settings)
		}
	}
	// WinRM
	expectedWinRM := map[string]interface{}{
		"accelerator":       "kvm",
		"boot_wait":         "10s",
		"communicator":      "winrm",
		"disk_cache":        "writeback",
		"disk_compression":  true,
		"disk_discard":      "ignore",
		"disk_image":        true,
		"disk_interface":    "ide",
		"disk_size":         40000,
		"format":            "ovf",
		"headless":          true,
		"http_directory":    "http",
		"http_port_max":     9000,
		"http_port_min":     8000,
		"iso_checksum":      "ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
		"iso_checksum_type": "sha256",
		"iso_target_path":   "isocache",
		"iso_url":           "http://releases.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
		"net_device":        "i82551",
		"output_directory":  "out/dir",
		"qemu_binary":       "qemu-system-x86_64",
		"skip_compaction":   true,
		"type":              "qemu",
		"winrm_host":        "host",
		"winrm_password":    "vagrant",
		"winrm_port":        22,
		"winrm_timeout":     "10m",
		"winrm_username":    "vagrant",
		"winrm_use_ssl":     true,
		"winrm_insecure":    true,
	}
	testAllBuilders.BaseURL = "http://releases.ubuntu.com/"
	settings, err = testAllBuildersWinRM.createQEMU("qemu")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if !reflect.DeepEqual(settings, expectedWinRM) {
			t.Errorf("Expected %#v, got %#v", expectedWinRM, settings)
		}
	}
}

func TestCreateVirtualboxISO(t *testing.T) {
	expected := map[string]interface{}{
		"boot_command": []string{
			"<bs>",
			"<del>",
			"<enter><return>",
			"<esc>",
		},
		"boot_wait": "5s",
		"disk_size": 20000,
		"export_opts": []string{
			"opt1",
		},
		"floppy_files": []string{
			"disk1",
		},
		"format":                 "ovf",
		"guest_additions_mode":   "upload",
		"guest_additions_path":   "path/to/additions",
		"guest_additions_sha256": "89dac78769b26f8facf98ce85020a605b7601fec1946b0597e22ced5498b3597",
		"guest_additions_url":    "file://guest-additions",
		"guest_os_type":          "Ubuntu_64",
		"hard_drive_interface":   "ide",
		"headless":               true,
		"http_directory":         "http",
		"http_port_max":          9000,
		"http_port_min":          8000,
		"iso_checksum":           "ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
		"iso_checksum_type":      "sha256",
		"iso_interface":          "ide",
		"iso_urls": []string{
			"http://releases.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
			"http://2.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
		},
		"output_directory":  "out/dir",
		"shutdown_command":  "echo 'shutdown -P now' > /tmp/shutdown.sh; echo 'vagrant'|sudo -S sh '/tmp/shutdown.sh'",
		"shutdown_timeout":  "5m",
		"ssh_host_port_max": 40,
		"ssh_host_port_min": 22,
		"ssh_password":      "vagrant",
		"ssh_username":      "vagrant",
		"type":              "virtualbox-iso",
		"vboxmanage": [][]string{
			[]string{
				"modifyvm",
				"{{.Name}}",
				"--cpus",
				"1",
			},
			[]string{
				"modifyvm",
				"{{.Name}}",
				"--memory",
				"4096",
			},
		},
		"vboxmanage_post": [][]string{
			[]string{
				"modifyvm",
				"{{.Name}}",
				"--something",
				"value",
			},
		},
		"virtualbox_version_file": ".vbox_version",
		"vm_name":                 "test-vb-iso",
	}
	testAllBuilders.BaseURL = "http://releases.ubuntu.com/"
	settings, err := testAllBuilders.createVirtualBoxISO("virtualbox-iso")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if !reflect.DeepEqual(settings, expected) {
			t.Errorf("Expected %#v, got %#v", expected, settings)
		}
	}
	// ssh
	expectedSSH := map[string]interface{}{
		"boot_wait":                    "5s",
		"communicator":                 "ssh",
		"disk_size":                    20000,
		"format":                       "ovf",
		"guest_additions_mode":         "upload",
		"guest_additions_path":         "path/to/additions",
		"guest_additions_sha256":       "89dac78769b26f8facf98ce85020a605b7601fec1946b0597e22ced5498b3597",
		"guest_additions_url":          "file://guest-additions",
		"guest_os_type":                "Ubuntu_64",
		"hard_drive_interface":         "ide",
		"headless":                     true,
		"http_directory":               "http",
		"http_port_max":                9000,
		"http_port_min":                8000,
		"iso_checksum":                 "ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
		"iso_checksum_type":            "sha256",
		"iso_interface":                "ide",
		"iso_url":                      "http://releases.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
		"output_directory":             "out/dir",
		"shutdown_command":             "echo 'shutdown -P now' > /tmp/shutdown.sh; echo 'vagrant'|sudo -S sh '/tmp/shutdown.sh'",
		"shutdown_timeout":             "5m",
		"ssh_bastion_host":             "bastion.host",
		"ssh_bastion_port":             2222,
		"ssh_bastion_username":         "packer",
		"ssh_bastion_password":         "packer",
		"ssh_bastion_private_key_file": "secret",
		"ssh_disable_agent":            true,
		"ssh_handshake_attempts":       10,
		"ssh_host_port_max":            40,
		"ssh_host_port_min":            22,
		"ssh_host":                     "127.0.0.1",
		"ssh_password":                 "vagrant",
		"ssh_port":                     22,
		"ssh_private_key_file":         "key/path",
		"ssh_pty":                      true,
		"ssh_username":                 "vagrant",
		"ssh_timeout":                  "10m",
		"type":                         "virtualbox-iso",
		"virtualbox_version_file": ".vbox_version",
		"vm_name":                 "test-vb-iso",
	}
	testAllBuildersSSH.BaseURL = "http://releases.ubuntu.com/"
	settings, err = testAllBuildersSSH.createVirtualBoxISO("virtualbox-iso")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if !reflect.DeepEqual(settings, expectedSSH) {
			t.Errorf("Expected %#v, got %#v", expectedSSH, settings)
		}
	}

	// winrm communicator
	expectedWinRM := map[string]interface{}{
		"boot_wait":               "5s",
		"communicator":            "winrm",
		"disk_size":               20000,
		"format":                  "ovf",
		"guest_additions_mode":    "upload",
		"guest_additions_path":    "path/to/additions",
		"guest_additions_sha256":  "89dac78769b26f8facf98ce85020a605b7601fec1946b0597e22ced5498b3597",
		"guest_additions_url":     "file://guest-additions",
		"guest_os_type":           "Ubuntu_64",
		"hard_drive_interface":    "ide",
		"headless":                true,
		"http_directory":          "http",
		"http_port_max":           9000,
		"http_port_min":           8000,
		"iso_checksum":            "ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
		"iso_checksum_type":       "sha256",
		"iso_interface":           "ide",
		"iso_url":                 "http://releases.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
		"output_directory":        "out/dir",
		"shutdown_command":        "echo 'shutdown -P now' > /tmp/shutdown.sh; echo 'vagrant'|sudo -S sh '/tmp/shutdown.sh'",
		"shutdown_timeout":        "5m",
		"type":                    "virtualbox-iso",
		"winrm_host":              "host",
		"winrm_password":          "vagrant",
		"winrm_port":              22,
		"winrm_timeout":           "10m",
		"winrm_username":          "vagrant",
		"winrm_use_ssl":           true,
		"winrm_insecure":          true,
		"virtualbox_version_file": ".vbox_version",
		"vm_name":                 "test-vb-iso",
	}
	testAllBuildersWinRM.BaseURL = "http://releases.ubuntu.com/"
	settings, err = testAllBuildersWinRM.createVirtualBoxISO("virtualbox-iso")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if !reflect.DeepEqual(settings, expectedWinRM) {
			t.Errorf("Expected %#v, got %#v", expectedWinRM, settings)
		}
	}
}

func TestCreateVirtualboxOVF(t *testing.T) {
	expected := map[string]interface{}{
		"boot_command": []string{
			"<bs>",
			"<del>",
			"<enter><return>",
			"<esc>",
		},
		"boot_wait": "5s",
		"export_opts": []string{
			"opt1",
		},
		"floppy_files": []string{
			"disk1",
		},
		"format":                 "ovf",
		"guest_additions_mode":   "upload",
		"guest_additions_path":   "path/to/additions",
		"guest_additions_sha256": "89dac78769b26f8facf98ce85020a605b7601fec1946b0597e22ced5498b3597",
		"guest_additions_url":    "file://guest-additions",
		"headless":               true,
		"http_directory":         "http",
		"http_port_max":          9000,
		"http_port_min":          8000,
		"import_flags": []string{
			"--eula-accept",
		},
		"import_opts":          "keepallmacs",
		"output_directory":     "out/dir",
		"shutdown_command":     "echo 'shutdown -P now' > /tmp/shutdown.sh; echo 'vagrant'|sudo -S sh '/tmp/shutdown.sh'",
		"shutdown_timeout":     "5m",
		"source_path":          "virtualbox-ovf/source.ova",
		"ssh_host_port_max":    40,
		"ssh_host_port_min":    22,
		"ssh_skip_nat_mapping": true,
		"ssh_username":         "vagrant",
		"type":                 "virtualbox-ovf",
		"vboxmanage": [][]string{
			[]string{
				"modifyvm",
				"{{.Name}}",
				"--cpus",
				"1",
			},
			[]string{
				"modifyvm",
				"{{.Name}}",
				"--memory",
				"4096",
			},
		},
		"vboxmanage_post": [][]string{
			[]string{
				"modifyvm",
				"{{.Name}}",
				"--something",
				"value",
			},
		},
		"virtualbox_version_file": ".vbox_version",
		"vm_name":                 "test-vb-ovf",
	}
	testAllBuilders.Files = make(map[string]string)
	settings, err := testAllBuilders.createVirtualBoxOVF("virtualbox-ovf")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(settings, expected) {
			t.Errorf("Expected %#v, got %#v", expected, settings)
		}
	}
	// ssh
	expectedSSH := map[string]interface{}{
		"boot_wait":                    "5s",
		"communicator":                 "ssh",
		"format":                       "ovf",
		"guest_additions_mode":         "upload",
		"guest_additions_path":         "path/to/additions",
		"guest_additions_sha256":       "89dac78769b26f8facf98ce85020a605b7601fec1946b0597e22ced5498b3597",
		"guest_additions_url":          "file://guest-additions",
		"headless":                     true,
		"http_directory":               "http",
		"http_port_max":                9000,
		"http_port_min":                8000,
		"import_opts":                  "keepallmacs",
		"output_directory":             "out/dir",
		"shutdown_command":             "echo 'shutdown -P now' > /tmp/shutdown.sh; echo 'vagrant'|sudo -S sh '/tmp/shutdown.sh'",
		"shutdown_timeout":             "5m",
		"source_path":                  "virtualbox-ovf/source.ova",
		"ssh_bastion_host":             "bastion.host",
		"ssh_bastion_port":             2222,
		"ssh_bastion_username":         "packer",
		"ssh_bastion_password":         "packer",
		"ssh_bastion_private_key_file": "secret",
		"ssh_disable_agent":            true,
		"ssh_handshake_attempts":       10,
		"ssh_host_port_max":            40,
		"ssh_host_port_min":            22,
		"ssh_host":                     "127.0.0.1",
		"ssh_password":                 "vagrant",
		"ssh_port":                     22,
		"ssh_private_key_file":         "key/path",
		"ssh_pty":                      true,
		"ssh_skip_nat_mapping":         true,
		"ssh_username":                 "vagrant",
		"ssh_timeout":                  "10m",
		"type":                         "virtualbox-ovf",
		"virtualbox_version_file": ".vbox_version",
		"vm_name":                 "test-vb-ovf",
	}
	testAllBuildersSSH.BaseURL = "http://releases.ubuntu.com/"
	settings, err = testAllBuildersSSH.createVirtualBoxOVF("virtualbox-ovf")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if !reflect.DeepEqual(settings, expectedSSH) {
			t.Errorf("Expected %#v, got %#v", expectedSSH, settings)
		}
	}

	// winrm communicator
	expectedWinRM := map[string]interface{}{
		"boot_wait":               "5s",
		"communicator":            "winrm",
		"format":                  "ovf",
		"guest_additions_mode":    "upload",
		"guest_additions_path":    "path/to/additions",
		"guest_additions_sha256":  "89dac78769b26f8facf98ce85020a605b7601fec1946b0597e22ced5498b3597",
		"guest_additions_url":     "file://guest-additions",
		"headless":                true,
		"http_directory":          "http",
		"http_port_max":           9000,
		"http_port_min":           8000,
		"import_opts":             "keepallmacs",
		"output_directory":        "out/dir",
		"shutdown_command":        "echo 'shutdown -P now' > /tmp/shutdown.sh; echo 'vagrant'|sudo -S sh '/tmp/shutdown.sh'",
		"shutdown_timeout":        "5m",
		"source_path":             "virtualbox-ovf/source.ova",
		"type":                    "virtualbox-ovf",
		"winrm_host":              "host",
		"winrm_password":          "vagrant",
		"winrm_port":              22,
		"winrm_timeout":           "10m",
		"winrm_username":          "vagrant",
		"winrm_use_ssl":           true,
		"winrm_insecure":          true,
		"virtualbox_version_file": ".vbox_version",
		"vm_name":                 "test-vb-ovf",
	}
	testAllBuildersWinRM.BaseURL = "http://releases.ubuntu.com/"
	settings, err = testAllBuildersWinRM.createVirtualBoxOVF("virtualbox-ovf")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err.Error())
	} else {
		if !reflect.DeepEqual(settings, expectedWinRM) {
			t.Errorf("Expected %#v, got %#v", expectedWinRM, settings)
		}
	}
}

func TestCreateVMWareISO(t *testing.T) {
	expected := map[string]interface{}{
		"boot_command": []string{
			"<bs>",
			"<del>",
			"<enter><return>",
			"<esc>",
		},
		"boot_wait":    "5s",
		"communicator": "none",
		"disk_additional_size": []int{
			10000,
		},
		"disk_size":    20000,
		"disk_type_id": "1",
		"floppy_files": []string{
			"disk1",
		},
		"fusion_app_path":   "/Applications/VMware Fusion.app",
		"guest_os_type":     "Ubuntu_64",
		"headless":          true,
		"http_directory":    "http",
		"http_port_max":     9000,
		"http_port_min":     8000,
		"iso_checksum":      "ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
		"iso_checksum_type": "sha256",
		"iso_target_path":   "../isocache/",
		"iso_urls": []string{
			"http://releases.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
			"http://2.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
		},
		"output_directory":        "out/dir",
		"remote_cache_datastore":  "datastore1",
		"remote_cache_directory":  "packer_cache",
		"remote_datastore":        "datastore1",
		"remote_host":             "remoteHost",
		"remote_password":         "rpassword",
		"remote_private_key_file": "secret",
		"remote_type":             "esx5",
		"shutdown_command":        "echo 'shutdown -P now' > /tmp/shutdown.sh; echo 'vagrant'|sudo -S sh '/tmp/shutdown.sh'",
		"shutdown_timeout":        "5m",
		"skip_compaction":         true,
		"ssh_username":            "vagrant",
		"tools_upload_flavor":     "linux",
		"tools_upload_path":       "{{.Flavor}}.iso",
		"type":                    "vmware-iso",
		"version":                 "9",
		"vmx_data": map[string]string{
			"cpuid.coresPerSocket": "1",
			"memsize":              "1024",
			"numvcpus":             "1",
		},
		"vmx_data_post": map[string]string{
			"something": "value",
		},
		"vm_name":           "packer-BUILDNAME",
		"vmdk_name":         "packer",
		"vmx_template_path": "template/path",
		"vnc_port_max":      6000,
		"vnc_port_min":      5900,
	}

	testAllBuilders.BaseURL = "http://releases.ubuntu.com/"
	settings, err := testAllBuilders.createVMWareISO("vmware-iso")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(settings, expected) {
			t.Errorf("Expected %#v, got %#v", expected, settings)
		}
	}
	// SSH
	expectedSSH := map[string]interface{}{
		"boot_wait":                    "5s",
		"communicator":                 "ssh",
		"disk_size":                    20000,
		"disk_type_id":                 "1",
		"fusion_app_path":              "/Applications/VMware Fusion.app",
		"guest_os_type":                "Ubuntu_64",
		"headless":                     true,
		"http_directory":               "http",
		"http_port_max":                9000,
		"http_port_min":                8000,
		"iso_checksum":                 "ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
		"iso_checksum_type":            "sha256",
		"iso_target_path":              "../isocache/",
		"iso_url":                      "http://releases.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
		"output_directory":             "out/dir",
		"remote_cache_datastore":       "datastore1",
		"remote_cache_directory":       "packer_cache",
		"remote_datastore":             "datastore1",
		"remote_host":                  "remoteHost",
		"remote_password":              "rpassword",
		"remote_private_key_file":      "secret",
		"remote_type":                  "esx5",
		"shutdown_command":             "echo 'shutdown -P now' > /tmp/shutdown.sh; echo 'vagrant'|sudo -S sh '/tmp/shutdown.sh'",
		"shutdown_timeout":             "5m",
		"skip_compaction":              true,
		"tools_upload_flavor":          "linux",
		"tools_upload_path":            "{{.Flavor}}.iso",
		"type":                         "vmware-iso",
		"ssh_bastion_host":             "bastion.host",
		"ssh_bastion_port":             2222,
		"ssh_bastion_username":         "packer",
		"ssh_bastion_password":         "packer",
		"ssh_bastion_private_key_file": "secret",
		"ssh_disable_agent":            true,
		"ssh_handshake_attempts":       10,
		"ssh_host":                     "127.0.0.1",
		"ssh_password":                 "vagrant",
		"ssh_port":                     22,
		"ssh_private_key_file":         "key/path",
		"ssh_pty":                      true,
		"ssh_username":                 "vagrant",
		"ssh_timeout":                  "10m",
		"version":                      "9",
		"vm_name":                      "packer-BUILDNAME",
		"vmdk_name":                    "packer",
		"vmx_template_path":            "template/path",
		"vnc_port_max":                 6000,
		"vnc_port_min":                 5900,
	}

	testAllBuildersSSH.BaseURL = "http://releases.ubuntu.com/"
	settings, err = testAllBuildersSSH.createVMWareISO("vmware-iso")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(settings, expectedSSH) {
			t.Errorf("Expected %#v, got %#v", expectedSSH, settings)
		}
	}
	// winrm
	expectedWinRM := map[string]interface{}{
		"boot_wait":               "5s",
		"communicator":            "winrm",
		"disk_size":               20000,
		"disk_type_id":            "1",
		"fusion_app_path":         "/Applications/VMware Fusion.app",
		"guest_os_type":           "Ubuntu_64",
		"headless":                true,
		"http_directory":          "http",
		"http_port_max":           9000,
		"http_port_min":           8000,
		"iso_checksum":            "ababb88a492e08759fddcf4f05e5ccc58ec9d47fa37550d63931d0a5fa4f7388",
		"iso_checksum_type":       "sha256",
		"iso_target_path":         "../isocache/",
		"iso_url":                 "http://releases.ubuntu.com/14.04/ubuntu-14.04.1-server-amd64.iso",
		"output_directory":        "out/dir",
		"remote_cache_datastore":  "datastore1",
		"remote_cache_directory":  "packer_cache",
		"remote_datastore":        "datastore1",
		"remote_host":             "remoteHost",
		"remote_password":         "rpassword",
		"remote_private_key_file": "secret",
		"remote_type":             "esx5",
		"shutdown_command":        "echo 'shutdown -P now' > /tmp/shutdown.sh; echo 'vagrant'|sudo -S sh '/tmp/shutdown.sh'",
		"shutdown_timeout":        "5m",
		"skip_compaction":         true,
		"tools_upload_flavor":     "linux",
		"tools_upload_path":       "{{.Flavor}}.iso",
		"type":                    "vmware-iso",
		"version":                 "9",
		"vm_name":                 "packer-BUILDNAME",
		"vmdk_name":               "packer",
		"vmx_template_path":       "template/path",
		"vnc_port_max":            6000,
		"vnc_port_min":            5900,
		"winrm_host":              "host",
		"winrm_password":          "vagrant",
		"winrm_port":              22,
		"winrm_timeout":           "10m",
		"winrm_username":          "vagrant",
		"winrm_use_ssl":           true,
		"winrm_insecure":          true,
	}

	testAllBuildersWinRM.BaseURL = "http://releases.ubuntu.com/"
	settings, err = testAllBuildersWinRM.createVMWareISO("vmware-iso")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(settings, expectedWinRM) {
			t.Errorf("Expected %#v, got %#v", expectedWinRM, settings)
		}
	}
}

func TestCreateVMWareVMX(t *testing.T) {
	expected := map[string]interface{}{
		"boot_command": []string{
			"<bs>",
			"<del>",
			"<enter><return>",
			"<esc>",
		},
		"boot_wait": "5s",
		"floppy_files": []string{
			"disk1",
		},
		"fusion_app_path":  "/Applications/VMware Fusion.app",
		"headless":         true,
		"http_directory":   "http",
		"http_port_max":    9000,
		"http_port_min":    8000,
		"output_directory": "out/dir",
		"shutdown_command": "echo 'shutdown -P now' > /tmp/shutdown.sh; echo 'vagrant'|sudo -S sh '/tmp/shutdown.sh'",
		"shutdown_timeout": "5m",
		"skip_compaction":  false,
		"source_path":      "vmware-vmx/source.vmx",
		"ssh_username":     "vagrant",
		"type":             "vmware-vmx",
		"vmx_data": map[string]string{
			"cpuid.coresPerSocket": "1",
			"memsize":              "1024",
			"numvcpus":             "1",
		},
		"vmx_data_post": map[string]string{
			"something": "value",
		},
		"vm_name":      "packer-BUILDNAME",
		"vnc_port_max": 6000,
		"vnc_port_min": 5900,
	}

	settings, err := testAllBuilders.createVMWareVMX("vmware-vmx")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(settings, expected) {
			t.Errorf("Expected %#v, got %#v", expected, settings)
		}
	}

	expectedSSH := map[string]interface{}{
		"boot_wait":                    "5s",
		"communicator":                 "ssh",
		"fusion_app_path":              "/Applications/VMware Fusion.app",
		"headless":                     true,
		"http_directory":               "http",
		"http_port_max":                9000,
		"http_port_min":                8000,
		"output_directory":             "out/dir",
		"shutdown_command":             "echo 'shutdown -P now' > /tmp/shutdown.sh; echo 'vagrant'|sudo -S sh '/tmp/shutdown.sh'",
		"shutdown_timeout":             "5m",
		"skip_compaction":              false,
		"source_path":                  "vmware-vmx/source.vmx",
		"ssh_bastion_host":             "bastion.host",
		"ssh_bastion_port":             2222,
		"ssh_bastion_username":         "packer",
		"ssh_bastion_password":         "packer",
		"ssh_bastion_private_key_file": "secret",
		"ssh_disable_agent":            true,
		"ssh_handshake_attempts":       10,
		"ssh_host":                     "127.0.0.1",
		"ssh_password":                 "vagrant",
		"ssh_port":                     22,
		"ssh_private_key_file":         "key/path",
		"ssh_pty":                      true,
		"ssh_username":                 "vagrant",
		"ssh_timeout":                  "10m",
		"type":                         "vmware-vmx",
		"vm_name":                      "packer-BUILDNAME",
		"vnc_port_max":                 6000,
		"vnc_port_min":                 5900,
	}

	settings, err = testAllBuildersSSH.createVMWareVMX("vmware-vmx")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(settings, expectedSSH) {
			t.Errorf("Expected %#v, got %#v", expectedSSH, settings)
		}
	}
	// WinRM
	expectedWinRM := map[string]interface{}{
		"boot_wait":        "5s",
		"communicator":     "winrm",
		"fusion_app_path":  "/Applications/VMware Fusion.app",
		"headless":         true,
		"http_directory":   "http",
		"http_port_max":    9000,
		"http_port_min":    8000,
		"output_directory": "out/dir",
		"shutdown_command": "echo 'shutdown -P now' > /tmp/shutdown.sh; echo 'vagrant'|sudo -S sh '/tmp/shutdown.sh'",
		"shutdown_timeout": "5m",
		"skip_compaction":  false,
		"source_path":      "vmware-vmx/source.vmx",
		"type":             "vmware-vmx",
		"vm_name":          "packer-BUILDNAME",
		"vnc_port_max":     6000,
		"vnc_port_min":     5900,
		"winrm_host":       "host",
		"winrm_password":   "vagrant",
		"winrm_port":       22,
		"winrm_timeout":    "10m",
		"winrm_username":   "vagrant",
		"winrm_use_ssl":    true,
		"winrm_insecure":   true,
	}

	settings, err = testAllBuildersWinRM.createVMWareVMX("vmware-vmx")
	if err != nil {
		t.Errorf("Expected error to be nil, got %q", err)
	} else {
		if !reflect.DeepEqual(settings, expectedWinRM) {
			t.Errorf("Expected %#v, got %#v", expectedWinRM, settings)
		}
	}
}

func TestDeepCopyMapStringBuilder(t *testing.T) {
	cpy := DeepCopyMapStringBuilderC(testDistroDefaults.Templates[Ubuntu].Builders)
	if !reflect.DeepEqual(cpy["common"], testDistroDefaults.Templates[Ubuntu].Builders["common"]) {
		t.Errorf("Expected %#v, got %#v", testDistroDefaults.Templates[Ubuntu].Builders["common"], cpy["common"])
	}
}

func TestProcessAMIBlockDeviceMappings(t *testing.T) {
	r := newRawTemplate()
	_, err := r.processAMIBlockDeviceMappings([]string{})
	if err == nil {
		t.Errorf("expected error, got none")
	} else {
		expected := "ami_block_device_mappings: \"\": not in a supported format"
		if err.Error() != expected {
			t.Errorf("got %q; expected %q", err, expected)
		}
	}

	expected := []map[string]interface{}{
		{
			"delete_on_termination": true,
			"device_name":           "/dev/sdb",
			"encrypted":             true,
			"iops":                  1000,
			"no_device":             false,
			"snapshot_id":           "SNAPSHOT",
			"virtual_name":          "/ephemeral0",
			"volume_type":           "io1",
			"volume_size":           20,
		},
		{
			"device_name":  "/dev/sdc",
			"iops":         500,
			"virtual_name": "/ephemeral1",
			"volume_type":  "io1",
			"volume_size":  10,
		},
	}
	// test using array of block mappings: []map[string]interface{}
	mappings := []map[string]interface{}{
		{
			"delete_on_termination": true,
			"device_name":           "/dev/sdb",
			"encrypted":             true,
			"iops":                  1000,
			"no_device":             false,
			"snapshot_id":           "SNAPSHOT",
			"virtual_name":          "/ephemeral0",
			"volume_type":           "io1",
			"volume_size":           20,
		},
		{
			"device_name":  "/dev/sdc",
			"iops":         500,
			"virtual_name": "/ephemeral1",
			"volume_type":  "io1",
			"volume_size":  10,
		},
	}
	// expected is the same as mappings
	ret, err := r.processAMIBlockDeviceMappings(mappings)
	if err != nil {
		t.Errorf("got %q, expected no error", err)
	} else {
		if !reflect.DeepEqual(expected, ret) {
			t.Errorf("Got %#v; want %#v", ret, expected)
		}
	}
	// test using array of block mappings: [][]string
	mappingsSlice := [][]string{
		[]string{
			"delete_on_termination=true",
			"device_name=/dev/sdb",
			"encrypted=true",
			"iops=1000",
			"no_device=false",
			"snapshot_id=SNAPSHOT",
			"virtual_name=/ephemeral0",
			"volume_type=io1",
			"volume_size=20",
		},
		[]string{
			"device_name=/dev/sdc",
			"iops=500",
			"virtual_name=/ephemeral1",
			"volume_type=io1",
			"volume_size=10",
		},
	}
	ret, err = r.processAMIBlockDeviceMappings(mappingsSlice)
	if err != nil {
		t.Errorf("got %q, expected no error", err)
	} else {
		if !reflect.DeepEqual(expected, ret) {
			t.Errorf("Got %#v; want %#v", ret, expected)
		}
	}
}

func TestCommandFromSlice(t *testing.T) {
	tests := []struct {
		lines    []string
		expected string
	}{
		{[]string{}, ""},
		{[]string{"hello"}, "hello"},
		{[]string{"hello \\ ", "world \\ ", "!"}, "hello world !"},
		{[]string{"hello \\ ", "world  ", "!"}, "hello world"},
		{[]string{"sudo -i -n ec2-bundle-vol \\ ",
			"	-k {{.KeyPath}} \\ ",
			"  -u {{.AccountId}} \\ ",
			" -c {{.CertPath}} \\ ",
			"   -r {{.Architecture}} \\",
			"	-e {{.PrivatePath}}/* \\ ",
			"-d {{.Destination}} \\ ",
			"	-p {{.Prefix}} \\ ",
			"	  --batch \\  ",
			"	--no-filter	"},
			"sudo -i -n ec2-bundle-vol -k {{.KeyPath}} -u {{.AccountId}} -c {{.CertPath}} -r {{.Architecture}} -e {{.PrivatePath}}/* -d {{.Destination}} -p {{.Prefix}} --batch --no-filter",
		},
	}
	for i, test := range tests {
		ret := commandFromSlice(test.lines)
		if ret != test.expected {
			t.Errorf("%d: got %s want %s", i, ret, test.expected)
		}
	}
}

func TestStringIsCommandFilename(t *testing.T) {
	tests := []struct {
		val      string
		expected bool
	}{
		{"", false},
		{"command", false},
		{".command", false},
		{"test", false},
		{"test.", false},
		{"test.command", true},
		{"a.test.command", true},
		{"a\\b\\test.command", true},
	}

	for i, test := range tests {
		b := stringIsCommandFilename(test.val)
		if b != test.expected {
			t.Errorf("%d: got %t, want %t", i, b, test.expected)
		}
	}
}
