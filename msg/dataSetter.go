package msg

import "fmt"

var _ Processor = (*DataSetter)(nil)

// DataSetter sets data to any message which is logged.
type DataSetter struct {

	// Key is the key of the data to set.
	Key string

	// Value is the value of the data to set.
	Value []interface{}
}

// NewDataSetter creates a new data setter message processor.
func NewDataSetter(key string, value ...interface{}) *DataSetter {
	return &DataSetter{
		Key:   key,
		Value: value,
	}
}

// Process will set the data of the given message and return the message.
func (ds *DataSetter) Process(msg *Message) *Message {
	if (ds != nil) && (msg != nil) {
		msg.Add(ds.Key, ds.Value...)
	}
	return msg
}

// String gets the string for this message processor.
func (ds *DataSetter) String() string {
	if ds == nil {
		return nilStr
	}
	return fmt.Sprint("DataSetter(", ds.Key, ": ", fmt.Sprint(ds.Value...), ")")
}
