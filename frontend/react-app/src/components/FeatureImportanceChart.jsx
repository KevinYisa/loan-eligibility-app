import React, { useEffect, useState } from 'react';
import { BarChart, Bar, XAxis, YAxis, Tooltip, ResponsiveContainer, Cell } from 'recharts';
import axios from 'axios';

const FeatureImportanceChart = () => {
  const [data, setData] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState('');

  useEffect(() => {
    const fetchData = async () => {
      try {
        const apiBase = import.meta.env.VITE_API_URL || 'http://localhost:8080';
        const response = await axios.get(`${apiBase}/features`);
        // Sort descending for display
        const sorted = response.data.sort((a, b) => b.importance - a.importance);
        setData(sorted);
      } catch (err) {
        setError('Failed to load feature importance');
      } finally {
        setLoading(false);
      }
    };
    fetchData();
  }, []);

  if (loading) return <p>Loading feature importance...</p>;
  if (error) return <p style={{ color: 'red' }}>{error}</p>;

  // Define colors
  const colors = ['#0088FE', '#00C49F', '#FFBB28', '#FF8042', '#FF6384', '#9966FF', '#FFB347', '#66B3FF'];

  return (
    <div style={{ width: '100%', height: 300, marginTop: '1.5rem' }}>
      <h3>Feature Importance</h3>
      <ResponsiveContainer width="100%" height="100%">
        <BarChart layout="vertical" data={data} margin={{ top: 5, right: 30, left: 100, bottom: 5 }}>
          <XAxis type="number" />
          <YAxis type="category" dataKey="feature" width={100} tick={{ fontSize: 12 }} />
          <Tooltip formatter={(value) => (value * 100).toFixed(2) + '%'} />
          <Bar dataKey="importance">
            {data.map((entry, index) => (
              <Cell key={`cell-${index}`} fill={colors[index % colors.length]} />
            ))}
          </Bar>
        </BarChart>
      </ResponsiveContainer>
    </div>
  );
};

export default FeatureImportanceChart;