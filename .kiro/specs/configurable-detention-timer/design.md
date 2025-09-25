# Design Document

## Overview

This feature will modify the existing detention email system to accept a configurable cron schedule through the `DETENTION_EMAIL_SCHEDULE` environment variable. The implementation will maintain backward compatibility by using the current hardcoded schedule as a fallback when the environment variable is not set or contains an invalid value.

## Architecture

The solution will be implemented entirely within the existing `main.go` file where the cron job is currently configured. The design follows the existing patterns in the codebase for environment variable handling and error management.

### Current Implementation
- Hardcoded cron schedule: `"0 12 * * 1-5"` (12:00 PM, Monday-Friday)
- Uses PocketBase's cron package for scheduling
- Calls `tasks.HandleDetentionReportSend(app)` when triggered

### Proposed Changes
- Read `DETENTION_EMAIL_SCHEDULE` environment variable at startup
- Validate the cron expression format
- Use the configured schedule or fall back to default
- Add appropriate logging for configuration and error states

## Components and Interfaces

### Environment Variable
- **Name**: `DETENTION_EMAIL_SCHEDULE`
- **Format**: Standard 5-field cron expression (minute hour day month weekday)
- **Example**: `"0 12 * * 1-5"` (current default)
- **Default**: `"0 12 * * 1-5"` when not set or invalid

### Modified Functions
- **Location**: `pocketbase/main.go` in the `OnBeforeServe` handler
- **Changes**: 
  - Add environment variable reading logic
  - Add cron expression validation
  - Add logging for configuration status
  - Maintain existing error handling for scheduler.Add()

### Validation Logic
The cron package's `scheduler.Add()` method already provides validation - if an invalid cron expression is passed, it returns an error. This existing validation will be leveraged rather than implementing custom validation.

## Data Models

No new data models are required. The feature only involves:
- Reading a string environment variable
- Passing the string to the existing cron scheduler
- Using existing error handling patterns

## Error Handling

### Invalid Cron Expression
- **Detection**: The `scheduler.Add()` method will return an error for invalid expressions
- **Response**: Log the error with details about the invalid expression
- **Fallback**: Retry with the default schedule `"0 12 * * 1-5"`
- **Logging**: Clear error message indicating fallback to default

### Missing Environment Variable
- **Detection**: `os.Getenv()` returns empty string when variable is not set
- **Response**: Use default schedule without logging an error (normal operation)
- **Logging**: Optional info message indicating use of default schedule

### Scheduler Errors
- **Existing Behavior**: Current code already handles scheduler.Add() errors by returning them
- **Enhanced Behavior**: Add more descriptive error messages that include the attempted cron expression

## Testing Strategy

### Manual Testing
1. **Default Behavior**: Start application without `DETENTION_EMAIL_SCHEDULE` set
   - Verify default schedule is used
   - Verify no error messages are logged

2. **Valid Custom Schedule**: Set `DETENTION_EMAIL_SCHEDULE="0 9 * * 1-5"`
   - Verify custom schedule is used
   - Verify confirmation message is logged

3. **Invalid Schedule**: Set `DETENTION_EMAIL_SCHEDULE="invalid"`
   - Verify error is logged
   - Verify fallback to default schedule occurs
   - Verify application continues running

4. **Edge Cases**: Test various cron expressions
   - Empty string: `DETENTION_EMAIL_SCHEDULE=""`
   - Malformed expressions: `DETENTION_EMAIL_SCHEDULE="0 25 * * *"` (invalid hour)
   - Different valid schedules: `DETENTION_EMAIL_SCHEDULE="30 14 * * 2,4"` (2:30 PM Tue/Thu)

### Integration Testing
- Verify the detention email functionality continues to work with custom schedules
- Confirm that changing the schedule doesn't affect the email content or recipient logic
- Test that the cron job actually fires at the configured times

### Logging Verification
- Confirm appropriate log messages are generated for each scenario
- Verify log messages are clear and actionable for administrators
- Ensure no sensitive information is logged

## Implementation Notes

### Code Location
The changes will be made in the `app.OnBeforeServe().Add()` function in `pocketbase/main.go`, specifically around lines 61-70 where the current cron job is configured.

### Backward Compatibility
- Existing deployments without the environment variable will continue to work unchanged
- The default schedule remains exactly the same as the current hardcoded value
- No changes to the detention report functionality itself

### Environment Variable Naming
The chosen name `DETENTION_EMAIL_SCHEDULE` follows the existing pattern in the codebase (e.g., `MANAGEBAC_API`) and clearly indicates its purpose.

### Logging Strategy
- Use `fmt.Printf()` for informational messages (consistent with existing code)
- Use `fmt.Errorf()` for error conditions (consistent with existing error handling)
- Include the actual cron expression in log messages for debugging