package common

type Platform interface {
	ExpandPath(path string) string
	IsAdmin() bool
	ClearRecycleBin() error
	GetCommonPaths() map[string]string
}
