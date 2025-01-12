package storage

import "testing"

func assertInt(t *testing.T, v *Int, expected int) {
	if v.kind() != intKind {
		t.Errorf("Int kind expected=%d ('INT'). got=%d", intKind, v.kind())
	}

	if v.Value != expected {
		t.Errorf("Int value expected=%d. got=%d", expected, v.Value)
	}
}

func assertBool(t *testing.T, v *Bool, expected bool) {
	if v.kind() != boolKind {
		t.Errorf("Bool kind expected=%d ('BOOL'). got=%d", boolKind, v.kind())
	}

	if v.Value != expected {
		t.Errorf("Bool value expected=%t. got=%t", expected, v.Value)
	}
}

func assertInt64(t *testing.T, v *Int64, expected int64) {
	if v.kind() != int64Kind {
		t.Errorf("Int64 kind expected=%d ('INT64'). got=%d", intKind, v.kind())
	}

	if v.Value != expected {
		t.Errorf("Int64 value expected=%d. got=%d", expected, v.Value)
	}
}

func assertFloat(t *testing.T, v *Float, expected float64) {
	if v.kind() != floatKind {
		t.Errorf("Float kind expected=%d ('FLOAT'). got=%d", intKind, v.kind())
	}

	if v.Value != expected {
		t.Errorf("Float value expected=%f. got=%f", expected, v.Value)
	}
}
func assertString(t *testing.T, v *String, expected string) {
	if v.kind() != stringKind {
		t.Errorf("Float kind expected=%d ('STRING'). got=%d", intKind, v.kind())
	}

	if v.Value != expected {
		t.Errorf("Float value expected=%q. got=%q", expected, v.Value)
	}
}
