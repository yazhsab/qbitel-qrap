import React, { useEffect, useState, useCallback } from "react";

interface Assessment {
  id: string;
  name: string;
  organization_id: string;
  status: string;
  overall_risk: string | null;
  risk_score: number;
  created_at: string;
  [key: string]: unknown;
}

export const AssessmentsPage: React.FC = () => {
  const [assessments, setAssessments] = useState<Assessment[]>([]);
  const [loading, setLoading] = useState(true);

  const fetchAssessments = useCallback(async () => {
    setLoading(true);
    try {
      const res = await fetch("/api/v1/assessments?limit=50");
      const data = await res.json();
      setAssessments(data.assessments ?? []);
    } catch {
      setAssessments([]);
    } finally {
      setLoading(false);
    }
  }, []);

  useEffect(() => {
    fetchAssessments();
  }, [fetchAssessments]);

  return (
    <div className="qtn-assessments-page">
      <div className="qtn-card">
        <div className="qtn-card__header">
          <h2>Risk Assessments</h2>
          <p>Manage quantum risk assessments for your organisation</p>
          <button className="qtn-btn qtn-btn--primary" onClick={fetchAssessments}>
            Refresh
          </button>
        </div>
        {loading ? (
          <div className="qtn-loading">Loading...</div>
        ) : assessments.length === 0 ? (
          <p>No assessments found. Create one using the API.</p>
        ) : (
          <table className="qtn-table">
            <thead>
              <tr>
                <th>Name</th>
                <th>Status</th>
                <th>Risk Level</th>
                <th>Score</th>
                <th>Created</th>
              </tr>
            </thead>
            <tbody>
              {assessments.map((a) => (
                <tr key={a.id}>
                  <td>{a.name}</td>
                  <td>{a.status}</td>
                  <td>{a.overall_risk ?? "-"}</td>
                  <td>{a.risk_score}</td>
                  <td>{a.created_at}</td>
                </tr>
              ))}
            </tbody>
          </table>
        )}
      </div>
    </div>
  );
};
