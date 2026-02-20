import React, { useEffect, useState } from "react";

interface Stats {
  totalAssessments: number;
  completedAssessments: number;
  totalOrganizations: number;
  criticalFindings: number;
}

export const DashboardPage: React.FC = () => {
  const [stats, setStats] = useState<Stats>({
    totalAssessments: 0,
    completedAssessments: 0,
    totalOrganizations: 0,
    criticalFindings: 0,
  });
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    const fetchStats = async () => {
      try {
        const [assessmentsRes, orgsRes] = await Promise.all([
          fetch("/api/v1/assessments?limit=1"),
          fetch("/api/v1/organizations?limit=1"),
        ]);

        const assessmentsData = await assessmentsRes.json();
        const orgsData = await orgsRes.json();

        setStats({
          totalAssessments: assessmentsData.total_count ?? 0,
          completedAssessments: assessmentsData.total_count ?? 0,
          totalOrganizations: orgsData.total_count ?? 0,
          criticalFindings: 0,
        });
      } catch {
        // API not available
      } finally {
        setLoading(false);
      }
    };

    fetchStats();
  }, []);

  if (loading) {
    return <div className="qtn-loading">Loading dashboard...</div>;
  }

  return (
    <div className="qtn-dashboard">
      <h1>Risk Assessment Dashboard</h1>
      <div className="qtn-dashboard__grid">
        <div className="qtn-card">
          <h3>Organizations</h3>
          <div className="qtn-stat">{stats.totalOrganizations}</div>
        </div>
        <div className="qtn-card">
          <h3>Total Assessments</h3>
          <div className="qtn-stat">{stats.totalAssessments}</div>
        </div>
        <div className="qtn-card">
          <h3>Completed</h3>
          <div className="qtn-stat">{stats.completedAssessments}</div>
        </div>
        <div className="qtn-card">
          <h3>Critical Findings</h3>
          <div className="qtn-stat">{stats.criticalFindings}</div>
        </div>
      </div>
    </div>
  );
};
