package util

import "testing"

func TestToUnderscore(t *testing.T) {
	str1 := "abcDef"
	t.Log(str1, "=>", ToUnderscore(str1))
	str2 := "AbcDef"
	t.Log(str2, "=>", ToUnderscore(str2))
	str3 := "AbcDefG"
	t.Log(str3, "=>", ToUnderscore(str3))
	str4 := "abc_Def"
	t.Log(str4, "=>", ToUnderscore(str4))
	str5 := "abcDEf"
	t.Log(str5, "=>", ToUnderscore(str5))
	str6 := "abcDEF"
	t.Log(str6, "=>", ToUnderscore(str6))
}

func TestToCamelCase(t *testing.T) {
	str1 := "abc_def"
	t.Log(str1, "=>", ToCamelCase(str1))
	str2 := "_abc_def"
	t.Log(str2, "=>", ToCamelCase(str2))
	str3 := "abc_def_"
	t.Log(str3, "=>", ToCamelCase(str3))
	str4 := "abc_de_f"
	t.Log(str4, "=>", ToCamelCase(str4))
	str5 := "a_bc_def"
	t.Log(str5, "=>", ToCamelCase(str5))
	str6 := "abc_def__ghi"
	t.Log(str6, "=>", ToCamelCase(str6))
}
