package web

type Args struct {
	Db databaser
}

func NewWith(args *Args) *Web {
	return &Web{db: args.Db}
}
