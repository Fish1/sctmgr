package gemanager

import (
	"archive/tar"
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"os"
	"strings"
)

const COMPATIBILITY_TOOLS_DIRECTORY = "/.steam/steam/compatibilitytools.d/"
const GE_GITHUB_URL = "https://api.github.com/repos/GloriousEggroll/proton-ge-custom/releases"

type Asset struct {
	Url                string `json:"url"`
	BrowserDownloadUrl string `json:"browser_download_url"`
	Id                 int    `json:"id"`
	Name               string `json:"name"`
	Size               int    `json:"size"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

type Release struct {
	Url             string  `json:"url"`
	HtmlUrl         string  `json:"html_url"`
	AssetsUrl       string  `json:"assets_url"`
	UploadUrl       string  `json:"upload_url"`
	TarballUrl      string  `json:"tarball_url"`
	ZipballUrl      string  `json:"zipball_url"`
	Id              int     `json:"id"`
	TagName         string  `json:"tag_name"`
	TargetCommitish string  `json:"target_commitish"`
	CreatedAt       string  `json:"created_at"`
	UpdatedAt       string  `json:"updated_at"`
	Assets          []Asset `json:"assets"`
}

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

type LocalRelease struct {
	Name string
	Path string
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

func Install(url string) (string, error) {
	err := downloadTempFile(url)
	if err != nil {
		return "", err
	}
	filename, err := extractTempFile()
	if err != nil {
		return "", err
	}
	return filename, nil
}

func downloadTempFile(url string) error {
	res, err := http.Get(url)
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

	file, err := os.Open("temp.tar.gz")
	if err != nil {
		return "", err
	}

	reader := tar.NewReader(file)
	var header *tar.Header
	for header, err = reader.Next(); err == nil; header, err = reader.Next() {
		switch header.Typeflag {
		case tar.TypeDir:
			if err := os.Mkdir(header.Name, 0755); err != nil {
				return "", err
			}
		case tar.TypeReg:
			outfile, err := os.Create(header.Name)
			if err != nil {
				return "", err
			}
			if _, err := io.Copy(outfile, reader); err != nil {
				outfile.Close()
				return "", err
			}
			if err := outfile.Close(); err != nil {
				return "", err
			}
		default:
			return "", errors.New("unknown tar type")
		}
	}

	return "", nil
}
