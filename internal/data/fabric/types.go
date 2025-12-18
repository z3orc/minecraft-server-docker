package fabric

type installerResponse []installer

type installer struct {
	Url     string
	Maven   string
	Version string
	Stable  bool
}

type loaderResponse []loaderResponseItem

type loaderResponseItem struct {
	Loader loader
}

type loader struct {
	Separator string
	Build     int
	Maven     string
	Version   string
	Stable    bool
}
