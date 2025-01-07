package users

type User struct {
	UID         uint32 `json:"uid"`
	GID         uint32 `json:"gid"`
	UIDSigned   int32  `json:"uid_signed"`
	GIDSigned   int32  `json:"gid_signed"`
	Username    string `json:"username"`
	Description string `json:"description"`
	Directory   string `json:"directory"`
	Shell       string `json:"shell"`
	UUID        string `json:"uuid"`
	Type        string `json:"type"`
}

var wellKnownSids = []string{
	"S-1-5-1",
	"S-1-5-2",
	"S-1-5-3",
	"S-1-5-4",
	"S-1-5-6",
	"S-1-5-7",
	"S-1-5-8",
	"S-1-5-9",
	"S-1-5-10",
	"S-1-5-11",
	"S-1-5-12",
	"S-1-5-13",
	"S-1-5-18",
	"S-1-5-19",
	"S-1-5-20",
	"S-1-5-21",
	"S-1-5-32",
}

func getRoamingProfileSids
