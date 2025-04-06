import React, { useState } from "react";
import UploadModal from "../modal/modal";

export default function Navbar({ onResults, onMeta }) {
  const [query, setQuery] = useState("");
  const [showFilters, setShowFilters] = useState(false);
  const [limit, setLimit] = useState(10);
  const [skip, setSkip] = useState(0);
  const [sortOrder, setSortOrder] = useState("desc");
  const [customQuery, setCustomQuery] = useState("");
  const [showUploadModal, setShowUploadModal] = useState(false);

  const handleSubmit = async (e) => {
    e.preventDefault();
    const url = query.trim()
      ? `http://localhost:9000/getDoc?query=${query}&start=${skip}&limit=${limit}`
      : `http://localhost:9000/get`;
    try {
      const res = await fetch(url);
      const data = await res.json();
      if (data.length > 0 && data[0].meta) {
        onMeta(data[0].meta);
        onResults(data.slice(1));
      } else {
        onMeta(null);
        onResults([]);
      }
    } catch (err) {
      console.error("Failed to fetch:", err);
      onMeta(null);
      onResults([]);
    }
  };

  const handleFilterSubmit = async () => {
    try {
      const parsedQuery = customQuery ? JSON.parse(customQuery) : {};
      const body = {
        query: parsedQuery,
        limit,
        start: skip,
        sort: sortOrder === "asc" ? "ascending" : "descending",
        matchType: "AND",
      };

      const res = await fetch("http://localhost:9000/search", {
        method: "POST",
        headers: {
          "Content-Type": "application/json",
        },
        body: JSON.stringify(body),
      });

      const data = await res.json();
      if (data.length > 0 && data[0].meta) {
        onMeta(data[0].meta);
        onResults(data.slice(1));
      } else {
        onMeta(null);
        onResults([]);
      }
    } catch (err) {
      console.error("Filter search failed:", err);
      onMeta(null);
      onResults([]);
    }
  };

  return (
    <>
      <nav className="navbar navbar-expand-lg bg-body-tertiary">
        <div className="container-fluid">
          <a className="navbar-brand" href="#">Go-Sage</a>
          <form className="d-flex" role="search" onSubmit={handleSubmit}>
            {!showFilters && (
              <>
                <input
                  className="form-control me-2"
                  type="search"
                  placeholder="Search"
                  value={query}
                  onChange={(e) => setQuery(e.target.value)}
                />
                <button className="btn btn-outline-success me-2" type="submit">
                  Search
                </button>
              </>
            )}
            <button
              type="button"
              className="btn btn-outline-secondary me-2"
              onClick={() => setShowFilters(!showFilters)}
            >
              {showFilters ? "Close Filters" : "Filters"}
            </button>
            <button
              type="button"
              className="btn btn-outline-primary"
              onClick={() => setShowUploadModal(true)}
            >
              Upload
            </button>
          </form>
        </div>
      </nav>

      {showFilters && (
        <div className="bg-light p-3 border-top border-bottom">
          <div className="container-fluid d-flex flex-wrap gap-3 align-items-center">
            <div className="form-group">
              <label className="form-label">Sort</label>
              <select
                className="form-select"
                value={sortOrder}
                onChange={(e) => setSortOrder(e.target.value)}
              >
                <option value="desc">Descending</option>
                <option value="asc">Ascending</option>
              </select>
            </div>

            <div className="form-group">
              <label className="form-label">Limit</label>
              <input
                type="number"
                className="form-control"
                value={limit}
                onChange={(e) => setLimit(Number(e.target.value))}
                min="1"
              />
            </div>

            <div className="form-group">
              <label className="form-label">Skip</label>
              <input
                type="number"
                className="form-control"
                value={skip}
                onChange={(e) => setSkip(Number(e.target.value))}
                min="0"
              />
            </div>

            <div className="form-group flex-grow-1">
              <label className="form-label">Custom Query</label>
              <input
                type="text"
                className="form-control"
                placeholder='{"key": "value"}'
                value={customQuery}
                onChange={(e) => setCustomQuery(e.target.value)}
              />
            </div>

            <button
              type="button"
              className="btn btn-primary align-self-end"
              onClick={handleFilterSubmit}
            >
              Find
            </button>
          </div>
        </div>
      )}

      {showUploadModal && <UploadModal onClose={() => setShowUploadModal(false)} />}
    </>
  );
}
