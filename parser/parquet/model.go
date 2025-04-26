package parquet

type LogEntry struct {
	MsgId          string `parquet:"name=MsgId, type=BYTE_ARRAY, convertedtype=UTF8" json:"MsgId"`
	PartitionId    int64  `parquet:"name=PartitionId, type=INT64" json:"PartitionId"`
	Timestamp      string `parquet:"name=Timestamp, type=BYTE_ARRAY, convertedtype=UTF8" json:"Timestamp"`
	Hostname       string `parquet:"name=Hostname, type=BYTE_ARRAY, convertedtype=UTF8" json:"Hostname"`
	Priority       int32  `parquet:"name=Priority, type=INT32" json:"Priority"`
	Facility       int32  `parquet:"name=Facility, type=INT32" json:"Facility"`
	FacilityString string `parquet:"name=FacilityString, type=BYTE_ARRAY, convertedtype=UTF8" json:"FacilityString"`
	Severity       int32  `parquet:"name=Severity, type=INT32" json:"Severity"`
	SeverityString string `parquet:"name=SeverityString, type=BYTE_ARRAY, convertedtype=UTF8" json:"SeverityString"`
	AppName        string `parquet:"name=AppName, type=BYTE_ARRAY, convertedtype=UTF8" json:"AppName"`
	ProcId         string `parquet:"name=ProcId, type=BYTE_ARRAY, convertedtype=UTF8" json:"ProcId"`
	Message        string `parquet:"name=Message, type=BYTE_ARRAY, convertedtype=UTF8" json:"Message"`
	MessageRaw     string `parquet:"name=MessageRaw, type=BYTE_ARRAY, convertedtype=UTF8" json:"MessageRaw"`
	StructuredData string `parquet:"name=StructuredData, type=BYTE_ARRAY, convertedtype=UTF8" json:"StructuredData"`
	Tag            string `parquet:"name=Tag, type=BYTE_ARRAY, convertedtype=UTF8" json:"Tag"`
	Sender         string `parquet:"name=Sender, type=BYTE_ARRAY, convertedtype=UTF8" json:"Sender"`
	Groupings      string `parquet:"name=Groupings, type=BYTE_ARRAY, convertedtype=UTF8" json:"Groupings"`
	Event          string `parquet:"name=Event, type=BYTE_ARRAY, convertedtype=UTF8" json:"Event"`
	EventId        string `parquet:"name=EventId, type=BYTE_ARRAY, convertedtype=UTF8" json:"EventId"`
	NanoTimeStamp  string `parquet:"name=NanoTimeStamp, type=BYTE_ARRAY, convertedtype=UTF8" json:"NanoTimeStamp"`
	Namespace      string `parquet:"name=namespace, type=BYTE_ARRAY, convertedtype=UTF8" json:"Namespace"`
}
