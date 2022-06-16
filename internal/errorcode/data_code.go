package errorcode

type DataCode struct {
	data string
}

func NewDataCode(data string) *DataCode {
	return &DataCode{data: data}
}

func (dataCode *DataCode) getData() string {
	return dataCode.data
}

func (dataCode *DataCode) Error() string {
	return dataCode.getData()
}
