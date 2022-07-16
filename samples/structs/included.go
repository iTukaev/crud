package main

import "errors"

type IncData struct {
	i int
}

func (i *IncData) check(m *MyCoolStruct) (string, error) {
	if i == nil {
		return "", errors.New("is nill")
	}
	if m.privateField == 1 {
		return "", errors.New("Я художник, я так вижу")
	}
	return m.PublicField + ", сойдет", nil
}
