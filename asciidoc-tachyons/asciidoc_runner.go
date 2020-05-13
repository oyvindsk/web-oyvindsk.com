package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"
)

func main() {

	err := generateAll("test-1")
	if err != nil {
		fmt.Println(err)
		return
	}

}

func generateAll(dirpath string) error {

	err := os.Chdir(dirpath)
	if err != nil {
		return fmt.Errorf("generateAll: %s", err)
	}

	// create the html, the main content for the website
	// it is not meant to be served directly, it be augmentet by postprocessing and go templates later
	err = runDoctorHTML("content.adoc", "content.html")
	if err != nil {
		return fmt.Errorf("generateAll: %s", err)

	}

	// Also try to generate a pdf
	// this is more involved since it:
	// - needs another prog
	// - must go through docbook first (at least with this setup)
	// - needs the metadata beacuse the pdf must be complete, no header or footer is added (no template like html)

	blogmetadata, err := loadMetadata("metadata.toml")
	if err != nil {
		return fmt.Errorf("generateAll: %s", err)
	}

	err = runDoctorDocbook("content.adoc", "full.xml", blogmetadata)
	if err != nil {
		return fmt.Errorf("generateAll: %s", err)

	}

	err = runOriginalPDF("full.xml")
	if err != nil {
		return fmt.Errorf("generateAll: %s", err)
	}

	return nil

}

func runDoctorHTML(inputpath, outputpath string) error {

	// -s : No header, footer, css etc
	// -a compat-mode: behave more like asciidoc
	cmd := exec.Command("asciidoctor", "-s", "-a", "compat-mode", inputpath, "-o", outputpath)

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runDoctorHTML: %s", err)
	}

	if len(stdoutStderr) != 0 {
		return fmt.Errorf("runDoctorHTML: cmd output: %s", stdoutStderr)
	}

	return nil
}

func runDoctorDocbook(inputpath, outputpath string, metadata BlogMetadata) error {

	cmd := exec.Command(
		"asciidoctor",
		"-a", "compat-mode", // behave more like asciidoc
		"-b", "docbook", // docbook xml,  not html

		// Set asciidoc variables so theyll be in the docbook, and eventually, the pdf
		"-a", "doctitle="+metadata.Postmeta.Title,
		"-a", "subtitleforpdf="+metadata.Postmeta.Subtitle,
		"-a", "author="+metadata.Postmeta.Author,
		"-a", "revdate="+metadata.Postmeta.Date.String(),

		inputpath,
		"-o", outputpath,
	)

	fmt.Println(cmd.String())

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runDoctorDocbook: %s", err)
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
		return fmt.Errorf("runOriginalPDF: %s", err)
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
		return fmt.Errorf("runOriginalHTML: %s", err)
	}

	if len(stdoutStderr) != 0 {
		return fmt.Errorf("runOriginalHTML: cmd output: %s", stdoutStderr)
	}

	return nil
}
