package logberry


type Context struct {
	Root *Root
	Class ContextClass
	Label string
}


//----------------------------------------------------------------------
//----------------------------------------------------------------------
func NewContext(root *Root, class ContextClass, label string) *Context {
	return &Context {
		Root: root,
		Class: class,
		Label: label,
	}
}

// toString
