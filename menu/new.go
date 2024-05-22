package menu

import (
	"strings"

	"github.com/charmbracelet/bubbles/spinner"
	tea "github.com/charmbracelet/bubbletea"
	"github.com/fish1/sctmgr/gemgr"
)

func (m Model) Init() tea.Cmd {
	return m.spinner.Tick
}

func New() Model {
	remoteChan := make(chan gemgr.RemoteReleaseResponse)
	localChan := make(chan gemgr.LocalReleaseResponse)

	go gemgr.RemoteReleases(remoteChan)
	go gemgr.LocalReleases(localChan)

	remote := <-remoteChan

	if remote.Err != nil {
		panic(remote.Err)
	}

	local := <-localChan

	if local.Err != nil {
		panic(local.Err)
	}

	choices := []Choice{}
	selections := make(map[int]Selection)

	for i, release := range remote.Releases {

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

		for _, local := range local.Releases {
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
	m.spinner.Spinner = spinner.Points

	return m
}
