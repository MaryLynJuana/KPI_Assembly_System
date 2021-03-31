package tested-binary

import "testing"

func TestHello(t *testing.T) {
	helloEmpty = hello("")
	if helloEmpty != "Hello, World!" {
		t.Errorf("hello(\"\") failed, expected: Hello, World!, got: ", helloEmpty)
	}

	helloName = hello("Valera")
	if helloName != "Hello, Valera!" {
		t.Errorf("hello(\"Valera\") failed, expected: Hello, Valera!, got: ", helloName)
	}
}
