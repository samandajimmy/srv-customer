package constant

type SubjectType int8

const (
	UserSubjectType = SubjectType(iota)
	SystemSubjectType
)

type ModifierRole = string

const (
	UserModifierRole   = ModifierRole("USER")
	AdminModifierRole  = ModifierRole("ADMIN")
	SystemModifierRole = ModifierRole("SYSTEM")
)
