package main

import (
	"fmt"
	"os"
	"os/exec"
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

	// Also try to generate a pdf. This uses another asciidoctor program, see the README
	err = runDoctorPDF(asciidocFilename, pdfFilename)
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

// TODO :
//   - Test and compare to the old
//   - Use metadata, like the old one? Passed in as "variables"
//     "github.com/oyvindsk/web-oyvindsk.com/internal/metadata"
func runDoctorPDF(inputpath, outputpath string) error {

	cmd := exec.Command("asciidoctor-pdf", "-b", "pdf", inputpath, "-o", outputpath)

	stdoutStderr, err := cmd.CombinedOutput()
	if err != nil {
		return fmt.Errorf("runDoctorPDF: %w, cmd output: %s", err, stdoutStderr)
	}

	if len(stdoutStderr) != 0 {
		return fmt.Errorf("runDoctorPDF: cmd output: %s", stdoutStderr)
	}

	return nil
}
