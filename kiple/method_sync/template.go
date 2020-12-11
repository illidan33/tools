package method_sync

import (
	"errors"
	"fmt"
	"github.com/dave/dst"
	"github.com/dave/dst/decorator"
	"tools/common"
	"tools/gen"
	"github.com/mohae/deepcopy"
	"go/format"
	"go/token"
	"os"
)

type KipleTemplatemethodsync struct {
	InterfaceName string

	gen.GenTemplate
	gen.TemplateModel
}

func (tm *KipleTemplatemethodsync) copyField(field *dst.Field) *dst.Field {
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

func (tm *KipleTemplatemethodsync) copyKipleNode(ffv *dst.FuncDecl) *dst.Field {
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
func (tm *KipleTemplatemethodsync) FindSourceInterfaceNode(node *dst.File) (interfaceNode *dst.InterfaceType, err error) {
	for _, decl := range node.Decls {
		if declv, ok := decl.(*dst.GenDecl); ok && declv.Tok == token.TYPE {
			if len(declv.Specs) == 0 {
				err = errors.New("FindInterfaceAndFillMethods - GenDecl has no Specs")
				return
			}
			if typespec, ok := declv.Specs[0].(*dst.TypeSpec); ok && typespec.Name.Name == tm.InterfaceName {
				if interfaceNode, ok = typespec.Type.(*dst.InterfaceType); ok {
					break
				}
			}
		}
	}
	if interfaceNode == nil {
		err = fmt.Errorf("Interface %s not found", tm.InterfaceName)
		return
	}
	return
}

func (tm *KipleTemplatemethodsync) FindInterfaceMethods(node *dst.File, interfaceNode *dst.InterfaceType) error {
	if interfaceNode == nil {
		return fmt.Errorf("Interface %s not found", tm.InterfaceName)
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

	interfaceNode.Methods.List = append(interfaceNode.Methods.List, newList...)

	return nil
}

func (tm *KipleTemplatemethodsync) ParseToFile(dstFilePath string, node *dst.File) (err error) {
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
