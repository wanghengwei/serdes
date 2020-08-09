package main

import "github.com/pkg/errors"

func (w *Serializer) SerializeFoo(v *Foo) error {
	var err error
	err = w.Serializeint(&v.A)
	if err != nil {
		return errors.WithStack(err)
	}

	len := len(v.B)
	err = w.Serializeint(&len)
	if err != nil {
		return err
	}
	for i := range v.B {
		err = w.Serializestring(&v.B[i])
		if err != nil {
			return err
		}
	}
	for k, v := range v.C {
		err = w.Serializeint(&k)
		if err != nil {
			return err
		}
		err = w.Serializestring(&v)
		if err != nil {
			return err
		}
	}
	err = w.SerializeBar(&v.D)
	if err != nil {
		return err
	}
	for i := range v.E {
		err = w.SerializeBar(&v.E[i])
		if err != nil {
			return err
		}
	}
	for k, v := range v.F {
		err = w.Serializestring(&k)
		if err != nil {
			return err
		}
		err = w.SerializeBar(&v)
		if err != nil {
			return err
		}
	}

	return nil
}
func (w *Serializer) SerializeBar(v *Bar) error {
	var err error
	err = w.Serializestring(&v.A)
	if err != nil {
		return err
	}
	err = w.Serializebool(&v.B)
	if err != nil {
		return err
	}

	return nil
}
