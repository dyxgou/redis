package storage

import "testing"

func assertInt(t *testing.T, v *Int, expected int) {
	if v.Kind() != IntKind {
		t.Errorf("Int kind expected=%d ('INT'). got=%d", IntKind, v.Kind())
	}

	if v.Value != expected {
		t.Errorf("Int value expected=%d. got=%d", expected, v.Value)
	}
}

func assertBool(t *testing.T, v *Bool, expected bool) {
	if v.Kind() != BoolKind {
		t.Errorf("Bool kind expected=%d ('BOOL'). got=%d", BoolKind, v.Kind())
	}

	if v.Value != expected {
		t.Errorf("Bool value expected=%t. got=%t", expected, v.Value)
	}
}

func assertInt64(t *testing.T, v *Int64, expected int64) {
	if v.Kind() != Int64Kind {
		t.Errorf("Int64 kind expected=%d ('INT64'). got=%d", Int64Kind, v.Kind())
	}

	if v.Value != expected {
		t.Errorf("Int64 value expected=%d. got=%d", expected, v.Value)
	}
}

func assertFloat(t *testing.T, v *Float, expected float64) {
	if v.Kind() != FloatKind {
		t.Errorf("Float kind expected=%d ('FLOAT'). got=%d", FloatKind, v.Kind())
	}

	if v.Value != expected {
		t.Errorf("Float value expected=%f. got=%f", expected, v.Value)
	}
}
func assertString(t *testing.T, v *String, expected string) {
	if v.Kind() != StringKind {
		t.Errorf("Float kind expected=%d ('STRING'). got=%d", StringKind, v.Kind())
	}

	if v.Value != expected {
		t.Errorf("Float value expected=%q. got=%q", expected, v.Value)
	}
}
