package configer

import (
	"strings"
)

type StringSlice []string

func (ss StringSlice) String() string {
	switch len(ss) {
	case 0:
		return ""
	case 1:
		return ss[0]
	}

	var sb strings.Builder
	sb.WriteString(ss[0])
	for i := 1; i < len(ss); i++ {
		sb.WriteRune(',')
		sb.WriteString(ss[i])
	}

	return sb.String()
}

func (ss *StringSlice) Set(val string) error {
	*ss = append(*ss, val)
	return nil
}
