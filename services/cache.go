package services

type DiskCache struct {
	BasePath string
}

func NewDiskCache(basePath string) *DiskCache {
	return &DiskCache{
		BasePath: basePath,
	}
}
