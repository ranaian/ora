// Copyright 2014 Rana Ian. All rights reserved.
// Use of this source code is governed by The MIT License
// found in the accompanying LICENSE file.

package ora

/*
#include <oci.h>
#include "version.h"
*/
import "C"
import (
	"unsafe"
)

type bndFloat32Ptr struct {
	stmt      *Stmt
	ocibnd    *C.OCIBind
	ociNumber [1]C.OCINumber
	value     *float32
	nullp
}

func (bnd *bndFloat32Ptr) bind(value *float32, position int, stmt *Stmt) error {
	bnd.stmt = stmt
	bnd.value = value
	bnd.nullp.Set(value == nil)
	if value != nil {
		if err := bnd.stmt.ses.srv.env.OCINumberFromFloat(&bnd.ociNumber[0], float64(*value), 4); err != nil {
			return err
		}
	}
	r := C.OCIBINDBYPOS(
		bnd.stmt.ocistmt, //OCIStmt      *stmtp,
		&bnd.ocibnd,
		bnd.stmt.ses.srv.env.ocierr,         //OCIError     *errhp,
		C.ub4(position),                     //ub4          position,
		unsafe.Pointer(&bnd.ociNumber[0]),   //void         *valuep,
		C.LENGTH_TYPE(C.sizeof_OCINumber),   //sb8          value_sz,
		C.SQLT_VNU,                          //ub2          dty,
		unsafe.Pointer(bnd.nullp.Pointer()), //void         *indp,
		nil,           //ub2          *alenp,
		nil,           //ub2          *rcodep,
		0,             //ub4          maxarr_len,
		nil,           //ub4          *curelep,
		C.OCI_DEFAULT) //ub4          mode );
	if r == C.OCI_ERROR {
		return bnd.stmt.ses.srv.env.ociError()
	}
	return nil
}

func (bnd *bndFloat32Ptr) setPtr() error {
	if bnd.nullp.IsNull() {
		return nil
	}
	val, err := bnd.stmt.ses.srv.env.OCINumberToFloat(&bnd.ociNumber[0], 4)
	*bnd.value = float32(val)
	return err
}

func (bnd *bndFloat32Ptr) close() (err error) {
	defer func() {
		if value := recover(); value != nil {
			err = errR(value)
		}
	}()

	stmt := bnd.stmt
	bnd.stmt = nil
	bnd.ocibnd = nil
	bnd.value = nil
	bnd.nullp.Free()
	stmt.putBnd(bndIdxFloat32Ptr, bnd)
	return nil
}
