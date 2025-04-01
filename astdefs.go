package main

// AST defs from https://bytedance.larkoffice.com/wiki/LTfKwv88aiVFZBkZDYDcCzlBn6e

// Below defined by hob

// 有一个 map[PkgPath]... ， 所以为了能 json marshal 必须是 string
type PkgPath string

type FileLine struct {
	file string
	line int
}

const (
	TypeKindStruct    = iota
	TypeKindInterface = iota
	TypeKindTypedef   = iota
)

type TypeKind struct {
	kind int
}

// Below copied from lark doc
type Identity struct {
	ModPath string // ModPath is the module which the package belongs to
	PkgPath        // Import Path of the package
	Name    string // Unique Name of declaration (FunctionName, TypeName.MethodName, InterfaceName<TypeName>.MethodName, or TypeName)
}

// Repository
type Repository struct {
	Name    string             `json:"id"` // module name
	Modules map[string]*Module // module name => Library
}

type Module struct {
	Name         string               // go module name
	Dir          string               // relative path to repo
	Packages     map[PkgPath]*Package // pkage import path => Package
	Dependencies map[string]string    `json:",omitempty"` // module name => module_path@version
	Files        map[string]*File     `json:",omitempty"` // relative path => file info
}

type File struct {
	Name    string   // file name
	Imports []string // imported symbols (in language spec)
}

type Package struct {
	IsMain bool
	PkgPath
	Functions    map[string]*Function // Function name (may be {{func}} or {{struct.method}}) => Function
	Types        map[string]*Type     // type name => type define
	Vars         map[string]*Var      // var name => var define
	CompressData *string              `json:"compress_data,omitempty"` // package summary
}

type Function struct {
	Exported bool // If public

	IsMethod bool // If the function is a method
	Identity      // unique identity in a repo
	FileLine
	Content string // codes of the function, including functiion signature and body

	Receiver *Receiver  `json:",omitempty"` // Method receiver
	Params   []Identity `json:",omitempty"` // function parameters, key is the parameter name
	Results  []Identity `json:",omitempty"` // function results, key is the result name or type name

	// call to in-the-project functions, key is {{pkgAlias.funcName}} or {{funcName}}
	FunctionCalls []Identity `json:",omitempty"`

	// call to internal methods,
	// NOTICE: method name may be duplicated, so we collect according to the SEQUENCE of APPEARANCE
	MethodCalls []Identity `json:",omitempty"`

	Types      []Identity `json:",omitempty"` // types used in the function
	GlobalVars []Identity `json:",omitempty"` // global vars used in the function

	// func llm compress result
	CompressData *string `json:"compress_data,omitempty"`
}

type Receiver struct {
	IsPointer bool
	Type      Identity
	Name      string
}

type Type struct {
	Exported bool // if the struct is exported

	TypeKind // type Kind: Struct / Interface / Typedef
	Identity // unique id in a repo
	FileLine
	Content string // struct declaration content

	// field type (not include basic types), type name => type id
	SubStruct []Identity `json:",omitempty"`

	// inline field type (not include basic types)
	InlineStruct []Identity `json:",omitempty"`

	// methods defined on the Struct, not including inlined type's method
	Methods map[string]Identity `json:",omitempty"`

	// Implemented interfaces
	Implements []Identity `json:",omitempty"`

	// functions defined in fields, key is type name, val is the function Signature
	// FieldFunctions map[string]string

	CompressData *string `json:"compress_data,omitempty"` // struct llm compress result
}

type Var struct {
	IsExported bool
	IsConst    bool
	Identity
	FileLine
	Type    *Identity `json:",omitempty"`
	Content string

	CompressData *string `json:"compress_data,omitempty"`
}
