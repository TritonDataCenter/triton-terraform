package main

// MockResourceData can be used to test functions that take ResourceData
// TODO: add assertions on this to make sure we're not leaving the real thing in
// a bad state
type MockResourceData struct {
	ID        string
	Attrs     map[string]interface{}
	IsPartial bool
	Changes   []string
}

func NewMockResourceData(id string, attrs map[string]interface{}) *MockResourceData {
	return &MockResourceData{id, attrs, false, []string{}}
}

func (d *MockResourceData) Get(key string) interface{} {
	v, _ := d.GetOk(key)
	return v
}

func (d *MockResourceData) GetOk(key string) (interface{}, bool) {
	if key == "" {
		return d.Attrs, true
	}

	v, ok := d.Attrs[key]
	return v, ok
}

func (d *MockResourceData) Set(key string, value interface{}) error {
	d.Attrs[key] = value
	return nil
}

func (d *MockResourceData) Change(key string, value interface{}) error {
	d.Attrs[key] = value
	d.Changes = append(d.Changes, key)
	return nil
}

func (d *MockResourceData) Id() string {
	return d.ID
}

func (d *MockResourceData) SetId(id string) {
	d.ID = id
}

func (d *MockResourceData) Partial(state bool) {
	d.IsPartial = state
}

func (d *MockResourceData) SetPartial(field string) {
	changes := []string{}
	for _, change := range d.Changes {
		if change != field {
			changes = append(changes, change)
		}
	}
	d.Changes = changes
}

func (d *MockResourceData) HasChange(field string) bool {
	for _, change := range d.Changes {
		if change == field {
			return true
		}
	}

	return false
}
