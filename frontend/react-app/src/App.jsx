import React from 'react';
import PredictorForm from './components/PredictorForm';
import FeatureImportanceChart from './components/FeatureImportanceChart';
import MetricsDisplay from './components/MetricsDisplay';
import ContactForm from './components/ContactForm';
import './App.css';

function App() {
  return (
    <div className="react-app">
      <h2>Loan Eligibility Predictor</h2>
      <PredictorForm />
      <FeatureImportanceChart />
      <MetricsDisplay />
      
      <hr style={{ margin: '2rem 0', borderColor: 'rgba(255,255,255,0.1)' }} />
      
      <h2>Contact Us</h2>
      <ContactForm />
    </div>
  );
}

export default App;