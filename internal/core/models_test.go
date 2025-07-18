package core

import "testing"

func TestPhpType_String(t *testing.T) {
	cases := []struct {
		name     string
		pt       phpType
		expected string
	}{
		{"plain", phpType{Name: "int"}, "int"},
		{"array", phpType{Name: "int", IsArray: true}, "array"},
		{"nullable", phpType{Name: "int", IsNull: true}, "?int"},
		{"nullable array", phpType{Name: "int", IsArray: true, IsNull: true}, "array"},
	}

	for _, tc := range cases {
		if got := tc.pt.String(); got != tc.expected {
			t.Errorf("phpType.String() (%s) = %q, want %q", tc.name, got, tc.expected)
		}
	}
}

func TestPhpType_IsDateTimeImmutable(t *testing.T) {
	pt := phpType{Name: "\\DateTimeImmutable"}
	if !pt.IsDateTimeImmutable() {
		t.Errorf("IsDateTimeImmutable() should be true for \\DateTimeImmutable")
	}

	pt = phpType{Name: "int"}
	if pt.IsDateTimeImmutable() {
		t.Errorf("IsDateTimeImmutable() should be false for int")
	}
}

func TestPhpType_IsInt(t *testing.T) {
	pt := phpType{Name: "int"}
	if !pt.IsInt() {
		t.Errorf("IsInt() should be true for int")
	}

	pt = phpType{Name: "string"}
	if pt.IsInt() {
		t.Errorf("IsInt() should be false for string")
	}
}

func TestPhpType_IsFloat(t *testing.T) {
	pt := phpType{Name: "float"}
	if !pt.IsFloat() {
		t.Errorf("IsFloat() should be true for float")
	}

	pt = phpType{Name: "int"}
	if pt.IsFloat() {
		t.Errorf("IsFloat() should be false for int")
	}
}

func TestPhpType_IsString(t *testing.T) {
	pt := phpType{Name: "string"}
	if !pt.IsString() {
		t.Errorf("IsString() should be true for string")
	}

	pt = phpType{Name: "int"}
	if pt.IsString() {
		t.Errorf("IsString() should be false for int")
	}
}

func TestQueryValue_Type(t *testing.T) {
	pt := phpType{Name: "int"}
	qv := QueryValue{Name: "foo", Typ: pt}
	if got := qv.Type(); got != "int" {
		t.Errorf("QueryValue.Type() = %q, want %q", got, "int")
	}

	mc := &ModelClass{Name: "Bar"}
	qv = QueryValue{Name: "bar", Struct: mc}
	if got := qv.Type(); got != "Bar" {
		t.Errorf("QueryValue.Type() = %q, want %q", got, "Bar")
	}
}

func TestQueryValue_IsStruct(t *testing.T) {
	qv := QueryValue{Name: "foo"}
	if qv.IsStruct() {
		t.Errorf("IsStruct() should be false when Struct is nil")
	}

	mc := &ModelClass{Name: "Bar"}
	qv.Struct = mc
	if !qv.IsStruct() {
		t.Errorf("IsStruct() should be true when Struct is set")
	}
}
