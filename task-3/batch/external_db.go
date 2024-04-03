package batch

type ExternalDB interface {
	SaveBatch(data []string, size int)
}
