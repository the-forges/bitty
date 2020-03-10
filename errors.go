package bitty

/*
	Copyright 2020 IBM

	Licensed under the Apache License, Version 2.0 (the "License");
	you may not use this file except in compliance with the License.
	You may obtain a copy of the License at

		http://www.apache.org/licenses/LICENSE-2.0

	Unless required by applicable law or agreed to in writing, software
	distributed under the License is distributed on an "AS IS" BASIS,
	WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
	See the License for the specific language governing permissions and
	limitations under the License.
*/

import "github.com/pkg/errors"

var (
	ErrUnitSymbolNotSupported       = errors.New("unit symbol not supported")
	ErrUnitSymbolNotSupportedf      = string(ErrUnitSymbolNotSupported.Error() + ": %s")
	ErrUnitSymbolEmptyNotSupportedf = errors.Errorf(ErrUnitSymbolNotSupported.Error() + ": empty symbol")
	ErrUnitExponentNotSupported     = errors.New("unit exponent not supported")
	ErrUnitExponentNotSupportedf    = string(ErrUnitExponentNotSupported.Error() + ": %s")
	ErrUnitStandardNotSupported     = errors.New("unit standard not supported")
	ErrUnitStandardNotSupportedf    = string(ErrUnitStandardNotSupported.Error() + ": %s")
)

func NewErrUnitSymbolNotSupported(s UnitSymbol) error {
	if s == "" {
		return ErrUnitSymbolEmptyNotSupportedf
	}
	return errors.Errorf(ErrUnitSymbolNotSupportedf, s)
}

func NewErrUnitStandardNotSupported(s UnitStandard) error {
	return errors.Errorf(ErrUnitSymbolNotSupportedf, s)
}
