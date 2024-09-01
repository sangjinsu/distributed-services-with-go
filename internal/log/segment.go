package log

type segment struct {
	store      *store
	index      *index
	baseOffset int64
	nextOffset int64
	config     Config
}
