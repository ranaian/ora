// Copyright 2014 Rana Ian. All rights reserved.
// Use of this source code is governed by The MIT License
// found in the accompanying LICENSE file.

package ora

/*
#include <stdlib.h>
#include <oci.h>
#include "version.h"
*/
import "C"
import "unsafe"

type bndInt64Ptr struct {
	stmt      *Stmt
	ocibnd    *C.OCIBind
	ociNumber [1]C.OCINumber
	value     *int64
	nullp
}

func (bnd *bndInt64Ptr) bind(value *int64, position int, stmt *Stmt) error {
	//bnd.stmt.logF(_drv.cfg.Log.Stmt.Bind, "Int64Ptr.bind(%d) value=%#v => number=%#v", position, value, bnd.ociNumber[0])
	bnd.stmt = stmt
	bnd.value = value
	bnd.nullp.Set(value == nil)
	if value != nil {
		if err := bnd.stmt.ses.srv.env.OCINumberFromInt(&bnd.ociNumber[0], *value, 8); err != nil {
			return err
		}
		bnd.stmt.logF(_drv.cfg.Log.Stmt.Bind,
			"Int64Ptr.bind(%d) value=%#v => number=%#v", position, value, bnd.ociNumber[0])
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

func (bnd *bndInt64Ptr) setPtr() error {
	if bnd.nullp.IsNull() {
		return nil
	}
	var err error
	*bnd.value, err = bnd.stmt.ses.srv.env.OCINumberToInt(&bnd.ociNumber[0], 8)
	return err
}

func (bnd *bndInt64Ptr) close() (err error) {
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
	stmt.putBnd(bndIdxInt64Ptr, bnd)
	return nil
}
