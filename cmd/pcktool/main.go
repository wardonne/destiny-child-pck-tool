package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gopi-frame/collection/list"
	"github.com/gopi-frame/console"
	"github.com/spf13/cobra"
	"github.com/wardonne/distiny-child-pck-tool"
	"os"
	"path/filepath"
	"runtime/debug"
	"strings"
	"sync"
	"time"
)

var ErrSkip = errors.New("skip")

func main() {
	wd, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	kernel := console.NewKernel()
	kernel.Command.Flags().StringP("source", "s", "", "source directory or file")
	kernel.Command.Flags().StringP("target", "t", filepath.Join(wd, "output"), "target directory to store the unpacked files")
	kernel.Command.Flags().BoolP("dry-run", "d", false, "dry run mode")
	kernel.Command.Flags().BoolP("live2d", "l", false, "rename and format the unpacked files to live2d model")
	kernel.Command.Flags().BoolP("force", "f", false, "force to overwrite the existing files, otherwise skip the existing files")
	kernel.Command.Flags().BoolP("verbose", "v", false, "verbose mode")
	kernel.Command.Flags().IntP("concurrency", "c", 1, "concurrency")
	kernel.Command.Run = func(cmd *cobra.Command, args []string) {
		var data struct {
			Source      string `flag:"source"`
			Target      string `flag:"target"`
			DryRun      bool   `flag:"dry-run"`
			Live2D      bool   `flag:"live2d"`
			Force       bool   `flag:"force"`
			Verbose     bool   `flag:"verbose"`
			Concurrency int    `flag:"concurrency"`
		}
		input := console.NewInput(cmd.Context(), cmd.Flags())
		output := console.GetOutput(cmd.Context())
		if err := input.Unmarshal(&data); err != nil {
			output.Errorf("Can not unmarshal input: %v", err)
			return
		}
		if data.Verbose {
			output = output.WithMode(output.GetMode().Append(console.OutputModeDebug))
		}

		if data.Source == "" {
			output.Error("Empty directory or file")
			return
		}
		stat, err := os.Stat(data.Source)
		if err != nil {
			output.Errorf("Source directory or file not found: %v", err)
			return
		}
		output.Debugf("Input.source: %v", data.Source)
		output.Debugf("Input.target: %v", data.Target)
		output.Debugf("Input.dry-run: %v", data.DryRun)
		output.Debugf("Input.live2d: %v", data.Live2D)
		output.Debugf("Input.force: %v", data.Force)
		output.Debugf("Input.verbose: %v", data.Verbose)
		output.Debugf("Input.concurrency: %v", data.Concurrency)
		var files []string
		if stat.IsDir() {
			files, err = filepath.Glob(filepath.Join(data.Source, "*.pck"))
			if err != nil {
				output.Errorf("Failed to glob source directory: %v", err)
				return
			}
		} else {
			if filepath.Ext(data.Source) != ".pck" {
				output.Errorf("Source file must be a .pck file")
				return
			}
			files = []string{data.Source}
		}
		if len(files) == 0 {
			output.Info("No .pck files found in source directory")
			return
		}
		output.Debugf("found %d files to unpack", len(files))
		start := time.Now()
		wg := new(sync.WaitGroup)
		if data.Concurrency <= 0 {
			data.Concurrency = 1
		}
		ch := make(chan string, data.Concurrency)
		success := list.NewList[string]()
		failed := list.NewList[string]()
		skipped := list.NewList[string]()
		for _, file := range files {
			wg.Add(1)
			ch <- file
			go func(file string, wg *sync.WaitGroup) {
				defer func() {
					if err := recover(); err != nil {
						if e, ok := err.(error); ok {
							if errors.Is(e, ErrSkip) {
								skipped.Lock()
								skipped.Push(file)
								skipped.Unlock()
								output.Warnf("Skipping %s", file)
							} else {
								if data.Verbose {
									debug.PrintStack()
								}
								output.Errorf("%v", err)
							}
						} else {
							if data.Verbose {
								debug.PrintStack()
							}
							output.Errorf("%v", err)
						}
					} else {
						success.Lock()
						success.Push(file)
						success.Unlock()
					}
					wg.Done()
					<-ch
				}()
				output.Debugf("Unpacking %s", file)
				pack, err := pcktool.Unpack(file, output)
				if err != nil {
					panic(fmt.Sprintf("Failed to unpack %s: %v", file, err))
					return
				}
				if data.Live2D {
					model, err := pcktool.GenerateLive2D(pack)
					if err != nil {
						panic(fmt.Sprintf("Failed to generate live2d model for %s: %v", file, err))
						return
					}
					jsonBytes, err := json.MarshalIndent(model, "", "  ")
					if err != nil {
						panic(fmt.Sprintf("Failed to marshal live2d model for %s: %v", file, err))
						return
					}
					output.Debug("Live2d model:")
					for _, line := range strings.Split(string(jsonBytes), "\n") {
						output.Debugf("%s", line)
					}
				}
				for index, entry := range pack.Entries {
					output.Debugf("File %d/%d: %s", index+1, pack.FileCount, entry.Filename)
				}
				if !data.DryRun {
					targetDir := filepath.Join(data.Target, filepath.Base(file)[:len(filepath.Base(file))-len(filepath.Ext(file))])
					if _, err := os.Stat(targetDir); os.IsNotExist(err) {
						if err := os.MkdirAll(targetDir, 0755); err != nil {
							panic(fmt.Sprintf("Failed to create target directory %s: %v", targetDir, err))
							return
						}
					} else if err != nil {
						panic(fmt.Sprintf("Failed to check target directory %s: %v", targetDir, err))
						return
					} else {
						if !data.Force {
							panic(ErrSkip)
							return
						}
					}
					for index, entry := range pack.Entries {
						output.Debugf("Saving %d/%d: %s", index+1, pack.FileCount, entry.Filename)
						targetPath := filepath.Join(targetDir, entry.Filename)
						if _, err := os.Stat(filepath.Dir(targetPath)); os.IsNotExist(err) {
							if err := os.MkdirAll(filepath.Dir(targetPath), 0755); err != nil {
								panic(fmt.Sprintf("Failed to create target directory %s: %v", filepath.Dir(targetPath), err))
								return
							}
						} else if err != nil {
							panic(fmt.Sprintf("Failed to check target directory %s: %v", filepath.Dir(targetPath), err))
							return
						}
						if err := os.WriteFile(targetPath, entry.Content, 0644); err != nil {
							panic(fmt.Sprintf("Failed to write file %s: %v", entry.Filename, err))
							return
						}
						output.Debugf("Saved %s", entry.Filename)
					}
				}
				output.Debugf("Unpaced %s", file)
			}(file, wg)
		}
		wg.Wait()
		if failed.Count() > 0 {
			output.Warnf("Unpacked %d files in %s [ %02d succeed | %02d failed | %02d skipped ]", failed.Count(), time.Since(start), success.Count(), failed.Count(), skipped.Count())
		}
		output.Successf("Unpacked %d files in %s [ %02d succeed | %02d skipped ]", len(files), time.Since(start), success.Count(), skipped.Count())
	}
	_ = kernel.Run()
}
