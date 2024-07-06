package gemgr

import (
	"archive/tar"
	"bytes"
	"compress/gzip"
	"context"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"io"
	"log/slog"
	"net/http"
	"os"
	"path/filepath"
	"time"
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
			results = append(results, LocalRelease{
				Name: directory.Name(),
				Path: homeDirectory + COMPATIBILITY_TOOLS_DIRECTORY + directory.Name(),
			})
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

func Install(remote string, ctx context.Context) (string, error) {
	var bytes []byte
	hasher := sha1.New()
	hasher.Write(bytes)
	filename := filepath.Join("/tmp/", base64.URLEncoding.EncodeToString(hasher.Sum([]byte(remote)))+".tar.gz")

	err := downloadTempFile(remote, filename, ctx)
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

func downloadTempFile(remote string, destination string, ctx context.Context) error {
	bodyReader := bytes.NewBuffer([]byte(""))
	req, err := http.NewRequestWithContext(ctx, "GET", remote, bodyReader)
	if err != nil {
		return err
	}

	client := http.Client{
		Timeout: 30 * time.Second,
	}

	res, err := client.Do(req)
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
		// info := header.FileInfo()

		slog.Debug(header.Linkname)

		switch header.Typeflag {
		case tar.TypeDir:
			err := os.MkdirAll(path, os.FileMode(header.Mode))
			if err != nil {
				return "", err
			}
			if resultpath == "" {
				resultpath = path
			}

		case tar.TypeReg:
			err := os.MkdirAll(filepath.Dir(path), os.ModePerm)
			if err != nil {
				return "", err
			}

			file, err := os.Create(path)
			if err != nil {
				return "", err
			}
			defer file.Close()

			_, err = io.Copy(file, reader)
			if err != nil {
				return "", err
			}

			err = os.Chmod(path, os.FileMode(header.Mode))
			if err != nil {
				return "", err
			}

			/*
				outfile, err := os.OpenFile(path, os.O_CREATE|os.O_TRUNC|os.O_WRONLY, info.Mode())
				if err != nil {
					return "", err
				}
				defer outfile.Close()
				_, err = io.Copy(outfile, reader)
				if err != nil {
					outfile.Close()
					return "", err
				}
			*/

		case tar.TypeSymlink:
			os.Symlink(header.Linkname, path)
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
