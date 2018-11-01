package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	utils "github.com/maorfr/helm-plugin-utils/pkg"
	helm_restore "github.com/maorfr/helm-restore/pkg"

	"github.com/mholt/archiver"
	"github.com/spf13/cobra"
)

var (
	releaseName     string
	tillerNamespace string
	label           string
	file            string
	restore         bool
)

func main() {
	cmd := &cobra.Command{
		Use:   "backup [flags] NAMESPACE",
		Short: "backup/restore releases in a namespace to/from a file",
		RunE:  run,
	}

	f := cmd.Flags()
	f.StringVar(&tillerNamespace, "tiller-namespace", "kube-system", "namespace of Tiller")
	f.StringVarP(&label, "label", "l", "OWNER=TILLER", "label to select tiller resources by")
	f.StringVar(&file, "file", "", "file name to use (.tgz file). If not provided - will use <namespace>.tgz")
	f.BoolVarP(&restore, "restore", "r", false, "restore instead of backup")

	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func run(cmd *cobra.Command, args []string) error {
	namespace := args[0]
	if restore {
		if err := Restore(namespace); err != nil {
			return err
		}
	} else {
		if err := Backup(namespace); err != nil {
			return err
		}
	}

	return nil
}

// Backup performs a backup of all releases from provided namespace
func Backup(namespace string) error {
	storage := utils.GetTillerStorage(tillerNamespace)
	inReleases, err := utils.ListReleaseNamesInNamespace(namespace)
	if err != nil {
		return err
	}
	backupCmd := []string{
		"kubectl",
		"--namespace", tillerNamespace,
		"get", storage,
		"-l", label,
		"-l", "NAME in (" + inReleases + ")",
		"-o", "yaml",
	}
	output := utils.Execute(backupCmd)

	manifestsFileName := "manifests.yaml"
	releasesFileName := "releases"
	tarGzName := getTarGzFileName(namespace)
	os.Remove(manifestsFileName)
	os.Remove(releasesFileName)
	os.Remove(tarGzName)
	if err := ioutil.WriteFile(manifestsFileName, output, 0644); err != nil {
		return err
	}
	if err := ioutil.WriteFile(releasesFileName, []byte(inReleases), 0644); err != nil {
		return err
	}
	err = archiver.TarGz.Make(tarGzName, []string{manifestsFileName, releasesFileName})
	if err != nil {
		return err
	}
	os.Remove(manifestsFileName)
	os.Remove(releasesFileName)
	fmt.Printf("backup of namespace \"%s\" to file %s complete (found releases: %s)\n", namespace, tarGzName, strings.Replace(inReleases, ",", ", ", -1))
	return nil
}

// Restore performs a restore of all releases to provided namespace
func Restore(namespace string) error {
	untarDir := "restore"
	manifestsFileName := untarDir + "/manifests.yaml"
	releasesFileName := untarDir + "/releases"
	os.RemoveAll(untarDir)
	tarGzName := getTarGzFileName(namespace)
	if err := archiver.TarGz.Open(tarGzName, untarDir); err != nil {
		return err
	}
	releasesToRestore, err := ioutil.ReadFile(releasesFileName)
	if err != nil {
		return err
	}
	restoreCmd := []string{
		"kubectl",
		"--namespace", tillerNamespace,
		"apply", "-f", manifestsFileName,
	}
	output := utils.Execute(restoreCmd)
	fmt.Print((string)(output))

	label += ",STATUS=DEPLOYED"
	for _, r := range strings.Split((string)(releasesToRestore), ",") {
		helm_restore.Restore(r, tillerNamespace, label)
	}

	os.RemoveAll(untarDir)
	fmt.Printf("restore file %s to namespace \"%s\" complete (found releases: %s)\n", tarGzName, namespace, strings.Replace((string)(releasesToRestore), ",", ", ", -1))
	return nil
}

func getTarGzFileName(namespace string) string {
	tarGzName := file
	if tarGzName == "" {
		tarGzName = namespace
	}
	if !strings.HasSuffix(tarGzName, ".tgz") {
		tarGzName = tarGzName + ".tgz"
	}

	return tarGzName
}
