package app

import (
	"testing"
)

// TODO rewrite with refactor
func TestCreate(t *testing.T) {
	p := PackerTemplate{}
	//	var scripts map[string]string
	//	b := BuildInf{BuildName: "test build"}
	//	i := IODirInf{}
	//	err := p.create(i, b, scripts)
	_ = p
	//	if err != nil {
	//		t.Errorf("Expected no error, got %q", err)
	//	}
	/*
		i = IODirInf{HTTPDir: "http"}
		err = p.create(i, b, scripts)
		if err == nil {
			t.Error("Expected \"HTTPSrcDir directory not set\", error was nil")
		} else {
			if err.Error() != "HTTPSrcDir directory not set" {
				t.Errorf("Expected \"HTTPSrcDir directory not set\", got %q", err.Error())
			}
		}

		i = IODirInf{HTTPDir: "http", HTTPSrcDir: "../test_files/http/"}
		err = p.create(i, b, scripts)
		if err == nil {
			t.Error("Expected \"output directory not set\", error was nil")
		} else {
			if err.Error() != "output directory not set" {
				t.Errorf("Expected \"output directory not set\", got %q", err.Error())
			}
		}

		i = IODirInf{HTTPDir: "http", HTTPSrcDir: "../test_files/http/", OutDir: "../test_files/out/"}
		err = p.create(i, b, scripts)
		if err == nil {
			t.Error("Expected \"SrcDir directory not set\", error was nil")
		} else {
			if err.Error() != "SrcDir directory not set" {
				t.Errorf("Expected \"SrcDir directory not set\", got %q", err.Error())
			}
		}

		i = IODirInf{HTTPDir: "http", HTTPSrcDir: "../test_files/http/", OutDir: "../test_files/out/", SrcDir: "../test_files/"}
		err = p.create(i, b, scripts)
		if err == nil {
			t.Error("Expected \"ScriptsDir directory not set\", error was nil")
		} else {
			if err.Error() != "ScriptsDir directory not set" {
				t.Errorf("Expected \"ScriptsDir directory not set\", got %q", err.Error())
			}
		}

		i = IODirInf{HTTPDir: "http", HTTPSrcDir: "../test_files/http/", OutDir: "../test_files/out/", SrcDir: "../test_files/", ScriptsDir: "scripts"}
		err = p.create(i, b, scripts)
		if err == nil {
			t.Error("Expected \"	ScriptsSrcDir directory not set\", error was nil")
		} else {
			if err.Error() != "ScriptsSrcDir directory not set" {
				t.Errorf("Expected \"ScriptsSrcDir directory not set\", got %q", err.Error())
			}
		}

		i = IODirInf{HTTPDir: "http", HTTPSrcDir: "../test_files/http/", OutDir: "../test_files/out/build/", SrcDir: "../test_files/", ScriptsDir: "scripts", ScriptsSrcDir: "../test_files/scripts/"}
		Scripts["../test_files/scripts"] := []string{"cleanup_test.sh", "setup_test.sh", "not_there.sh", "missing.sh", "test_file.sh"}
		err = p.create(i, b, Scripts)
		if err == nil {
			t.Error("Expected \"open ../test_files/scripts/test_file.sh: no such file or directory\", error was nil")
		} else {
			if err.Error() != "open ../test_files/scripts/test_file.sh: no such file or directory" {
				t.Errorf("Expected \"open ../test_files/scripts/test_file.sh: no such file or directory\", got %q", err.Error())
			}
		}
	*/
}
