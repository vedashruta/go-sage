import React, { useState } from "react";

export default function UploadModal({ onClose }) {
  const [file, setFile] = useState(null);
  const [uploadResult, setUploadResult] = useState(null);

  const handleFileUpload = async () => {
    if (!file) return;
    const formData = new FormData();
    formData.append("file", file);

    try {
      const res = await fetch("http://localhost:9000/upload", {
        method: "POST",
        body: formData,
      });

      if (!res.ok) throw new Error("Upload failed");

      const result = await res.json();
      setUploadResult(result);
    } catch (err) {
      console.error("Upload error:", err);
      setUploadResult({ error: "Upload failed" });
    }
  };

  return (
    <div className="modal d-block" tabIndex="-1" role="dialog" style={{ background: "rgba(0,0,0,0.5)" }}>
      <div className="modal-dialog" role="document">
        <div className="modal-content">
          <div className="modal-header">
            <h5 className="modal-title">Upload .csv or .parquet File</h5>
            <button type="button" className="btn-close" onClick={onClose}></button>
          </div>
          <div className="modal-body">
            <input
              type="file"
              accept=".csv,.parquet"
              className="form-control"
              onChange={(e) => setFile(e.target.files[0])}
            />
            {uploadResult && (
              <div className="mt-3">
                {uploadResult.error ? (
                  <p className="text-danger">Error: {uploadResult.error}</p>
                ) : (
                  <>
                    <p className="text-success">{uploadResult.message}</p>
                    <p>Documents Parsed: {uploadResult.documentsParsed}</p>
                    <p>Time Taken: {uploadResult.duration}</p>
                  </>
                )}
              </div>
            )}
          </div>
          <div className="modal-footer">
            <button onClick={handleFileUpload} className="btn btn-primary">Upload</button>
            <button onClick={onClose} className="btn btn-secondary">Close</button>
          </div>
        </div>
      </div>
    </div>
  );
}
