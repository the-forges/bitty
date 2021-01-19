# Bitty

[![Build Status](https://travis-ci.org/the-forges/bitty.svg?branch=main)](https://travis-ci.org/the-forges/bitty) [![Go Report Card](https://goreportcard.com/badge/github.com/the-forges/bitty)](https://goreportcard.com/report/github.com/the-forges/bitty) [![codecov](https://codecov.io/gh/the-forges/bitty/branch/main/graph/badge.svg?token=QLM3JEIUFU)](https://codecov.io/gh/the-forges/bitty)

Bitty is a memory unit conversion library that makes working with multiple sizes and SI/IEC standards straight forward. It is based on unit ecapsulation, immutability, plugability, and testability: the idea that each unit should know how to operate with other valid units idempotently; every unit function returns a new unit or value, instead of changing itself; each unit implements interfaces so that it's easy to plug in new unit types (which also makes testing new unit types straight forward).

## Features

### Standards Compliance

- [x] Full IEC Binary notation [SI 9th edition (page 145)](https://www.bipm.org/utils/common/pdf/si-brochure/SI-Brochure-9.pdf) compliance
- [ ] Full SI Decimal notation [SI 9th edition (page 145)](https://www.bipm.org/utils/common/pdf/si-brochure/SI-Brochure-9.pdf) compliance (**partial**)

### Mathematics

- [x] Adding different units against each other
- [x] Subtracting units from each other
- [ ] Multiplying units against each other
- [ ] Dividing units from each other

### Conversions

- [x] Unit conversions
- [x] Standard conversions

### Helpers

- [x] Unit parsing
- [ ] Finding a specific unit
