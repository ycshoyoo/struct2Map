package struct2Map

import (
	"encoding/json"
	"fmt"
	"testing"
)

type T1 struct {
	Name []*T2     `json:"name"`
	Age  int       `json:"age"`
	Qw   []*string `json:"qw"`
}

type T2 struct {
	N1 *string `json:"n1, omitempty"`
}

func TestStructToMap(t *testing.T) {
	var t2 T2
	m, _ := StructToMap(t2, "json")
	bys, _ := json.Marshal(m)
	fmt.Println(string(bys))
}

type benchmarkUser struct {
	Name    string `json:"name"`
	Age     int    `json:"age"`
	Address string `json:"address"`
	Contact string `json:"contact"`
}

func newBenchmarkUser() benchmarkUser {
	return benchmarkUser{
		Name:    "name",
		Age:     18,
		Address: "github address",
		Contact: "github contact",
	}
}

func BenchmarkStructToMapByJson(b *testing.B) {
	user := newBenchmarkUser()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		data, _ := json.Marshal(&user)
		m := make(map[string]interface{})
		json.Unmarshal(data, &m)
	}
}

func BenchmarkStructToMapByToMap(b *testing.B) {
	user := newBenchmarkUser()
	tag := "json"
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		StructToMap(&user, tag)
	}
}
