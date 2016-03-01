package main

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestSplitEnvironKeyValue(t *testing.T) {
	var pair, k, v string
	var err error
	Convey("Split environment key-value pairs", t, func() {
		pair = "key1=value1"
		k, v, err = splitEnvironKeyValue(pair)
		So(err, ShouldBeNil)
		So(k, ShouldEqual, "key1")
		So(v, ShouldEqual, "value1")

		pair = "key1=value1=alsovalue"
		k, v, err = splitEnvironKeyValue(pair)
		So(err, ShouldBeNil)
		So(k, ShouldEqual, "key1")
		So(v, ShouldNotEqual, "value1")
		So(v, ShouldEqual, "value1=alsovalue")
	})
}

func TestAppendEnvironment(t *testing.T) {
	var origEnviron, newEnviron []string
	var err error
	Convey("Append environment", t, func() {
		origEnviron = []string{"key1=value1", "key2=value2"}
		newEnviron, err = appendEnvironment(origEnviron)
		So(err, ShouldBeNil)
		So(len(newEnviron), ShouldEqual, len(origEnviron))
		So(newEnviron, ShouldResemble, origEnviron)

		origEnviron = []string{"key1=value1", "key2=value2"}
		newEnviron, err = appendEnvironment(origEnviron, "KEY2=value2new")
		So(err, ShouldBeNil)
		So(len(newEnviron), ShouldEqual, len(origEnviron))
		So(newEnviron, ShouldNotResemble, origEnviron)
		So(newEnviron[0], ShouldEqual, origEnviron[0])
		So(newEnviron[1], ShouldNotEqual, origEnviron[1])
		So(newEnviron[1], ShouldEqual, "key2=value2new") //The "original key" is kept

		origEnviron = []string{"key1=value1", "key2=value2"}
		newEnviron, err = appendEnvironment(origEnviron, "KEY2=value2new", "Key3=value3")
		So(err, ShouldBeNil)
		So(len(newEnviron), ShouldNotEqual, len(origEnviron))
		So(len(newEnviron), ShouldEqual, 3)
		So(newEnviron, ShouldNotResemble, origEnviron)
		So(newEnviron[0], ShouldEqual, origEnviron[0])
		So(newEnviron[1], ShouldNotEqual, origEnviron[1])
		So(newEnviron[1], ShouldEqual, "key2=value2new") //The "original key" is kept
		So(newEnviron[2], ShouldEqual, "Key3=value3")
	})
}

func _testsPhaseMapContains(m map[string]nodeData, s string) bool {
	_, ok := m[s]
	return ok
}

func TestDeleteVariablesFromPhasesMap(t *testing.T) {
	Convey("Delete variables from phases map", t, func() {
		m1 := map[string]nodeData{"variables": nodeData{}, "variables1": nodeData{}, "myvariables": nodeData{}}
		deleteVariablesFromPhasesMap(m1)
		So(len(m1), ShouldEqual, 2)
		So(_testsPhaseMapContains(m1, "variables"), ShouldBeFalse)
		So(_testsPhaseMapContains(m1, "variables1"), ShouldBeTrue)
		So(_testsPhaseMapContains(m1, "myvariables"), ShouldBeTrue)

		m2 := map[string]nodeData{"VARIABLES": nodeData{}, "VARIABLES1": nodeData{}, "MYVARIABLES": nodeData{}}
		deleteVariablesFromPhasesMap(m2)
		So(len(m2), ShouldEqual, 2)
		So(_testsPhaseMapContains(m2, "VARIABLES"), ShouldBeFalse)
		So(_testsPhaseMapContains(m2, "VARIABLES1"), ShouldBeTrue)
		So(_testsPhaseMapContains(m2, "MYVARIABLES"), ShouldBeTrue)

		m3 := map[string]nodeData{"Variables": nodeData{}, "Variables1": nodeData{}, "MyVariables": nodeData{}}
		deleteVariablesFromPhasesMap(m3)
		So(len(m3), ShouldEqual, 2)
		So(_testsPhaseMapContains(m3, "Variables"), ShouldBeFalse)
		So(_testsPhaseMapContains(m3, "Variables1"), ShouldBeTrue)
		So(_testsPhaseMapContains(m3, "MyVariables"), ShouldBeTrue)
	})
}

func TestReplaceVariables(t *testing.T) {
	var origStr, replacedStr string
	var variables map[string]string
	Convey("Replace variables", t, func() {
		origStr = "$COMMAND_VAR & echo test123"
		variables = map[string]string{"COMMAND_VAR": "where python & echo Hallo"}

		replacedStr = replaceVariables(origStr, variables)
		So(replacedStr, ShouldEqual, "where python & echo Hallo & echo test123")
	})
}
