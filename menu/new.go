package menu

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fish1/ge-downloader/gemanager"
)

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func New() Model {
	remote, err := gemanager.RemoteReleases()
	if err != nil {
		panic(err)
	}

	local, err := gemanager.LocalReleases()
	if err != nil {
		panic(err)
	}

	choices := []Choice{}
	selections := make(map[int]Selection)

	for i, release := range remote {

		downloadUrl := ""
		for _, asset := range release.Assets {
			if strings.Contains(asset.Name, "GE-Proton") && strings.Contains(asset.Name, ".tar.gz") {
				downloadUrl = asset.BrowserDownloadUrl
			}
		}

		choice := Choice{
			name:        release.TagName,
			downloadUrl: downloadUrl,
		}

		choices = append(choices, choice)

		for _, local := range local {
			if local.Name == release.TagName {
				selections[i] = Selection{
					status:   Idle,
					localUri: local.Path,
				}
				break
			}
		}
	}

	m := Model{
		header:   "Choose a GE release to install",
		choices:  choices,
		selected: selections,
	}

	m.spinner = spinner.New()
	m.spinner.Spinner = spinner.Globe

	return m
}
