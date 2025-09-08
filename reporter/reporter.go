package reporter

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/yansol0/blacklight/tester"
	"github.com/yansol0/blacklight/utils"
)

func WriteReports(results tester.Results, outdir string) error {
	if err := os.MkdirAll(outdir, 0755); err != nil {
		return err
	}

	unauthPath := filepath.Join(outdir, "unauth_report.txt")
	authPath := filepath.Join(outdir, "auth_report.txt")
	bypassPath := filepath.Join(outdir, "bypass_report.txt")

	if err := writeList(unauthPath, results.Unauth); err != nil {
		return err
	}
	if err := writeList(authPath, results.Auth); err != nil {
		return err
	}
	if err := writeList(bypassPath, results.Bypass); err != nil {
		return err
	}

	utils.LogSuccess(fmt.Sprintf("Reports written to %s", outdir))
	utils.LogInfo(fmt.Sprintf("%d endpoints to be tested for IDOR", len(results.IDORCandidates)))
	return nil
}

func writeList(path string, entries [][2]string) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	for _, entry := range entries {
		line := fmt.Sprintf("%s [%s]\n", entry[0], entry[1])
		if _, err := f.WriteString(line); err != nil {
			return err
		}
	}
	return nil
}
