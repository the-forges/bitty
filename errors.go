package main

import "github.com/pkg/errors"

var (
	ErrUnitSymbolNotSupported       = errors.New("unit symbol not supported")
	ErrUnitSymbolNotSupportedf      = string(ErrUnitSymbolNotSupported.Error() + ": %s")
	ErrUnitSymbolEmptyNotSupportedf = errors.Errorf(ErrUnitSymbolNotSupported.Error() + ": empty symbol")
)

func NewErrUnitSymbolNotSupported(s UnitSymbol) error {
	if s == "" {
		return ErrUnitSymbolEmptyNotSupportedf
	}
	return errors.Errorf(ErrUnitSymbolNotSupportedf, s)
}
