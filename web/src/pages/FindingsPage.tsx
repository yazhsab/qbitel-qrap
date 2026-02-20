import React, { useEffect, useState, useCallback } from "react";

interface Finding {
  id: string;
  assessment_id: string;
  category: string;
  risk_level: string;
  title: string;
  affected_asset: string;
  current_algorithm: string | null;
  recommended_algorithm: string | null;
  discovered_at: string;
  [key: string]: unknown;
}

export const FindingsPage: React.FC = () => {
  const [findings, setFindings] = useState<Finding[]>([]);
  const [loading, setLoading] = useState(true);
  const [assessmentId, setAssessmentId] = useState("");

  const fetchFindings = useCallback(async () => {
    if (!assessmentId) {
      setFindings([]);
      setLoading(false);
      return;
    }
    setLoading(true);
    try {
      const res = await fetch(
        `/api/v1/findings?assessment_id=${assessmentId}&limit=50`,
      );
      const data = await res.json();
      setFindings(data.findings ?? []);
    } catch {
      setFindings([]);
    } finally {
      setLoading(false);
    }
  }, [assessmentId]);

  useEffect(() => {
    fetchFindings();
  }, [fetchFindings]);

  return (
    <div className="qtn-findings-page">
      <div className="qtn-card">
        <div className="qtn-card__header">
          <h2>Findings</h2>
          <p>Browse findings from risk assessments</p>
          <button className="qtn-btn qtn-btn--primary" onClick={fetchFindings}>
            Refresh
          </button>
        </div>
        <div style={{ marginBottom: "1rem" }}>
          <label htmlFor="assessment-id">Assessment ID: </label>
          <input
            id="assessment-id"
            type="text"
            placeholder="Enter assessment UUID"
            value={assessmentId}
            onChange={(e) => setAssessmentId(e.target.value)}
            style={{ padding: "0.25rem 0.5rem", width: "320px" }}
          />
        </div>
        {loading ? (
          <div className="qtn-loading">Loading...</div>
        ) : findings.length === 0 ? (
          <p>
            {assessmentId
              ? "No findings for this assessment."
              : "Enter an assessment ID to view findings."}
          </p>
        ) : (
          <table className="qtn-table">
            <thead>
              <tr>
                <th>Title</th>
                <th>Category</th>
                <th>Risk</th>
                <th>Asset</th>
                <th>Current</th>
                <th>Recommended</th>
              </tr>
            </thead>
            <tbody>
              {findings.map((f) => (
                <tr key={f.id}>
                  <td>{f.title}</td>
                  <td>{f.category}</td>
                  <td>{f.risk_level}</td>
                  <td>{f.affected_asset}</td>
                  <td>{f.current_algorithm ?? "-"}</td>
                  <td>{f.recommended_algorithm ?? "-"}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
};
