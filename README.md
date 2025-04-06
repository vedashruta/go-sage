# 🦊 GO Sage — In-Memory Search Engine in Golang

**GO Sage** is a blazing fast, in-memory full-text search engine written in Golang. It processes large datasets (like 24k+ CSV records) with stemming, stop-word filtering, and indexing — all designed for sub-second search latency.

---

## 📐 Architecture & Design Choices

- **In-Memory Reverse Index**  
  Built using `map[string][]docID`, mapping stemmed, normalized tokens to document IDs.

- **Normalization Pipeline**

  - Lowercasing
  - Punctuation removal
  - Tokenization (by whitespace)
  - Stop-word removal
  - Porter stemming

- **Storage**

  - Indexed documents stored in `model.store`
  - Recent docs tracked via ordered `[]docID` array

- **Concurrency**

  - `ParseCSVFromDir()` uses 10 goroutines for parallel CSV parsing
  - Thread-safe doc counter via `sync/atomic`

- **Search Features**
  - Filter by field
  - Limit results (`FindOptions`)
  - Get latest documents (`Get`)
  - Track total records, matched count, returned count, and time taken

---

## ⚡ Performance Tuning Strategies

- Bounded goroutine pool for parsing
- Fast tokenizer & stemmer
- In-place string manipulation to avoid memory bloat
- Avoiding deferred file closures inside loops
- Atomic counters for low-overhead shared metrics

## Features Implemented

☑️ Stemmed indexing & stop-word filtering</br>
☑️ Concurrency-safe CSV parsing</br>
☑️ Fast document retrieval</br>
☑️ Meta search stats (matched/returned/total/time)</br>
☑️ Frontend (React) to visualize results</br>

## 🗂️ Project Structure

GO-Sage/</br>
── backend/</br>
&nbsp;├── main.go</br>
&nbsp;├── model.go</br>
&nbsp;├── search.go</br>
&nbsp;├── parser.go</br>
&nbsp;├── helper/</br>
&nbsp;├── search_test.go</br>
&nbsp;├── benchmark_test.go</br>
&nbsp;└── go.mod</br>
</br>
── frontend/</br>
&nbsp;├── public/</br>
&nbsp;├── src/</br>
&nbsp;&nbsp;&nbsp;├── App.js</br>
&nbsp;&nbsp;&nbsp;├── components/</br>
&nbsp;&nbsp;&nbsp;└── index.js</br>
&nbsp;├─ package.json</br>
&nbsp;└── .env</br>
&nbsp;
── README.md</br>

## ⚙️ How to Build & Run

### 🔨 Build

### Backend (Go)
```bash
cd backend
go mod tidy
go build -o server .
./server
```

### Frontend (React)
```bash
cd frontend
npm install
npm run dev
```

## 🧪 Testing & Benchmark

### ✅ Unit Tests

```bash
go test -v ./...
```

### 📊 Benchmarking

```bash
go test -bench=. -benchmem
```

## Future Enhancements
Export/load index from disk (WIP)</br>
Support phrase matching & advanced queries (WIP)</br>
