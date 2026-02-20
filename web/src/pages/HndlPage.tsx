import React, { useState } from "react";

interface HndlResult {
  algorithm: string;
  estimated_break_year: number;
  data_shelf_life_years: number;
  risk_window_years: number;
  is_at_risk: boolean;
  urgency: string;
}

const ALGORITHMS = [
  "RSA-2048",
  "RSA-3072",
  "RSA-4096",
  "ECDSA-P256",
  "ECDSA-P384",
  "Ed25519",
  "X25519",
  "ML-KEM-768",
  "ML-DSA-65",
];

export const HndlPage: React.FC = () => {
  const [algorithm, setAlgorithm] = useState("RSA-2048");
  const [shelfLife, setShelfLife] = useState(10);
  const [result, setResult] = useState<HndlResult | null>(null);
  const [loading, setLoading] = useState(false);

  const calculate = async () => {
    setLoading(true);
    try {
      const res = await fetch("/api/v1/hndl", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          algorithm,
          data_shelf_life_years: shelfLife,
        }),
      });
      const data = await res.json();
      setResult(data);
    } catch {
      setResult(null);
    } finally {
      setLoading(false);
    }
  };

  return (
    <div className="qtn-hndl-page">
      <div className="qtn-card">
        <h2>Harvest Now, Decrypt Later Analysis</h2>
        <p>Estimate HNDL risk based on algorithm and data shelf life</p>
        <div style={{ display: "flex", gap: "1rem", marginBottom: "1rem" }}>
          <div>
            <label htmlFor="algorithm">Algorithm: </label>
            <select
              id="algorithm"
              value={algorithm}
              onChange={(e) => setAlgorithm(e.target.value)}
            >
              {ALGORITHMS.map((a) => (
                <option key={a} value={a}>
                  {a}
                </option>
              ))}
            </select>
          </div>
          <div>
            <label htmlFor="shelf-life">Data shelf life (years): </label>
            <input
              id="shelf-life"
              type="number"
              min={1}
              max={50}
              value={shelfLife}
              onChange={(e) => setShelfLife(Number(e.target.value))}
              style={{ width: "60px" }}
            />
          </div>
          <button
            className="qtn-btn qtn-btn--primary"
            onClick={calculate}
            disabled={loading}
          >
            {loading ? "Calculating..." : "Calculate"}
          </button>
        </div>

        {result && (
          <div className="qtn-hndl-result">
            <div className="qtn-card">
              <h3>Result</h3>
              <table style={{ width: "100%", borderCollapse: "collapse" }}>
                <tbody>
                  <tr>
                    <td style={{ padding: "0.5rem", fontWeight: "bold" }}>Algorithm</td>
                    <td style={{ padding: "0.5rem" }}>{result.algorithm}</td>
                  </tr>
                  <tr>
                    <td style={{ padding: "0.5rem", fontWeight: "bold" }}>
                      Estimated break year
                    </td>
                    <td style={{ padding: "0.5rem" }}>{result.estimated_break_year}</td>
                  </tr>
                  <tr>
                    <td style={{ padding: "0.5rem", fontWeight: "bold" }}>
                      Risk window (years)
                    </td>
                    <td style={{ padding: "0.5rem" }}>{result.risk_window_years}</td>
                  </tr>
                  <tr>
                    <td style={{ padding: "0.5rem", fontWeight: "bold" }}>At risk?</td>
                    <td style={{ padding: "0.5rem" }}>
                      {result.is_at_risk ? "YES" : "NO"}
                    </td>
                  </tr>
                  <tr>
                    <td style={{ padding: "0.5rem", fontWeight: "bold" }}>Urgency</td>
                    <td style={{ padding: "0.5rem" }}>{result.urgency}</td>
                  </tr>
                </tbody>
              </table>
            </div>
          </div>
        )}
      </div>
    </div>
  );
};
