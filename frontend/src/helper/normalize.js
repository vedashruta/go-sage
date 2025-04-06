const DEFAULT_KEYS = [
  "AppName",
  "Event",
  "EventId",
  "Facility",
  "FacilityString",
  "Groupings",
  "Hostname",
  "Message",
  "MessageRaw",
  "MsgId",
  "NanoTimeStamp",
  "PartitionId",
  "Priority",
  "ProcId",
  "Sender",
  "Severity",
  "SeverityString",
  "StructuredData",
  "Tag",
  "Timestamp",
  "namespace",
];

export function normalizeResults(data) {
  return data.map((item) => {
    const parsed = {};
    DEFAULT_KEYS.forEach((key) => {
      parsed[key] = item[key] !== undefined ? item[key] : "";
    });
    return parsed;
  });
}
