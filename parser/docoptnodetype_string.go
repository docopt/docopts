// Code generated by "stringer -type=DocoptNodeType"; DO NOT EDIT.

package docopt_language

import "strconv"

func _() {
	// An "invalid array index" compiler error signifies that the constant values have changed.
	// Re-run the stringer command to generate them again.
	var x [1]struct{}
	_ = x[Unmatched_node - -1]
	_ = x[Root-2]
	_ = x[Prologue-3]
	_ = x[Prologue_node-4]
	_ = x[Usage_section-5]
	_ = x[Usage-6]
	_ = x[Usage_line-7]
	_ = x[Prog_name-8]
	_ = x[Usage_short_option-9]
	_ = x[Usage_long_option-10]
	_ = x[Usage_argument-11]
	_ = x[Usage_unmatched_punct-12]
	_ = x[Usage_command-13]
	_ = x[Usage_optional_group-14]
	_ = x[Usage_required_group-15]
	_ = x[Group_alternative-16]
	_ = x[Free_section-17]
	_ = x[Section_name-18]
	_ = x[Section_node-19]
	_ = x[Options_section-20]
	_ = x[Options_node-21]
	_ = x[Option_line-22]
	_ = x[Option_short-23]
	_ = x[Option_long-24]
	_ = x[Option_argument-25]
	_ = x[Option_alternative_group-26]
	_ = x[Option_description-27]
	_ = x[Description_node-28]
}

const (
	_DocoptNodeType_name_0 = "Unmatched_node"
	_DocoptNodeType_name_1 = "RootProloguePrologue_nodeUsage_sectionUsageUsage_lineProg_nameUsage_short_optionUsage_long_optionUsage_argumentUsage_unmatched_punctUsage_commandUsage_optional_groupUsage_required_groupGroup_alternativeFree_sectionSection_nameSection_nodeOptions_sectionOptions_nodeOption_lineOption_shortOption_longOption_argumentOption_alternative_groupOption_descriptionDescription_node"
)

var (
	_DocoptNodeType_index_1 = [...]uint16{0, 4, 12, 25, 38, 43, 53, 62, 80, 97, 111, 132, 145, 165, 185, 202, 214, 226, 238, 253, 265, 276, 288, 299, 314, 338, 356, 372}
)

func (i DocoptNodeType) String() string {
	switch {
	case i == -1:
		return _DocoptNodeType_name_0
	case 2 <= i && i <= 28:
		i -= 2
		return _DocoptNodeType_name_1[_DocoptNodeType_index_1[i]:_DocoptNodeType_index_1[i+1]]
	default:
		return "DocoptNodeType(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
