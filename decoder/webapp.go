package decoder

type Binder interface {
	Bind(c webappv2.Context, i interface{}) error
}

type webappBinder struct {
}
