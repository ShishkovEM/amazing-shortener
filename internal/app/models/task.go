package models

type DeletionTask struct {
	UserID      uint32
	UrlToDelete string
}

func NewDeletionTask(uid uint32, URL string) *DeletionTask {
	return &DeletionTask{UserID: uid, UrlToDelete: URL}
}
