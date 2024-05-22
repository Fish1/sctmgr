package gemgr

type Asset struct {
	Url                string `json:"url"`
	BrowserDownloadUrl string `json:"browser_download_url"`
	Id                 int    `json:"id"`
	Name               string `json:"name"`
	Size               int    `json:"size"`
	CreatedAt          string `json:"created_at"`
	UpdatedAt          string `json:"updated_at"`
}

type RemoteRelease struct {
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

type LocalRelease struct {
	Name string
	Path string
}
