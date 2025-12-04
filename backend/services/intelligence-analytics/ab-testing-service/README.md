# A/B Testing Service

Experimentation framework for feature testing and optimization.

## Purpose
Manages A/B tests, assigns users to variants, and tracks experiment results.

## Technology Stack
- **Database**: PostgreSQL (experiment config, assignments)
- **API**: gRPC

## Key Features
- ✅ Create and manage experiments
- ✅ Multi-variant testing (A/B/C/D)
- ✅ Traffic splitting
- ✅ Consistent user assignment
- ✅ Statistical significance testing
- ✅ Experiment results tracking
- ✅ Feature flags integration

## API
- `CreateExperiment`: Setup new test
- `AssignVariant`: Get user's variant
- `GetExperimentResults`: View test results
