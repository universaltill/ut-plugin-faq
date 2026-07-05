// pkgtool is the self-contained packaging helper for the FAQ plugin.
// It validates the manifest against both the marketplace and POS host
// contracts, stages a per-target manifest, and emits release metadata,
// so packaging never depends on external CLIs being installed.
//
// Packaging is single-target: one artifact per os/arch pair. Architecture-
// independent ("universal") artifacts are defined in the release contract but
// not produced by this plugin, which ships a compiled binary.
//
// Usage:
//
//	go run ./tools/pkgtool validate [-manifest path] [-root path]
//	go run ./tools/pkgtool version [-manifest path]
//	go run ./tools/pkgtool id [-manifest path]
//	go run ./tools/pkgtool stage-manifest -out path -os OS -arch ARCH [-manifest path]
//	go run ./tools/pkgtool release-json -out path -os OS -arch ARCH -artifact NAME [-manifest path]
package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fatal("usage: pkgtool <validate|version|id|stage-manifest|release-json> [flags]")
	}
	cmd, args := os.Args[1], os.Args[2:]

	fs := flag.NewFlagSet(cmd, flag.ExitOnError)
	manifestPath := fs.String("manifest", "src/manifest/manifest.json", "path to manifest.json")
	root := fs.String("root", ".", "repo root for resource/version checks ('' to skip)")
	out := fs.String("out", "", "output file path")
	targetOS := fs.String("os", "", "target OS (e.g. linux, darwin)")
	targetArch := fs.String("arch", "", "target architecture (e.g. amd64, arm64)")
	artifact := fs.String("artifact", "", "artifact file name")
	if err := fs.Parse(args); err != nil {
		fatal(err.Error())
	}

	m, err := LoadManifest(*manifestPath)
	if err != nil {
		fatal(err.Error())
	}

	switch cmd {
	case "validate":
		if errs := ValidateManifest(m, *root); len(errs) > 0 {
			fmt.Fprintln(os.Stderr, "manifest validation failed:")
			for _, e := range errs {
				fmt.Fprintf(os.Stderr, "  - %s\n", e)
			}
			os.Exit(1)
		}
		if *root != "" {
			if err := CheckVersionAlignment(m, filepath.Join(*root, "package.json")); err != nil {
				fatal(err.Error())
			}
		}
		fmt.Printf("manifest OK: %s %s\n", m.ID, m.Version)

	case "version":
		fmt.Println(m.Version)

	case "id":
		fmt.Println(m.ID)

	case "stage-manifest":
		requireFlags(map[string]string{"out": *out, "os": *targetOS, "arch": *targetArch})
		if !m.SupportsTarget(*targetOS, *targetArch) {
			fatal(fmt.Sprintf("target %s/%s is not in manifest supported_architectures %v", *targetOS, *targetArch, m.SupportedArch))
		}
		if err := writeStagedManifest(*manifestPath, *out, *targetOS+"/"+*targetArch); err != nil {
			fatal(err.Error())
		}

	case "release-json":
		requireFlags(map[string]string{"out": *out, "os": *targetOS, "arch": *targetArch, "artifact": *artifact})
		rel := map[string]string{
			"plugin_id":     m.ID,
			"version":       m.Version,
			"build_time":    time.Now().UTC().Format(time.RFC3339),
			"target_os":     *targetOS,
			"target_arch":   *targetArch,
			"artifact_name": *artifact,
		}
		if sha := gitSHA(); sha != "" {
			rel["git_sha"] = sha
		}
		if err := writeJSON(*out, rel); err != nil {
			fatal(err.Error())
		}

	default:
		fatal("unknown command: " + cmd)
	}
}

// writeStagedManifest copies the manifest with device_arch rewritten for the
// packaging target, preserving all other fields (including unknown ones).
func writeStagedManifest(src, dst, deviceArch string) error {
	data, err := os.ReadFile(src)
	if err != nil {
		return err
	}
	var raw map[string]any
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}
	raw["device_arch"] = deviceArch
	return writeJSON(dst, raw)
}

func writeJSON(path string, v any) error {
	data, err := json.MarshalIndent(v, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile(path, append(data, '\n'), 0o644)
}

func gitSHA() string {
	out, err := exec.Command("git", "rev-parse", "HEAD").Output()
	if err != nil {
		return ""
	}
	return strings.TrimSpace(string(out))
}

func requireFlags(flags map[string]string) {
	for name, value := range flags {
		if value == "" {
			fatal("-" + name + " is required")
		}
	}
}

func fatal(msg string) {
	fmt.Fprintln(os.Stderr, "pkgtool: "+msg)
	os.Exit(1)
}
