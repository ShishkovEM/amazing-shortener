package workerpool

type DeletionTask struct {
	userID      uint32
	urlToDelete string
}

func NewDeletionTask(uid uint32, URL string) *DeletionTask {
	return &DeletionTask{userID: uid, urlToDelete: URL}
}
