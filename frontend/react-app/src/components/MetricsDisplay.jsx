import React, { useEffect, useState } from 'react';
import axios from 'axios';

const MetricsDisplay = () => {
  const [metrics, setMetrics] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchData = async () => {
      const apiBase = import.meta.env.VITE_API_URL || '';
      
      if (!apiBase) {
        console.error("CRITICAL: VITE_API_URL is not defined in the environment variables!");
      }

      try {
        const response = await axios.get(`${apiBase}/metrics`);
        setMetrics(response.data);
      } catch (err) {
        setError('Failed to load metrics');
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  if (loading) return <p>Loading metrics...</p>;
  if (error) return <p style={{ color: 'red' }}>{error}</p>;
  if (!metrics) return null;

  return (
    <div style={{ marginTop: '1.5rem' }}>
      <h3>Model Performance</h3>
      <div style={{ display: 'grid', gridTemplateColumns: 'repeat(auto-fit, minmax(120px, 1fr))', gap: '0.5rem' }}>
        <div><strong>Accuracy:</strong> {(metrics.accuracy * 100).toFixed(2)}%</div>
        <div><strong>Precision:</strong> {(metrics.precision * 100).toFixed(2)}%</div>
        <div><strong>Recall:</strong> {(metrics.recall * 100).toFixed(2)}%</div>
        <div><strong>F1:</strong> {(metrics.f1 * 100).toFixed(2)}%</div>
        <div><strong>ROC AUC:</strong> {(metrics.roc_auc * 100).toFixed(2)}%</div>
      </div>
    </div>
  );
};

export default MetricsDisplay;