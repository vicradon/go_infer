# Git Hooks

This project uses custom git hooks that are tracked in the repository.

## Installing Hooks

After cloning the repository, install the hooks:

```bash
make install-hooks
```

Or run the install script directly:

```bash
./scripts/install-hooks.sh
```

## Available Hooks

### pre-push
Automatically creates a new version tag before pushing to main/master.

- Increments the patch version (e.g., `v0.0.1` → `v0.0.2`)
- Creates and pushes the new tag
- Triggers the GitHub Actions release workflow

## Adding New Hooks

1. Create a new hook file in the `hooks/` directory
2. Make it executable: `chmod +x hooks/your-hook`
3. Run `make install-hooks` to install it

## Troubleshooting

If hooks aren't running:
- Ensure hooks are installed: `make install-hooks`
- Check hook permissions: `ls -la .git/hooks/`
- Verify the hook is executable: `chmod +x .git/hooks/your-hook`
