import React, { useState } from 'react';
import axios from 'axios';

const PredictorForm = () => {
  const [formData, setFormData] = useState({
    no_of_dependents: 0,
    education: 0,          // 1 = Graduate, 0 = Not Graduate
    self_employed: 0,      // 1 = Yes, 0 = No
    income_annum: '',
    loan_amount: '',
    loan_term: 12,
    cibil_score: '',
    total_assets: ''
  });

  const [result, setResult] = useState(null);
  const [loading, setLoading] = useState(false);
  const [error, setError] = useState('');

  const handleChange = (e) => {
    const { name, value } = e.target;
    setFormData(prev => ({ ...prev, [name]: value }));
  };

  const handleSubmit = async (e) => {
    e.preventDefault();
    setLoading(true);
    setError('');
    setResult(null);

    // Build request payload (ensure numeric fields are numbers)
    const payload = {
      no_of_dependents: parseInt(formData.no_of_dependents),
      education: parseInt(formData.education),
      self_employed: parseInt(formData.self_employed),
      income_annum: parseFloat(formData.income_annum),
      loan_amount: parseFloat(formData.loan_amount),
      loan_term: parseInt(formData.loan_term),
      cibil_score: parseFloat(formData.cibil_score),
      total_assets: parseFloat(formData.total_assets)
    };

    try {
      // Use environment variable for API base URL (fallback to localhost for dev)
      const apiBase = import.meta.env.VITE_API_URL || 'http://localhost:8080';
      const response = await axios.post(`${apiBase}/predict`, payload);
      setResult(response.data);
    } catch (err) {
      setError(err.response?.data || 'Failed to get prediction');
    } finally {
      setLoading(false);
    }
  };

  return (
    <div style={{ padding: '1rem 0' }}>
      <form onSubmit={handleSubmit} style={{ display: 'grid', gridTemplateColumns: '1fr 1fr', gap: '1rem' }}>
        <div>
          <label>Number of Dependents</label>
          <input type="number" name="no_of_dependents" value={formData.no_of_dependents} onChange={handleChange} min="0" required />
        </div>
        <div>
          <label>Education (1=Graduate, 0=Not Graduate)</label>
          <input type="number" name="education" value={formData.education} onChange={handleChange} min="0" max="1" required />
        </div>
        <div>
          <label>Self Employed (1=Yes, 0=No)</label>
          <input type="number" name="self_employed" value={formData.self_employed} onChange={handleChange} min="0" max="1" required />
        </div>
        <div>
          <label>Annual Income</label>
          <input type="number" name="income_annum" value={formData.income_annum} onChange={handleChange} required />
        </div>
        <div>
          <label>Loan Amount</label>
          <input type="number" name="loan_amount" value={formData.loan_amount} onChange={handleChange} required />
        </div>
        <div>
          <label>Loan Term (months)</label>
          <input type="number" name="loan_term" value={formData.loan_term} onChange={handleChange} min="1" required />
        </div>
        <div>
          <label>CIBIL Score</label>
          <input type="number" name="cibil_score" value={formData.cibil_score} onChange={handleChange} min="300" max="900" required />
        </div>
        <div>
          <label>Total Assets</label>
          <input type="number" name="total_assets" value={formData.total_assets} onChange={handleChange} required />
        </div>
        <div style={{ gridColumn: 'span 2', textAlign: 'center' }}>
          <button type="submit" className="button primary" disabled={loading}>
            {loading ? 'Checking...' : 'Check Eligibility'}
          </button>
        </div>
      </form>

      {error && <p style={{ color: 'red', marginTop: '1rem' }}>{error}</p>}

      {result && (
        <div style={{ marginTop: '1.5rem', padding: '1rem', border: '1px solid #ccc', borderRadius: '4px' }}>
          <h3>Result</h3>
          <p><strong>Eligible:</strong> {result.eligible ? '✅ Yes' : '❌ No'}</p>
          <p><strong>Probability:</strong> {(result.probability * 100).toFixed(2)}%</p>
          <p><strong>Risk Level:</strong> {result.risk}</p>
        </div>
      )}
    </div>
  );
};

export default PredictorForm;