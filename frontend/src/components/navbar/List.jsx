import React from "react";

export default function List({ results, meta }) {
  if (results === null) return null;
  if (!results || results.length === 0)
    return <p className="text-center mt-4">No results found.</p>;

  const keys = Object.keys(results[0]);
  const reorderedKeys = ["MsgId", ...keys.filter(k => k !== "MsgId")];

  return (
    <div style={{ overflowX: "auto" }}>
      {meta && (
        <div className="alert alert-info text-center my-3">
          <p className="mb-1">
            Matched Records: <strong>{meta.matchedRecords}</strong>
          </p>
          <p className="mb-0">
            Total Time: <strong>{meta.totalTime}</strong>
          </p>
        </div>
      )}

      <table className="table table-bordered table-striped">
        <thead className="table-dark">
          <tr>
            {reorderedKeys.map((key) => (
              <th key={key}>{key}</th>
            ))}
          </tr>
        </thead>
        <tbody>
          {results.map((item, index) => (
            <tr key={index}>
              {reorderedKeys.map((key) => (
                <td
                  key={key}
                  style={{
                    maxWidth: "250px",
                    overflow: "auto",
                    whiteSpace: "nowrap",
                  }}
                >
                  {String(item[key])}
                </td>
              ))}
            </tr>
          ))}
        </tbody>
      </table>

      <div className="text-center">
        <button
          className="btn btn-danger mt-3"
          onClick={() => window.location.reload()}
        >
          Reset
        </button>
      </div>
    </div>
  );
}
