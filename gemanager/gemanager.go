package gemanager

import (
	"archive/tar"
	"compress/gzip"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

const COMPATIBILITY_TOOLS_DIRECTORY = "/.steam/steam/compatibilitytools.d/"
const GE_GITHUB_URL = "https://api.github.com/repos/GloriousEggroll/proton-ge-custom/releases"

func RemoteReleases() ([]Release, error) {
	res, err := http.Get(GE_GITHUB_URL)
	if err != nil {
		return []Release{}, err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return []Release{}, err
	}

	var releases []Release
	err = json.Unmarshal(body, &releases)
	if err != nil {
		return []Release{}, err
	}

	return releases, nil
}

func LocalReleases() ([]LocalRelease, error) {
	results := []LocalRelease{}

	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return []LocalRelease{}, err
	}

	directories, err := os.ReadDir(homeDirectory + COMPATIBILITY_TOOLS_DIRECTORY)
	if err != nil {
		return []LocalRelease{}, err
	}

	for _, directory := range directories {
		if directory.IsDir() {
			files, err := os.ReadDir(homeDirectory + COMPATIBILITY_TOOLS_DIRECTORY + directory.Name())
			if err != nil {
				return []LocalRelease{}, err
			}
			for _, file := range files {
				if file.IsDir() == false && file.Name() == "version" {
					folder := homeDirectory + COMPATIBILITY_TOOLS_DIRECTORY + directory.Name()
					path := folder + "/" + file.Name()
					buffer, err := os.ReadFile(path)
					if err != nil {
						return []LocalRelease{}, err
					}
					fileData := strings.Split(string(buffer), " ")
					if len(fileData) != 2 {
						return []LocalRelease{}, errors.New("unable to read version file")
					}
					tag := strings.Trim(fileData[1], "\n")
					results = append(results, LocalRelease{
						Name: tag,
						Path: folder,
					})
				}
			}
		}
	}

	return results, nil
}

func Delete(localUri string) error {
	err := os.RemoveAll(localUri)
	if err != nil {
		return err
	}

	return nil
}

func Install(remote string) (string, error) {
	err := downloadTempFile(remote)
	if err != nil {
		return "", err
	}

	path, err := extractTempFile()
	if err != nil {
		return "", err
	}

	err = cleanupTempFile()
	if err != nil {
		return "", err
	}

	return path, nil
}

func downloadTempFile(remote string) error {
	res, err := http.Get(remote)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	file, err := os.Create("temp.tar.gz")
	if err != nil {
		return err
	}

	_, err = file.Write(body)
	if err != nil {
		return err
	}

	return nil
}

func extractTempFile() (string, error) {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	file, err := os.Open("temp.tar.gz")
	if err != nil {
		return "", err
	}

	stream, err := gzip.NewReader(file)
	if err != nil {
		return "", err
	}

	var resultpath string

	reader := tar.NewReader(stream)
	var header *tar.Header
	for header, err = reader.Next(); err == nil; header, err = reader.Next() {
		path := filepath.Join(homeDirectory, COMPATIBILITY_TOOLS_DIRECTORY, header.Name)
		switch header.Typeflag {
		case tar.TypeDir:
			err := os.Mkdir(path, 0755)
			if err != nil {
				return "", err
			}
			if resultpath == "" {
				resultpath = path
			}
		case tar.TypeReg, tar.TypeSymlink:
			outfile, err := os.Create(path)
			if err != nil {
				return "", err
			}
			_, err = io.Copy(outfile, reader)
			if err != nil {
				outfile.Close()
				return "", err
			}
			err = outfile.Close()
			if err != nil {
				return "", err
			}
		default:
			return "", errors.New("unknown tar type: " + header.Name + " " + string(header.Typeflag))
		}
	}

	return resultpath, nil
}

func cleanupTempFile() error {
	err := os.Remove("temp.tar.gz")
	if err != nil {
		return err
	}
	return nil
}
