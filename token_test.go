package token

import (
	"bytes"
	"database/sql"
	"database/sql/driver"
	"encoding"
	"encoding/gob"
	"encoding/json"
	"testing"
)

const testTokenLength = 128

func TestNew(t *testing.T) {
	token, err := New(testTokenLength)
	if err != nil {
		t.Fatal(err)
	}

	if returnedTokenLength := len(token); returnedTokenLength != testTokenLength {
		t.Errorf("expected %d bytes, got %d", testTokenLength, returnedTokenLength)
	}
}

func TestNewSequential(t *testing.T) {
	token1, err := New(testTokenLength)
	if err != nil {
		t.Fatal(err)
	}

	token2, err := New(testTokenLength)
	if err != nil {
		t.Fatal(err)
	}

	if bytes.Equal(token1, token2) {
		t.Errorf("sequential tokens were equal, %s", token1)
	}
}

func TestImplementsEncodingBinaryMarshaler(t *testing.T) {
	token := Token{}
	if _, ok := interface{}(token).(encoding.BinaryMarshaler); !ok {
		t.Error("Token does not implement encoding.BinaryMarshaler")
	}
}

func TestImplementsEncodingBinaryUnmarshaler(t *testing.T) {
	token := &Token{}
	if _, ok := interface{}(token).(encoding.BinaryUnmarshaler); !ok {
		t.Error("*Token does not implement encoding.BinaryUnmarshaler")
	}
}

func TestMarshalBinary(t *testing.T) {
	token, err := New(testTokenLength)
	if err != nil {
		t.Fatal(err)
	}

	tokenBinary, err := token.MarshalBinary()
	if err != nil {
		t.Fatal(err)
	}

	returnedToken := Token{}
	if err := returnedToken.UnmarshalBinary(tokenBinary); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(token, returnedToken) {
		t.Errorf("expected %s, got %s", token, returnedToken)
	}
}

func TestImplementsEncodingTextMarshaler(t *testing.T) {
	token := Token{}
	if _, ok := interface{}(token).(encoding.TextMarshaler); !ok {
		t.Error("Token does not implement encoding.TextMarshaler")
	}
}

func TestImplementsEncodingTextUnmarshaler(t *testing.T) {
	token := &Token{}
	if _, ok := interface{}(token).(encoding.TextUnmarshaler); !ok {
		t.Error("*Token does not implement encoding.TextUnmarshaler")
	}
}

func TestMarshalText(t *testing.T) {
	token, err := New(testTokenLength)
	if err != nil {
		t.Fatal(err)
	}

	tokenText, err := token.MarshalText()
	if err != nil {
		t.Fatal(err)
	}

	returnedToken := Token{}
	if err := returnedToken.UnmarshalText(tokenText); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(token, returnedToken) {
		t.Errorf("expected %s, got %s", token, returnedToken)
	}
}

func TestGobEncode(t *testing.T) {
	token, err := New(testTokenLength)
	if err != nil {
		t.Fatal(err)
	}

	buffer := bytes.NewBuffer([]byte{})
	if err := gob.NewEncoder(buffer).Encode(token); err != nil {
		t.Fatal(err)
	}

	returnedToken := Token{}
	if err := gob.NewDecoder(buffer).Decode(&returnedToken); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(token, returnedToken) {
		t.Errorf("expected %s, got %s", token, returnedToken)
	}
}

func TestJSONEncode(t *testing.T) {
	token, err := New(testTokenLength)
	if err != nil {
		t.Fatal(err)
	}

	buffer := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(buffer).Encode(token); err != nil {
		t.Fatal(err)
	}

	returnedToken := Token{}
	if err := json.NewDecoder(buffer).Decode(&returnedToken); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(token, returnedToken) {
		t.Errorf("expected %s, got %s", token, returnedToken)
	}
}

func TestImplementsDriverValuer(t *testing.T) {
	token := Token{}
	if _, ok := interface{}(token).(driver.Valuer); !ok {
		t.Error("Token does not implement driver.Valuer")
	}
}

func TestImplementsSQLScanner(t *testing.T) {
	token := &Token{}
	if _, ok := interface{}(token).(sql.Scanner); !ok {
		t.Error("*Token does not implement sql.Scanner")
	}
}

func TestValueScan(t *testing.T) {
	token, err := New(testTokenLength)
	if err != nil {
		t.Fatal(err)
	}

	tokenValue, err := token.Value()
	if err != nil {
		t.Fatal(err)
	}

	if !driver.IsValue(tokenValue) {
		t.Fatal("invalid driver value")
	}

	returnedToken := Token{}
	if err := returnedToken.Scan(tokenValue); err != nil {
		t.Fatal(err)
	}

	if !bytes.Equal(token, returnedToken) {
		t.Errorf("expected %s, got %s", token, returnedToken)
	}
}
