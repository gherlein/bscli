module bscli/test

go 1.21

// This module is for integration tests only
// It doesn't import the main bscli module to avoid circular dependencies
// Instead, it runs the bscli binary as an external command