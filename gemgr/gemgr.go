package gemgr

import (
	"archive/tar"
	"compress/gzip"
	"crypto/sha1"
	"encoding/base64"
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

type RemoteReleaseResponse struct {
	Releases []RemoteRelease
	Err      error
}

func RemoteReleases(releasesChan chan RemoteReleaseResponse) {
	res, err := http.Get(GE_GITHUB_URL)
	if err != nil {
		releasesChan <- RemoteReleaseResponse{[]RemoteRelease{}, err}
		return
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		releasesChan <- RemoteReleaseResponse{[]RemoteRelease{}, err}
		return
	}

	var releases []RemoteRelease
	err = json.Unmarshal(body, &releases)
	if err != nil {
		releasesChan <- RemoteReleaseResponse{[]RemoteRelease{}, err}
		return
	}

	releasesChan <- RemoteReleaseResponse{releases, nil}
}

type LocalReleaseResponse struct {
	Releases []LocalRelease
	Err      error
}

func LocalReleases(releaseChan chan LocalReleaseResponse) {
	results := []LocalRelease{}

	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		releaseChan <- LocalReleaseResponse{[]LocalRelease{}, err}
		return
	}

	directories, err := os.ReadDir(homeDirectory + COMPATIBILITY_TOOLS_DIRECTORY)
	if err != nil {
		releaseChan <- LocalReleaseResponse{[]LocalRelease{}, err}
		return
	}

	for _, directory := range directories {
		if directory.IsDir() {
			files, err := os.ReadDir(homeDirectory + COMPATIBILITY_TOOLS_DIRECTORY + directory.Name())
			if err != nil {
				releaseChan <- LocalReleaseResponse{[]LocalRelease{}, err}
				return
			}
			for _, file := range files {
				if file.IsDir() == false && file.Name() == "version" {
					folder := homeDirectory + COMPATIBILITY_TOOLS_DIRECTORY + directory.Name()
					path := folder + "/" + file.Name()
					buffer, err := os.ReadFile(path)
					if err != nil {
						releaseChan <- LocalReleaseResponse{[]LocalRelease{}, err}
						return
					}
					fileData := strings.Split(string(buffer), " ")
					if len(fileData) != 2 {
						releaseChan <- LocalReleaseResponse{[]LocalRelease{}, err}
						return
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

	releaseChan <- LocalReleaseResponse{results, nil}
}

func Delete(localUri string) error {
	err := os.RemoveAll(localUri)
	if err != nil {
		return err
	}

	return nil
}

func Install(remote string) (string, error) {

	var bytes []byte
	hasher := sha1.New()
	hasher.Write(bytes)
	filename := filepath.Join("/tmp/", base64.URLEncoding.EncodeToString(hasher.Sum([]byte(remote)))+".tar.gz")

	err := downloadTempFile(remote, filename)
	if err != nil {
		return "", err
	}

	path, err := extractTempFile(filename)
	if err != nil {
		return "", err
	}

	err = cleanupTempFile(filename)
	if err != nil {
		return "", err
	}

	return path, nil
}

func downloadTempFile(remote string, destination string) error {
	res, err := http.Get(remote)
	if err != nil {
		return err
	}

	body, err := io.ReadAll(res.Body)
	if err != nil {
		return err
	}

	file, err := os.Create(destination)
	if err != nil {
		return err
	}

	_, err = file.Write(body)
	if err != nil {
		return err
	}

	return nil
}

func extractTempFile(source string) (string, error) {
	homeDirectory, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	file, err := os.Open(source)
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

func cleanupTempFile(source string) error {
	err := os.Remove(source)
	if err != nil {
		return err
	}
	return nil
}
