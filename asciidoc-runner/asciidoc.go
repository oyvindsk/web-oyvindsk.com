package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/oyvindsk/web-oyvindsk.com/internal/metadata"
)

func generate(dirpath string) error {

	fmt.Printf("Generating all static files in folder %q\n", dirpath)

	err := os.Chdir(dirpath)
	if err != nil {
		return fmt.Errorf("generateAll: %w", err)
	}

	// create the html, the main content for the website
	// it is not meant to be served directly, it be augmentet by postprocessing and go templates later
	err = runDoctorHTML(asciidocFilename, htmlFilename)
	if err != nil {
		return fmt.Errorf("generateAll: %w", err)

	}

	// Also try to generate a pdf
	// this is more involved since it:
	// - needs another prog
	// - must go through docbook first (at least with this setup)
	// - needs the metadata beacuse the pdf must be complete, no header or footer is added (no template like in the html)

	metadata, err := metadata.Fromfile(metadataFilename)
	if err != nil {
		return fmt.Errorf("generateAll: %w", err)
	}

	err = runDoctorDocbook(asciidocFilename, docBookFilename, metadata)
	if err != nil {
		return fmt.Errorf("generateAll: %w", err)

	}

	err = runOriginalPDF(docBookFilename)
	if err != nil {
		return fmt.Errorf("generateAll: %w", err)
	}

	// Remove the DocBook file if we suceeded with the pdf
	err = os.Remove(docBookFilename)
	if err != nil {
		return fmt.Errorf("generateAll: %w", err)
	}

	fmt.Println("Done!")
	return nil
}

func runDoctorHTML(inputpath, outputpath string) error {

	// -s : No header, footer, css etc
	// -a compat-mode: behave more like asciidoc
	cmd := exec.Command("asciidoctor", "-s", "-a", "compat-mode", inputpath, "-o", outputpath)

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runDoctorHTML: %w, cmd output: %s", err, stdoutStderr)
	}

	if len(stdoutStderr) != 0 {
		return fmt.Errorf("runDoctorHTML: cmd output: %s", stdoutStderr)
	}

	return nil
}

func runDoctorDocbook(inputpath, outputpath string, meta metadata.M) error {

	cmd := exec.Command(
		"asciidoctor",
		"-a", "compat-mode", // behave more like asciidoc
		"-b", "docbook", // docbook xml,  not html

		// Set asciidoc variables so theyll be in the docbook, and eventually, the pdf
		"-a", "doctitle="+meta.Title,
		"-a", "subtitleforpdf="+meta.BlogSubtitle,
		"-a", "author="+meta.Author,
		"-a", "revdate="+meta.Date.String(),

		inputpath,
		"-o", outputpath,
	)

	// fmt.Println(cmd.String())

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runDoctorDocbook: %w, cmd output: %s", err, stdoutStderr)
	}

	if len(stdoutStderr) != 0 {
		return fmt.Errorf("runDoctorDocbook: cmd output: %s", stdoutStderr)
	}

	return nil
}

func runOriginalPDF(inputpath string) error {

	if !strings.HasSuffix(inputpath, ".xml") {
		return fmt.Errorf("runOriginalPDF: inputpath should be a docbook xml file and end in .xml")
	}

	// Options for the pdf backend,
	// http://www.methods.co.nz/asciidoc/faq.html#_how_can_i_customize_pdf_files_generated_by_dblatex
	dblatexOpts := " -P doc.layout=\"coverpage mainmatter\" -P doc.publisher.show=0 -P latex.output.revhistory=0"
	cmd := exec.Command("a2x", "-d", "article", "--dblatex-opts", dblatexOpts, inputpath)

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runOriginalPDF: %w, cmd output: %s", err, stdoutStderr)
	}

	if len(stdoutStderr) != 0 {
		return fmt.Errorf("runOriginalPDF: cmd output: %s", stdoutStderr)
	}

	return nil
}

func runOriginalHTML(inputpath, outputpath string) error {

	// -s : No header, footer, css etc
	cmd := exec.Command("asciidoc", "-s", "-o", outputpath, inputpath)

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runOriginalHTML: %w, cmd output: %s", err, stdoutStderr)
	}

	if len(stdoutStderr) != 0 {
		return fmt.Errorf("runOriginalHTML: cmd output: %s", stdoutStderr)
	}

	return nil
}
