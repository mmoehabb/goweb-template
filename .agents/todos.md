# Todo List

Use this file to track tasks across sessions. Update it at the start of each session and after completing tasks.

## Format

- `pending` - Not started
- `in_progress` - Currently working on
- `completed` - Finished successfully
- `cancelled` - No longer needed

## Example Tasks (Template Reference)

### Adding a New User Profile Feature

- [ ] Create handler types: `handlers/user/types.go` (extend Credentials)
- [ ] Create profile handler: `handlers/user/profile.go`
- [ ] Create profile form template: `ui/forms/profile_templ.go`
- [ ] Create profile page template: `pages/user/profile_templ.go`
- [ ] Add profile GET route in `main.go`
- [ ] Add profile POST route in `main.go`
- [ ] Create profile DB model: `db/users/profile.go`
- [ ] Add profile query functions: `db/users/queries.go`
- [ ] Add seed for profile table: `db/db.go`

### Adding a New Feature (Generic)

- [ ] Create handler package: `handlers/<feature>/`
- [ ] Create handler file: `handlers/<feature>/handler.go`
- [ ] Create types file: `handlers/<feature>/types.go`
- [ ] Create validators file: `handlers/<feature>/validators.go` (optional)
- [ ] Create form template: `ui/forms/<form>_templ.go`
- [ ] Create page template: `pages/<page>_templ.go`
- [ ] Add route in `main.go`
- [ ] Create DB package: `db/<feature>/`
- [ ] Create DB model: `db/<feature>/model.go`
- [ ] Create DB queries: `db/<feature>/queries.go`
- [ ] Add seed to `db/db.go` (if new table)
- [ ] Run `templ generate`
- [ ] Test the feature
- [ ] Run lint/typecheck

---

## Active Tasks

<!-- Add current tasks here -->

## Completed Tasks

<!-- Record completed work here -->