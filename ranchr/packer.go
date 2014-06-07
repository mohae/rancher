package ranchr


import (
	_"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	_"reflect"
	"strconv"
	_"strings"
	_"time"
)

type packerer interface {
	mergeSettings([]string)
}

type builderer interface {
	mergeVMSettings([]string)
}

type PackerTemplate struct {
	Description      string                 `json:"description"`
	MinPackerVersion string                 `json:"min_packer_version"`
	Builders         []interface{}          `json:"builders"`
	PostProcessors   []interface{}          `json:"post-processors"`
	Provisioners     []interface{}          `json:"provisioners"`
	Variables        map[string]interface{} `json:"variables"`
}

func (p *PackerTemplate) TemplateToFileJSON(i IODirInf, b BuildInf, scripts []string) error {

	if i.OutDir == "" {
		err := errors.New("ranchr.TemplateToFileJSON: output directory for " + b.BuildName + " not set")
		Log.Error(err.Error())
		return err
	}

	if i.SrcDir == "" {
		err := errors.New("ranchr.TemplateToFileJSON: SrcDir directory for " + b.BuildName + " not set")
		Log.Error(err.Error())
		return err
	}

	if i.ScriptsDir == "" {
		err := errors.New("ranchr.TemplateToFileJSON: ScriptsDir directory for " + b.BuildName + " not set")
		Log.Error(err.Error())
		return err
	}

	// if the OutDir exists, archive it before processing
	archiveDir(i.OutDir)

	destDir := i.OutDir + "scripts/"
	var errCnt, okCnt int
	for _, script := range scripts {

		if wB, err := copyFile(i.ScriptsDir, destDir, script); err != nil {
			Log.Error(err.Error())
			errCnt++
		} else {
			Log.Info(strconv.FormatInt(wB, 10) + " Bytes were copied from " + i.ScriptsDir + script + " to " + destDir + script)
			okCnt++
		}
	}

	if errCnt > 0 {
		fmt.Printf("Copy of scripts for build, %s, had %s errors. There were %s scripts that were copied without error.", b.BuildName, string(errCnt), string(okCnt))
		Log.Error("Copy of scripts for build, " + b.BuildName + ", had " + strconv.Itoa(errCnt) + " errors. There were " + strconv.Itoa(okCnt) + " scripts that were copied without error.")
	} else {
		Log.Info(strconv.Itoa(okCnt) + " scripts were successfully copied.")
	}
	
	if err := os.MkdirAll(i.OutDir + "http", os.FileMode(0766)); err != nil {
		Log.Error(err.Error())
		return err
	}

	// Write it out as JSON
	tplJSON, err := json.MarshalIndent(p, "", "\t")
	if err != nil {
		Log.Error("Marshalling of the Packer Template failed: " + err.Error())
		return err
	}
	
	f, err := os.Create(i.OutDir + b.Name + ".json")
	if err != nil {
		Log.Error(err.Error())
		return err
	}
	defer f.Close()
	
	_, err = io.WriteString(f,  string(tplJSON[:]))
	if err != nil {
		Log.Error(err.Error())
		return err
	}

	return nil
}
