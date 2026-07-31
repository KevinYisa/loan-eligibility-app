# scripts/preprocess.py

import pandas as pd
import numpy as np
from sklearn.model_selection import train_test_split
import json
import os

# ------------------ 1. Load Data ------------------
df = pd.read_csv('../datasets/loan_approval_dataset.csv')
print(f"Original shape: {df.shape}")

# ------------------ 2. Drop Unnecessary Columns ------------------
# loan_id is unique and useless for prediction

# Remove leading/trailing whitespace from column names
df.columns = df.columns.str.strip()
df.drop('loan_id', axis=1, inplace=True)

for col in df.select_dtypes(include='object'):
    df[col] = df[col].str.strip()
    
# ------------------ 3. Handle Negative Asset Values ------------------
# Some values are negative (e.g., -100000). This is likely a data error.
# Replace negative values with 0 (assuming they mean no asset/debt).
asset_cols = ['residential_assets_value', 'commercial_assets_value', 
              'luxury_assets_value', 'bank_asset_value']

for col in asset_cols:
    df[col] = df[col].apply(lambda x: max(0, x))  # set negatives to 0

# ------------------ 4. Feature Engineering: Aggregate Assets ------------------
# We only want to ask the user for ONE asset value (Total Assets)
df['total_assets'] = df[asset_cols].sum(axis=1)

# Drop the individual asset columns
df.drop(columns=asset_cols, inplace=True)

# ------------------ 5. Encode Categorical Variables ------------------
# Education: Graduate -> 1, Not Graduate -> 0
df['education'] = df['education'].map({'Graduate': 1, 'Not Graduate': 0})

# Self Employed: Yes -> 1, No -> 0
df['self_employed'] = df['self_employed'].map({'Yes': 1, 'No': 0})

# Target: Approved -> 1, Rejected -> 0
df['loan_status'] = df['loan_status'].map({'Approved': 1, 'Rejected': 0})

# ------------------ 6. Check for Missing Values ------------------
print("\nMissing values per column:")
print(df.isnull().sum())

# ------------------ 7. Define Final Feature Set (Matches Frontend) ------------------
feature_cols = [
    'no_of_dependents', 
    'education', 
    'self_employed', 
    'income_annum', 
    'loan_amount', 
    'loan_term', 
    'cibil_score', 
    'total_assets'
]

X = df[feature_cols]
y = df['loan_status']

print(f"\nFinal Features ({len(feature_cols)}): {feature_cols}")

# ------------------ 8. Save Cleaned Full Dataset ------------------
# Ensure processed directory exists
os.makedirs('../processed', exist_ok=True)
df.to_csv('../processed/cleaned_data.csv', index=False)

# ------------------ 9. Train / Validation / Test Split ------------------
# Split 70% train, 30% temp
X_train, X_temp, y_train, y_temp = train_test_split(
    X, y, test_size=0.3, random_state=42, stratify=y
)

# Split temp into 50% validation, 50% test (i.e., 15% each of total)
X_val, X_test, y_val, y_test = train_test_split(
    X_temp, y_temp, test_size=0.5, random_state=42, stratify=y_temp
)

print(f"\nTrain size: {len(X_train)}, Val size: {len(X_val)}, Test size: {len(X_test)}")

# ------------------ 10. Save Splits ------------------
X_train.to_csv('../processed/X_train.csv', index=False)
X_val.to_csv('../processed/X_val.csv', index=False)
X_test.to_csv('../processed/X_test.csv', index=False)

y_train.to_csv('../processed/y_train.csv', index=False)
y_val.to_csv('../processed/y_val.csv', index=False)
y_test.to_csv('../processed/y_test.csv', index=False)

# ------------------ 11. Export Metadata for Go Backend ------------------
os.makedirs('../models', exist_ok=True)
with open('../models/feature_names.json', 'w') as f:
    json.dump(feature_cols, f)

print("\nPreprocessing complete! All files saved in ../processed/ and ../models/")