package main

type Foo struct {
	A int
	B []string
	C map[int]string
	D Bar
	E []Bar
	F map[string]Bar
}

type Bar struct {
	A string
	B bool
}
