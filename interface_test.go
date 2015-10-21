package main

// MockResourceData can be used to test functions that take ResourceData
type MockResourceData struct {
	ID    string
	Attrs map[string]interface{}
}

func NewMockResourceData(id string, attrs map[string]interface{}) *MockResourceData {
	return &MockResourceData{id, attrs}
}

func (d *MockResourceData) Get(key string) interface{} {
	return d.Attrs[key]
}

func (d *MockResourceData) SetId(id string) {
	d.ID = id
}

func (d *MockResourceData) Set(key string, value interface{}) error {
	d.Attrs[key] = value
	return nil
}
