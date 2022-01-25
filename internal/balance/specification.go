package balance

type Specification interface {
	ToSQLClosure() string
}

type ForUpdateSpec struct{}

func (s ForUpdateSpec) ToSQLClosure() string {
	return "FOR UPDATE"
}
