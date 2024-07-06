package gemgrscreen

import (
	"strings"

	"github.com/charmbracelet/bubbles/list"
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

	items := []list.Item{}
	for _, release := range remote.Releases {

		remotePath := ""
		for _, asset := range release.Assets {
			if strings.Contains(asset.Name, "GE-Proton") && strings.Contains(asset.Name, ".tar.gz") {
				remotePath = asset.BrowserDownloadUrl
			}
		}

		localPath := ""
		for _, local := range local.Releases {
			if local.Name == release.TagName {
				localPath = local.Path
			}
		}

		status := None
		if localPath != "" {
			status = Downloaded
		}

		geitem := GEItem{
			name:       release.TagName,
			localPath:  localPath,
			remotePath: remotePath,
			status:     status,
		}
		items = append(items, item(geitem))
	}

	s := spinner.New()
	s.Spinner = spinner.Points

	list := list.New(items, itemDelegate{spinner: &s}, 40, 20)
	list.Title = "Glorious Eggroll Releases"

	m := Model{
		list:    list,
		spinner: s,
	}

	return m
}
