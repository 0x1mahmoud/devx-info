package utils

import (
	"strings"

	"cuelang.org/go/cue"
	"cuelang.org/go/cue/build"
	"cuelang.org/go/cue/cuecontext"
	cueload "cuelang.org/go/cue/load"
)

func LoadInstances(configDir string) []*build.Instance {
	buildConfig := &cueload.Config{
		Dir:     configDir,
		Overlay: map[string]cueload.Source{},
	}
	return cueload.Instances([]string{}, buildConfig)
}

func LoadProject(configDir string) cue.Value {
	buildConfig := &cueload.Config{
		Dir:     configDir,
		Overlay: map[string]cueload.Source{},
	}
	instances := cueload.Instances([]string{}, buildConfig)

	ctx := cuecontext.New()

	return ctx.BuildInstance(instances[0])
}

func GetLastPathFragement(value cue.Value) string {
	selector := value.Path().Selectors()
	return selector[len(selector)-1].String()
}

func GetComments(value cue.Value) string {
	comments := value.Doc()
	result := ""
	for _, comment := range comments {
		result += comment.Text()
	}
	return strings.ReplaceAll(result, "\n", " ")
}

func HasComments(value cue.Value) bool {
	comments := value.Doc()

	return len(comments) > 0
}

func Walk(v cue.Value, before func(cue.Value) bool, after func(cue.Value)) {
	switch v.Kind() {
	case cue.StructKind:
		if before != nil && !before(v) {
			return
		}
		fieldIter, _ := v.Fields(cue.All())
		for fieldIter.Next() {
			Walk(fieldIter.Value(), before, after)
		}
	case cue.ListKind:
		if before != nil && !before(v) {
			return
		}
		valueIter, _ := v.List()
		for valueIter.Next() {
			Walk(valueIter.Value(), before, after)
		}
	default:
		if before != nil {
			before(v)
		}
	}
	if after != nil {
		after(v)
	}
}

func RemoveMeta(value cue.Value) (cue.Value, error) {
	result := value.Context().CompileString("_")

	iter, err := value.Fields()
	if err != nil {
		return result, err
	}

	for iter.Next() {
		v := iter.Value()
		if !strings.HasPrefix(v.Path().String(), "$") {
			result = result.FillPath(v.Path(), v)
		}
	}

	return result, nil
}
