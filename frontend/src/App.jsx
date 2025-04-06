import React, { useState } from "react";
import Navbar from "./components/navbar/Navbar";
import List from "./components/navbar/List";
import { parseDurationToSeconds } from "./helper/time";

export default function App() {
  const [results, setResults] = useState([]);
  const [meta, setMeta] = useState(null);

  return (
    <div className="container mt-4">
      <Navbar onResults={setResults} onMeta={setMeta} />

      {meta && (
        <div className="bg-white border rounded p-3 text-center my-3">
          <strong>Search Stats:</strong>
          <br />
          <span>Total Records: {meta.totalRecords}</span> |{" "}
          <span>Matched Records: {meta.matchedRecords}</span> |{" "}
          <span>Returned Records: {meta.returnedRecords}</span> |{" "}
          <span>Total Time: {parseDurationToSeconds(meta.totalTime)} sec</span>
        </div>
      )}

      <List results={results} />
    </div>
  );
}
