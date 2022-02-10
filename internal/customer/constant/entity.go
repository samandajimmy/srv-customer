package constant

type SubjectType = int8

const (
	UserSubjectType = SubjectType(iota)
)

type ModifierRole = string

const (
	UserModifierRole  = ModifierRole("USER")
	AdminModifierRole = ModifierRole("ADMIN")
)
