# Initialization

## Initial Configuration

### NAS Service

After installing the NAS service:

1. Register the device with the central control server; after registration, the NAS will obtain a `NAT.ID`.
2. Enable "Connection Password"; the system will generate a secure random password for you.
3. Enable "Auto Start".

### Obsidian Plugin

After installing the Obsidian plugin, in the settings:

1. Enter the "Unique Identifier," which is the `NAT.ID` obtained when registering the NAS with the central control server.
2. Enter the "Connection Password."

## Vault Configuration (Important)

Please follow the appropriate section based on your situation:

- If you have never used this service:
  - I store my `Obsidian Vault` on the `NAS`.
    - Follow the G4 Plan.
  - I only have an `Obsidian Vault` locally.
    - Follow the G1 Plan.
- If you have used any component of this service:
  - I already have a `Vault` on the `NAS Server`.
    - Follow the G2 Plan.
  - I already have a synchronized `Vault` in my local `Obsidian`.
    - Follow the G3 Plan.

### G1 Plan

No additional actions are required. Complete the connection as described earlier and use it normally.

### G2 Plan

Since you already have a `NAS Vault`, you need to confirm:

- Your local `Obsidian` is a brand new `Vault`.
- Files and data in your local `Obsidian Vault` can be overwritten.

> Q: What if I have a local `Obsidian Vault` with important data that cannot be overwritten? <br/>
> A: Create a new `Obsidian Vault`.

After confirming, complete the connection as described earlier. Upon successful connection, open the command palette in `Obsidian` and select `BetaX NAS Sync: Manual Update`.

### G3 Plan

If you already have a local `Obsidian Vault` and do not have a `NAS Vault`, you can proceed to the [G1 Plan](#G1 Plan) to push data to the new `NAS Server`.

If you have a `NAS Server` with `Vault` data, determine if the local `Vault` and `NAS Vault` are consistent. If they are consistent, you can proceed to the [G1 Plan](#G1 Plan). Both will synchronize differences based on timestamp.

If they are inconsistent, you need to clear the end that is not needed, preferably by deleting the entire folder and recreating it.

After clearing, complete the connection as described earlier. Upon successful connection, open the command palette in `Obsidian` and select `BetaX NAS Sync: Manual Update`.

### G4 Plan

It seems you previously used `NAS` for simple file sharing to store your `Obsidian Vault`.

First, deploy the `NAS Server` and map your `Obsidian Vault` on the `NAS` to the container's `/app/vault` directory.

Next, create a `.synclog` file in the `Obsidian Vault` directory; this file has no extension.

After creating the file, open it with a notepad or another editor and enter a 10-digit timestamp in seconds, then save the changes.

After completing these steps, follow the initial connection instructions. Upon successful connection, open the command palette in `Obsidian` and select `BetaX NAS Sync: Manual Update` to automatically pull the `Obsidian Vault` from the `NAS` to the local machine.
