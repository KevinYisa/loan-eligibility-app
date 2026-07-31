# scripts/train.py

import pandas as pd
import numpy as np
from sklearn.linear_model import LogisticRegression
from sklearn.preprocessing import StandardScaler
from sklearn.metrics import (
    accuracy_score, precision_score, recall_score, f1_score,
    roc_auc_score, confusion_matrix
)
import json
import os
from datetime import datetime

# Paths
PROCESSED_DIR = '../processed/'
MODELS_DIR = '../models/'
FEATURES_FILE = '../models/feature_names.json'

# Ensure models directory exists
os.makedirs(MODELS_DIR, exist_ok=True)

# Load feature names
with open(FEATURES_FILE, 'r') as f:
    feature_names = json.load(f)

print("Feature order:", feature_names)

# Load splits
X_train = pd.read_csv(os.path.join(PROCESSED_DIR, 'X_train.csv'))
X_val = pd.read_csv(os.path.join(PROCESSED_DIR, 'X_val.csv'))
X_test = pd.read_csv(os.path.join(PROCESSED_DIR, 'X_test.csv'))
y_train = pd.read_csv(os.path.join(PROCESSED_DIR, 'y_train.csv')).squeeze('columns')
y_val = pd.read_csv(os.path.join(PROCESSED_DIR, 'y_val.csv')).squeeze('columns')
y_test = pd.read_csv(os.path.join(PROCESSED_DIR, 'y_test.csv')).squeeze('columns')

# Ensure columns are in the correct order
X_train = X_train[feature_names]
X_val = X_val[feature_names]
X_test = X_test[feature_names]

# Check class distribution
print("Train class distribution:\n", y_train.value_counts())

# Scale features
scaler = StandardScaler()
X_train_scaled = scaler.fit_transform(X_train)
X_val_scaled = scaler.transform(X_val)
X_test_scaled = scaler.transform(X_test)

# Save scaler parameters (mean and scale)
scaler_params = {
    'mean': scaler.mean_.tolist(),
    'scale': scaler.scale_.tolist()
}
with open(os.path.join(MODELS_DIR, 'scaler_params.json'), 'w') as f:
    json.dump(scaler_params, f, indent=2)

# Train Logistic Regression
# We'll use balanced class_weight to handle potential imbalance
model = LogisticRegression(
    C=1.0,
    solver='liblinear',  # good for small datasets and gives coefficients
    max_iter=1000,
    class_weight='balanced',
    random_state=42
)
model.fit(X_train_scaled, y_train)

# Predictions on test set
y_pred = model.predict(X_test_scaled)
y_proba = model.predict_proba(X_test_scaled)[:, 1]  # probability of class 1 (Approved)

# Metrics
accuracy = accuracy_score(y_test, y_pred)
precision = precision_score(y_test, y_pred)
recall = recall_score(y_test, y_pred)
f1 = f1_score(y_test, y_pred)
roc_auc = roc_auc_score(y_test, y_proba)
conf_matrix = confusion_matrix(y_test, y_pred).tolist()

metrics = {
    'accuracy': round(accuracy, 4),
    'precision': round(precision, 4),
    'recall': round(recall, 4),
    'f1': round(f1, 4),
    'roc_auc': round(roc_auc, 4),
    'confusion_matrix': conf_matrix
}

# Save metrics
with open(os.path.join(MODELS_DIR, 'metrics.json'), 'w') as f:
    json.dump(metrics, f, indent=2)

# Extract weights
bias = float(model.intercept_[0])
coefficients = model.coef_[0].tolist()  # list in feature order

weights = {
    'bias': bias,
    'coefficients': dict(zip(feature_names, coefficients))
}
# Also save as list for easy Go loading (optional, but we'll keep both)
weights_list = {
    'bias': bias,
    'coefficients_list': coefficients  # order matches feature_names
}
with open(os.path.join(MODELS_DIR, 'weights.json'), 'w') as f:
    json.dump(weights, f, indent=2)

# Feature importance (absolute value of coefficients)
importance = {name: abs(coef) for name, coef in zip(feature_names, coefficients)}
# Normalize to sum to 1 for better visualization
total = sum(importance.values())
importance_normalized = {k: v/total for k, v in importance.items()}
feature_importance = [{'feature': k, 'importance': v} for k, v in importance_normalized.items()]
feature_importance.sort(key=lambda x: x['importance'], reverse=True)

with open(os.path.join(MODELS_DIR, 'feature_importance.json'), 'w') as f:
    json.dump(feature_importance, f, indent=2)

# Save metadata
metadata = {
    'model_type': 'LogisticRegression',
    'training_date': datetime.now().isoformat(),
    'features': feature_names,
    'num_features': len(feature_names),
    'class_weights': 'balanced',
    'threshold': 0.5,
    'scaler': 'StandardScaler'
}
with open(os.path.join(MODELS_DIR, 'metadata.json'), 'w') as f:
    json.dump(metadata, f, indent=2)

print("Training complete. All model artifacts saved to", MODELS_DIR)
print("Metrics:", metrics)