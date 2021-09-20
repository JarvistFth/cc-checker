// Copyright 2020 Google LLC
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// https://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

// Package utils contains various utility functions.
package utils

import (
	"cc-checker/logger"
	"go/build"
	"go/types"
	"golang.org/x/tools/go/callgraph"
	"golang.org/x/tools/go/ssa"
	"strings"
)

var log = logger.GetLogger()

// Dereference returns the underlying type of a pointer.
// If the input is not a pointer, then the type of the input is returned.
func Dereference(t types.Type) types.Type {
	for {
		tt, ok := t.Underlying().(*types.Pointer)
		if !ok {
			return t
		}
		t = tt.Elem()
	}
}

// DecomposeType returns the path and name of a Named type
// Returns empty strings if the type is not *types.Named
func DecomposeType(t types.Type) (path, name string) {
	n, ok := t.(*types.Named)
	if !ok {
		return
	}

	if pkg := n.Obj().Pkg(); pkg != nil {
		path = pkg.Path()
	}

	return path, n.Obj().Name()
}

// DecomposeField returns the decomposed type of the
// struct containing the field, as well as the field's name.
// If the referenced struct's type is not a named type,
// the type path and name will both be empty strings.
func DecomposeField(t types.Type, field int) (typePath, typeName, fieldName string) {
	deref := Dereference(t)
	typePath, typeName = DecomposeType(deref)
	fieldName = deref.Underlying().(*types.Struct).Field(field).Name()
	return
}

// UnqualifiedName returns the name of the given type, without the qualifying
// prefix containing the package in which it was declared.
// Example: for a type named T declared in package p, the returned string will
// be just `T` instead of `p.T`.
func UnqualifiedName(t *types.Var) string {
	return types.TypeString(t.Type(), func(*types.Package) string { return "" })
}

// DecomposeFunction returns the path, receiver, and name strings of a ssa.Function.
// For functions that have no receiver, returns an empty string for recv.
// For shared functions (wrappers and error.Error), returns an empty string for path.
// Panics if provided a nil argument.
func DecomposeFunction(f *ssa.Function) (path, recv, name string) {
	if f.Pkg != nil {
		path = f.Pkg.Pkg.Path()
	}
	name = f.Name()
	if recvVar := f.Signature.Recv(); recvVar != nil {
		recv = UnqualifiedName(recvVar)
	}

	log.Debugf("path:%s, recv:%s, name:%s", path, recv, name)

	return
}

func DecomposeAbstractMethod(callCom *ssa.CallCommon) (path, recv, name string) {
	if callCom.Method.Pkg() != nil {
		path = callCom.Method.Pkg().Path()
	}

	typename := callCom.Value.Type().String()
	relativepkgs := strings.Split(typename, "/")
	objname := relativepkgs[len(relativepkgs)-1]
	recv = objname
	name = callCom.Method.Name()
	return
}
func FindMethod(prog *ssa.Program, ssaPkg *ssa.Package, name string) *ssa.Function {
	var ret *ssa.Function
	for _, member := range ssaPkg.Members {
		if ty, ok := member.(*ssa.Type); ok {
			t := ty.Type()
			p := types.NewPointer(t)
			initselt := prog.MethodSets.MethodSet(t).Lookup(ssaPkg.Pkg, name)
			initselp := prog.MethodSets.MethodSet(p).Lookup(ssaPkg.Pkg, name)
			if initselt != nil {
				ret = prog.LookupMethod(t, ssaPkg.Pkg, name)
			}
			if initselp != nil {
				ret = prog.LookupMethod(p, ssaPkg.Pkg, name)
			}
			if ret == nil {
				continue
			} else {
				break
			}
		}
	}
	return ret
}

func FindMethodByType(prog *ssa.Program, ssaPkg *ssa.Package, t types.Type, name string) *ssa.Function {
	var ret *ssa.Function
	for _, impotedPkg := range ssaPkg.Pkg.Imports() {
		p := types.NewPointer(t)
		initselt := prog.MethodSets.MethodSet(t).Lookup(impotedPkg, name)
		initselp := prog.MethodSets.MethodSet(p).Lookup(impotedPkg, name)
		log.Info(prog.MethodSets.MethodSet(t).String())

		if initselt != nil {
			ret = prog.LookupMethod(t, impotedPkg, name)
		}
		if initselp != nil {
			ret = prog.LookupMethod(p, impotedPkg, name)
		}
	}

	return ret
}

func FindMethodWithAllPkgs(prog *ssa.Program, ssaPkgs []*ssa.Package, t types.Type, name string) *ssa.Function {
	var ret *ssa.Function
	p := types.NewPointer(t)
	for _, ssaPkg := range ssaPkgs {
		selectiont := prog.MethodSets.MethodSet(t).Lookup(ssaPkg.Pkg, name)
		selectionp := prog.MethodSets.MethodSet(p).Lookup(ssaPkg.Pkg, name)
		if selectiont != nil {
			ret = prog.LookupMethod(t, ssaPkg.Pkg, name)
		}
		if selectionp != nil {
			ret = prog.LookupMethod(p, ssaPkg.Pkg, name)
		}
		if ret != nil {
			log.Infof("find method %s in pkg: %s", name, ssaPkg)
			return ret
		}
	}
	return ret
}

func FindInvokeMethod(prog *ssa.Program, mainPkg *ssa.Package) (*ssa.Function, *ssa.Function) {
	var initf *ssa.Function
	var invokef *ssa.Function
	for _, member := range mainPkg.Members {
		if ty, ok := member.(*ssa.Type); ok {
			t := ty.Type()
			p := types.NewPointer(t)
			initselt := prog.MethodSets.MethodSet(t).Lookup(mainPkg.Pkg, "Invoke")
			initselp := prog.MethodSets.MethodSet(p).Lookup(mainPkg.Pkg, "Invoke")
			if initselt != nil {
				initf = prog.LookupMethod(t, mainPkg.Pkg, "Init")
				invokef = prog.LookupMethod(t, mainPkg.Pkg, "Invoke")
			}
			if initselp != nil {
				initf = prog.LookupMethod(p, mainPkg.Pkg, "Init")
				invokef = prog.LookupMethod(p, mainPkg.Pkg, "Invoke")
			}
			if initf == nil || invokef == nil {
				continue
			} else {
				break
			}
		}
	}
	return initf, invokef
}

func ReverseNewSlice(s []interface{}) []interface{} {
	t := make([]interface{}, len(s))
	j := len(s) - 1
	for i, _ := range s {
		t[i] = s[j]
		j -= 1
	}
	return t
}

func IsSynthetic(edge *callgraph.Edge) bool {
	return edge.Caller.Func.Pkg == nil || edge.Callee.Func.Synthetic != ""
}

func InStd(node *callgraph.Node) bool {
	pkg, _ := build.Import(node.Func.Pkg.Pkg.Path(), "", 0)
	return pkg.Goroot
}

func InFabric(node *callgraph.Node) bool {
	return strings.Contains(node.Func.Pkg.Pkg.Path(), "fabric")
}
