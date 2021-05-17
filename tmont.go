package tmont

import (
	"fmt"
	"io/fs"
	"os"
	"time"

	ht "html/template"
)

type TMonitor struct {
	htmlTemplates map[string]*ht.Template

	lastUpdated map[string]fs.FileInfo
	files       []string
	lock        chan struct{}
	errMsg      chan error
}

func New(files ...string) *TMonitor {
	tmon := &TMonitor{
		htmlTemplates: make(map[string]*ht.Template),
		lastUpdated:   make(map[string]fs.FileInfo),
		files:         files,
	}
	tmon.lock = make(chan struct{}, 1)
	tmon.lock <- struct{}{}

	tmon.errMsg = make(chan error)

	for _, f := range files {
		if err := tmon.populateTemplate(f); err != nil {
			return nil
		}

		stat, err := os.Stat(f)
		if err != nil {
			return nil
		}
		tmon.lastUpdated[f] = stat
	}

	return tmon
}

func (tmon *TMonitor) Watch() {
	ticker := time.NewTicker(1 * time.Second)
	go func() {
		for {
			select {
			case msg := <-tmon.errMsg:
				fmt.Printf("Error reloading: %v\n", msg)
			case <-ticker.C:
				tmon.checkChanged()
			}
		}
	}()
}

func (tmon *TMonitor) recompile(file string) error {
	t, err := ht.ParseFiles(file)
	if err != nil {
		return err
	}
	<-tmon.lock
	tmon.htmlTemplates[file] = t
	tmon.lock <- struct{}{}
	return nil
}

func (tmon *TMonitor) Get(name string) *ht.Template {
	var tpl *ht.Template
	<-tmon.lock
	tpl = tmon.htmlTemplates[name]
	tmon.lock <- struct{}{}
	return tpl
}

func (tmon *TMonitor) populateTemplate(filename string) error {
	t, err := ht.ParseFiles(filename)
	if err != nil {
		return err
	}
	tmon.htmlTemplates[filename] = t
	return nil
}

func (tmon *TMonitor) checkChanged() error {
	for _, f := range tmon.files {
		stat, err := os.Stat(f)
		if err != nil {
			return err
		}

		if changed(tmon.lastUpdated[f], stat) {
			fmt.Printf("[*] Recompiling %s\n", f)
			tmon.recompile(f)
			tmon.lastUpdated[f] = stat
		}
	}
	return nil
}

func changed(first, second fs.FileInfo) bool {
	if first.Size() != second.Size() || first.ModTime() != second.ModTime() {
		return true
	}
	return false
}
