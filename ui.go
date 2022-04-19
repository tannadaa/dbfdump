package main

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/data/binding"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/layout"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
)

func showUI(out, dir, skip string) {
	a := app.New()
	w := a.NewWindow("DBF Dump")

	boundInput := binding.BindString(&dir)
	boundSkip := binding.BindString(&skip)
	boundOutput := binding.BindString(&out)

	inputDialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
		if uri == nil {
			return
		}
		_ = boundInput.Set(uri.Path())
	}, w)

	skipEntry := widget.NewEntryWithData(boundSkip)
	skipEntry.Validator = nil

	outputDialog := dialog.NewFolderOpen(func(uri fyne.ListableURI, err error) {
		if uri == nil {
			return
		}
		_ = boundOutput.Set(uri.Path())
	}, w)

	exportButton := widget.NewButton("Export", func() {
		inDir, _ := boundInput.Get()
		outDir, _ := boundOutput.Get()
		skips, _ := boundSkip.Get()
		err := processDBF(inDir, outDir, skipFiles(skips))
		if err != nil {
			// this is a hack to get the error to be multiline
			dialog.ShowError(fmt.Errorf(strings.ReplaceAll(err.Error(), ": ", "\n")), w)
		} else {
			open(outDir)
		}
	})
	exportButton.Importance = widget.HighImportance

	formContainer := container.New(layout.NewFormLayout(),
		widget.NewLabel("Input Directory"),
		container.New(layout.NewFormLayout(),
			widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
				inputDialog.Show()
			}), widget.NewLabelWithData(boundInput)),
		widget.NewLabel("Skip"),
		skipEntry,
		widget.NewLabel("Output Directory"),
		container.New(layout.NewFormLayout(),
			widget.NewButtonWithIcon("", theme.FolderOpenIcon(), func() {
				outputDialog.Show()
			}), widget.NewLabelWithData(boundOutput)),
		exportButton,
		widget.NewLabel(""),
	)

	w.SetContent(formContainer)
	w.Resize(fyne.NewSize(600, 600))
	w.ShowAndRun()
}
