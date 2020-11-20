package dao_sync

import (
	"errors"
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"github.com/illidan33/tools/common"
	"github.com/illidan33/tools/gen"
	"github.com/mohae/deepcopy"
	"go/format"
	"go/token"
	"os"
)

type KipleTemplateDaoSync struct {
	InterfaceName string

	gen.GenTemplate
	gen.TemplateModel
}

func (tm *KipleTemplateDaoSync) copyField(field *dst.Field) *dst.Field {
	newField := dst.Field{
		Names: []*dst.Ident{},
		Type:  deepcopy.Copy(field.Type).(dst.Expr),
		Tag:   deepcopy.Copy(field.Tag).(*dst.BasicLit),
		Decs: dst.FieldDecorations{
			NodeDecs: dst.NodeDecs{},
			Type:     dst.Decorations{},
		},
	}
	if field.Names != nil {
		for _, name := range field.Names {
			newField.Names = append(newField.Names, dst.NewIdent(name.Name))
		}
	}
	return &newField
}

func (tm *KipleTemplateDaoSync) copyKipleNode(ffv *dst.FuncDecl) *dst.Field {
	ffnew := dst.Field{}
	ffnew.Names = []*dst.Ident{dst.NewIdent(ffv.Name.Name)}
	newType := dst.FuncType{
		Func: false,
		Params: &dst.FieldList{
			Opening: true,
			List:    []*dst.Field{},
			Closing: true,
			Decs:    dst.FieldListDecorations{},
		},
		Results: &dst.FieldList{
			Opening: true,
			List:    []*dst.Field{},
			Closing: true,
			Decs:    dst.FieldListDecorations{},
		},
		Decs: dst.FuncTypeDecorations{},
	}
	for i := 0; i < len(ffv.Type.Results.List); i++ {
		newField := tm.copyField(ffv.Type.Results.List[i])
		newType.Results.List = append(newType.Results.List, newField)
	}
	for i := 0; i < len(ffv.Type.Params.List); i++ {
		newField := tm.copyField(ffv.Type.Params.List[i])
		newType.Params.List = append(newType.Params.List, newField)
	}
	ffnew.Type = &newType
	return &ffnew
}

func (tm *KipleTemplateDaoSync) FindInterfaceMethods(node *dst.File) (*dst.File, error) {
	var interfaceNode *dst.InterfaceType
	for _, decl := range node.Decls {
		if declv, ok := decl.(*dst.GenDecl); ok && declv.Tok == token.TYPE {
			if len(declv.Specs) == 0 {
				return node, errors.New("FindInterfaceAndFillMethods - GenDecl has no Specs")
			}
			if typespec, ok := declv.Specs[0].(*dst.TypeSpec); ok && typespec.Name.Name == tm.InterfaceName {
				if interfaceNode, ok = typespec.Type.(*dst.InterfaceType); ok {
					break
				}
			}
		}
	}
	if interfaceNode == nil {
		return node, fmt.Errorf("Interface %s not found", tm.InterfaceName)
	}

	newList := make([]*dst.Field, 0)
	for i := 0; i < len(node.Decls); i++ {
		decl := node.Decls[i]
		if ffv, ok := decl.(*dst.FuncDecl); ok && ffv.Recv != nil {
			if len(ffv.Recv.List) == 0 || ffv.Recv.List[0].Type == nil {
				continue
			}
			if ffvse, ok := (ffv.Recv.List[0].Type).(*dst.StarExpr); ok && ffvse.X.(*dst.Ident).Name == tm.ModelName {
				if common.IsUpperLetter(rune(ffv.Name.Name[0])) {
					fmt.Println(ffv.Name.Name)
					ffnew := tm.copyKipleNode(ffv)
					newList = append(newList, ffnew)
				}
			}
			if ffvse, ok := (ffv.Recv.List[0].Type).(*dst.Ident); ok && ffvse.Name == tm.ModelName {
				if common.IsUpperLetter(rune(ffv.Name.Name[0])) {
					fmt.Println(ffv.Name.Name)
					ffnew := tm.copyKipleNode(ffv)
					newList = append(newList, ffnew)
				}
			}
		}
	}

	interfaceNode.Methods.List = newList

	return node, nil
}

func (tm *KipleTemplateDaoSync) ParseToFile(dstFilePath string, node *dst.File) (err error) {
	fset, af, e := decorator.RestoreFile(node)
	if e != nil {
		err = e
		return
	}

	var file *os.File
	file, err = os.OpenFile(dstFilePath, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0644)
	if err != nil {
		return err
	}
	defer file.Close()
	err = format.Node(file, fset, af)

	return
}
