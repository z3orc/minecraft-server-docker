package management

type WhitelistEntry struct {
	UUID string `json:"uuid"`
	Name string `json:"name"`
	//IgnoresPlayerLimit bool   `json:"ignoresPlayerLimit"`
}

type Whitelist []WhitelistEntry
