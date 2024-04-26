package src_const

const (
	StatusCodeInvalid = 601

	StatusCodeMongoCreateError = 701
	StatusCodeMongoReadError   = 702
	StatusCodeMongoUpdateError = 703
	StatusCodeMongoDeleteError = 704
)

const (
	CreateErr     = "1-"
	ReadErr       = "2-"
	UpdateErr     = "3-"
	DeleteErr     = "4-"
	InvalidErr    = "5-"
	ExistedErr    = "6-"
	NotExistedErr = "7-"
)

const (
	ElementErr_Request = "1-"
)

const (
	InternalError = "100"
	GrpcSchollErr = "99"
)

const (
	ServiceErr_Request = "3-"
)
